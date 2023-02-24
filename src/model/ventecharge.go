/*
*****************************************************************************

	Chargement pour une livraison de plaquettes

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2020-03-13 11:10:29+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"time"
	//	"strconv"

	//	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	//"fmt"
)

type VenteCharge struct {
	Id            int
	IdLivraison   int `db:"id_livraison"`
	IdChargeur    int `db:"id_chargeur"`
	IdConducteur  int `db:"id_conducteur"`
	IdProprioutil int `db:"id_proprioutil"`
	IdTas         int `db:"id_tas"`
	Qte           float64
	DateCharge    time.Time
	TypeCout      string // G (global) ou D (détail)
	// coût global
	GlPrix    float64
	GlTVA     float64
	GlDatePay time.Time
	// coût détaillé, main d'oeuvre
	MoNHeure  float64
	MoPrixH   float64
	MoTVA     float64
	MoDatePay time.Time
	// coût détaillé, outil
	OuPrix    float64
	OuTVA     float64
	OuDatePay time.Time
	//
	Notes string
	// Pas stocké en base
	IdVente     int
	Livraison   *VenteLivre
	Chargeur    *Acteur
	Conducteur  *Acteur
	Proprioutil *Acteur
	Tas         *Tas
}

// ************************** Get *******************************

func GetVenteCharge(db *sqlx.DB, id int) (vc *VenteCharge, err error) {
	vc = &VenteCharge{}
	query := "select * from ventecharge where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(vc)
	if err != nil {
		return vc, werr.Wrapf(err, "Erreur query : "+query)
	}
	return vc, nil
}

func GetVenteChargeFull(db *sqlx.DB, id int) (vc *VenteCharge, err error) {
	vc, err = GetVenteCharge(db, id)
	if err != nil {
		return vc, werr.Wrapf(err, "Erreur appel GetVenteCharge()")
	}
	err = vc.ComputeChargeur(db)
	if err != nil {
		return vc, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeChargeur()")
	}
	err = vc.ComputeConducteur(db)
	if err != nil {
		return vc, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeConducteur()")
	}
	err = vc.ComputeProprioutil(db)
	if err != nil {
		return vc, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeProprioutil()")
	}
	err = vc.ComputeTas(db)
	if err != nil {
		return vc, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeTas()")
	}
	return vc, nil
}

// ************************** Compute *******************************

func (vc *VenteCharge) ComputeChargeur(db *sqlx.DB) (err error) {
	if vc.IdChargeur == 0 {
		return nil // pas de chargeur (mais conducteur et proprioutil)
	}
	if vc.Chargeur != nil {
		return nil // déjà calculé
	}
	vc.Chargeur, err = GetActeur(db, vc.IdChargeur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (vc *VenteCharge) ComputeConducteur(db *sqlx.DB) (err error) {
	if vc.IdConducteur == 0 {
		return nil // pas de conducteur ni proprioutil (mais un chargeur)
	}
	if vc.Conducteur != nil {
		return nil // déjà calculé
	}
	vc.Conducteur, err = GetActeur(db, vc.IdConducteur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (vc *VenteCharge) ComputeProprioutil(db *sqlx.DB) (err error) {
	if vc.IdProprioutil == 0 {
		return nil // pas de conducteur ni proprioutil (mais un chargeur)
	}
	if vc.Proprioutil != nil {
		return nil // déjà calculé
	}
	vc.Proprioutil, err = GetActeur(db, vc.IdProprioutil)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (vc *VenteCharge) ComputeLivraison(db *sqlx.DB) (err error) {
	if vc.Livraison != nil {
		return nil // déjà calculé
	}
	vc.Livraison, err = GetVenteLivre(db, vc.IdLivraison)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetVenteLivre()")
	}
	return nil
}

func (vc *VenteCharge) ComputeTas(db *sqlx.DB) (err error) {
	if vc.Tas != nil {
		return nil // déjà calculé
	}
	vc.Tas, err = GetTasFull(db, vc.IdTas)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetTasFull()")
	}
	return nil
}

func (vc *VenteCharge) ComputeIdVente(db *sqlx.DB) (err error) {
	if vc.IdVente != 0 {
		return nil // déjà calculé
	}
	err = vc.ComputeLivraison(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel VenteCharge.ComputeLivraison()")
	}
	vc.IdVente = vc.Livraison.IdVente
	return nil
}

// ************************** CRUD *******************************

func InsertVenteCharge(db *sqlx.DB, vc *VenteCharge) (id int, err error) {
	query := `insert into ventecharge(
        id_livraison,
        id_chargeur,
        id_conducteur,
        id_proprioutil,
        id_tas,
        qte,
        datecharge,
        typecout,
        glprix,
        gltva,
        gldatepay,
        ouprix,
        outva,
        oudatepay,
        monheure,
        moprixh,
        motva,
        modatepay,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) returning id`
	err = db.QueryRow(
		query,
		vc.IdLivraison,
		vc.IdChargeur,
		vc.IdConducteur,
		vc.IdProprioutil,
		vc.IdTas,
		vc.Qte,
		vc.DateCharge,
		vc.TypeCout,
		vc.GlPrix,
		vc.GlTVA,
		vc.GlDatePay,
		vc.OuPrix,
		vc.OuTVA,
		vc.OuDatePay,
		vc.MoNHeure,
		vc.MoPrixH,
		vc.MoTVA,
		vc.MoDatePay,
		vc.Notes).Scan(&id)
	if err != nil {
		return 0, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

func UpdateVenteCharge(db *sqlx.DB, vc *VenteCharge) (err error) {
	query := `update ventecharge set(
        id_livraison,
        id_chargeur,
        id_conducteur,
        id_proprioutil,
        id_tas,
        qte,
        datecharge,
        typecout,
        glprix,
        gltva,
        gldatepay,
        ouprix,
        outva,              
        oudatepay,
        monheure,
        moprixh,
        motva,
        modatepay,
        notes                          
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) where id=$20`
	_, err = db.Exec(
		query,
		vc.IdLivraison,
		vc.IdChargeur,
		vc.IdConducteur,
		vc.IdProprioutil,
		vc.IdTas,
		vc.Qte,
		vc.DateCharge,
		vc.TypeCout,
		vc.GlPrix,
		vc.GlTVA,
		vc.GlDatePay,
		vc.OuPrix,
		vc.OuTVA,
		vc.OuDatePay,
		vc.MoNHeure,
		vc.MoPrixH,
		vc.MoTVA,
		vc.MoDatePay,
		vc.Notes,
		vc.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeleteVenteCharge(db *sqlx.DB, id int) (err error) {
	// rétablit le stock du tas concerné par le chargement
	// avant de supprimer le chargement
	vc, err := GetVenteCharge(db, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetVenteCharge()")
	}
	err = vc.ComputeTas(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ComputeTas()")
	}
	err = vc.Tas.ModifierStock(db, vc.Qte) // Ajoute des plaquettes au tas
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ModifierStock()")
	}
	// delete le chargement
	query := "delete from ventecharge where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
