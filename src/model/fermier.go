/******************************************************************************
    Fermiers = exploitants agricoles, issus de la base SCTL

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2021-02-08 09:09:36+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type Fermier struct {
	Id        int // IdExploitant de la base SCTL
	Nom       string
	Prenom    string
	Adresse   string
	Cp        string
	Ville     string
	Tel       string
	Email     string
	Parcelles []*Parcelle
}

// ************************** Nom *******************************

func (f *Fermier) String() string {
	if f.Prenom == "" {
		return f.Nom
	}
	return f.Prenom + " " + f.Nom
}

// Renvoie Nom + Prenom, adapté aux besoins de autocomplete
func (f *Fermier) NomAutocomplete() string {
	if f.Prenom == "" {
		return f.Nom
	}
	return f.Nom + " " + f.Prenom
}

// ************************** Divers *******************************

func CountFermiers(db *sqlx.DB) int {
	var count int
	_ = db.QueryRow("select count(*) from fermier").Scan(&count)
	return count
}

// ************************** Compute *******************************

// Calcule le champ Parcelles d'un fermier
func (f *Fermier) ComputeParcelles(db *sqlx.DB) error {
	if len(f.Parcelles) != 0 {
		return nil // déjà calculé
	}
	query := `select * from parcelle where id in(
	    select id_parcelle from parcelle_fermier where id_fermier=$1)
	`
	err := db.Select(&f.Parcelles, query, f.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

// ************************** Get one *******************************

// Renvoie un Fermier à partir de son id.
// Ne contient que les champs de la table fermier.
// Les autres champs ne sont pas remplis.
func GetFermier(db *sqlx.DB, id int) (*Fermier, error) {
	f := &Fermier{}
	query := "select * from fermier where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(f)
	if err != nil {
		return f, werr.Wrapf(err, "Erreur query : "+query)
	}
	return f, err
}

// ************************** Get many *******************************

// Renvoie une liste de Fermiers triés en utilisant un champ de la table
// @param field    Champ de la table fermier utilisé pour le tri
func GetSortedFermiers(db *sqlx.DB, field string) ([]*Fermier, error) {
	fermiers := []*Fermier{}
	query := "select * from fermier where id<>0 order by " + field
	err := db.Select(&fermiers, query)
	if err != nil {
		return fermiers, werr.Wrapf(err, "Erreur query : "+query)
	}
	return fermiers, err
}

// Renvoie des Fermiers à partir d'un lieu-dit.
// Utilise les parcelles pour faire le lien
// Ne contient que les champs de la table fermier.
// Les autres champs ne sont pas remplis.
// Utilisé par ajax
func GetFermiersFromLieudit(db *sqlx.DB, idLieudit int) ([]*Fermier, error) {
	fermiers := []*Fermier{}
	query := `
	    select * from fermier where id in(
            select distinct id_fermier from parcelle_fermier where id_parcelle in(
                select id_parcelle from parcelle_lieudit where id_lieudit=$1
            )
        ) order by nom`
	err := db.Select(&fermiers, query, idLieudit)
	return fermiers, err
}

// Renvoie des Fermiers à partir d'une UG.
// Utilise les parcelles pour faire le lien
// Ne contient que les champs de la table fermier.
// Les autres champs ne sont pas remplis.
// Utilisé par ajax
func GetFermiersFromCodeUG(db *sqlx.DB, codeUG string) ([]*Fermier, error) {
	fermiers := []*Fermier{}
	query := `
	    select * from fermier where id in(
            select distinct id_fermier from parcelle_fermier where id_parcelle in(
                select id_parcelle from parcelle_ug where id_ug in(
                    select id from ug where code=$1
                )
            )
        ) order by nom`
	err := db.Select(&fermiers, query, codeUG)
	return fermiers, err
}

// Renvoie des Fermiers à partir du début de leurs noms.
// Ne contient que les champs de la table fermier.
// Les autres champs ne sont pas remplis.
// Utilisé par ajax
func GetFermiersAutocomplete(db *sqlx.DB, str string) ([]*Fermier, error) {
	fermiers := []*Fermier{}
	query := "select * from fermier where nom ilike '" + str + "%'"
	err := db.Select(&fermiers, query)
	if err != nil {
		return fermiers, werr.Wrapf(err, "Erreur query : "+query)
	}
	return fermiers, nil
}
