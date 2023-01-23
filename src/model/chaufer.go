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
	IdUG         int `db:"id_ug"`
	DateChantier time.Time
	Exploitation string // de 1 à 5
	Essence      string
	Volume       float64
	Unite        string // stères ou map
	Notes        string
	// Pas stocké dans la table
	Fermier        *Fermier
	UG             *UG
	LiensParcelles []*ChantierParcelle
}

// ************************** Nom *******************************

func (chantier *Chaufer) String() string {
	if chantier.Fermier == nil {
		panic("Erreur dans le code - Le fermier d'un chantier chauffage fermier doit être calculé avant d'appeler String()")
	}
	return chantier.Fermier.String() + " " + tiglib.DateFr(chantier.DateChantier)
}

func (chantier *Chaufer) FullString() string {
	return "Chantier chauffage fermier " + chantier.String()
}

// ************************** Get *******************************

// Renvoie un chantier chauffage fermier
// contenant uniquement les données stockées en base
func GetChaufer(db *sqlx.DB, idChantier int) (*Chaufer, error) {
	chantier := &Chaufer{}
	query := "select * from chaufer where id=$1"
	row := db.QueryRowx(query, idChantier)
	err := row.StructScan(chantier)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur query : "+query)
	}
	return chantier, nil
}

// Renvoie un chantier chauffage fermier contenant :
//      - les données stockées dans la table
//      - Fermier
//      - UG
//      - LiensParcelles
func GetChauferFull(db *sqlx.DB, idChantier int) (*Chaufer, error) {
	chantier, err := GetChaufer(db, idChantier)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Chaufer()")
	}
	err = chantier.ComputeFermier(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Chaufer.ComputeFermier()")
	}
	err = chantier.ComputeUG(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Chaufer.ComputeUG()")
	}
	err = chantier.ComputeLiensParcelles(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Chaufer.ComputeLiensParcelles()")
	}
	return chantier, nil
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
		chantier, err := GetChauferFull(db, tmp2.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetChauferFull()")
		}
		res = append(res, chantier)
	}
	return res, nil
}

// ************************** Compute *******************************

func (chantier *Chaufer) ComputeFermier(db *sqlx.DB) error {
	if chantier.Fermier != nil {
		return nil
	}
	var err error
	chantier.Fermier, err = GetFermier(db, chantier.IdFermier)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetFermier()")
	}
	return nil
}

func (chantier *Chaufer) ComputeUG(db *sqlx.DB) error {
	if chantier.UG != nil {
		return nil
	}
	var err error
	chantier.UG, err = GetUG(db, chantier.IdUG)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetUG()")
	}
	return nil
}

func (chantier *Chaufer) ComputeLiensParcelles(db *sqlx.DB) error {
	if len(chantier.LiensParcelles) != 0 {
		return nil
	}
	var err error
	query := "select * from chantier_parcelle where type_chantier='chaufer' and id_chantier=$1"
	err = db.Select(&chantier.LiensParcelles, query, chantier.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for i, lien := range chantier.LiensParcelles {
		chantier.LiensParcelles[i].Parcelle, err = GetParcelle(db, lien.IdParcelle)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel LienParcelle()")
		}
	}
	return nil
}

// ************************** CRUD *******************************

func InsertChaufer(db *sqlx.DB, chantier *Chaufer) (int, error) {
	query := `insert into chaufer(
        id_fermier,
        id_ug,
        datechantier,
        exploitation,
        essence,
        volume,
        unite,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		chantier.IdFermier,
		chantier.IdUG,
		chantier.DateChantier,
		chantier.Exploitation,
		chantier.Essence,
		chantier.Volume,
		chantier.Unite,
		chantier.Notes).Scan(&id)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	query = "insert into chantier_parcelle values($1,$2,$3,$4,$5)"
	for _, lien := range chantier.LiensParcelles {
		_, err = db.Exec(
			query,
			id,
			lien.IdParcelle,
			lien.Entiere,
			lien.Surface,
			"chaufer")
		if err != nil {
			return id, werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	return id, nil
}

func UpdateChaufer(db *sqlx.DB, chantier *Chaufer) error {
	query := `update chaufer set(
        id_fermier,
        id_ug,
        datechantier,
        exploitation,
        essence,
        volume,
        unite,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8) where id=$9`
	_, err := db.Exec(
		query,
		chantier.IdFermier,
		chantier.IdUG,
		chantier.DateChantier,
		chantier.Exploitation,
		chantier.Essence,
		chantier.Volume,
		chantier.Unite,
		chantier.Notes,
		chantier.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from chantier_parcelle where type_chantier='chaufer' and id_chantier=$1"
	_, err = db.Exec(query, chantier.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "insert into chantier_parcelle values($1,$2,$3,$4,$5)"
	for _, lien := range chantier.LiensParcelles {
		_, err = db.Exec(
			query,
			lien.IdChantier,
			lien.IdParcelle,
			lien.Entiere,
			lien.Surface,
			"chaufer")
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	return nil
}

func DeleteChaufer(db *sqlx.DB, id int) error {
	query := "delete from chantier_parcelle where type_chantier='chaufer' and id_chantier=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from chaufer where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
