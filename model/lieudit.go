/******************************************************************************
    Lieux-dits

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2019-11-07 10:07:45+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type Lieudit struct {
	Id        int // IdLieuDit de la base SCTL
	Nom       string
	Parcelles []*Parcelle
	Communes  []*Commune
}

// ************************** Get one *******************************

// Renvoie un Lieudit à partir de son id.
// Ne contient que les champs de la table lieudit.
// Les autres champs ne sont pas remplis.
func GetLieudit(db *sqlx.DB, id int) (ld *Lieudit, err error) {
	ld = &Lieudit{}
	query := "select * from lieudit where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(ld)
	if err != nil {
		return ld, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ld, nil
}

// Renvoie un Lieudit à partir de son nom.
// Ne contient que les champs de la table lieudit.
// Les autres champs ne sont pas remplis.
func GetLieuditByNom(db *sqlx.DB, nom string) (ld *Lieudit, err error) {
	ld = &Lieudit{}
	query := "select * from lieudit where nom=$1"
	row := db.QueryRowx(query, nom)
	err = row.StructScan(ld)
	if err != nil {
		return ld, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ld, nil
}

// ************************** Get many *******************************

// Renvoie des Lieudit à partir du début du nom.
// Les mots comme LE LA LES DE DU D' ne sont pas pris en compte.
func GetLieuditsAutocomplete(db *sqlx.DB, str string) (lds []*Lieudit, err error) {
	lds = []*Lieudit{}
	query := "select id,nom from lieudit_mot where mot ilike '" + str + "%'"
	err = db.Select(&lds, query)
	if err != nil {
		return lds, werr.Wrapf(err, "Erreur query : "+query)
	}
	return lds, nil
}

// Renvoie des Lieudit à partir d'un code UG.
// Utilise les parcelles pour faire le lien
// Ne contient que les champs de la table lieudit
// + le champ Communes.
// Les autres champs ne sont pas remplis.
func GetLieuditsFromCodeUG(db *sqlx.DB, codeUG string) (lds []*Lieudit, err error) {
	lds = []*Lieudit{}
	query := `
	    select * from lieudit where id in(
            select id_lieudit from parcelle_lieudit where id_parcelle in(
                select id_parcelle from parcelle_ug where id_ug in(
                    select id from ug where code=$1
                )
            )
        )`
	err = db.Select(&lds, query, codeUG)
	if err != nil {
		return lds, werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, ld := range lds {
		err = ld.ComputeCommunes(db)
		if err != nil {
			return lds, werr.Wrapf(err, "Erreur appel ComputeCommunes()")
		}
	}
	return lds, nil
}

// ************************** Compute *******************************

// Remplit le champ Parcelles d'un Lieudit
func (ld *Lieudit) ComputeParcelles(db *sqlx.DB) (err error) {
	if len(ld.Parcelles) != 0 {
		return nil // déjà calculé
	}
	query := `
	    select * from parcelle where id in(
            select id_parcelle from parcelle_lieudit where id_lieudit=$1
        ) order by code`
	err = db.Select(&ld.Parcelles, query, ld.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

// Remplit le champ Communes d'un Lieudit
func (ld *Lieudit) ComputeCommunes(db *sqlx.DB) (err error) {
	if len(ld.Communes) != 0 {
		return nil // déjà calculé
	}
	query := `  
	    select * from commune where id in(
            select id_commune from commune_lieudit where id_lieudit=$1
        ) order by nom`
	err = db.Select(&ld.Communes, query, ld.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
