/******************************************************************************
    Rangement de plaquettes dans un hangar à plaquettes
    = Déchargement, accompagne les transports de plaquettes
    Se fait avec un télescopique

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-01-22 02:19:38+01:00+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"time"

	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type PlaqRange struct {
	Id            int
	IdChantier    int `db:"id_chantier"`
	IdTas         int `db:"id_tas"`
	IdRangeur     int `db:"id_rangeur"`
	IdConducteur  int `db:"id_conducteur"`
	IdProprioutil int `db:"id_proprioutil"`
	DateRange     time.Time
	TypeCout      string // G (global) ou D (détail)
	// coût global
	GlPrix    float64 // HT
	GlTVA     float64
	GlDatePay time.Time
	// coût détail - conducteur
	CoPrixH   float64 // HT
	CoNheure  float64
	CoTVA     float64
	CoDatePay time.Time
	// coût détail - outil
	OuPrix    float64 // HT
	OuTVA     float64
	OuDatePay time.Time
	//
	Notes string
	// Pas stocké en base
	Chantier    *Plaq
	Tas         *Tas
	Rangeur     *Acteur
	Conducteur  *Acteur
	Proprioutil *Acteur
}

// ************************** Get *******************************

func GetPlaqRange(db *sqlx.DB, id int) (*PlaqRange, error) {
	pr := &PlaqRange{}
	query := "select * from plaqrange where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(pr)
	if err != nil {
		return pr, werr.Wrapf(err, "Erreur query : "+query)
	}
	return pr, nil
}

// ************************** Compute *******************************

func (pr *PlaqRange) ComputeTas(db *sqlx.DB) error {
	var err error
	pr.Tas, err = GetTasFull(db, pr.IdTas)
	return err
}

func (pr *PlaqRange) ComputeRangeur(db *sqlx.DB) error {
	if pr.IdRangeur == 0 {
		return nil
	}
	var err error
	pr.Rangeur, err = GetActeur(db, pr.IdRangeur)
	return err
}

func (pr *PlaqRange) ComputeConducteur(db *sqlx.DB) error {
	var err error
	pr.Conducteur, err = GetActeur(db, pr.IdConducteur)
	return err
}

func (pr *PlaqRange) ComputeProprioutil(db *sqlx.DB) error {
	if pr.IdProprioutil == 0 {
		return nil
	}
	var err error
	pr.Proprioutil, err = GetActeur(db, pr.IdProprioutil)
	return err
}

// ************************** CRUD *******************************

func InsertPlaqRange(db *sqlx.DB, pr *PlaqRange) (int, error) {
	query := `insert into plaqrange(
        id_chantier,
        id_tas,
        id_rangeur,
        id_conducteur,
        id_proprioutil,
        daterange,
        typecout,
        glprix,
        gltva,
        gldatepay,
        conheure,
        coprixh,
        cotva,
        codatepay,
        ouprix,
        outva,
        oudatepay,
        notes)
        values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		pr.IdChantier,
		pr.IdTas,
		pr.IdRangeur,
		pr.IdConducteur,
		pr.IdProprioutil,
		pr.DateRange,
		pr.TypeCout,
		pr.GlPrix,
		pr.GlTVA,
		pr.GlDatePay,
		pr.CoNheure,
		pr.CoPrixH,
		pr.CoTVA,
		pr.CoDatePay,
		pr.OuPrix,
		pr.OuTVA,
		pr.OuDatePay,
		pr.Notes).Scan(&id)
	return id, err
}

func UpdatePlaqRange(db *sqlx.DB, pr *PlaqRange) error {
	query := `update plaqrange set(
        id_chantier,
        id_tas,
        id_rangeur,
        id_conducteur,
        id_proprioutil,
        daterange,
        typecout,
        glprix,
        gltva,
        gldatepay,
        conheure,
        coprixh,              
        cotva,
        codatepay,
        ouprix,
        outva,
        oudatepay,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) where id=$19`
	_, err := db.Exec(
		query,
		pr.IdChantier,
		pr.IdTas,
		pr.IdRangeur,
		pr.IdConducteur,
		pr.IdProprioutil,
		pr.DateRange,
		pr.TypeCout,
		pr.GlPrix,
		pr.GlTVA,
		pr.GlDatePay,
		pr.CoNheure,
		pr.CoPrixH,
		pr.CoTVA,
		pr.CoDatePay,
		pr.OuPrix,
		pr.OuTVA,
		pr.OuDatePay,
		pr.Notes,
		pr.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeletePlaqRange(db *sqlx.DB, id int) error {
	query := "delete from plaqrange where id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
