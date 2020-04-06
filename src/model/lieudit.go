/******************************************************************************
    Lieux-dits

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-11-07 10:07:45+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type Lieudit struct {
	Id        int
	Nom       string
	Parcelles []*Parcelle
	Communes  []*Commune
}

// ************************** Get *******************************
// Renvoie un Lieudit à partir de son id.
// Ne contient que les champs de la table lieudit.
// Les autres champs ne sont pas remplis.
func GetLieudit(db *sqlx.DB, id int) (*Lieudit, error) {
	ld := &Lieudit{}
	query := "select * from lieudit where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(ld)
	if err != nil {
		return ld, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ld, err
}

// Renvoie un Lieudit à partir de son nom.
// Ne contient que les champs de la table lieudit.
// Les autres champs ne sont pas remplis.
func GetLieuditByNom(db *sqlx.DB, nom string) (*Lieudit, error) {
	ld := &Lieudit{}
	query := "select * from lieudit where nom=$1"
	row := db.QueryRowx(query, nom)
	err := row.StructScan(ld)
	if err != nil {
		return ld, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ld, nil
}

// Renvoie des Lieudit à partir du début du nom.
// Les mots comme LE LA LES DE DU D' ne sont pas pris en compte.
func GetLieuDitAutocomplete(db *sqlx.DB, str string) ([]*Lieudit, error) {
	lds := []*Lieudit{}
	query := "select id,nom from lieudit_mot where mot ilike '" + str + "%'"
	err := db.Select(&lds, query)
	if err != nil {
		return lds, werr.Wrapf(err, "Erreur query : "+query)
	}
	return lds, nil
}

// ************************** Compute *******************************

// Remplit le champ Parcelles d'un Lieudit
func (ld *Lieudit) ComputeParcelles(db *sqlx.DB) error {
	query := "select id_parcelle from parcelle_lieudit where id_lieudit=$1"
	rows, err := db.Query(query, ld.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	defer rows.Close()
	var idP int
	for rows.Next() {
		if err := rows.Scan(&idP); err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
		parcelle, err := GetParcelle(db, idP)
		if err != nil {
			return werr.Wrapf(err, "Erreur GetParcelle()")
		}
		ld.Parcelles = append(ld.Parcelles, parcelle)
	}
	err = rows.Err()
	if err != nil {
		return werr.Wrapf(err, "Erreur rows.Next()")
	}
	return nil
}

// Remplit le champ Communes d'un Lieudit
func (ld *Lieudit) ComputeCommunes(db *sqlx.DB) error {
	query := "select id_commune from commune_lieudit where id_lieudit=$1"
	rows, err := db.Query(query, ld.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	defer rows.Close()
	var idC int
	for rows.Next() {
		if err := rows.Scan(&idC); err != nil {
			return werr.Wrapf(err, "Erreur query : %s"+query)
		}
		commune, err := GetCommune(db, idC)
		if err != nil {
			return werr.Wrapf(err, "Erreur GetCommune()")
		}
		ld.Communes = append(ld.Communes, commune)
	}
	err = rows.Err()
	if err != nil {
		return werr.Wrapf(err, "Erreur rows.Next()")
	}
	return nil
}
