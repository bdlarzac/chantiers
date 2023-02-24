/*
*****************************************************************************

	Fermiers = exploitants agricoles, issus de la base SCTL

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2021-02-08 09:09:36+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Fermier struct {
	Id      int // IdExploitant de la base SCTL
	Nom     string
	Prenom  string
	Adresse string
	Cp      string
	Ville   string
	Tel     string
	Email   string
	// Pas stocké en base
	Parcelles []*Parcelle
}

// ************************** Nom *******************************

func (f *Fermier) String() string {
	return strings.TrimSpace(f.Prenom + " " + f.Nom)
}

// ************************** Divers *******************************

func CountFermiers(db *sqlx.DB) (count int) {
	_ = db.QueryRow("select count(*) from fermier").Scan(&count)
	return count
}

// ************************** Compute *******************************

/*
*

	Calcule le champ Parcelles d'un fermier

*
*/
func (f *Fermier) ComputeParcelles(db *sqlx.DB) (err error) {
	if len(f.Parcelles) != 0 {
		return nil // déjà calculé
	}
	query := `select * from parcelle where id in(
            select id_parcelle from parcelle_fermier where id_fermier=$1
	    ) order by code`
	err = db.Select(&f.Parcelles, query, f.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

// ************************** Get one *******************************

/*
*

	Renvoie un Fermier à partir de son id.
	Ne contient que les champs de la table fermier.
	Les autres champs ne sont pas remplis.

*
*/
func GetFermier(db *sqlx.DB, id int) (f *Fermier, err error) {
	f = &Fermier{}
	query := "select * from fermier where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(f)
	if err != nil {
		return f, werr.Wrapf(err, "Erreur query : "+query)
	}
	return f, nil
}

// ************************** Get many *******************************

/*
*

	Renvoie une liste de Fermiers triés en utilisant un champ de la table
	@param field    Champ de la table fermier utilisé pour le tri

*
*/
func GetSortedFermiers(db *sqlx.DB, field string) (fermiers []*Fermier, err error) {
	fermiers = []*Fermier{}
	query := "select * from fermier where id<>0 order by " + field
	err = db.Select(&fermiers, query)
	if err != nil {
		return fermiers, werr.Wrapf(err, "Erreur query : "+query)
	}
	return fermiers, nil
}

/*
// TODO supprimer si toujours inutile
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
	if err != nil {
		return fermiers, werr.Wrapf(err, "Erreur query : "+query)
	}
	return fermiers, nil
}
*/
/**
    Renvoie des fermiers à partir d'un id UG.
    Contient les champs de la table fermier.
    Les autres champs ne sont pas remplis.
    @param      strIdsUGs   Chaîne contenant les ids séparés par des virgules. ex : "1, 34, 87"
**/
func GetFermiersFromIdsUGs(db *sqlx.DB, strIdsUGs string) (fermiers []*Fermier, err error) {
	fermiers = []*Fermier{}
	query := `
	    select * from fermier where id in(
            select id_fermier from parcelle_fermier where id_parcelle in(
                select id_parcelle from parcelle_ug where id_ug in(` + strIdsUGs + `)
            )
        )`
	err = db.Select(&fermiers, query)
	if err != nil {
		return fermiers, werr.Wrapf(err, "Erreur query : "+query)
	}
	return fermiers, nil
}

/*
*

	Renvoie des Fermiers à partir d'une UG.
	Utilise les parcelles pour faire le lien
	Ne contient que les champs de la table fermier.
	Les autres champs ne sont pas remplis.
	Utilisé par ajax

*
*/
func GetFermiersFromIdUG(db *sqlx.DB, idUG int) ([]*Fermier, error) {
	fermiers := []*Fermier{}
	query := `
	    select * from fermier where id in(
            select distinct id_fermier from parcelle_fermier where id_parcelle in(
                select id_parcelle from parcelle_ug where id_ug=$1
            )
        ) order by nom`
	err := db.Select(&fermiers, query, idUG)
	if err != nil {
		return fermiers, werr.Wrapf(err, "Erreur query : "+query)
	}
	return fermiers, nil
}

// ************************** CRUD *******************************

/*
*

	Pas un insert habituel car id est fourni

*
*/
func InsertFermier(db *sqlx.DB, f *Fermier) (err error) {
	query := `insert into fermier(
	    id,
	    nom,
	    prenom,
	    adresse,
	    cp,
	    ville,
	    tel,
	    email
    ) values($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err = db.Exec(
		query,
		f.Id,
		f.Nom,
		f.Prenom,
		f.Adresse,
		f.Cp,
		f.Ville,
		f.Tel,
		f.Email)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func UpdateFermier(db *sqlx.DB, f *Fermier) (err error) {
	query := `update fermier set(
	    nom,
	    prenom,
	    adresse,
	    cp,
	    ville,
	    tel,
	    email
    ) = ($1,$2,$3,$4,$5,$6,$7) where id=$8`
	_, err = db.Exec(
		query,
		f.Nom,
		f.Prenom,
		f.Adresse,
		f.Cp,
		f.Ville,
		f.Tel,
		f.Email,
		f.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
