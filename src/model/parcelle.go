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
	//"fmt"
)

type Parcelle struct {
	Id             int
	IdProprietaire int `db:"id_proprietaire"`
	Code           string
	Surface        float32
	// Pas en base
	Proprietaire *Acteur
	Lieudits     []*Lieudit
	Exploitants  []*Acteur
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
	query := "select * from parcelle where id in (select id_parcelle from parcelle_lieudit where id_lieudit=$1) order by code"
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

// Remplit le champ Exploitants d'une parcelle
func (p *Parcelle) ComputeExploitants(db *sqlx.DB) error {
	query := "select id_sctl_exploitant from parcelle_exploitant where id_parcelle=$1"
	rows, err := db.Query(query, p.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	defer rows.Close()
	var idE int
	for rows.Next() {
		if err := rows.Scan(&idE); err != nil {
			return werr.Wrapf(err, "Erreur rows.Scan sur query : "+query)
		}
		exploitant, err := GetActeurByIdSctl(db, idE)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetActeurByIdSctl()")
		}
		p.Exploitants = append(p.Exploitants, exploitant)
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

// ************************** Activité *******************************

// Renvoie les activités ayant eu lieu sur une parcelle.
// Ordre chronologique inverse
// Ne renvoie que des infos pour afficher la liste, pas les activités réelles.
// @todo A FINIR - faire ugs d'abord - CONFIRMER SI DOIT ETRE FAIT
// A PRIORI PAS BESOIN DE LE FAIRE
func (p *Parcelle) GetActivitesByDate(db *sqlx.DB) ([]*ParcelleActivite, error) {
	res := []*ParcelleActivite{}
	var err error
	var query string
	//
	// Chantiers plaquettes
	//
	list1 := []Plaq{}
	query = "select * from plaq where id_ug in(select id_ug from parcelle_ug where id_parcelle=$1)"
	err = db.Select(&list1, query, p.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	return res, nil
}