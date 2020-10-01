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
	IdConducteur   int `db:"id_conducteur"`
	IdProprioutil  int `db:"id_proprioutil"`
	DateTrans      time.Time
	Qte            float64
	TypeCout       string // G (global) C (camion) ou T (tracteur)
	// Coût global
	GlPrix    float64 // prix total
	GlTVA     float64
	GlDatePay time.Time
	// Concerne le conducteur
	ConNheure  float64
	ConPrixH   float64 // prix ht / heure
	ConTVA     float64
	ConDatePay time.Time
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
	Conducteur   *Acteur
	Proprioutil  *Acteur
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
    if pt.IdTransporteur == 0{
        return nil
    }
	var err error
	pt.Transporteur, err = GetActeur(db, pt.IdTransporteur)
	return err
}

func (pt *PlaqTrans) ComputeConducteur(db *sqlx.DB) error {
    if pt.IdConducteur == 0{
        return nil
    }
	var err error
	pt.Conducteur, err = GetActeur(db, pt.IdConducteur)
	return err
}

func (pt *PlaqTrans) ComputeProprioutil(db *sqlx.DB) error {
    if pt.IdProprioutil == 0{
        return nil
    }
	var err error
	pt.Proprioutil, err = GetActeur(db, pt.IdProprioutil)
	return err
}

// ************************** CRUD *******************************

func InsertPlaqTrans(db *sqlx.DB, pt *PlaqTrans) (int, error) {
	// Mise à jour du stock du tas
	err := pt.ComputeTas(db)
	if err != nil {
		return 0, err
	}
	err = pt.Tas.ModifierStock(db, pt.Qte) // Ajoute plaquettes au tas
	if err != nil {
		return 0, err
	}
	query := `insert into plaqtrans(
        id_chantier,
        id_tas,
        id_transporteur,
        id_conducteur,
        id_proprioutil,
        datetrans,
        qte,
        typecout,
        glprix,
        gltva,
        gldatepay,
        connheure,
        conprixh,
        contva,
        condatepay,
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
        values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25) returning id`
	id := int(0)
	err = db.QueryRow(
		query,
		pt.IdChantier,
		pt.IdTas,
		pt.IdTransporteur,
		pt.IdConducteur,
		pt.IdProprioutil,
		pt.DateTrans,
		pt.Qte,
		pt.TypeCout,
		pt.GlPrix,
		pt.GlTVA,
		pt.GlDatePay,
		pt.ConNheure,
		pt.ConPrixH,
		pt.ConTVA,
		pt.ConDatePay,
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
	// Mise à jour du stock du tas
	// Enlève la qté du transport avant update transport
	// puis ajoute qté après update transport
	// Attention, le tas avant update n'est pas forcément le même que le tas après update
	// (cas où plusieurs tas pour un chantier plaquette et changement de tas lors de update transport)
	ptAvant, err := GetPlaqTrans(db, pt.Id)
	if err != nil {
		return err
	}
	err = ptAvant.ComputeTas(db)
	if err != nil {
		return err
	}
	err = ptAvant.Tas.ModifierStock(db, -ptAvant.Qte) // Retire des plaquettes au tas
	if err != nil {
		return err
	}
	//
	err = pt.ComputeTas(db)
	if err != nil {
		return err
	}
	err = pt.Tas.ModifierStock(db, pt.Qte) // Ajoute des plaquettes au tas
	if err != nil {
		return err
	}
	// clés étrangères pouvant être nulles
	/* 
	var idTransporteur interface{}
	if pt.IdTransporteur != 0{
	    idTransporteur = pt.IdTransporteur
	}
	var idConducteur interface{}
	if pt.IdConducteur != 0{
	    idConducteur = pt.IdConducteur
	}
	var idProprioutil interface{}
	if pt.IdProprioutil != 0{
	    idProprioutil = pt.IdProprioutil
	}
	*/
	//
	query := `update plaqtrans set(
        id_chantier,
        id_tas,
        id_transporteur,
        id_conducteur,
        id_proprioutil,
        datetrans,
        qte,
        typecout,
        glprix,
        gltva,
        gldatepay,
        connheure,
        conprixh,
        contva,
        condatepay,
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
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25) where id=$26`
	_, err = db.Exec(
		query,
		pt.IdChantier,
		pt.IdTas,
		pt.IdTransporteur,
		pt.IdConducteur,
		pt.IdProprioutil,
		pt.DateTrans,
		pt.Qte,
		pt.TypeCout,
		pt.GlPrix,
		pt.GlTVA,
		pt.GlDatePay,
		pt.ConNheure,
		pt.ConPrixH,
		pt.ConTVA,
		pt.ConDatePay,
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
	// Enlève le stock du tas concerné par le transport
	// avant de supprimer le transport
	pt, err := GetPlaqTrans(db, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetPlaqTrans()")
	}
	err = pt.ComputeTas(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ComputeTas()")
	}
	err = pt.Tas.ModifierStock(db, -pt.Qte) // Retire des plaquettes au tas
	if err != nil {
		return err
	}
	// delete le transport
	query := "delete from plaqtrans where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
