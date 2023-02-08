/******************************************************************************
    Chaufer = Chantier Bois de chauffage fermiers

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2020-02-04 18:55:10+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type Chaufer struct {
	Id           int
	IdFermier    int `db:"id_fermier"`
	DateChantier time.Time
	Exploitation string // de 1 à 5
	Essence      string
	Volume       float64
	Unite        string // stères ou map
	Notes        string
	// Pas stocké dans la table
	Fermier        *Fermier
	UGs            []*UG
	LiensParcelles []*ChantierParcelle
}

// ************************** Nom *******************************

func (ch *Chaufer) String() string {
	if ch.Fermier == nil {
		panic("Erreur dans le code - Le fermier d'un chantier chauffage fermier doit être calculé avant d'appeler String()")
	}
	return ch.Fermier.String() + " " + tiglib.DateFr(ch.DateChantier)
}

func (ch *Chaufer) FullString() string {
	return "Chantier chauffage fermier " + ch.String()
}

// ************************** Get *******************************

// Renvoie un chantier chauffage fermier
// contenant uniquement les données stockées en base
func GetChaufer(db *sqlx.DB, idChantier int) (*Chaufer, error) {
	ch := &Chaufer{}
	query := "select * from chaufer where id=$1"
	row := db.QueryRowx(query, idChantier)
	err := row.StructScan(ch)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ch, nil
}

// Renvoie un chantier chauffage fermier contenant :
//      - les données stockées dans la table
//      - Fermier
//      - UG
//      - LiensParcelles
func GetChauferFull(db *sqlx.DB, idChantier int) (*Chaufer, error) {
	ch, err := GetChaufer(db, idChantier)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chaufer()")
	}
	err = ch.ComputeFermier(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chaufer.ComputeFermier()")
	}
	err = ch.ComputeUGs(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chaufer.ComputeUgs()")
	}
	err = ch.ComputeLiensParcelles(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Chaufer.ComputeLiensParcelles()")
	}
	return ch, nil
}

// Renvoie la liste des années ayant des chantiers chauffage fermier,
// @param exclude   Année à exclure du résultat
func GetChauferDifferentYears(db *sqlx.DB, exclude string) ([]string, error) {
	res := []string{}
	list := []time.Time{}
	query := "select datechantier from chaufer order by datechantier desc"
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

// Renvoie la liste des chantiers chauffage fermier pour une année donnée,
// triés par ordre chronologique inverse.
// Chaque chantier contient les mêmes champs que ceux renvoyés par GetChauferFull()
func GetChaufersOfYear(db *sqlx.DB, annee string) ([]*Chaufer, error) {
	res := []*Chaufer{}
	type ligne struct {
		Id           int
		DateChantier time.Time
	}
	tmp1 := []*ligne{}
	// select aussi datecontrat au lieu de seulement id pour pouvoir faire le order by
	query := "select id,datechantier from chaufer where extract(year from datechantier)=$1 order by datechantier"
	err := db.Select(&tmp1, query, annee)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, tmp2 := range tmp1 {
		ch, err := GetChauferFull(db, tmp2.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetChauferFull()")
		}
		res = append(res, ch)
	}
	return res, nil
}

// ************************** Compute *******************************

func (ch *Chaufer) ComputeFermier(db *sqlx.DB) error {
	if ch.Fermier != nil {
		return nil
	}
	var err error
	ch.Fermier, err = GetFermier(db, ch.IdFermier)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetFermier()")
	}
	return nil
}

func (ch *Chaufer) ComputeUGs(db *sqlx.DB) (err error) {
	if len(ch.UGs) != 0 {
		return nil // déjà calculé
	}
	ch.UGs, err = computeUGsOfChantier(db, "chaufer", ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeUGsOfChantier")
	}
	return nil
}

func (ch *Chaufer) ComputeLiensParcelles(db *sqlx.DB) (err error) {
	if len(ch.LiensParcelles) != 0 {
		return nil
	}
	ch.LiensParcelles, err = computeLiensParcellesOfChantier(db, "chaufer", ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeLiensParcellesOfChantier")
	}
	return nil
}

// ************************** CRUD *******************************

func InsertChaufer(db *sqlx.DB, ch *Chaufer, idsUG []int) (idChantier int, err error) {
	query := `insert into chaufer(
        id_fermier,
        datechantier,
        exploitation,
        essence,
        volume,
        unite,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7) returning id`
	idChantier = int(0)
	err = db.QueryRow(
		query,
		ch.IdFermier,
		ch.DateChantier,
		ch.Exploitation,
		ch.Essence,
		ch.Volume,
		ch.Unite,
		ch.Notes).Scan(&idChantier)
	if err != nil {
		return idChantier, werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// insert associations avec UGs, Parcelles
	//
	err = insertLiensChantierUG(db, "chaufer", idChantier, idsUG)
	if err != nil {
		return idChantier, werr.Wrapf(err, "Erreur appel insertLiensChantierUG()")
	}
	//
	err = insertLiensChantierParcelle(db, "chaufer", idChantier, ch.LiensParcelles)
	if err != nil {
		return idChantier, werr.Wrapf(err, "Erreur appel insertLiensChantierParcelle()")
	}
	return idChantier, nil
}

func UpdateChaufer(db *sqlx.DB, ch *Chaufer, idsUG []int) (err error) {
	query := `update chaufer set(
        id_fermier,
        datechantier,
        exploitation,
        essence,
        volume,
        unite,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7) where id=$8`
	_, err = db.Exec(
		query,
		ch.IdFermier,
		ch.DateChantier,
		ch.Exploitation,
		ch.Essence,
		ch.Volume,
		ch.Unite,
		ch.Notes,
		ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// update associations avec UGs, Parcelles
	//
	err = updateLiensChantierUG(db, "chaufer", ch.Id, idsUG)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel updateLiensChantierUG()")
	}
	//
	err = updateLiensChantierParcelle(db, "chaufer", ch.Id, ch.LiensParcelles)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel updateLiensChantierParcelle()")
	}
	//
	return nil
}

func DeleteChaufer(db *sqlx.DB, id int) (err error) {
	//
	// delete associations avec UGs, Parcelles
	//
    err = deleteLiensChantierUG(db, "chaufer", id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierUG()")
	}
	//
    err = deleteLiensChantierParcelle(db, "chaufer", id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierParcelle()")
	}
	//
	// delete le chantier, fait à la fin pour respecter les clés étrangères
	//
	query := "delete from chaufer where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
