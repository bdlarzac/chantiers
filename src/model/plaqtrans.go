/******************************************************************************
    Transport de plaquettes, depuis un chantier vers un lieu de stockage
    Soit par camion, soit par tracteur + benne

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-01-15 02:37:25+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"time"
)

type PlaqTrans struct {
	Id             int
	IdChantier     int `db:"id_chantier"`
	IdTas          int `db:"id_tas"`
	IdTransporteur int `db:"id_transporteur"`
	DateTrans      time.Time
	Qte            float64
	TypeCout       string // G (global) C (camion) ou T (tracteur)
	// Coût global
	GlPrix    float64 // prix total
	GlTVA     float64
	GlDatePay time.Time
	// Concerne le transporteur
	TrNheure  float64
	TrPrixH   float64 // prix ht / heure
	TrTVA     float64
	TrDatePay time.Time
	// Transport camion
	CaNkm     float64
	CaPrixKm  float64 // prix ht / km
	CaTVA     float64
	CaDatePay time.Time
	// Transport tracteur + benne
	TbNbenne  int
	TbDuree   float64
	TbPrixH   float64 // prix ht / heure
	TbTVA     float64
	TbDatePay time.Time
	Notes     string
	// Pas stocké en base
	Chantier     *Plaq
	Tas          *Tas
	Transporteur *Acteur
}

// ************************** Get *******************************

func GetPlaqTrans(db *sqlx.DB, id int) (*PlaqTrans, error) {
	pt := &PlaqTrans{}
	query := "select * from plaqtrans where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(pt)
	if err != nil {
		return pt, werr.Wrapf(err, "Erreur query : "+query)
	}
	return pt, nil
}

// ************************** Compute *******************************

func (pt *PlaqTrans) ComputeTas(db *sqlx.DB) error {
	var err error
	pt.Tas, err = GetTasFull(db, pt.IdTas)
	return err
}

func (pt *PlaqTrans) ComputeTransporteur(db *sqlx.DB) error {
	var err error
	pt.Transporteur, err = GetActeur(db, pt.IdTransporteur)
	return err
}

// ************************** CRUD *******************************

func InsertPlaqTrans(db *sqlx.DB, pt *PlaqTrans) (int, error) {
	query := `insert into plaqtrans(
        id_chantier,
        id_tas,
        id_transporteur,
        datetrans,
        qte,
        typecout,
        glprix,
        gltva,
        gldatepay,
        trnheure,
        trprixh,
        trtva,
        trdatepay,
        cankm,
        caprixkm,
        catva,
        cadatepay,
        tbnbenne,
        tbduree,
        tbprixh,
        tbtva,
        tbdatepay,
        notes)
        values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		pt.IdChantier,
		pt.IdTas,
		pt.IdTransporteur,
		pt.DateTrans,
		pt.Qte,
		pt.TypeCout,
		pt.GlPrix,
		pt.GlTVA,
		pt.GlDatePay,
		pt.TrNheure,
		pt.TrPrixH,
		pt.TrTVA,
		pt.TrDatePay,
		pt.CaNkm,
		pt.CaPrixKm,
		pt.CaTVA,
		pt.CaDatePay,
		pt.TbNbenne,
		pt.TbDuree,
		pt.TbPrixH,
		pt.TbTVA,
		pt.TbDatePay,
		pt.Notes).Scan(&id)
	return id, err
}

func UpdatePlaqTrans(db *sqlx.DB, pt *PlaqTrans) error {
	query := `update plaqtrans set(
        id_chantier,
        id_tas,
        id_transporteur,
        datetrans,
        qte,
        typecout,
        glprix,
        gltva,
        gldatepay,
        trnheure,
        trprixh,
        trtva,
        trdatepay,
        cankm,
        caprixkm,
        catva,
        cadatepay,
        tbnbenne,
        tbduree,
        tbprixh,
        tbtva,
        tbdatepay,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23) where id=$24`
	_, err := db.Exec(
		query,
		pt.IdChantier,
		pt.IdTas,
		pt.IdTransporteur,
		pt.DateTrans,
		pt.Qte,
		pt.TypeCout,
		pt.GlPrix,
		pt.GlTVA,
		pt.GlDatePay,
		pt.TrNheure,
		pt.TrPrixH,
		pt.TrTVA,
		pt.TrDatePay,
		pt.CaNkm,
		pt.CaPrixKm,
		pt.CaTVA,
		pt.CaDatePay,
		pt.TbNbenne,
		pt.TbDuree,
		pt.TbPrixH,
		pt.TbTVA,
		pt.TbDatePay,
		pt.Notes,
		pt.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeletePlaqTrans(db *sqlx.DB, id int) error {
	query := "delete from plaqtrans where id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
