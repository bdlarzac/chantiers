/*
*****************************************************************************

	Parcelles

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2019-11-07, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type Parcelle struct {
	Id             int // IdParcelle de la base SCTL
	IdProprietaire int `db:"id_proprietaire"`
	Code           string
	Surface        float64
	IdCommune      int `db:"id_commune"`
	// Pas en base
	Proprietaire *Acteur
	Commune      *Commune
	Lieudits     []*Lieudit
	Fermiers     []*Fermier
	UGs          []*UG
}

// Sert à afficher la liste des activités sur une parcelle.
// Contient les infos utilisées pour l'affichage, pas les activités.
// @todo supprimer si finalement pas utilisé
type ParcelleActivite struct {
	Date        time.Time
	URL         string // URL de la page de l'activité concernée
	NomActivite string
}

func (p *Parcelle) String() string {
	return p.Code
}

// ************************** Get one *******************************

/*
Renvoie une Parcelle à partir de son id.
*/
func GetParcelle(db *sqlx.DB, id int) (p *Parcelle, err error) {
	p = &Parcelle{}
	query := "select * from parcelle where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(p)
	if err != nil {
		return p, werr.Wrapf(err, "Erreur query : "+query)
	}
	return p, nil
}

/*
Renvoie une Parcelle à partir de son code et de l'id de la commune.
(id de la commune nécessaire car le code est unique au sein de la commune,
donc plusieurs parcelles avec le même code existent en base).
*/
func GetParcelleFromCodeAndCommuneId(db *sqlx.DB, codeParcelle string, idCommune int) (p *Parcelle, err error) {
	p = &Parcelle{}
	query := "select * from parcelle where code=$1 and id_commune=$2"
	row := db.QueryRowx(query, codeParcelle, strconv.Itoa(idCommune))
	err = row.StructScan(p)
	switch {
	case err == sql.ErrNoRows:
		return p, nil
	case err != nil:
		return p, werr.Wrapf(err, "Erreur query : "+query)
	}
	if err != nil {
		return p, werr.Wrapf(err, "Erreur query : "+query)
	}
	return p, nil
}

// ************************** Get many *******************************

/*
Utilisé par ajax
@param  idsUG  string, par ex : "12,432,35"
*/
func GetParcellesFromIdsUGs(db *sqlx.DB, idsUG string) (result []*Parcelle, err error) {
	query := `select * from parcelle where id in(select id_parcelle from parcelle_ug where id_ug in(` + idsUG + `)) order by code`
	err = db.Select(&result, query)
	return result, nil
}

// ************************** Compute *******************************

func (p *Parcelle) ComputeProprietaire(db *sqlx.DB) (err error) {
	p.Proprietaire, err = GetActeur(db, p.IdProprietaire)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (p *Parcelle) ComputeLieudits(db *sqlx.DB) (err error) {
	if len(p.Lieudits) != 0 {
		return nil // déjà calculé
	}
	query := `
	    select * from lieudit where id in(
            select id_lieudit from parcelle_lieudit where id_parcelle=$1
        ) order by nom`
	err = db.Select(&p.Lieudits, query, p.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func (p *Parcelle) ComputeCommune(db *sqlx.DB) (err error) {
	if p.Commune != nil {
		return nil // déjà calculé
	}
	p.Commune, err = GetCommune(db, p.IdCommune)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetCommune()")
	}
	return nil
}

func (p *Parcelle) ComputeFermiers(db *sqlx.DB) (err error) {
	if len(p.Fermiers) != 0 {
		return nil // déjà calculé
	}
	query := `
	    select * from fermier where id in(
	        select id_fermier from parcelle_fermier where id_parcelle=$1
	    ) order by nom`
	err = db.Select(&p.Fermiers, query, p.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func (p *Parcelle) ComputeUGs(db *sqlx.DB) (err error) {
	query := `
	    select * from ug where id in(
	        select id_ug from parcelle_ug where id_parcelle=$1
	    ) order by code`
	err = db.Select(&p.UGs, query, p.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
