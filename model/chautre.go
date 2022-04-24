/******************************************************************************
    Chautre = Chantiers Autres valorisations
    Bois vendu sur pied à des particuliers, faisant l'objet d'une facturation par BDL

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2020-02-04 19:32:43+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
	//"fmt"
)

type Chautre struct {
	Id            int
	IdAcheteur    int `db:"id_acheteur"`
	TypeVente     string
	TypeValo      string
	DateContrat   time.Time
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
	UGs      []*UG
	Lieudits []*Lieudit
	Fermiers []*Fermier
	Acheteur *Acteur
}

// ************************** Nom *******************************

func (ch *Chautre) String() string {
	if ch.Acheteur == nil {
		panic("Erreur dans le code - L'acheteur d'un chantier autre valorisation doit être calculé avant d'appeler String()")
	}
	return LabelValorisation(ch.TypeValo) + " " + ch.Acheteur.String() + " " + tiglib.DateFr(ch.DateContrat)
}

func (ch *Chautre) FullString() string {
	return "Chantier autre valorisation " + ch.String()
}

// ************************** Get *******************************

// Renvoie un chantier bois sur pied
// contenant uniquement les données stockées en base
func GetChautre(db *sqlx.DB, idChantier int) (*Chautre, error) {
	ch := &Chautre{}
	query := "select * from chautre where id=$1"
	row := db.QueryRowx(query, idChantier)
	err := row.StructScan(ch)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ch, nil
}

// Renvoie un chantier bois sur pied contenant :
//      - les données stockées dans la table
//      - Acheteur
//      - les lieux-dits
//      - les UGs
//      - les fermiers
func GetChautreFull(db *sqlx.DB, idChantier int) (*Chautre, error) {
	ch, err := GetChautre(db, idChantier)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre()")
	}
	err = ch.ComputeAcheteur(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre.ComputeAcheteur()")
	}
	err = ch.ComputeLieudits(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre.ComputeLieuDits()")
	}
	err = ch.ComputeUGs(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre.ComputeUGs()")
	}
	err = ch.ComputeFermiers(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chautre.ComputeFermiers()")
	}
	return ch, nil
}

// Renvoie la liste des années ayant des chantiers bois sur pied,
// @param exclude   Année à exclure du résultat
func GetChautreDifferentYears(db *sqlx.DB, exclude string) ([]string, error) {
	res := []string{}
	list := []time.Time{}
	query := "select datecontrat from chautre order by datecontrat desc"
	err := db.Select(&list, query)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, d := range list {
		y := strconv.Itoa(d.Year())
		if !tiglib.InArrayString(y, res) && y != exclude {
			res = append(res, y)
		}
	}
	return res, nil
}

// Renvoie la liste des chantiers bois sur pied pour une année donnée,
// triés par ordre chronologique inverse.
// Chaque chantier contient les mêmes champs que ceux renvoyés par GetChautreFull()
func GetChautresOfYear(db *sqlx.DB, annee string) ([]*Chautre, error) {
	res := []*Chautre{}
	type ligne struct {
		Id          int
		DateContrat time.Time
	}
	tmp1 := []*ligne{}
	query := "select id,datecontrat from chautre where extract(year from datecontrat)=$1 order by datecontrat"
	err := db.Select(&tmp1, query, annee)
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

func (ch *Chautre) ComputeAcheteur(db *sqlx.DB) error {
	if ch.Acheteur != nil {
		return nil
	}
	var err error
	ch.Acheteur, err = GetActeur(db, ch.IdAcheteur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (ch *Chautre) ComputeUGs(db *sqlx.DB) error {
	if len(ch.UGs) != 0 {
		return nil // déjà calculé
	}
	query := `select * from ug where id in(
	    select id_ug from chantier_ug where type_chantier='chautre' and id_chantier=$1
    )`
	err := db.Select(&ch.UGs, query, &ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func (ch *Chautre) ComputeLieudits(db *sqlx.DB) error {
	if len(ch.Lieudits) != 0 {
		return nil // déjà calculé
	}
	query := `select * from lieudit where id in(
	    select id_lieudit from chantier_lieudit where type_chantier='chautre' and id_chantier=$1
    )`
	err := db.Select(&ch.Lieudits, query, &ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func (ch *Chautre) ComputeFermiers(db *sqlx.DB) error {
	if len(ch.Fermiers) != 0 {
		return nil // déjà calculé
	}
	query := `select * from fermier where id in(
	    select id_fermier from chantier_fermier where type_chantier='chautre' and id_chantier=$1
    )`
	err := db.Select(&ch.Fermiers, query, &ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

// ************************** CRUD *******************************

func InsertChautre(db *sqlx.DB, ch *Chautre, idsUG, idsLieudit, idsFermier []int) (int, error) {
	query := `insert into chautre(
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
        numfacture,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
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
		ch.NumFacture,
		ch.Notes).Scan(&id)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// UGs
	//
	query = `insert into chantier_ug(
        type_chantier,
        id_chantier,
        id_ug) values($1,$2,$3)`
	for _, idUG := range idsUG {
		_, err = db.Exec(
			query,
			"chautre",
			id,
			idUG)
		if err != nil {
			return id, werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	//
	// Lieudits
	//
	query = `insert into chantier_lieudit(
        type_chantier,
        id_chantier,
        id_lieudit) values($1,$2,$3)`
	for _, idLieudit := range idsLieudit {
		_, err = db.Exec(
			query,
			"chautre",
			id,
			idLieudit)
		if err != nil {
			return id, werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	//
	// Fermiers
	//
	query = `insert into chantier_fermier(
        type_chantier,
        id_chantier,
        id_fermier) values($1,$2,$3)`
	for _, idFermier := range idsFermier {
		_, err = db.Exec(
			query,
			"chautre",
			id,
			idFermier)
		if err != nil {
			return id, werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	//
	return id, nil
}

func UpdateChautre(db *sqlx.DB, ch *Chautre, idsUG, idsLieudit, idsFermier []int) error {
	query := `update chautre set(
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
        numfacture,
        notes    
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) where id=$15`
	_, err := db.Exec(
		query,
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
		ch.NumFacture,
		ch.Notes,
		ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// UGs
	//
	query = "delete from chantier_ug where type_chantier='chautre' and id_chantier=$1"
	_, err = db.Exec(query, ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = `insert into chantier_ug(
        type_chantier,
        id_chantier,
        id_ug) values($1,$2,$3)`
	for _, idUG := range idsUG {
		_, err = db.Exec(
			query,
			"chautre",
			ch.Id,
			idUG)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	//
	// Lieudits
	//
	query = "delete from chantier_lieudit where type_chantier='chautre' and id_chantier=$1"
	_, err = db.Exec(query, ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = `insert into chantier_lieudit(
        type_chantier,
        id_chantier,
        id_lieudit) values($1,$2,$3)`
	for _, idLieudit := range idsLieudit {
		_, err = db.Exec(
			query,
			"chautre",
			ch.Id,
			idLieudit)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	//
	// Fermiers
	//
	query = "delete from chantier_fermier where type_chantier='chautre' and id_chantier=$1"
	_, err = db.Exec(query, ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = `insert into chantier_fermier(
        type_chantier,
        id_chantier,
        id_fermier) values($1,$2,$3)`
	for _, idFermier := range idsFermier {
		_, err = db.Exec(
			query,
			"chautre",
			ch.Id,
			idFermier)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	//
	return nil
}

func DeleteChautre(db *sqlx.DB, id int) error {
	//
	// delete UGs, Lieudits, Fermiers associés à ce chantier
	//
	var query string
	var err error
	query = "delete from chantier_ug where type_chantier='chautre' and id_chantier=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	query = "delete from chantier_lieudit where type_chantier='chautre' and id_chantier=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	query = "delete from chantier_fermier where type_chantier='chautre' and id_chantier=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// delete le chantier
	//
	query = "delete from chautre where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
