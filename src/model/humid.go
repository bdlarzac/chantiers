/******************************************************************************
    Mesures d'humidité effectuées sur les tas de plaquettes

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-12-20 16:30:41+01:00, Thierry Graff : Creation
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

type Humid struct {
	Id         int
	IdTas      int `db:"id_tas"`
	Valeur     float64
	DateMesure time.Time
	Notes      string
	// pas stocké en base
	IdsMesureurs []int
	Mesureurs    []*Acteur
	Tas          *Tas
}

// ************************** Get *******************************
// Renvoie une mesure d'humidité contenant
// - les données stockées en base
// - les mesureurs
// - le stockage
func GetHumidFull(db *sqlx.DB, idMesure int) (*Humid, error) {
	humid := &Humid{}
	query := "select * from humid where id=$1"
	row := db.QueryRowx(query, idMesure)
	err := row.StructScan(humid)
	if err != nil {
		return humid, werr.Wrapf(err, "Erreur query : "+query)
	}
	// mesureurs
	query2 := "select id_acteur from humid_acteur where id_humid=$1"
	rows, err := db.Query(query2, humid.Id)
	if err != nil {
		return humid, werr.Wrapf(err, "Erreur query DB : "+query2)
	}
	defer rows.Close()
	var idActeur int
	for rows.Next() {
		err = rows.Scan(&idActeur)
		if err != nil {
			return humid, werr.Wrapf(err, "Erreur row scan id acteur pour mesure humidité")
		}
		acteur, err := GetActeur(db, idActeur)
		if err != nil {
			return humid, werr.Wrapf(err, "Erreur récup acteur pour mesure humidité")
		}
		humid.Mesureurs = append(humid.Mesureurs, acteur)
	}
	// tas
	humid.Tas, err = GetTasFull(db, humid.IdTas)
	if err != nil {
		return humid, werr.Wrapf(err, "Erreur appel GetTas()")
	}
	return humid, nil
}

// Renvoie la liste des années ayant des mesures d'humidité,
// @param exclude   Année à exclure du résultat
func GetHumidDifferentYears(db *sqlx.DB, exclude string) ([]string, error) {
	res := []string{}
	list := []time.Time{}
	query := "select datemesure from humid order by datemesure desc"
	err := db.Select(&list, query)
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

// Renvoie la liste des mesures d'humidité pour une année donnée,
// triés par ordre chronologique inverse.
// Chaque mesure contient les mêmes champs que ceux renvoyés par GetHumidFull()
func GetHumidsOfYear(db *sqlx.DB, annee string) ([]*Humid, error) {
	res := []*Humid{}
	type ligne struct {
		Id         int
		DateMesure time.Time
	}
	tmp1 := []*ligne{}
	query := "select id,datemesure from humid where extract(year from datemesure)=$1 order by datemesure desc"
	err := db.Select(&tmp1, query, annee)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, tmp2 := range tmp1 {
		mesure, err := GetHumidFull(db, tmp2.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetHumidFull()")
		}
		res = append(res, mesure)
	}
	return res, nil
}

// ************************** CRUD *******************************

func InsertHumid(db *sqlx.DB, humid *Humid) (int, error) {
	query := `insert into humid(
        id_tas,
        valeur,
        datemesure,
        notes
        ) values($1,$2,$3,$4) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		humid.IdTas,
		humid.Valeur,
		humid.DateMesure,
		humid.Notes).Scan(&id)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "insert into humid_acteur values($1,$2)"
	for _, idMesureur := range humid.IdsMesureurs {
		_ = db.QueryRow(query, id, idMesureur)
	}
	return id, nil
}

func UpdateHumid(db *sqlx.DB, humid *Humid) error {
	query := `update humid set(
        id_tas,
        valeur,
        datemesure,
        notes
        ) = ($1,$2,$3,$4) where id=$5`
	_, err := db.Exec(
		query,
		humid.IdTas,
		humid.Valeur,
		humid.DateMesure,
		humid.Notes,
		humid.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from humid_acteur where id_humid=$1"
	_, err = db.Exec(query, humid.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "insert into humid_acteur values($1,$2)"
	for _, idMesureur := range humid.IdsMesureurs {
		_ = db.QueryRow(query, humid.Id, idMesureur)
	}
	return nil
}

func DeleteHumid(db *sqlx.DB, id int) error {
	query := "delete from humid_acteur where id_humid=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from humid where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
