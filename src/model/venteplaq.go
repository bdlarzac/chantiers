/*
*****************************************************************************

	Vente de plaquettes, depuis un lieu de stockage
	Note: La string stockée dans FactureLivraisonUnite ne vient pas d'une enum.
	Donc les valeurs "km" et "map" sont codées en dur
	- dans le js de venteplaq-form.html
	- dans venteplaq-show.html
	- dans venteplaq-list.html
	La valeur est insérée directement en base dans control.ventePlaqForm2var()

	Pour afficher sur la facture :
	Si FactureLivraisonUnite = map, on utilise Qte
	Si FactureLivraisonUnite = km, on utilise FactureLivraisonNbKm

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2020-01-22 02:56:23+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"strconv"
	"time"

	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type VentePlaq struct {
	Id            int
	IdClient      int `db:"id_client"`
	IdFournisseur int `db:"id_fournisseur"`
	PUHT          float64
	TVA           float64
	DateVente     time.Time
	DatePaiement  time.Time
	// Facture
	NumFacture            string
	DateFacture           time.Time
	FactureLivraison      bool
	FactureLivraisonPUHT  float64
	FactureLivraisonTVA   float64
	FactureNotes          bool
	FactureLivraisonUnite string  // voir note dans commentaire de la classe
	FactureLivraisonNbKm  float64 // voir note dans commentaire de la classe
	//
	Notes string
	// Pas stocké en base
	Qte         float64 // maps
	Client      *Acteur
	Fournisseur *Acteur
	Livraisons  []*VenteLivre
	Chantiers   []*Plaq
}

// ************************** Manipulation Quantité *******************************

// @param   qte en maps
func (vp *VentePlaq) ModifierQte(db *sqlx.DB, qte float64) {
	vp.Qte += qte
}

// ************************** Nom *******************************

func (vp *VentePlaq) String() string {
	if vp.Client == nil {
		panic("Erreur dans le code - Le client d'une vente plaquettes doit être calculé avant d'appeler String()")
	}
	return vp.Client.String() + " " + tiglib.DateFr(vp.DateVente)
}

func (vp *VentePlaq) FullString() string {
	return "Vente " + vp.String()
}

func (vp *VentePlaq) StringSansDate() string {
	if vp.Client == nil {
		panic("Erreur dans le code - Le client d'une vente plaquettes doit être calculé avant d'appeler String()")
	}
	return "Vente " + vp.Client.String()
}

// ************************** Get one *******************************

func GetVentePlaq(db *sqlx.DB, id int) (*VentePlaq, error) {
	vp := &VentePlaq{}
	query := "select * from venteplaq where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(vp)
	if err != nil {
		return vp, werr.Wrapf(err, "Erreur query : "+query)
	}
	return vp, nil
}

func GetVentePlaqFull(db *sqlx.DB, id int) (*VentePlaq, error) {
	vp, err := GetVentePlaq(db, id)
	if err != nil {
		return vp, werr.Wrapf(err, "Erreur appel GetVentePlaq()")
	}
	err = vp.ComputeQte(db)
	if err != nil {
		return vp, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeQte()")
	}
	err = vp.ComputeClient(db)
	if err != nil {
		return vp, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeClient()")
	}
	err = vp.ComputeFournisseur(db)
	if err != nil {
		return vp, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeFournisseur()")
	}
	err = vp.ComputeLivraisons(db)
	if err != nil {
		return vp, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeLivraisons()")
	}
	err = vp.ComputeChantiers(db)
	if err != nil {
		return vp, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeChantiers()")
	}
	return vp, nil
}

// ************************** Get many *******************************

// Renvoie la liste des années ayant des ventes de plaquettes,
// @param   exclude   Année à exclure du résultat
// @return  Liste de string au format YYYY
func GetVentePlaqDifferentYears(db *sqlx.DB, exclude string) ([]string, error) {
	res := []string{}
	list := []time.Time{}
	query := "select datevente from venteplaq order by datevente"
	err := db.Select(&list, query)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, d := range list {
		y := strconv.Itoa(d.Year())
		if !tiglib.InArray(y, res) && y != exclude {
			res = append(res, y)
		}
	}
	return res, nil
}

// Renvoie la liste des ventes de plaquettes pour une année donnée,
// triés par ordre chronologique inverse.
// Chaque vente contient les mêmes champs que ceux renvoyés par GetVentePlaqFull()
func GetVentePlaqsOfYear(db *sqlx.DB, annee string) ([]*VentePlaq, error) {
	res := []*VentePlaq{}
	type ligne struct {
		Id        int
		DateVente time.Time
	}
	tmp1 := []*ligne{}
	// select aussi datevente au lieu de seulement id pour pouvoir faire le order by
	query := "select id,datevente from venteplaq where extract(year from datevente)=$1 order by datevente"
	err := db.Select(&tmp1, query, annee)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, tmp2 := range tmp1 {
		v, err := GetVentePlaqFull(db, tmp2.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetVentePlaqFull()")
		}
		res = append(res, v)
	}
	return res, nil
}

// Renvoie la liste des ventes de plaquettes de date 1 à date 2,
// triés par ordre chronologique.
// Chaque vente contient les mêmes champs que ceux renvoyés par GetVentePlaqFull()
func GetVentePlaqsOfPeriod(db *sqlx.DB, date1, date2 time.Time) ([]*VentePlaq, error) {
	res := []*VentePlaq{}
	query := "select * from venteplaq where datevente>=$1 and datevente<=$2 order by datevente"
	err := db.Select(&res, query, date1, date2)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for i, tmp := range res {
		v, err := GetVentePlaqFull(db, tmp.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetVentePlaqFull()")
		}
		res[i] = v
	}
	return res, nil
}

// Renvoie la liste des ventes de plaquettes pour un client donné, situé entre 2 dates,
// Chaque chantier contient les mêmes champs que ceux renvoyés par GetVentePlaqFull()
func GetVentePlaqsOfClient(db *sqlx.DB, idClient int, dateDebut, dateFin time.Time) ([]*VentePlaq, error) {
	res := []*VentePlaq{}
	query := "select * from venteplaq where id_client=$1 and datevente>=$2 and datevente<=$3 order by datevente"
	err := db.Select(&res, query, idClient, dateDebut, dateFin)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, vp := range res {
		err := vp.ComputeQte(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel vp.ComputeQte()")
		}
		err = vp.ComputeLivraisons(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeLivraisons()")
		}
	}
	return res, nil
}

// ************************** Compute *******************************

func (vp *VentePlaq) ComputeQte(db *sqlx.DB) error {
	var qtes []float64
	query := `select qte from ventecharge where id_livraison in(
                select id from ventelivre where id_vente=$1
            )`
	err := db.Select(&qtes, query, vp.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	vp.Qte = 0
	for _, qte := range qtes {
		vp.Qte += qte
	}
	return nil
}

func (vp *VentePlaq) ComputeClient(db *sqlx.DB) error {
	var err error
	vp.Client, err = GetActeur(db, vp.IdClient)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur() pour VentePlaq.ComputeClient()")
	}
	return nil
}

func (vp *VentePlaq) ComputeFournisseur(db *sqlx.DB) error {
	var err error
	vp.Fournisseur, err = GetActeur(db, vp.IdFournisseur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur() pour VentePlaq.ComputeFournisseur()")
	}
	return nil
}

func (vp *VentePlaq) ComputeLivraisons(db *sqlx.DB) error {
	query := "select id from ventelivre where id_vente=$1 order by datelivre"
	idsLivraison := []int{}
	err := db.Select(&idsLivraison, query, &vp.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, idL := range idsLivraison {
		vl, err := GetVenteLivreFull(db, idL)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetVenteLivreFull()")
		}
		vp.Livraisons = append(vp.Livraisons, vl)
	}
	return nil
}

func (vp *VentePlaq) ComputeChantiers(db *sqlx.DB) error {
	ids := []int{}
	query := `select distinct id_chantier from tas where id in (
                  select id_tas from ventecharge where id_livraison in(
                      select id from ventelivre where id_vente in(
                          select id from venteplaq where id=$1
                      )
                  )
              )`
	err := db.Select(&ids, query, &vp.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, idChantier := range ids {
		chantier, err := GetPlaq(db, idChantier)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetPlaq()")
		}
		// Ajoute lieu-dit pour avoir le nom du chantier
		err = chantier.ComputeLieudits(db)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
		}
		vp.Chantiers = append(vp.Chantiers, chantier)
	}
	return nil
}

// ************************** CRUD *******************************

func InsertVentePlaq(db *sqlx.DB, vp *VentePlaq) (int, error) {
	query := `insert into venteplaq(
        id_client,
        id_fournisseur,
        puht,
        tva,
        datevente,
        datepaiement,
        numfacture,
        datefacture,
        facturelivraison,
        facturelivraisonpuht,
        facturelivraisonunite,
        facturelivraisonnbkm,
        facturelivraisontva,
        facturenotes,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		vp.IdClient,
		vp.IdFournisseur,
		vp.PUHT,
		vp.TVA,
		vp.DateVente,
		vp.DatePaiement,
		vp.NumFacture,
		vp.DateFacture,
		vp.FactureLivraison,
		vp.FactureLivraisonPUHT,
		vp.FactureLivraisonUnite,
		vp.FactureLivraisonNbKm,
		vp.FactureLivraisonTVA,
		vp.FactureNotes,
		vp.Notes).Scan(&id)
	return id, err
}

func UpdateVentePlaq(db *sqlx.DB, vp *VentePlaq) error {
	query := `update venteplaq set(
        id_client,
        id_fournisseur,
        puht,
        tva,
        datevente,
        datepaiement,
        numfacture,
        datefacture,
        facturelivraison,
        facturelivraisonpuht,
        facturelivraisonunite,
        facturelivraisonnbkm,
        facturelivraisontva,
        facturenotes,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) where id=$16`
	_, err := db.Exec(
		query,
		vp.IdClient,
		vp.IdFournisseur,
		vp.PUHT,
		vp.TVA,
		vp.DateVente,
		vp.DatePaiement,
		vp.NumFacture,
		vp.DateFacture,
		vp.FactureLivraison,
		vp.FactureLivraisonPUHT,
		vp.FactureLivraisonUnite,
		vp.FactureLivraisonNbKm,
		vp.FactureLivraisonTVA,
		vp.FactureNotes,
		vp.Notes,
		vp.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeleteVentePlaq(db *sqlx.DB, id int) error {
	// delete les livraisons dépendant de cette vente
	idsLivraison := []int{}
	query := "select id from ventelivre where id_vente=$1"
	err := db.Select(&idsLivraison, query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, idL := range idsLivraison {
		err := DeleteVenteLivre(db, idL) // va aussi effacer les chargements
		if err != nil {
			return werr.Wrapf(err, "Erreur appel DeleteVenteLivre()")
		}
	}
	// delete la vente
	query = "delete from venteplaq where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
