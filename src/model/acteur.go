/*
*****************************************************************************

	Acteurs

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2019-11-08 10:43:40+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// valeurs du champ id
// déterminés à la création de la base
// voir manage/db-install/acteur.go
const ID_SCTL = 1
const ID_BDL = 2
const ID_GFA = 3

type Acteur struct {
	Id           int
	Nom          string
	CodesRole    []string
	Prenom       string
	Adresse1     string
	Adresse2     string
	Cp           string
	Ville        string
	Tel          string
	Mobile       string
	Email        string
	Bic          string
	Iban         string
	Siret        string
	Proprietaire bool
	Fournisseur  bool
	Actif        bool
	Notes        string
	// pas stocké en base
	Deletable bool
	Parcelles []*Parcelle
}

// Sert à afficher la liste des activités d'un acteur (acteur-show.html).
// Contient les infos utilisées pour l'affichage, pas les activités.
// Distinct de model.Activite car prend en compte tous les rôles possibles des acteurs
// Pourrait être supprimé et remplacé par model.Activite
// Mais obligerait à rendre model.Activite plus complexe => laissé en l'état pour l'instant.
type ActeurActivite struct {
	Date        time.Time
	Role        string
	URL         string    // URL de la page de l'activité concernée
	NomActivite string
	Quantite    float64
	Unite       string    // pour Quantite
}

// ************************** Structure *******************************

// IsDeletable indique si un acteur peut être supprimé, ou s'il doit être marqué comme inactif
// cf règles de gestion dans cahier des charges
func (a *Acteur) IsDeletable(db *sqlx.DB) (res bool, err error) {
    queries := []string{
        // mis en premier car le plus fréquent
        "select count(*) from chautre where id_acheteur=$1",
        "select count(*) from venteplaq where id_client=$1",
        //
        "select count(*) from plaqop where id_acteur=$1",
        //
        "select count(*) from plaqtrans where id_transporteur=$1",
        "select count(*) from plaqtrans where id_conducteur=$1",
        "select count(*) from plaqtrans where id_proprioutil=$1",
        //
        "select count(*) from plaqrange where id_rangeur=$1",
        "select count(*) from plaqrange where id_conducteur=$1",
        "select count(*) from plaqrange where id_proprioutil=$1",
        //
        "select count(*) from ventelivre where id_livreur=$1",
        "select count(*) from ventelivre where id_conducteur=$1",
        "select count(*) from ventelivre where id_proprioutil=$1",
        //
        "select count(*) from ventecharge where id_chargeur=$1",
        "select count(*) from ventecharge where id_conducteur=$1",
        "select count(*) from ventecharge where id_proprioutil=$1",
        //
        "select count(*) from humid_acteur where id_acteur=$1",
    }
    var count int
    for _, query := range(queries){
        err = db.QueryRow(query, a.Id).Scan(&count)
        if err != nil {
            return false, werr.Wrapf(err, fmt.Sprint("Erreur query: %s\n Acteur %d: %s", query, a.Id, a.String()))
        }
        if count != 0 {
            return false, nil
        }
    }
    return true, nil
}

// ************************** Nom *******************************

func (a *Acteur) String() string {
	return strings.TrimSpace(a.Prenom + " " + a.Nom)
}

// ************************** Divers *******************************

func CountActeurs(db *sqlx.DB) (count int) {
	_ = db.QueryRow("select count(*) from acteur").Scan(&count)
	return count
}

// ************************** Get one *******************************

// Renvoie un Acteur à partir de son id.
// Ne contient que les champs de la table acteur.
// Les autres champs ne sont pas remplis.
func GetActeur(db *sqlx.DB, id int) (a *Acteur, err error) {
	a = &Acteur{}
	query := "select * from acteur where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(a)
	if err != nil {
		return a, werr.Wrapf(err, "Erreur query : "+query)
	}
	return a, err
}

/*
Renvoie un acteur avec les codes de ses rôles
*/
func GetActeurFull(db *sqlx.DB, id int) (a *Acteur, err error) {
	a, err = GetActeur(db, id)
	if err != nil {
		return nil, werr.Wrapf(err, "Erreur Appel GetActeur()")
	}
	a.ComputeCodesRole(db)
	return a, nil
}

// ************************** Get many *******************************

// GetSortedActeurs renvoie une liste d'Acteurs triés en utilisant un champ de la table.
// TODO En fait, toujours utilisée en triant par nom, on pourrait supprimer le param field
// @param field    Champ de la table acteur utilisé pour le tri
func GetSortedActeurs(db *sqlx.DB, field string) (acteurs []*Acteur, err error) {
	acteurs = []*Acteur{}
	query := "select * from acteur where id<>0 order by " + field
	err = db.Select(&acteurs, query)
	if err != nil {
		return acteurs, werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, a := range acteurs {
		err = a.ComputeCodesRole(db)
		if err != nil {
			return acteurs, werr.Wrapf(err, "Erreur appel ComputeCodesRole()")
		}
		a.Deletable, err = a.IsDeletable(db)
		if err != nil {
			return acteurs, werr.Wrapf(err, "Erreur appel Acteur.IsDeletable()")
		}
	}
	return acteurs, nil
}

/*
Renvoie une liste d'Acteurs ayant un rôle donné.
Les acteurs ne contiennent que les champs de la table.
Les acteurs sont triés par nom.
@param code_role    Code d'un rôle utilisateur

////////////// pas encore utilisée //////////////////
*/
func GetActeursByRole(db *sqlx.DB, code_role string) (acteurs []*Acteur, err error) {
	acteurs = []*Acteur{}
	query := `select * from acteur id in(select id_acteur from acteur_role where code_role=$1) order by nom`
	err = db.Select(&acteurs, query, code_role)
	if err != nil {
		return acteurs, werr.Wrapf(err, "Erreur query : "+query)
	}
	return acteurs, nil
}

/*
Renvoie une liste d'Acteurs ayant comme rôle client PF ou client d'un chantier autres valorisations.
Les acteurs ne contiennent que les champs de la table.
Les acteurs sont triés par nom.
*/
func GetClients(db *sqlx.DB) (acteurs []*Acteur, err error) {
	acteurs = []*Acteur{}
	query := `select * from acteur where id in(
	    select id_acteur from acteur_role where code_role in(
	        'AVC-CH', 'AVC-BO', 'AVC-PL', 'AVC-PP', 'VPL-CL'
	        )
        ) order by nom`
	err = db.Select(&acteurs, query)
	if err != nil {
		return acteurs, werr.Wrapf(err, "Erreur query : "+query)
	}
	return acteurs, nil
}

/*
Renvoie les Acteurs dont le champ Fournisseur = true
( = les fournisseurs de plaquettes ; en pratique, en 2020, 1 seul fournisseur : BDL)
Ne contient que les champs de la table acteur.
Les autres champs ne sont pas remplis.
*/
func GetFournisseurs(db *sqlx.DB) (acteurs []*Acteur, err error) {
	////////////// remplacer par GetActeursByRole() //////////////
	acteurs = []*Acteur{}
	query := "select * from acteur where fournisseur"
	err = db.Select(&acteurs, query)
	if err != nil {
		return acteurs, werr.Wrapf(err, "Erreur query : "+query)
	}
	return acteurs, nil
}

/*
Renvoie les Acteurs ayant participé à une vente plaquettes en tant que client
Ne contient que les champs de la table acteur.
Les autres champs ne sont pas remplis.
////////////// remplacer par GetSortedActeursByRole() //////////////
*/
func GetClientsPlaquettes(db *sqlx.DB) (acteurs []*Acteur, err error) {
	acteurs = []*Acteur{}
	query := `select * from acteur where id in(
                select id_client from venteplaq
	        ) order by nom,prenom`
	err = db.Select(&acteurs, query)
	if err != nil {
		return acteurs, werr.Wrapf(err, "Erreur query : "+query)
	}
	return acteurs, nil
}

/*
Utilisé pour construire html datalist
*/
func GetListeActeurs(db *sqlx.DB) (res map[int]string, err error) {
	res = map[int]string{}
	acteurs := []*Acteur{}
	query := "select id,prenom,nom from acteur"
	err = db.Select(&acteurs, query)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, a := range acteurs {
		res[a.Id] = a.String()
	}
	return res, nil
}

/*
Renvoie les acteurs SCTL et GFA, marqué comme propriétaires
////////////// remplacer par GetSortedActeursByRole() //////////////
*/
func GetProprietaires(db *sqlx.DB) (res map[int]string, err error) {
	res = map[int]string{}
	acteurs := []*Acteur{}
	query := "select id,nom from acteur where proprietaire=true"
	err = db.Select(&acteurs, query)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, a := range acteurs {
		res[a.Id] = a.String()
	}
	return res, nil
}

// ************************** Compute *******************************

func (a *Acteur) ComputeCodesRole(db *sqlx.DB) (err error) {
	if len(a.CodesRole) != 0 {
		return nil // déjà calculé
	}
	query := `select code_role from acteur_role where id_acteur=$1`
	err = db.Select(&a.CodesRole, query, a.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	return nil

}

// ************************** Get activité *******************************

/*
Renvoie les activités auxquelles un acteur a participé.
Ordre chronologique inverse
Ne renvoie que des infos pour afficher la liste, pas les activités réelles.
*/
func (a *Acteur) GetActivitesByDate(db *sqlx.DB) (res []*ActeurActivite, err error) {
	res = []*ActeurActivite{}
	var query string
	//
	// Opérations simples pour chantiers plaquettes
	//
	list1 := []PlaqOp{}
	query = "select * from plaqop where id_acteur=$1"
	err = db.Select(&list1, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list1 {
		plaq, err := GetPlaq(db, elt.IdChantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaq()")
		}
		err = plaq.ComputeLieudits(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
		}
		new := &ActeurActivite{
			Date:        elt.DateDebut,
			Role:        elt.RoleName(),
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.IdChantier) + "/chantiers",
			NomActivite: plaq.FullString(),
            Quantite:    elt.Qte,
            Unite:       elt.Unite,
		}
		res = append(res, new)
	}
	//
	// Transport plateforme de chantier plaquette à lieu de stockage - transporteur (coût global)
	//
	list2 := []PlaqTrans{}
	query = "select * from plaqtrans where id_transporteur=$1"
	err = db.Select(&list2, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list2 {
		plaq, err := GetPlaq(db, elt.IdChantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaq()")
		}
		err = plaq.ComputeLieudits(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
		}
		new := &ActeurActivite{
			Date:        elt.DateTrans,
			Role:        "transporteur",
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.IdChantier) + "/chantiers",
			NomActivite: plaq.FullString(),
            Quantite:    elt.Qte,
            Unite:       "MA",
		}
		res = append(res, new)
	}
	//
	// Transport plateforme de chantier plaquette à lieu de stockage - conducteur
	//
	list2a := []PlaqTrans{}
	query = "select * from plaqtrans where id_conducteur=$1"
	err = db.Select(&list2a, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list2a {
		plaq, err := GetPlaq(db, elt.IdChantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaq()")
		}
		err = plaq.ComputeLieudits(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
		}
		new := &ActeurActivite{
			Date:        elt.DateTrans,
			Role:        "conducteur (transport)",
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.IdChantier) + "/chantiers",
			NomActivite: plaq.FullString(),
            Quantite:    elt.Qte,
            Unite:       "MA",
		}
		res = append(res, new)
	}
	//
	// Transport plateforme de chantier plaquette à lieu de stockage - propriétaire outil
	//
	list2b := []PlaqTrans{}
	query = "select * from plaqtrans where id_proprioutil=$1"
	err = db.Select(&list2b, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list2b {
		plaq, err := GetPlaq(db, elt.IdChantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaq()")
		}
		err = plaq.ComputeLieudits(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
		}
		new := &ActeurActivite{
			Date:        elt.DateTrans,
			Role:        "propriétaire outil (transport)",
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.IdChantier) + "/chantiers",
			NomActivite: plaq.FullString(),
            Quantite:    elt.Qte,
            Unite:       "MA",
		}
		res = append(res, new)
	}
	//
	// Rangement plaquettes suite au transport - rangeur (coût global)
	//
	list3 := []PlaqRange{}
	query = "select * from plaqrange where id_rangeur=$1"
	err = db.Select(&list3, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list3 {
		plaq, err := GetPlaq(db, elt.IdChantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaq()")
		}
		err = plaq.ComputeLieudits(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
		}
		new := &ActeurActivite{
			Date:        elt.DateRange,
			Role:        "rangeur",
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.IdChantier) + "/chantiers",
			NomActivite: plaq.FullString(),
            // Quantite:    ,
            // Unite:       ,
		}
		res = append(res, new)
	}
	//
	// Rangement plaquettes suite au transport - conducteur
	//
	list3a := []PlaqRange{}
	query = "select * from plaqrange where id_conducteur=$1"
	err = db.Select(&list3a, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list3a {
		plaq, err := GetPlaq(db, elt.IdChantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaq()")
		}
		err = plaq.ComputeLieudits(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
		}
		new := &ActeurActivite{
			Date:        elt.DateRange,
			Role:        "conducteur (rangement)",
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.IdChantier) + "/chantiers",
			NomActivite: plaq.FullString(),
            // Quantite:    ,
            // Unite:       ,
		}
		res = append(res, new)
	}
	//
	// Rangement plaquettes suite au transport - propriétaire outil
	//
	list3b := []PlaqRange{}
	query = "select * from plaqrange where id_proprioutil=$1"
	err = db.Select(&list3b, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list3b {
		plaq, err := GetPlaq(db, elt.IdChantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaq()")
		}
		err = plaq.ComputeLieudits(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
		}
		new := &ActeurActivite{
			Date:        elt.DateRange,
			Role:        "propriétaire outil (rangement)",
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.IdChantier) + "/chantiers",
			NomActivite: plaq.FullString(),
            // Quantite:    ,
            // Unite:       ,
		}
		res = append(res, new)
	}
	//
	// Livraison pour vente plaquette - livreur (coût global)
	//
	list4 := []VenteLivre{}
	query = "select * from ventelivre where id_livreur=$1"
	err = db.Select(&list4, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list4 {
		vp, err := GetVentePlaq(db, elt.IdVente)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetVentePlaq()")
		}
		err = vp.ComputeClient(db) // besoin de client pour calculer vp.FullString()
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeClient()")
		}
		err = elt.ComputeChargements(db) // besoin de chargements pour calculer quantité
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeChargements()")
		}
		new := &ActeurActivite{
			Date:        elt.DateLivre,
			Role:        "livreur",
			URL:         "/vente/" + strconv.Itoa(elt.IdVente),
			NomActivite: vp.FullString(),
            Quantite:    elt.Qte,
            Unite:       "MA",
		}
		res = append(res, new)
	}
	//
	// Livraison pour vente plaquette - conducteur
	//
	list4a := []VenteLivre{}
	query = "select * from ventelivre where id_conducteur=$1"
	err = db.Select(&list4a, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list4a {
		vp, err := GetVentePlaq(db, elt.IdVente)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetVentePlaq()")
		}
		err = vp.ComputeClient(db) // besoin de client pour calculer vp.FullString()
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeClient()")
		}
		err = elt.ComputeChargements(db) // besoin de chargements pour calculer quantité
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeChargements()")
		}
		new := &ActeurActivite{
			Date:        elt.DateLivre,
			Role:        "conducteur (livraison)",
			URL:         "/vente/" + strconv.Itoa(elt.IdVente),
			NomActivite: vp.FullString(),
            Quantite:    elt.Qte,
            Unite:       "MA",
		}
		res = append(res, new)
	}
	//
	// Livraison pour vente plaquette - propriétaire outil
	//
	list4b := []VenteLivre{}
	query = "select * from ventelivre where id_proprioutil=$1"
	err = db.Select(&list4b, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list4b {
		vp, err := GetVentePlaq(db, elt.IdVente)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetVentePlaq()")
		}
		err = vp.ComputeClient(db) // besoin de client pour calculer vp.FullString()
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeClient()")
		}
		err = elt.ComputeChargements(db) // besoin de chargements pour calculer vp.Quantité
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeChargements()")
		}
		new := &ActeurActivite{
			Date:        elt.DateLivre,
			Role:        "propriétaire outil (livraison)",
			URL:         "/vente/" + strconv.Itoa(elt.IdVente),
			NomActivite: vp.FullString(),
            Quantite:    elt.Qte,
            Unite:       "MA",
		}
		res = append(res, new)
	}
	//
	// Chargement pour livraison plaquette - chargeur (coût global)
	//
	list5 := []VenteCharge{}
	query = "select * from ventecharge where id_chargeur=$1"
	err = db.Select(&list5, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list5 {
		err = elt.ComputeIdVente(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeIdVente()")
		}
		// besoin de vente et de son client pour NomActivite
		vp, err := GetVentePlaq(db, elt.IdVente)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetVentePlaq()")
		}
		err = vp.ComputeClient(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel vente.ComputeClient()")
		}
		new := &ActeurActivite{
			Date:        elt.DateCharge,
			Role:        "chargeur",
			URL:         "/vente/" + strconv.Itoa(elt.IdVente),
			NomActivite: vp.FullString(),
            Quantite:    elt.Qte,
            Unite:       "MA",
		}
		res = append(res, new)
	}
	//
	// Chargement pour livraison plaquette - conducteur
	//
	list5a := []VenteCharge{}
	query = "select * from ventecharge where id_conducteur=$1"
	err = db.Select(&list5a, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list5a {
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeProprioutil()")
		}
		err = elt.ComputeIdVente(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeIdVente()")
		}
		// besoin de vente et de son client pour NomActivite
		vp, err := GetVentePlaq(db, elt.IdVente)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetVentePlaq()")
		}
		err = vp.ComputeClient(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel vente.ComputeClient()")
		}
		new := &ActeurActivite{
			Date:        elt.DateCharge,
			Role:        "conducteur (chargement)",
			URL:         "/vente/" + strconv.Itoa(elt.IdVente),
			NomActivite: vp.FullString(),
            Quantite:    elt.Qte,
            Unite:       "MA",
		}
		res = append(res, new)
	}
	//
	// Chargement pour livraison plaquette - propriétaire outil
	//
	list5b := []VenteCharge{}
	query = "select * from ventecharge where id_proprioutil=$1"
	err = db.Select(&list5b, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list5b {
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeConducteur()")
		}
		err = elt.ComputeIdVente(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeIdVente()")
		}
		// besoin de vente et de son client pour NomActivite
		vp, err := GetVentePlaq(db, elt.IdVente)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetVentePlaq()")
		}
		err = vp.ComputeClient(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel vente.ComputeClient()")
		}
		new := &ActeurActivite{
			Date:        elt.DateCharge,
			Role:        "propriétaire outil (chargement)",
			URL:         "/vente/" + strconv.Itoa(elt.IdVente),
			NomActivite: vp.FullString(),
            Quantite:    elt.Qte,
            Unite:       "MA",
		}
		res = append(res, new)
	}
	//
	// Client plaquette
	//
	list6 := []VentePlaq{}
	query = "select * from venteplaq where id_client=$1"
	err = db.Select(&list6, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list6 {
		elt.Client, err = GetActeur(db, elt.IdClient) // besoin pour appel elt.FullString()
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetActeur()")
		}
        err = elt.ComputeQte(db)
        if err != nil {
            return res, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeQte()")
        }
		new := &ActeurActivite{
			Date:        elt.DateVente,
			Role:        "client plaquettes",
			URL:         "/vente/" + strconv.Itoa(elt.Id),
			NomActivite: elt.FullString(),
            Quantite:    elt.Qte,
            Unite:       "MA",
		}
		res = append(res, new)
	}
	//
	// Chantier autres valorisations - acheteur
	//
	list8 := []Chautre{}
	query = "select * from chautre where id_acheteur=$1"
	err = db.Select(&list8, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list8 {
		elt.Acheteur = a
		err = elt.ComputeLieudits(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Chautre.ComputeLieudit()")
		}
		new := &ActeurActivite{
			Date:        elt.DateContrat,
			Role:        "acheteur chantier autres valorisations",
			URL:         "/chantier/autre/liste/" + strconv.Itoa(elt.DateContrat.Year()),
			NomActivite: elt.FullString(),
            Quantite:    elt.VolumeRealise,
            Unite:       elt.Unite,
		}
		res = append(res, new)
	}
	//
	// mesures d'humidité - mesureur
	//
	list10 := []Humid{}
	query = "select * from humid where id in(select id_humid from humid_acteur where id_acteur=$1)"
	err = db.Select(&list10, query, a.Id)
	for _, elt := range list10 {
		new := &ActeurActivite{
			Date:        elt.DateMesure,
			Role:        "mesureur",
			URL:         "/humidite/liste/" + strconv.Itoa(elt.DateMesure.Year()),
			NomActivite: "Mesure humidité",
            // Quantite:    ,
            // Unite:       ,
		}
		res = append(res, new)
	}
	// Tri par date
	sortedRes := make(acteurActiviteSlice, 0, len(res))
	for _, elt := range res {
		sortedRes = append(sortedRes, elt)
	}
	sort.Sort(sortedRes)
	//
	return sortedRes, nil
}

// Auxiliaires de GetActivitesByDate() pour trier par date
type acteurActiviteSlice []*ActeurActivite

func (p acteurActiviteSlice) Len() int           { return len(p) }
func (p acteurActiviteSlice) Less(i, j int) bool { return p[i].Date.After(p[j].Date) }
func (p acteurActiviteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// ************************** CRUD *******************************

func InsertActeur(db *sqlx.DB, acteur *Acteur) (id int, err error) {
	query := `insert into acteur(
        nom,
        prenom,
        adresse1,
        adresse2,
        cp,
        ville,
        tel,
        mobile,
        email,
        bic,
        iban,
        siret,
        proprietaire,
        fournisseur,
        actif,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) returning id`
	err = db.QueryRow(
		query,
		acteur.Nom,
		acteur.Prenom,
		acteur.Adresse1,
		acteur.Adresse2,
		acteur.Cp,
		acteur.Ville,
		acteur.Tel,
		acteur.Mobile,
		acteur.Email,
		acteur.Bic,
		acteur.Iban,
		acteur.Siret,
		acteur.Proprietaire,
		acteur.Fournisseur,
		acteur.Actif,
		acteur.Notes).Scan(&id)
	if err != nil {
		return 0, werr.Wrapf(err, "Erreur query : "+query)
	}
	err = insertLiensActeurRole(db, id, acteur.CodesRole)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur appel insertLiensActeurRole()")
	}
	return id, nil
}

func UpdateActeur(db *sqlx.DB, acteur *Acteur) (err error) {
	query := `update acteur set(
        nom,
        prenom,
        adresse1,
        adresse2,
        cp,
        ville,
        tel,
        mobile,
        email,
        bic,
        iban,
        siret,
        proprietaire,
        fournisseur,
        actif,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) where id=$17`
	_, err = db.Exec(
		query,
		acteur.Nom,
		acteur.Prenom,
		acteur.Adresse1,
		acteur.Adresse2,
		acteur.Cp,
		acteur.Ville,
		acteur.Tel,
		acteur.Mobile,
		acteur.Email,
		acteur.Bic,
		acteur.Iban,
		acteur.Siret,
		acteur.Proprietaire,
		acteur.Fournisseur,
		acteur.Actif,
		acteur.Notes,
		acteur.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	err = updateLiensActeurRole(db, acteur.Id, acteur.CodesRole)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel updateLiensActeurRole()")
	}
	return nil
}

func DeleteActeur(db *sqlx.DB, id int) (err error) {
	// peut-être ici protection pour savoir si Deletable = true
	// (la situation actuelle fait confiance à l'UI pour ne pas proposer delete sur acteur non deletable)
	query := "delete from acteur where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	err = deleteLiensActeurRole(db, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensActeurRole()")
	}
	return nil
}

//
// Fonctions auxiliares de InsertActeur(), UpdateActeur() et DeleteActeur()
//

func insertLiensActeurRole(db *sqlx.DB, idActeur int, codesRoles []string) (err error) {
	query := "insert into acteur_role values($1,$2)"
	for _, code := range codesRoles {
		_, err = db.Exec(
			query,
			idActeur,
			code)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	return nil
}

func deleteLiensActeurRole(db *sqlx.DB, idActeur int) (err error) {
	query := "delete from acteur_role where id_acteur=$1"
	_, err = db.Exec(query, idActeur)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func updateLiensActeurRole(db *sqlx.DB, idActeur int, codesRoles []string) (err error) {
	err = deleteLiensActeurRole(db, idActeur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensActeurRole() à partir de updateLiensActeurRole()")
	}
	//
	err = insertLiensActeurRole(db, idActeur, codesRoles)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel insertLiensActeurRole() à partir de updateLiensActeurRole()")
	}
	//
	return nil
}
