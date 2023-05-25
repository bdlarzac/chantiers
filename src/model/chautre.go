/*
*****************************************************************************

	Chautre = Chantiers Autres valorisations
	Bois vendu sur pied à des particuliers, faisant l'objet d'une facturation par BDL

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2020-02-04 19:32:43+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type Chautre struct {
	Id            int
	IdAcheteur    int `db:"id_acheteur"`
	Titre         string
	TypeVente     string
	TypeValo      string
	DateContrat   time.Time
	DatePaiement  time.Time
	Exploitation  string
	Essence       string
	VolumeContrat float64
	VolumeRealise float64
	Unite         string
	PUHT          float64
	TVA           float64
	DateFacture   time.Time
	NumFacture    string
	Notes         string
	// pas stocké en base
	UGs            []*UG
	LiensParcelles []*ChantierParcelle
	Lieudits       []*Lieudit
	Fermiers       []*Fermier
	Proprietaires  []*Acteur
	Acheteur       *Acteur
}

// Association code type vente => label
// Les codes correspondent aux valeurs stockées en base dans chautre.typevente
var ChautreTypeVenteMap = map[string]string{
	"NON": "Non spécifié",
	"BSP": "Bois-sur-pied",
	"BDR": "Bord-de-route",
	"LIV": "Livré",
}

// ************************** Nom *******************************

func (ch *Chautre) String() string {
	return ch.Titre
}

func (ch *Chautre) FullString() string {
	return "Chantier autre valorisation " + ch.String()
}

// ************************** Get *******************************

// Renvoie un chantier autres valorisations
// contenant uniquement les données stockées en base
func GetChautre(db *sqlx.DB, idChantier int) (ch *Chautre, err error) {
	ch = &Chautre{}
	query := "select * from chautre where id=$1"
	row := db.QueryRowx(query, idChantier)
	err = row.StructScan(ch)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ch, nil
}

/*
Renvoie un chantier autres valorisations contenant :
  - les données stockées dans la table
  - Acheteur
  - les UGs
  - les parcelles
  - les lieux-dits
  - les fermiers
*/
func GetChautreFull(db *sqlx.DB, idChantier int) (ch *Chautre, err error) {
	ch, err = GetChautre(db, idChantier)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel GetChautre()")
	}
	err = ch.ComputeAcheteur(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre.ComputeAcheteur()")
	}
	err = ch.ComputeUGs(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre.ComputeUGs()")
	}
	err = ch.ComputeLiensParcelles(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre.ComputeLiensParcelles()")
	}
	err = ch.ComputeLieudits(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre.ComputeLieuDits()")
	}
	err = ch.ComputeFermiers(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre.ComputeFermiers()")
	}
	err = ch.ComputeProprietaires(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre.ComputeProprietaires()")
	}
	return ch, nil
}

/*
Renvoie la liste des années ayant des chantiers autres valorisations,
@param exclude   Année à exclure du résultat
*/
func GetChautreDifferentYears(db *sqlx.DB, exclude string) (res []string, err error) {
	res = []string{}
	list := []time.Time{}
	query := "select datecontrat from chautre order by datecontrat desc"
	err = db.Select(&list, query)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, d := range list {
		y := strconv.Itoa(d.Year())
		if !tiglib.InArray(y, res) && y != exclude {
			res = append(res, y)
		}
	}
	return res, nil
}

/*
Renvoie la liste des chantiers autres valorisations pour une année donnée,
triés par ordre chronologique inverse.
Chaque chantier contient les mêmes champs que ceux renvoyés par GetChautreFull()
*/
func GetChautresOfYear(db *sqlx.DB, annee string) (res []*Chautre, err error) {
	res = []*Chautre{}
	type ligne struct {
		Id          int
		DateContrat time.Time
	}
	tmp1 := []*ligne{}
	query := "select id,datecontrat from chautre where extract(year from datecontrat)=$1 order by datecontrat"
	err = db.Select(&tmp1, query, annee)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, tmp2 := range tmp1 {
		ch, err := GetChautreFull(db, tmp2.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetChautreFull()")
		}
		res = append(res, ch)
	}
	return res, nil
}

// ************************** Compute *******************************

func (ch *Chautre) ComputeAcheteur(db *sqlx.DB) (err error) {
	if ch.Acheteur != nil {
		return nil // déjà calculé
	}
	ch.Acheteur, err = GetActeur(db, ch.IdAcheteur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (ch *Chautre) ComputeUGs(db *sqlx.DB) (err error) {
	if len(ch.UGs) != 0 {
		return nil // déjà calculé
	}
	ch.UGs, err = computeUGsOfChantier(db, "chautre", ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeUGsOfChantier")
	}
	return nil
}

func (ch *Chautre) ComputeLiensParcelles(db *sqlx.DB) (err error) {
	if len(ch.LiensParcelles) != 0 {
		return nil // déjà calculé
	}
	ch.LiensParcelles, err = computeLiensParcellesOfChantier(db, "chautre", ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeLiensParcellesOfChantier")
	}
	return nil
}

func (ch *Chautre) ComputeLieudits(db *sqlx.DB) (err error) {
	if len(ch.Lieudits) != 0 {
		return nil // déjà calculé
	}
	ch.Lieudits, err = computeLieuditsOfChantier(db, "chautre", ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeLieuditsOfChantier")
	}
	return nil
}

func (ch *Chautre) ComputeFermiers(db *sqlx.DB) (err error) {
	if len(ch.Fermiers) != 0 {
		return nil // déjà calculé
	}
	ch.Fermiers, err = computeFermiersOfChantier(db, "chautre", ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeFermiersOfChantier")
	}
	return nil
}

func (ch *Chautre) ComputeProprietaires(db *sqlx.DB) (err error) {
	if len(ch.Proprietaires) != 0 {
		return nil // déjà calculé
	}
	err = ch.ComputeUGs(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ch.ComputeUGs()")
	}
	idsProprios := []int{}
	for _, ug := range ch.UGs {
		err = ug.ComputeProprietaires(db)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel UG.ComputeProprietaires()")
		}
		for _, proprio := range ug.Proprietaires {
			if !tiglib.InArray(proprio.Id, idsProprios) {
				idsProprios = append(idsProprios, proprio.Id)
				ch.Proprietaires = append(ch.Proprietaires, proprio)
			}
		}
	}
	return nil
}

// ************************** CRUD *******************************

func InsertChautre(db *sqlx.DB, ch *Chautre, idsUG, idsLieudit, idsFermier []int) (idChantier int, err error) {
	query := `insert into chautre(
        titre,
        id_acheteur,
        typevente,
        typevalo,
        datecontrat,
        exploitation,
        essence,
        volumecontrat,
        volumerealise,
        unite,
        puht,
        tva,
        datefacture,
        datepaiement,
        numfacture,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) returning id`
	idChantier = int(0)
	err = db.QueryRow(
		query,
		ch.Titre,
		ch.IdAcheteur,
		ch.TypeVente,
		ch.TypeValo,
		ch.DateContrat,
		ch.Exploitation,
		ch.Essence,
		ch.VolumeContrat,
		ch.VolumeRealise,
		ch.Unite,
		ch.PUHT,
		ch.TVA,
		ch.DateFacture,
		ch.DatePaiement,
		ch.NumFacture,
		ch.Notes).Scan(&idChantier)
	if err != nil {
		return idChantier, werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// insert associations avec UGs, Parcelles, Lieudits, Fermiers
	//
	fmt.Printf("InsertChautre() - idsUG = %+v\n", idsUG)
	err = insertLiensChantierUG(db, "chautre", idChantier, idsUG)
	if err != nil {
		return idChantier, werr.Wrapf(err, "Erreur appel insertLiensChantierUG()")
	}
	//
	err = insertLiensChantierParcelle(db, "chautre", idChantier, ch.LiensParcelles)
	if err != nil {
		return idChantier, werr.Wrapf(err, "Erreur appel insertLiensChantierParcelle()")
	}
	//
	err = insertLiensChantierLieudit(db, "chautre", idChantier, idsLieudit)
	if err != nil {
		return idChantier, werr.Wrapf(err, "Erreur appel insertLiensChantierLieudit()")
	}
	//
	err = insertLiensChantierFermier(db, "chautre", idChantier, idsFermier)
	if err != nil {
		return idChantier, werr.Wrapf(err, "Erreur appel insertLiensChantierFermier()")
	}
	//
	return idChantier, nil
}

func UpdateChautre(db *sqlx.DB, ch *Chautre, idsUG, idsLieudit, idsFermier []int) (err error) {
	query := `update chautre set(
        titre,
        id_acheteur,
        typevente,
        typevalo,
        datecontrat,
        exploitation,
        essence,
        volumecontrat,
        volumerealise,
        unite,
        puht,
        tva,
        datefacture,
        datepaiement,
        numfacture,
        notes    
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) where id=$17`
	_, err = db.Exec(
		query,
		ch.Titre,
		ch.IdAcheteur,
		ch.TypeVente,
		ch.TypeValo,
		ch.DateContrat,
		ch.Exploitation,
		ch.Essence,
		ch.VolumeContrat,
		ch.VolumeRealise,
		ch.Unite,
		ch.PUHT,
		ch.TVA,
		ch.DateFacture,
		ch.DatePaiement,
		ch.NumFacture,
		ch.Notes,
		ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// update associations avec UGs, Parcelles, Lieudits, Fermiers
	//
	err = updateLiensChantierUG(db, "chautre", ch.Id, idsUG)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel updateLiensChantierUG()")
	}
	//
	err = updateLiensChantierParcelle(db, "chautre", ch.Id, ch.LiensParcelles)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel updateLiensChantierParcelle()")
	}
	//
	err = updateLiensChantierLieudit(db, "chautre", ch.Id, idsLieudit)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel updateLiensChantierLieudit()")
	}
	//
	err = updateLiensChantierFermier(db, "chautre", ch.Id, idsFermier)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel updateLiensChantierFermier()")
	}
	//
	return nil
}

func DeleteChautre(db *sqlx.DB, id int) (err error) {
	//
	// delete associations avec UGs, Parcelles, Lieudits, Fermiers
	//
	err = deleteLiensChantierUG(db, "chautre", id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierUG()")
	}
	//
	err = deleteLiensChantierParcelle(db, "chautre", id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierParcelle()")
	}
	//
	err = deleteLiensChantierLieudit(db, "chautre", id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierLieudit()")
	}
	//
	err = deleteLiensChantierFermier(db, "chautre", id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierFermier()")
	}
	//
	// delete le chantier, fait à la fin pour respecter les clés étrangères
	//
	query := "delete from chautre where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
