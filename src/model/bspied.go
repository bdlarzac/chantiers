/******************************************************************************
    BSPied = Chantier Bois sur pied
    Bois vendu sur pied à des particuliers, faisant l'objet d'une facturation par BDL

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-02-04 19:32:43+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
	//"fmt"
)

type BSPied struct {
	Id            int
	IdAcheteur    int `db:"id_acheteur"`
	IdLieudit     int `db:"id_lieudit"`
	IdUG          int `db:"id_ug"`
	Nom           string
	DateContrat   time.Time
	Exploitation  string
	Essence       string
	NStereContrat float64
	NStereCoupees float64
	PrixStere     float64
	TVA           float64
	DateFacture   time.Time
	NumFacture    string
	Notes         string
	// Stocké dans bspied_parcelle
	IdsParcelles []int
	// pas stocké en base
	Acheteur  *Acteur
	Lieudit   *Lieudit
	Parcelles []*Parcelle
	UG        *UG
}

// ************************** Nom *******************************

func (bsp *BSPied) String() string {
	if bsp.Lieudit == nil {
		panic("Erreur dans le code - Le lieu-dit d'un chantier bois sur pied doit être calculé avant d'appeler String()")
	}
	if bsp.Acheteur == nil {
		panic("Erreur dans le code - L'acheteur d'un chantier bois sur pied doit être calculé avant d'appeler String()")
	}
	return bsp.Acheteur.String() + " " + bsp.Lieudit.Nom + " " + tiglib.DateFr(bsp.DateContrat)
}

// ************************** Get *******************************

// Renvoie un chantier bois sur pied
// contenant uniquement les données stockées en base
func GetBSPied(db *sqlx.DB, idChantier int) (*BSPied, error) {
	chantier := &BSPied{}
	query := "select * from bspied where id=$1"
	row := db.QueryRowx(query, idChantier)
	err := row.StructScan(chantier)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur query : "+query)
	}
	return chantier, nil
}

// Renvoie un chantier bois sur pied contenant :
//      - les données stockées dans la table
//      - Acheteur
//      - Lieudit
//      - IdsParcelles
//      - Parcelles
//      - UG
func GetBSPiedFull(db *sqlx.DB, idChantier int) (*BSPied, error) {
	bsp, err := GetBSPied(db, idChantier)
	if err != nil {
		return bsp, werr.Wrapf(err, "Erreur appel GetBSPied()")
	}
	err = bsp.ComputeLieudit(db)
	if err != nil {
		return bsp, werr.Wrapf(err, "Erreur appel BSPied.ComputeLieudit()")
	}
	err = bsp.ComputeAcheteur(db)
	if err != nil {
		return bsp, werr.Wrapf(err, "Erreur appel BSPied.ComputeAcheteur()")
	}
	err = bsp.ComputeParcelles(db)
	if err != nil {
		return bsp, werr.Wrapf(err, "Erreur appel BSPied.ComputeParcelles()")
	}
	err = bsp.ComputeUG(db)
	if err != nil {
		return bsp, werr.Wrapf(err, "Erreur appel BSPied.ComputeUG()")
	}
	return bsp, nil
}

// Renvoie la liste des années ayant des chantiers bois sur pied,
// @param exclude   Année à exclure du résultat
func GetBSPiedDifferentYears(db *sqlx.DB, exclude string) ([]string, error) {
	res := []string{}
	list := []time.Time{}
	query := "select datecontrat from bspied order by datecontrat desc"
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

// Renvoie la liste des chantiers bois sur pied pour une année donnée,
// triés par ordre chronologique inverse.
// Chaque chantier contient les mêmes champs que ceux renvoyés par GetBSPiedFull()
func GetBSPiedsOfYear(db *sqlx.DB, annee string) ([]*BSPied, error) {
	res := []*BSPied{}
	type ligne struct {
		Id          int
		DateContrat time.Time
	}
	tmp1 := []*ligne{}
	// select aussi datecontrat au lieu de seulement id pour pouvoir faire le order by
	query := "select id,datecontrat from bspied where extract(year from datecontrat)=$1 order by datecontrat desc"
	err := db.Select(&tmp1, query, annee)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, tmp2 := range tmp1 {
		chantier, err := GetBSPiedFull(db, tmp2.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetBSPiedFull()")
		}
		res = append(res, chantier)
	}
	return res, nil
}

// ************************** Compute *******************************

func (bsp *BSPied) ComputeLieudit(db *sqlx.DB) error {
	if bsp.Lieudit != nil {
		return nil
	}
	var err error
	bsp.Lieudit, err = GetLieudit(db, bsp.IdLieudit)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetLieudit()")
	}
	return nil
}

func (bsp *BSPied) ComputeAcheteur(db *sqlx.DB) error {
	if bsp.Acheteur != nil {
		return nil
	}
	var err error
	bsp.Acheteur, err = GetActeur(db, bsp.IdAcheteur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

// Calcule à la fois Parcelles et IdsParcelles
func (bsp *BSPied) ComputeParcelles(db *sqlx.DB) error {
	if len(bsp.Parcelles) != 0 {
		return nil
	}
	var err error
	query := "select id_parcelle from bspied_parcelle where id_bspied=$1"
	idsP := []int{}
	err = db.Select(&idsP, query, bsp.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, idP := range idsP {
		p, err := GetParcelle(db, idP)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetParcelle()")
		}
		bsp.IdsParcelles = append(bsp.IdsParcelles, idP)
		bsp.Parcelles = append(bsp.Parcelles, p)
	}
	return nil
}

func (bsp *BSPied) ComputeUG(db *sqlx.DB) error {
	if bsp.UG != nil {
		return nil
	}
	var err error
	bsp.UG, err = GetUG(db, bsp.IdUG)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetUG()")
	}
	return nil
}

// ************************** CRUD *******************************

func InsertBSPied(db *sqlx.DB, bsp *BSPied) (int, error) {
	query := `insert into bspied(
        id_acheteur,
        id_lieudit,
        id_ug,
        datecontrat,
        exploitation,
        essence,
        nsterecontrat,
        nsterecoupees,
        prixstere,
        tva,
        datefacture,
        numfacture,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		bsp.IdAcheteur,
		bsp.IdLieudit,
		bsp.IdUG,
		bsp.DateContrat,
		bsp.Exploitation,
		bsp.Essence,
		bsp.NStereContrat,
		bsp.NStereCoupees,
		bsp.PrixStere,
		bsp.TVA,
		bsp.DateFacture,
		bsp.NumFacture,
		bsp.Notes).Scan(&id)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	query = "insert into bspied_parcelle values($1, $2)"
	for _, idP := range bsp.IdsParcelles {
		_ = db.QueryRow(query, id, idP)
	}
	return id, nil
}

func UpdateBSPied(db *sqlx.DB, bsp *BSPied) error {
	query := `update bspied set(
        id_acheteur,
        id_lieudit,
        id_ug,
        datecontrat,
        exploitation,
        essence,
        nsterecontrat,
        nsterecoupees,
        prixstere,
        tva,
        datefacture,
        numfacture,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) where id=$14`
	_, err := db.Exec(
		query,
		bsp.IdAcheteur,
		bsp.IdLieudit,
		bsp.IdUG,
		bsp.DateContrat,
		bsp.Exploitation,
		bsp.Essence,
		bsp.NStereContrat,
		bsp.NStereCoupees,
		bsp.PrixStere,
		bsp.TVA,
		bsp.DateFacture,
		bsp.NumFacture,
		bsp.Notes,
		bsp.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from bspied_parcelle where id_bspied=$1"
	_, err = db.Exec(query, bsp.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "insert into bspied_parcelle values($1,$2)"
	for _, idP := range bsp.IdsParcelles {
		_ = db.QueryRow(query, bsp.Id, idP)
	}
	return nil
}

func DeleteBSPied(db *sqlx.DB, id int) error {
	query := "delete from bspied_parcelle where id_bspied=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from bspied where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
