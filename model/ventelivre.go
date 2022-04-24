/******************************************************************************
    Livraison lors d'une vente de plaquettes

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2020-03-13 10:59:50+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type VenteLivre struct {
	Id            int
	IdVente       int `db:"id_vente"`
	IdLivreur     int `db:"id_livreur"`
	IdConducteur  int `db:"id_conducteur"`
	IdProprioutil int `db:"id_proprioutil"`
	DateLivre     time.Time
	TypeCout      string // G (global) ou D (détail)
	// coût global
	GlPrix    float64
	GlTVA     float64
	GlDatePay time.Time
	// coût main d'oeuvre
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
	Qte         float64 // somme des quantités des chargements
	Livreur     *Acteur
	Conducteur  *Acteur
	Proprioutil *Acteur
	Vente       *VentePlaq
	Chargements []*VenteCharge
}

// ************************** Nom *******************************

// Les livraisons sont les seules opérations composées (avec coût global / coût détaillé)
// ayant besoin de String(), car le nom apparaît dans le formulaire de new / update chargement.
func (vl *VenteLivre) String() string {
	nom := ""
	if vl.TypeCout == "G" {
		if vl.Livreur == nil {
			msg := "Erreur dans le code VenteLivre.String() : livreur doit être calculé avant d'appeler String()" +
				"\nventecharge id = " + strconv.Itoa(vl.Id)
			panic(msg)
		}
		nom = vl.Livreur.String()
	} else {
		if vl.Conducteur == nil || vl.Proprioutil == nil {
			msg := "Erreur dans le code VenteLivre.String() : conducteur et proprioutil doivent être calculés avant d'appeler String()" +
				"\nventecharge id = " + strconv.Itoa(vl.Id)
			panic(msg)
		}
		nom = vl.Conducteur.String() + " / " + vl.Proprioutil.String()
	}
	return nom + " " + tiglib.DateFr(vl.DateLivre)
}

func (vl *VenteLivre) FullString() string {
	return "Livraison plaquettes " + vl.String()
}

// ************************** Get *******************************

func GetVenteLivre(db *sqlx.DB, id int) (vl *VenteLivre, err error) {
	vl = &VenteLivre{}
	query := "select * from ventelivre where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(vl)
	if err != nil {
		return vl, werr.Wrapf(err, "Erreur query : "+query)
	}
	return vl, nil
}

func GetVenteLivreFull(db *sqlx.DB, id int) (vl *VenteLivre, err error) {
	vl, err = GetVenteLivre(db, id)
	if err != nil {
		return vl, werr.Wrapf(err, "Erreur appel GetVenteLivre()")
	}
	err = vl.ComputeLivreur(db)
	if err != nil {
		return vl, werr.Wrapf(err, "Erreur appel VenteLivre.ComputeLivreur()")
	}
	err = vl.ComputeConducteur(db)
	if err != nil {
		return vl, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeConducteur()")
	}
	err = vl.ComputeProprioutil(db)
	if err != nil {
		return vl, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeProprioutil()")
	}
	err = vl.ComputeChargements(db)
	if err != nil {
		return vl, werr.Wrapf(err, "Erreur appel VenteCharge.ComputeChargements()")
	}
	return vl, nil
}

// ************************** Compute *******************************

func (vl *VenteLivre) ComputeLivreur(db *sqlx.DB) (err error) {
	if vl.Livreur != nil {
		return nil
	}
	vl.Livreur, err = GetActeur(db, vl.IdLivreur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (vl *VenteLivre) ComputeConducteur(db *sqlx.DB) (err error) {
	if vl.Conducteur != nil {
		return nil
	}
	vl.Conducteur, err = GetActeur(db, vl.IdConducteur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (vl *VenteLivre) ComputeProprioutil(db *sqlx.DB) (err error) {
	if vl.Proprioutil != nil {
		return nil
	}
	vl.Proprioutil, err = GetActeur(db, vl.IdProprioutil)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

// Calcule à la fois les chargements et la quantité de la livraison
func (vl *VenteLivre) ComputeChargements(db *sqlx.DB) (err error) {
	if vl.Chargements != nil {
		return nil
	}
	query := "select id from ventecharge where id_livraison=$1 order by datecharge"
	idsCharge := []int{}
	err = db.Select(&idsCharge, query, &vl.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, idC := range idsCharge {
		vc, err := GetVenteChargeFull(db, idC)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetVenteChargeFull()")
		}
		vl.Chargements = append(vl.Chargements, vc)
		vl.Qte += vc.Qte // ICI, calcule aussi la quantité de la livraison
	}
	return nil
}

// ************************** CRUD *******************************

func InsertVenteLivre(db *sqlx.DB, vl *VenteLivre) (id int, err error) {
	query := `insert into ventelivre(
        id_vente,
        id_livreur,
        id_conducteur,
        id_proprioutil,
        datelivre,
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
	err = db.QueryRow(
		query,
		vl.IdVente,
		vl.IdLivreur,
		vl.IdConducteur,
		vl.IdProprioutil,
		vl.DateLivre,
		vl.TypeCout,
		vl.GlPrix,
		vl.GlTVA,
		vl.GlDatePay,
		vl.OuPrix,
		vl.OuTVA,
		vl.OuDatePay,
		vl.MoNHeure,
		vl.MoPrixH,
		vl.MoTVA,
		vl.MoDatePay,
		vl.Notes).Scan(&id)
	if err != nil {
		return 0, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

func UpdateVenteLivre(db *sqlx.DB, vl *VenteLivre) (err error) {
	query := `update ventelivre set(
        id_vente,
        id_livreur,
        id_conducteur,
        id_proprioutil,
        datelivre,
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
	_, err = db.Exec(
		query,
		vl.IdVente,
		vl.IdLivreur,
		vl.IdConducteur,
		vl.IdProprioutil,
		vl.DateLivre,
		vl.TypeCout,
		vl.GlPrix,
		vl.GlTVA,
		vl.GlDatePay,
		vl.OuPrix,
		vl.OuTVA,
		vl.OuDatePay,
		vl.MoNHeure,
		vl.MoPrixH,
		vl.MoTVA,
		vl.MoDatePay,
		vl.Notes,
		vl.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeleteVenteLivre(db *sqlx.DB, id int) (err error) {
	// delete les chargements dépendant de cette livraison
	idsCharge := []int{}
	query := "select id from ventecharge where id_livraison=$1"
	err = db.Select(&idsCharge, query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, idC := range idsCharge {
		// Attention ici ne pas faire directement delete ventecharge en base
		// car DeleteVenteCharge() gère le stock des tas associés
		err := DeleteVenteCharge(db, idC)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel DeleteVenteCharge()")
		}
	}
	// delete la livraison
	query = "delete from ventelivre where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
