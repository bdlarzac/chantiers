/*
Transport de plaquettes, depuis un chantier vers un lieu de stockage
Soit par camion, soit par tracteur + benne

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2020-01-15 02:37:25+01:00, Thierry Graff : Creation
*/
package model

import (
	"time"

	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
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
	PourcentPerte  float64 // différence bois sec / bois vert
	TypeCout       string  // G (global) C (camion) ou T (tracteur)
	// Coût global
	GlPrix    float64 // prix total
	GlTVA     float64
	GlDatePay time.Time
	// Concerne le conducteur
	CoNheure  float64
	CoPrixH   float64 // prix ht / heure
	CoTVA     float64
	CoDatePay time.Time
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

func GetPlaqTrans(db *sqlx.DB, id int) (pt *PlaqTrans, err error) {
	pt = &PlaqTrans{}
	query := "select * from plaqtrans where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(pt)
	if err != nil {
		return pt, werr.Wrapf(err, "Erreur query : "+query)
	}
	return pt, nil
}

// ************************** Compute *******************************

func (pt *PlaqTrans) ComputeTas(db *sqlx.DB) (err error) {
	if pt.Tas != nil {
		return nil // déjà calculé
	}
	pt.Tas, err = GetTasFull(db, pt.IdTas)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetTasFull()")
	}
	return nil
}

func (pt *PlaqTrans) ComputeTransporteur(db *sqlx.DB) (err error) {
	if pt.IdTransporteur == 0 {
		return nil // pas de transporteur (mais conducteur et proprioutil)
	}
	if pt.Transporteur != nil {
		return nil // déjà calculé
	}
	pt.Transporteur, err = GetActeur(db, pt.IdTransporteur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (pt *PlaqTrans) ComputeConducteur(db *sqlx.DB) (err error) {
	if pt.IdConducteur == 0 {
		return nil // pas de conducteur ni proprioutil (mais un transporteur)
	}
	if pt.Conducteur != nil {
		return nil // déjà calculé
	}
	pt.Conducteur, err = GetActeur(db, pt.IdConducteur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (pt *PlaqTrans) ComputeProprioutil(db *sqlx.DB) (err error) {
	if pt.IdProprioutil == 0 {
		return nil // pas de conducteur ni proprioutil (mais un transporteur)
	}
	if pt.Proprioutil != nil {
		return nil // déjà calculé
	}
	pt.Proprioutil, err = GetActeur(db, pt.IdProprioutil)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

// ************************** CRUD *******************************

func InsertPlaqTrans(db *sqlx.DB, pt *PlaqTrans) (id int, err error) {
	// Mise à jour du stock du tas
	err = pt.ComputeTas(db)
	if err != nil {
		return 0, werr.Wrapf(err, "Erreur appel ComputeTas()")
	}
	// Ajoute plaquettes au tas
	// Attention, transport en map vert et tas en map sec
	// => qté pour le tas = qté du transport - pourcentage de perte
	err = pt.Tas.ModifierStock(db, Vert2sec(pt.Qte, pt.PourcentPerte))
	if err != nil {
		return 0, werr.Wrapf(err, "Erreur appel PlaqTrans.Tas.ModifierStock()")
	}
	query := `insert into plaqtrans(
        id_chantier,
        id_tas,
        id_transporteur,
        id_conducteur,
        id_proprioutil,
        datetrans,
        qte,
        pourcentperte,
        typecout,
        glprix,
        gltva,
        gldatepay,
        conheure,
        coprixh,
        cotva,
        codatepay,
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
        values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26) returning id`
	err = db.QueryRow(
		query,
		pt.IdChantier,
		pt.IdTas,
		pt.IdTransporteur,
		pt.IdConducteur,
		pt.IdProprioutil,
		pt.DateTrans,
		pt.Qte,
		pt.PourcentPerte,
		pt.TypeCout,
		pt.GlPrix,
		pt.GlTVA,
		pt.GlDatePay,
		pt.CoNheure,
		pt.CoPrixH,
		pt.CoTVA,
		pt.CoDatePay,
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
	if err != nil {
		return 0, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

func UpdatePlaqTrans(db *sqlx.DB, pt *PlaqTrans) (err error) {
	// Mise à jour du stock du tas
	// Enlève la qté du transport avant update transport
	// puis ajoute qté après update transport
	// Attention, le tas avant update n'est pas forcément le même que le tas après update
	// (cas où plusieurs tas pour un chantier plaquette et changement de tas lors de update transport)
	ptAvant, err := GetPlaqTrans(db, pt.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetPlaqTrans()")
	}
	err = ptAvant.ComputeTas(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ptAvant.ComputeTas()")
	}
	err = ptAvant.Tas.ModifierStock(db, -Vert2sec(ptAvant.Qte, ptAvant.PourcentPerte)) // Retire des plaquettes au tas
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ptAvant.ModifierStock()")
	}
	//
	err = pt.ComputeTas(db)
	if err != nil {
		return err
	}
	err = pt.Tas.ModifierStock(db, Vert2sec(pt.Qte, pt.PourcentPerte)) // Ajoute des plaquettes au tas
	if err != nil {
		return werr.Wrapf(err, "Erreur appel PlaqTrans.ComputeTas()")
	}
	//
	query := `update plaqtrans set(
        id_chantier,
        id_tas,
        id_transporteur,
        id_conducteur,
        id_proprioutil,
        datetrans,
        qte,
        pourcentperte,
        typecout,
        glprix,
        gltva,
        gldatepay,
        conheure,
        coprixh,
        cotva,
        codatepay,
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
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26) where id=$27`
	_, err = db.Exec(
		query,
		pt.IdChantier,
		pt.IdTas,
		pt.IdTransporteur,
		pt.IdConducteur,
		pt.IdProprioutil,
		pt.DateTrans,
		pt.Qte,
		pt.PourcentPerte,
		pt.TypeCout,
		pt.GlPrix,
		pt.GlTVA,
		pt.GlDatePay,
		pt.CoNheure,
		pt.CoPrixH,
		pt.CoTVA,
		pt.CoDatePay,
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

func DeletePlaqTrans(db *sqlx.DB, id int) (err error) {
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
	err = pt.Tas.ModifierStock(db, -Vert2sec(pt.Qte, pt.PourcentPerte)) // Retire des plaquettes au tas
	if err != nil {
		return werr.Wrapf(err, "Erreur appel PlaqTrans.Tas.ModifierStock()")
	}
	// delete le transport
	query := "delete from plaqtrans where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
