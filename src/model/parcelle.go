/******************************************************************************
    Parcelles

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-11-07, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"time"

	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type Parcelle struct {
	Id             int
	IdProprietaire int `db:"id_proprietaire"`
	Code           string
	Surface        float32
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
func GetParcelle(db *sqlx.DB, id int) (*Parcelle, error) {
	p := &Parcelle{}
	query := "select * from parcelle where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(p)
	if err != nil {
		return p, werr.Wrapf(err, "Erreur query : "+query)
	}
	return p, err
}

// Renvoie des Parcelles à partir d'un lieu-dit.
// Ne contient que les champs de la table parcelle.
// Les autres champs ne sont pas remplis.
func GetParcellesFromLieudit(db *sqlx.DB, idLieudit int) ([]*Parcelle, error) {
	parcelles := []*Parcelle{}
	query := `select * from parcelle where id in(
	            select id_parcelle from parcelle_lieudit where id_lieudit=$1
	        ) order by code`
	err := db.Select(&parcelles, query, idLieudit)
	if err != nil {
		return parcelles, werr.Wrapf(err, "Erreur query : "+query)
	}
	return parcelles, nil
}

// ************************** Compute *******************************

// Remplit le champ Proprietaire d'une parcelle
func (p *Parcelle) ComputeProprietaire(db *sqlx.DB) error {
	var err error
	p.Proprietaire, err = GetActeur(db, p.IdProprietaire)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

// Remplit le champ Lieudits d'une parcelle
func (p *Parcelle) ComputeLieudits(db *sqlx.DB) error {
	if len(p.Lieudits) != 0 {
		return nil // déjà calculé
	}
	query := "select id_lieudit from parcelle_lieudit where id_parcelle=$1"
	rows, err := db.Query(query, p.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	defer rows.Close()
	var idLD int
	for rows.Next() {
		if err := rows.Scan(&idLD); err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
		lieudit, err := GetLieudit(db, idLD)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetLieudit()")
		}
		p.Lieudits = append(p.Lieudits, lieudit)
	}
	return nil
}

// Remplit le champ Communes d'une parcelle
func (p *Parcelle) ComputeCommunes(db *sqlx.DB) error {
	if len(p.Communes) != 0 {
		return nil // déjà calculé
	}
	query := `select * from commune where id in(
                select id_commune from commune_lieudit where id_lieudit in(
                    select id_lieudit from parcelle_lieudit where id_parcelle=$1
                )
            )`
	communes := []*Commune{}
	err := db.Select(&communes, query, p.Id)
	p.Communes = communes
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

// Remplit le champ Fermiers d'une parcelle
func (p *Parcelle) ComputeFermiers(db *sqlx.DB) error {
	if len(p.Fermiers) != 0 {
		return nil // déjà calculé
	}
	query := `select * from fermier where id in(
	        select id_fermier from parcelle_fermier where id_parcelle=$1
	    )`
	fermiers := []*Fermier{}
	err := db.Select(&fermiers, query, p.Id)
	p.Fermiers = fermiers
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

// Remplit le champ UGs d'une parcelle
func (p *Parcelle) ComputeUGs(db *sqlx.DB) error {
	query := "select id_ug from parcelle_ug where id_parcelle=$1"
	rows, err := db.Query(query, p.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	defer rows.Close()
	var idUG int
	for rows.Next() {
		if err := rows.Scan(&idUG); err != nil {
			return werr.Wrapf(err, "Erreur rows.Scan sur query : "+query)
		}
		ug, err := GetUG(db, idUG)
		if err != nil {
			return err
		}
		p.UGs = append(p.UGs, ug)
	}
	return nil
}
