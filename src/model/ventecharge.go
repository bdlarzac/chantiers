/******************************************************************************
    Chargement pour une livraison de plaquettes

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-03-13 11:10:29+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	//"strconv"
	"time"

	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	//"fmt"
)

type VenteCharge struct {
	Id          int
	IdLivraison int `db:"id_livraison"`
	IdChargeur  int `db:"id_chargeur"`
	IdTas       int `db:"id_tas"`
	Qte         float64
	DateCharge  time.Time
	TypeCout    string // G (global) ou D (détail)
	// coût global
	GlPrix    float64
	GlTVA     float64
	GlDatePay time.Time
	// coût détaillé, outil
	OuPrix    float64
	OuTVA     float64
	OuDatePay time.Time
	// coût détaillé, main d'oeuvre
	MoNHeure  float64
	MoPrixH   float64
	MoTVA     float64
	MoDatePay time.Time
	//
	Notes string
	// Pas stocké en base
	IdVente   int
	Livraison *VenteLivre
	Chargeur  *Acteur
	Tas       *Tas
}

// ************************** Nom *******************************

func (vc *VenteCharge) String() string {
	if vc.Chargeur == nil {
		panic("Erreur dans le code - Le chargeur d'un chargement doit être calculé avant d'appeler String()")
	}
	return vc.Chargeur.String() + " " + tiglib.DateFr(vc.DateCharge)
}

// ************************** Get *******************************

func GetVenteCharge(db *sqlx.DB, id int) (*VenteCharge, error) {
	vc := &VenteCharge{}
	query := "select * from ventecharge where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(vc)
	if err != nil {
		return vc, werr.Wrapf(err, "Erreur query : "+query)
	}
	return vc, nil
}

func GetVenteChargeFull(db *sqlx.DB, id int) (*VenteCharge, error) {
	vc, err := GetVenteCharge(db, id)
	if err != nil {
		return vc, werr.Wrapf(err, "Erreur appel GetVenteCharge()")
	}
	err = vc.ComputeChargeur(db)
	if err != nil {
		return vc, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeChargeur()")
	}
	err = vc.ComputeTas(db)
	if err != nil {
		return vc, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeTas()")
	}
	// @todo  aussi Livraison et IdVente ???
	return vc, err
}

// ************************** Compute *******************************

func (vc *VenteCharge) ComputeChargeur(db *sqlx.DB) error {
	if vc.Chargeur != nil {
		return nil
	}
	var err error
	vc.Chargeur, err = GetActeur(db, vc.IdChargeur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (vc *VenteCharge) ComputeLivraison(db *sqlx.DB) error {
	if vc.Livraison != nil {
		return nil
	}
	var err error
	vc.Livraison, err = GetVenteLivre(db, vc.IdLivraison)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetVenteLivre()")
	}
	return nil
}

func (vc *VenteCharge) ComputeTas(db *sqlx.DB) error {
	if vc.Tas != nil {
		return nil
	}
	var err error
	vc.Tas, err = GetTasFull(db, vc.IdTas)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetTasFull()")
	}
	return nil
}

func (vc *VenteCharge) ComputeIdVente(db *sqlx.DB) error {
	if vc.IdVente != 0 {
		return nil
	}
	var err error
	err = vc.ComputeLivraison(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel VenteCharge.ComputeLivraison()")
	}
	vc.IdVente = vc.Livraison.IdVente
	return nil
}

// ************************** CRUD *******************************

func InsertVenteCharge(db *sqlx.DB, vc *VenteCharge) (int, error) {
	query := `insert into ventecharge(
        id_livraison,
        id_chargeur,
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
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		vc.IdLivraison,
		vc.IdChargeur,
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
	return id, err
}

func UpdateVenteCharge(db *sqlx.DB, vc *VenteCharge) error {
	query := `update ventecharge set(
        id_livraison,
        id_chargeur,
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
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17) where id=$18`
	_, err := db.Exec(
		query,
		vc.IdLivraison,
		vc.IdChargeur,
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

func DeleteVenteCharge(db *sqlx.DB, id int) error {
	query := "delete from ventecharge where id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}