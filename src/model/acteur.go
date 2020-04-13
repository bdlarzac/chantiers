/******************************************************************************
    Acteurs

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-11-08 10:43:40+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"sort"
	"strconv"
	"time"

	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	//"fmt"
)

type Acteur struct {
	Id           int
	IdSctl       int `db:"id_sctl"`
	Nom          string
	Prenom       string
	Adresse1     string
	Adresse2     string
	Cp           string
	Ville        string
	Tel          string
	TelPortable  string
	Email        string
	Bic          string
	Iban         string
	Siret        string
	Proprietaire bool // à supprimer
	Fournisseur  bool
	Actif        bool
	Notes        string
	// pas stocké en base
	Deletable bool
	Parcelles []*Parcelle
}

// Sert à afficher la liste des activités d'un acteur.
// Contient les infos utilisées pour l'affichage, pas les activités.
type ActeurActivite struct {
	Date        time.Time
	Role        string
	URL         string // URL de la page de l'activité concernée
	NomActivite string
}

// ************************** Structure *******************************

func (a *Acteur) IsSctl() bool {
	return a.IdSctl != 0
}

// cf règles de gestion dans cahier des charges
func (a *Acteur) IsDeletable(db *sqlx.DB) (bool, error) {
	if a.IdSctl != 0 {
		return false, nil
	}
	act, err := a.GetActivitesByDate(db)
	if err != nil {
		return false, werr.Wrapf(err, "Erreur appel GetActivitesByDate()")
	}
	// supprimable si associé à aucune activité
	return len(act) == 0, nil
}

// ************************** Nom *******************************

func (a *Acteur) String() string {
	if a.Prenom == "" {
		return a.Nom
	}
	return a.Prenom + " " + a.Nom
}

// Comme autocomplete se fait par ne Nom, besoin de cet affichage
// particulier dans les input avec autocomplete sur le nom
func (a *Acteur) NomAutocomplete() string {
	if a.Prenom == "" {
		return a.Nom
	}
	return a.Nom + " " + a.Prenom
}

// ************************** Divers *******************************

func CountActeurs(db *sqlx.DB) int {
	var count int
	_ = db.QueryRow("select count(*) from acteur").Scan(&count)
	return count
}

// ************************** Get généraux *******************************

// Renvoie un Acteur à partir de son id.
// Ne contient que les champs de la table acteur.
// Les autres champs ne sont pas remplis.
func GetActeur(db *sqlx.DB, id int) (*Acteur, error) {
	a := &Acteur{}
	query := "select * from acteur where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(a)
	if err != nil {
		return a, werr.Wrapf(err, "Erreur query : "+query)
	}
	return a, err
}

// Renvoie une liste d'Acteurs triés en utilisant un champ de la table
// @param field    Champ de la table acteur utilisé pour le tri
func SortedActeurs(db *sqlx.DB, field string) ([]*Acteur, error) {
	acteurs := []*Acteur{}
	query := "select * from acteur order by " + field
	err := db.Select(&acteurs, query)
	if err != nil {
		return acteurs, werr.Wrapf(err, "Erreur query : "+query)
	}
	return acteurs, err
}

// Renvoie un Acteur à partir de son id sctl.
// Ne contient que les champs de la table acteur.
// Les autres champs ne sont pas remplis.
func GetActeurByIdSctl(db *sqlx.DB, id int) (*Acteur, error) {
	a := &Acteur{}
	query := "select * from acteur where id_sctl=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(a)
	if err != nil {
		return a, werr.Wrapf(err, "Erreur query : "+query)
	}
	return a, nil
}

// ************************** Get liés à l'activité *******************************

// Renvoie les Acteurs dont le champ Fournisseur = true
// ( = les fournisseurs de plaquettes ; en pratique, en 2020, 1 seul fournisseur : BDL)
// Ne contient que les champs de la table acteur.
// Les autres champs ne sont pas remplis.
func GetFournisseurs(db *sqlx.DB) ([]*Acteur, error) {
	acteurs := []*Acteur{}
	query := "select * from acteur where fournisseur"
	err := db.Select(&acteurs, query)
	return acteurs, err
}

// Renvoie les activités auxquelles un acteur a participé.
// Ordre chronologique inverse
// Ne renvoie que des infos pour afficher la liste, pas les activités réelles.
// Voir GetAffactureActivitesByDate(), qui renvoie les activités.
func (a *Acteur) GetActivitesByDate(db *sqlx.DB) ([]*ActeurActivite, error) {
	res := []*ActeurActivite{}
	var err error
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
		err = plaq.ComputeLieudit(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudit()")
		}
		new := &ActeurActivite{
			Date:        elt.DateOp,
			Role:        elt.RoleName(),
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.IdChantier) + "/chantiers",
			NomActivite: "Chantier plaquettes " + plaq.String()}
		res = append(res, new)
	}
	//
	// Transport plateforme de chantier plaquette à lieu de stockage - conducteur
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
		err = plaq.ComputeLieudit(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudit()")
		}
		new := &ActeurActivite{
			Date:        elt.DateTrans,
			Role:        "transporteur",
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.IdChantier) + "/chantiers",
			NomActivite: "Chantier plaquettes " + plaq.String()}
		res = append(res, new)
	}
	//
	// Rangement plaquettes suite au transport - conducteur
	//
	list3 := []PlaqRange{}
	query = "select * from plaqrange where id_conducteur=$1"
	err = db.Select(&list3, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list3 {
		plaq, err := GetPlaq(db, elt.IdChantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaq()")
		}
		err = plaq.ComputeLieudit(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudit()")
		}
		new := &ActeurActivite{
			Date:        elt.DateRange,
			Role:        "rangeur",
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.IdChantier) + "/chantiers",
			NomActivite: "Chantier plaquettes " + plaq.String()}
		res = append(res, new)
	}
	//
	// Livraison pour vente plaquette - livreur
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
		err = vp.ComputeClient(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeClient()")
		}
		new := &ActeurActivite{
			Date:        elt.DateLivre,
			Role:        "livreur",
			URL:         "/vente/" + strconv.Itoa(elt.IdVente),
			NomActivite: "Vente plaquettes " + vp.String()}
		res = append(res, new)
	}
	//
	// Chargement pour livraison plaquette - chargeur
	//
	list5 := []VenteCharge{}
	query = "select * from ventecharge where id_chargeur=$1"
	err = db.Select(&list5, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list5 {
		elt.Chargeur = a
		err = elt.ComputeIdVente(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeIdVente()")
		}
		new := &ActeurActivite{
			Date:        elt.DateCharge,
			Role:        "chargeur",
			URL:         "/vente/" + strconv.Itoa(elt.IdVente),
			NomActivite: "Chargement plaquettes " + elt.String()}
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
		elt.Client, err = GetActeur(db, elt.IdClient)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetActeur()")
		}
		new := &ActeurActivite{
			Date:        elt.DateVente,
			Role:        "client plaquettes",
			URL:         "/vente/" + strconv.Itoa(elt.Id),
			NomActivite: "Vente plaquettes " + elt.String()}
		res = append(res, new)
	}
	//
	// Chantier bois sur pied - client
	//
	list7 := []BSPied{}
	query = "select * from bspied where id_acheteur=$1"
	err = db.Select(&list7, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list7 {
		elt.Acheteur = a
		err = elt.ComputeLieudit(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel BSPied.ComputeLieudit()")
		}
		new := &ActeurActivite{
			Date:        elt.DateContrat,
			Role:        "client bois sur pied",
			URL:         "/chantier/bois-sur-pied/" + strconv.Itoa(elt.DateContrat.Year()),
			NomActivite: "Chantier bois sur pied " + elt.String()}
		res = append(res, new)
	}
	//
	// Chantier autres valorisations - client
	//
	list8 := []Chautre{}
	query = "select * from chautre where id_client=$1"
	err = db.Select(&list8, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list8 {
		elt.Client = a
		err = elt.ComputeLieudit(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Chautre.ComputeLieudit()")
		}
		new := &ActeurActivite{
			Date:        elt.DateContrat,
			Role:        "client bois sur pied",
			URL:         "/chantier/autre/liste/" + strconv.Itoa(elt.DateContrat.Year()),
			NomActivite: "Chantier " + elt.String()}
		res = append(res, new)
	}
	//
	// Chantier chauffage fermier - fermier
	//
	list9 := []Chaufer{}
	query = "select * from chaufer where id_fermier=$1"
	err = db.Select(&list9, query, a.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list9 {
		elt.Fermier = a
		new := &ActeurActivite{
			Date:        elt.DateChantier,
			Role:        "fermier",
			URL:         "/chantier/chauffage-fermier/liste/" + strconv.Itoa(elt.DateChantier.Year()),
			NomActivite: "Chantier chauffage fermier " + elt.String()}
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
		}
		res = append(res, new)
	}
	// tri par date
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

func (p acteurActiviteSlice) Len() int {
	return len(p)
}
func (p acteurActiviteSlice) Less(i, j int) bool {
	return p[i].Date.After(p[j].Date)
}
func (p acteurActiviteSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// ************************** Get utilisés par ajax *******************************

// Renvoie des Acteurs (exploitants ou fermiers) à partir d'un lieu-dit.
// Utilise les parcelles pour faire le lien
// Ne contient que les champs de la table acteur.
// Les autres champs ne sont pas remplis.
//
// ATTENTION les cas idsExploitants et idsParcelles vides doivent être traités
func GetFermiersFromLieudit(db *sqlx.DB, idLieudit int) ([]*Acteur, error) {
	acteurs := []*Acteur{}
	// parcelles
	idsParcelles := []int{}
	query := "select id_parcelle from parcelle_lieudit where id_lieudit=$1"
	err := db.Select(&idsParcelles, query, idLieudit)
	if err != nil {
		return acteurs, werr.Wrapf(err, "Erreur query : "+query)
	}
	if len(idsParcelles) == 0 {
		return acteurs, nil // empty res
	}
	// ids exploitants
	strIdsParcelles := tiglib.JoinInt(idsParcelles, ",")
	idsExploitants := []int{}
	query = "select distinct id_sctl_exploitant from parcelle_exploitant where id_parcelle in(" + strIdsParcelles + ")"
	err = db.Select(&idsExploitants, query)
	if err != nil {
		return acteurs, werr.Wrapf(err, "Erreur query : "+query)
	}
	if len(idsExploitants) == 0 {
		return acteurs, nil // empty res
	}
	// exploitants
	strIdsExploitants := tiglib.JoinInt(idsExploitants, ",")
	query = "select * from acteur where id_sctl in(" + strIdsExploitants + ") order by nom"
	err = db.Select(&acteurs, query)
	return acteurs, err
}

// Renvoie des Acteurs à partir du début de leurs noms.
// Ne contient que les champs de la table acteur.
// Les autres champs ne sont pas remplis.
func GetActeursAutocomplete(db *sqlx.DB, str string) ([]*Acteur, error) {
	acteurs := []*Acteur{}
	query := "select * from acteur where nom ilike '" + str + "%'"
	err := db.Select(&acteurs, query)
	if err != nil {
		return acteurs, werr.Wrapf(err, "Erreur query : "+query)
	}
	return acteurs, nil
}

// Renvoie des fermiers à partir du début de leurs noms.
// fermier = acteur associé à une ou plusieurs parcelles
// Ne contient que les champs de la table acteur.
// Les autres champs ne sont pas remplis.
func GetFermiersAutocomplete(db *sqlx.DB, str string) ([]*Acteur, error) {
	acteurs := []*Acteur{}
	query := `select * from acteur where nom ilike '` + str + `%' and id in(
        select id_sctl_exploitant from parcelle_exploitant
    )`
	err := db.Select(&acteurs, query)
	if err != nil {
		return acteurs, werr.Wrapf(err, "Erreur query : "+query)
	}
	return acteurs, nil
}

// Renvoie un Acteur à partir de son nom et de son prénom.
// Ne contient que les champs de la table acteur.
// Les autres champs ne sont pas remplis.
// @param  str Chaîne composée de nom + " " + prénom
func GetActeurByNomAutocomplete(db *sqlx.DB, str string) (*Acteur, error) {
	a := &Acteur{}
	query := "select * from acteur where nom || ' ' || prenom=$1"
	row := db.QueryRowx(query, str)
	err := row.StructScan(a)
	if err != nil {
		return a, werr.Wrapf(err, "Erreur query : "+query)
	}
	return a, nil
}

// ************************** CRUD *******************************

func InsertActeur(db *sqlx.DB, acteur *Acteur) (int, error) {
	query := `insert into acteur(
        nom,
        prenom,
        adresse1,
        adresse2,
        cp,
        ville,
        tel,
        telportable,
        email,
        bic,
        iban,
        siret,
        proprietaire,
        fournisseur,
        actif,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		acteur.Nom,
		acteur.Prenom,
		acteur.Adresse1,
		acteur.Adresse2,
		acteur.Cp,
		acteur.Ville,
		acteur.Tel,
		acteur.TelPortable,
		acteur.Email,
		acteur.Bic,
		acteur.Iban,
		acteur.Siret,
		acteur.Proprietaire,
		acteur.Fournisseur,
		acteur.Actif,
		acteur.Notes).Scan(&id)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

func UpdateActeur(db *sqlx.DB, acteur *Acteur) error {
	query := `update acteur set(
        nom,
        prenom,
        adresse1,
        adresse2,
        cp,
        ville,
        tel,
        telportable,
        email,
        bic,
        iban,
        siret,
        proprietaire,
        fournisseur,
        actif,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) where id=$17`
	_, err := db.Exec(
		query,
		acteur.Nom,
		acteur.Prenom,
		acteur.Adresse1,
		acteur.Adresse2,
		acteur.Cp,
		acteur.Ville,
		acteur.Tel,
		acteur.TelPortable,
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
	return nil
}

/*
func DeleteActeur(db *sqlx.DB, id int) error {
    // peut-être ici protection pour davoir si Deletable = true
    // (la situation actuelle fait confiance à l'UI pour ne pas proposer delete sur acteur non deletable)
    query := "delete from acteur where id=$1"
    _, err := db.Exec(query, id)
    if err != nil {
        return werr.Wrapf(err, "Erreur query : " + query)
    }
    return nil
}
*/
