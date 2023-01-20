/******************************************************************************
    Parcelles

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2019-11-07, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"time"

	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type Parcelle struct {
	Id             int // IdParcelle de la base SCTL
	IdProprietaire int `db:"id_proprietaire"`
	Code           string
	Surface        float32
	IdCommune      int `db:"id_commune"`
	// Pas en base
	Proprietaire *Acteur
	Lieudits     []*Lieudit
	Communes     []*Commune
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

// ************************** Get *******************************

// Renvoie une Parcelle contenant Id, Nom et surface.
// Les autres champs ne sont pas remplis.
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

// ************************** Compute *******************************

// Remplit le champ Proprietaire d'une parcelle
func (p *Parcelle) ComputeProprietaire(db *sqlx.DB) (err error) {
	p.Proprietaire, err = GetActeur(db, p.IdProprietaire)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

// Remplit le champ Lieudits d'une parcelle
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

// Remplit le champ Communes d'une parcelle
func (p *Parcelle) ComputeCommunes(db *sqlx.DB) (err error) {
	if len(p.Communes) != 0 {
		return nil // déjà calculé
	}
	query := `
	    select * from commune where id=$1`
	err = db.Select(&p.Communes, query, p.IdCommune)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

// Remplit le champ Fermiers d'une parcelle
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

// Remplit le champ UGs d'une parcelle
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
