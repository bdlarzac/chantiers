/******************************************************************************
    Livraison lors d'une vente de plaquettes

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-03-13 10:59:50+01:00, Thierry Graff : Creation
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

type VenteLivre struct {
	Id        int
	IdVente   int `db:"id_vente"`
	IdLivreur int `db:"id_livreur"`
	DateLivre time.Time
	TypeCout  string // G (global) ou D (détail)
	// coût global
	GlPrix float64
	GlTVA  float64
	// coût main d'oeuvre
	MoNHeure float64
	MoPrixH  float64
	MoTVA    float64
	//
	DatePay time.Time
	Notes   string
	// Pas stocké en base
	Quantite    float64
	Livreur     *Acteur
	Vente       *VentePlaq
	Chargements []*VenteCharge
}

// ************************** Nom *******************************

func (vl *VenteLivre) String() string {
	if vl.Livreur == nil {
		panic("Erreur dans le code - Le livreur d'une livraison plaquettes doit être calculé avant d'appeler String()")
	}
	return vl.Livreur.String() + " " + tiglib.DateFr(vl.DateLivre)
}

// ************************** Get *******************************

func GetVenteLivre(db *sqlx.DB, id int) (*VenteLivre, error) {
	vl := &VenteLivre{}
	query := "select * from ventelivre where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(vl)
	if err != nil {
		return vl, werr.Wrapf(err, "Erreur query : "+query)
	}
	return vl, nil
}

func GetVenteLivreFull(db *sqlx.DB, id int) (*VenteLivre, error) {
	vl, err := GetVenteLivre(db, id)
	if err != nil {
		return vl, werr.Wrapf(err, "Erreur appel GetVenteLivre()")
	}
	err = vl.ComputeLivreur(db)
	if err != nil {
		return vl, werr.Wrapf(err, "Erreur appel VenteLivre.ComputeLivreur()")
	}
	err = vl.ComputeChargements(db)
	return vl, err
}

// ************************** Compute *******************************

func (vl *VenteLivre) ComputeLivreur(db *sqlx.DB) error {
	if vl.Livreur != nil {
		return nil
	}
	var err error
	vl.Livreur, err = GetActeur(db, vl.IdLivreur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (vl *VenteLivre) ComputeChargements(db *sqlx.DB) error {
	query := "select id from ventecharge where id_livraison=$1 order by datecharge desc"
	idsCharge := []int{}
	err := db.Select(&idsCharge, query, &vl.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, idC := range idsCharge {
		vc, err := GetVenteChargeFull(db, idC)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetVenteChargeFull()")
		}
		vl.Chargements = append(vl.Chargements, vc)
	}
	return nil
}

// ************************** CRUD *******************************

func InsertVenteLivre(db *sqlx.DB, vl *VenteLivre) (int, error) {
	query := `insert into ventelivre(
        id_vente,
        id_livreur,
        datelivre,
        typecout,
        glprix,
        gltva,
        monheure,
        moprixh,
        motva,
        datepay,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, $11) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		vl.IdVente,
		vl.IdLivreur,
		vl.DateLivre,
		vl.TypeCout,
		vl.GlPrix,
		vl.GlTVA,
		vl.MoNHeure,
		vl.MoPrixH,
		vl.MoTVA,
		vl.DatePay,
		vl.Notes).Scan(&id)
	return id, err
}

func UpdateVenteLivre(db *sqlx.DB, vl *VenteLivre) error {
	query := `update ventelivre set(
        id_vente,
        id_livreur,
        datelivre,
        typecout,
        glprix,
        gltva,
        monheure,
        moprixh,
        motva,
        datepay,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) where id=$12`
	_, err := db.Exec(
		query,
		vl.IdVente,
		vl.IdLivreur,
		vl.DateLivre,
		vl.TypeCout,
		vl.GlPrix,
		vl.GlTVA,
		vl.MoNHeure,
		vl.MoPrixH,
		vl.MoTVA,
		vl.DatePay,
		vl.Notes,
		vl.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeleteVenteLivre(db *sqlx.DB, id int) error {
	// @todo delete VenteCharge
	query := "delete from ventelivre where id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
