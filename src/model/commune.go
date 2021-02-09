/******************************************************************************
    Communes

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-11-07, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type Commune struct {
	Id       int    // IdCommune de la base SCTL
	Nom      string
	NomCourt string
	Lieudits []*Lieudit
}

const N_COMMUNES = 13

// *********************************************************
/**
    Renvoie une Commune contenant Id et Nom.
    Les autres champs ne sont pas remplis.
**/
func GetCommune(db *sqlx.DB, id int) (*Commune, error) {
	commune := &Commune{}
	query := "select * from commune where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(commune)
	if err != nil {
		return commune, werr.Wrapf(err, "Incapable de fabriquer la commune")
	}
	return commune, err
}

// *********************************************************
/**
    Renvoie la liste de toutes les communes avec leurs lieux-dits
**/
func ListCommunesEtLieudits(db *sqlx.DB) ([]*Commune, error) {
	communes := make([]*Commune, N_COMMUNES)
	query := "select * from commune"
	rows, err := db.Queryx(query)
	if err != nil {
		return communes, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	defer rows.Close()
	queryLD := "select * from lieudit where lieudit.id in(select id_lieudit from commune_lieudit where id_commune=$1)"
	i := 0
	for rows.Next() {
		commune := &Commune{}
		err = rows.StructScan(commune)
		if err != nil {
			return communes, werr.Wrapf(err, "Erreur fabrication Commune")
		}
		rows2, err := db.Queryx(queryLD, commune.Id)
		if err != nil {
			return communes, werr.Wrapf(err, "Erreur query DB : "+query)
		}
		for rows2.Next() {
			lieudit := &Lieudit{}
			err = rows2.StructScan(lieudit)
			if err != nil {
				rows2.Close()
				return communes, werr.Wrapf(err, "Erreur fabrication LieuDit")
			}
			commune.Lieudits = append(commune.Lieudits, lieudit)
		}
		// voir si ok de Close comme ça dans la boucle
		rows2.Close()
		err = rows2.Err()
		if err != nil {
			return communes, werr.Wrapf(err, "Erreur rows2.Next()")
		}
		communes[i] = commune
		i++
	}
	err = rows.Err()
	if err != nil {
		return communes, werr.Wrapf(err, "Erreur rows.Next()")
	}
	return communes, nil
}
