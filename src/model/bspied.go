/******************************************************************************
    BSPied = Chantier Bois sur pied
    Bois vendu sur pied à des particuliers, faisant l'objet d'une facturation par BDL

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-02-04 19:32:43+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
	//"fmt"
)

type BSPied struct {
	Id            int
	IdAcheteur    int `db:"id_acheteur"`
	Nom           string
	DateContrat   time.Time
	Exploitation  string
	Essence       string
	NStereContrat float64
	NStereCoupees float64
	PrixStere     float64
	TVA           float64
	DateFacture   time.Time
	NumFacture    string
	Notes         string
	// pas stocké en base
	UGs        []*UG
	Lieudits   []*Lieudit
	Fermiers   []*Acteur
	Acheteur  *Acteur
}

// ************************** Nom *******************************

func (ch *BSPied) String() string {
	if len(ch.Lieudits) == 0 {
		panic("Erreur dans le code - Les lieux-dits d'un chantier bois-sur-pied doivent être calculés avant d'appeler String()")
	}
	if ch.Acheteur == nil {
		panic("Erreur dans le code - L'acheteur d'un chantier bois sur pied doit être calculé avant d'appeler String()")
	}
	var noms []string 
    for _, ld := range ch.Lieudits {
        noms = append(noms, ld.Nom)
    }
	return ch.Acheteur.String() + " " + strings.Join(noms, " - ") + " " + tiglib.DateFr(ch.DateContrat)
}

func (ch *BSPied) FullString() string {
	return "Chantier bois sur pied " + ch.String()
}

// ************************** Get *******************************

// Renvoie un chantier bois sur pied
// contenant uniquement les données stockées en base
func GetBSPied(db *sqlx.DB, idChantier int) (*BSPied, error) {
	chantier := &BSPied{}
	query := "select * from bspied where id=$1"
	row := db.QueryRowx(query, idChantier)
	err := row.StructScan(chantier)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur query : "+query)
	}
	return chantier, nil
}

// Renvoie un chantier bois sur pied contenant :
//      - les données stockées dans la table
//      - Acheteur
//      - les lieux-dits
//      - les UGs
//      - les fermiers
func GetBSPiedFull(db *sqlx.DB, idChantier int) (*BSPied, error) {
	ch, err := GetBSPied(db, idChantier)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel GetBSPied()")
	}
	err = ch.ComputeAcheteur(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel BSPied.ComputeAcheteur()")
	}
	err = ch.ComputeLieudits(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel BSPied.ComputeLieudits()")
	}
	err = ch.ComputeUGs(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel BSPied.ComputeUGs()")
	}
	err = ch.ComputeFermiers(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel BSPied.ComputeFermiers()")
	}
	return ch, nil
}

// Renvoie la liste des années ayant des chantiers bois sur pied,
// @param exclude   Année à exclure du résultat
func GetBSPiedDifferentYears(db *sqlx.DB, exclude string) ([]string, error) {
	res := []string{}
	list := []time.Time{}
	query := "select datecontrat from bspied order by datecontrat desc"
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
// Chaque chantier contient les mêmes champs que ceux renvoyés par GetBSPiedFull()
func GetBSPiedsOfYear(db *sqlx.DB, annee string) ([]*BSPied, error) {
	res := []*BSPied{}
	type ligne struct {
		Id          int
		DateContrat time.Time
	}
	tmp1 := []*ligne{}
	// select aussi datecontrat au lieu de seulement id pour pouvoir faire le order by
	query := "select id,datecontrat from bspied where extract(year from datecontrat)=$1 order by datecontrat desc"
	err := db.Select(&tmp1, query, annee)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, tmp2 := range tmp1 {
		chantier, err := GetBSPiedFull(db, tmp2.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetBSPiedFull()")
		}
		res = append(res, chantier)
	}
	return res, nil
}

// ************************** Compute *******************************

func (ch *BSPied) ComputeAcheteur(db *sqlx.DB) error {
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

func (ch *BSPied) ComputeUGs(db *sqlx.DB) error {
	if len(ch.UGs) != 0 {
		return nil // déjà calculé
	}
	query := `select * from ug where id in(
	    select id_ug from chantier_ug where type_chantier='bspied' and id_chantier=$1
    )`
	err := db.Select(&ch.UGs, query, &ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func (ch *BSPied) ComputeLieudits(db *sqlx.DB) error {
	if len(ch.Lieudits) != 0 {
		return nil // déjà calculé
	}
	query := `select * from lieudit where id in(
	    select id_lieudit from chantier_lieudit where type_chantier='bspied' and id_chantier=$1
    )`
	err := db.Select(&ch.Lieudits, query, &ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func (ch *BSPied) ComputeFermiers(db *sqlx.DB) error {
	if len(ch.Fermiers) != 0 {
		return nil // déjà calculé
	}
	query := `select * from acteur where id in(
	    select id_fermier from chantier_fermier where type_chantier='bspied' and id_chantier=$1
    )`
	err := db.Select(&ch.Fermiers, query, &ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

// ************************** CRUD *******************************

func InsertBSPied(db *sqlx.DB, ch *BSPied, idsUG, idsLieudit, idsFermier []int) (int, error) {
	query := `insert into bspied(
        id_acheteur,
        datecontrat,
        exploitation,
        essence,
        nsterecontrat,
        nsterecoupees,
        prixstere,
        tva,
        datefacture,
        numfacture,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		ch.IdAcheteur,
		ch.DateContrat,
		ch.Exploitation,
		ch.Essence,
		ch.NStereContrat,
		ch.NStereCoupees,
		ch.PrixStere,
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
	for _, idUG := range(idsUG){
        _, err = db.Exec(
            query,
            "bspied",
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
	for _, idLieudit := range(idsLieudit){
        _, err = db.Exec(
            query,
            "bspied",
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
	for _, idFermier := range(idsFermier){
        _, err = db.Exec(
            query,
            "bspied",
            id,
            idFermier)
        if err != nil {
            return id, werr.Wrapf(err, "Erreur query : "+query)
        }
    }
    //
	return id, nil
}

func UpdateBSPied(db *sqlx.DB, ch *BSPied, idsUG, idsLieudit, idsFermier []int) error {
	query := `update bspied set(
        id_acheteur,
        datecontrat,
        exploitation,
        essence,
        nsterecontrat,
        nsterecoupees,
        prixstere,
        tva,
        datefacture,
        numfacture,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) where id=$12`
	_, err := db.Exec(
		query,
		ch.IdAcheteur,
		ch.DateContrat,
		ch.Exploitation,
		ch.Essence,
		ch.NStereContrat,
		ch.NStereCoupees,
		ch.PrixStere,
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
	query = "delete from chantier_ug where type_chantier='bspied' and id_chantier=$1"
	_, err = db.Exec(query, ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = `insert into chantier_ug(
        type_chantier,
        id_chantier,
        id_ug) values($1,$2,$3)`
	for _, idUG := range(idsUG){
        _, err = db.Exec(
            query,
            "bspied",
            ch.Id,
            idUG)
        if err != nil {
            return werr.Wrapf(err, "Erreur query : "+query)
        }
    }
    //
	// Lieudits
	//
	query = "delete from chantier_lieudit where type_chantier='bspied' and id_chantier=$1"
	_, err = db.Exec(query, ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = `insert into chantier_lieudit(
        type_chantier,
        id_chantier,
        id_lieudit) values($1,$2,$3)`
	for _, idLieudit := range(idsLieudit){
        _, err = db.Exec(
            query,
            "bspied",
            ch.Id,
            idLieudit)
        if err != nil {
            return werr.Wrapf(err, "Erreur query : "+query)
        }
    }
    //
	// Fermiers
	//
	query = "delete from chantier_fermier where type_chantier='bspied' and id_chantier=$1"
	_, err = db.Exec(query, ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = `insert into chantier_fermier(
        type_chantier,
        id_chantier,
        id_fermier) values($1,$2,$3)`
	for _, idFermier := range(idsFermier){
        _, err = db.Exec(
            query,
            "bspied",
            ch.Id,
            idFermier)
        if err != nil {
            return werr.Wrapf(err, "Erreur query : "+query)
        }
    }
	//
	return nil
}

func DeleteBSPied(db *sqlx.DB, id int) error {
	//
	// delete UGs, Lieudits, Fermiers associés à ce chantier
	//
	var query string
	var err error
	query = "delete from chantier_ug where type_chantier='bspied' and id_chantier=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	query = "delete from chantier_lieudit where type_chantier='bspied' and id_chantier=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	query = "delete from chantier_fermier where type_chantier='bspied' and id_chantier=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// delete le chantier
	//
	query = "delete from bspied where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	return nil
}
