/******************************************************************************
    Chautre = Chantiers Autres valorisations
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

type Chautre struct {
	Id           int
	IdClient     int `db:"id_client"`
	IdLieudit    int `db:"id_lieudit"`
	IdUG         int `db:"id_ug"`
	TypeValo     string
	DateContrat  time.Time
	Exploitation string
	Essence      string
	Volume       float64
	Unite        string
	PUHT         float64
	TVA          float64
	DateFacture  time.Time
	NumFacture   string
	Notes        string
	// pas stocké en base
	Client  *Acteur
	Lieudit *Lieudit
	UG      *UG
}

// ************************** Nom *******************************

func (chantier *Chautre) String() string {
	if chantier.Client == nil {
		panic("Erreur dans le code - Le client d'un chantier autre valorisation doit être calculé avant d'appeler String()")
	}
	return LabelValorisation(chantier.TypeValo) + " " + chantier.Client.String() + " " + tiglib.DateFr(chantier.DateContrat)
}

func (chantier *Chautre) FullString() string {
	return "Chantier autre valorisation " + chantier.String()
}

// ************************** Get *******************************

// Renvoie un chantier bois sur pied
// contenant uniquement les données stockées en base
func GetChautre(db *sqlx.DB, idChantier int) (*Chautre, error) {
	chantier := &Chautre{}
	query := "select * from chautre where id=$1"
	row := db.QueryRowx(query, idChantier)
	err := row.StructScan(chantier)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur query : "+query)
	}
	return chantier, nil
}

// Renvoie un chantier bois sur pied contenant :
//      - les données stockées dans la table
//      - Client
//      - Lieudit
//      - UG
func GetChautreFull(db *sqlx.DB, idChantier int) (*Chautre, error) {
	chantier, err := GetChautre(db, idChantier)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Chautre()")
	}
	err = chantier.ComputeClient(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Chautre.ComputeClient()")
	}
	err = chantier.ComputeLieudit(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Chautre.ComputeLieuDit()")
	}
	err = chantier.ComputeUG(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Chautre.ComputeUG()")
	}
	return chantier, nil
}

// Renvoie la liste des années ayant des chantiers bois sur pied,
// @param exclude   Année à exclure du résultat
func GetChautreDifferentYears(db *sqlx.DB, exclude string) ([]string, error) {
	res := []string{}
	list := []time.Time{}
	query := "select datecontrat from chautre order by datecontrat desc"
	err := db.Select(&list, query)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, d := range list {
		y := strconv.Itoa(d.Year())
		if !tiglib.InArrayString(y, res) && y != exclude {
			res = append(res, y)
		}
	}
	return res, nil
}

// Renvoie la liste des chantiers bois sur pied pour une année donnée,
// triés par ordre chronologique inverse.
// Chaque chantier contient les mêmes champs que ceux renvoyés par GetChautreFull()
func GetChautresOfYear(db *sqlx.DB, annee string) ([]*Chautre, error) {
	res := []*Chautre{}
	type ligne struct {
		Id          int
		DateContrat time.Time
	}
	tmp1 := []*ligne{}
	query := "select id,datecontrat from chautre where extract(year from datecontrat)=$1 order by datecontrat desc"
	err := db.Select(&tmp1, query, annee)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, tmp2 := range tmp1 {
		chantier, err := GetChautreFull(db, tmp2.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetChautreFull()")
		}
		res = append(res, chantier)
	}
	return res, nil
}

// ************************** Compute *******************************

func (chantier *Chautre) ComputeClient(db *sqlx.DB) error {
	if chantier.Client != nil {
		return nil
	}
	var err error
	chantier.Client, err = GetActeur(db, chantier.IdClient)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (chantier *Chautre) ComputeLieudit(db *sqlx.DB) error {
	if chantier.Lieudit != nil {
		return nil
	}
	var err error
	chantier.Lieudit, err = GetLieudit(db, chantier.IdLieudit)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetLieudit()")
	}
	return nil
}

func (chantier *Chautre) ComputeUG(db *sqlx.DB) error {
	if chantier.UG != nil {
		return nil
	}
	var err error
	chantier.UG, err = GetUG(db, chantier.IdUG)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetUG()")
	}
	return nil
}

// ************************** CRUD *******************************

func InsertChautre(db *sqlx.DB, chantier *Chautre) (int, error) {
	query := `insert into chautre(
        id_client,
        id_lieudit,                                                        
        id_ug,
        typevalo,
        datecontrat,
        exploitation,
        essence,
        volume,
        unite,
        puht,
        tva,
        datefacture,
        numfacture,
        notes    
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		chantier.IdClient,
		chantier.IdLieudit,
		chantier.IdUG,
		chantier.TypeValo,
		chantier.DateContrat,
		chantier.Exploitation,
		chantier.Essence,
		chantier.Volume,
		chantier.Unite,
		chantier.PUHT,
		chantier.TVA,
		chantier.DateFacture,
		chantier.NumFacture,
		chantier.Notes).Scan(&id)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

func UpdateChautre(db *sqlx.DB, chantier *Chautre) error {
	query := `update chautre set(
        id_client,
        id_lieudit,
        id_ug,
        typevalo,
        datecontrat,
        exploitation,
        essence,
        volume,
        unite,
        puht,
        tva,
        datefacture,
        numfacture,
        notes    
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) where id=$15`
	_, err := db.Exec(
		query,
		chantier.IdClient,
		chantier.IdLieudit,
		chantier.IdUG,
		chantier.TypeValo,
		chantier.DateContrat,
		chantier.Exploitation,
		chantier.Essence,
		chantier.Volume,
		chantier.Unite,
		chantier.PUHT,
		chantier.TVA,
		chantier.DateFacture,
		chantier.NumFacture,
		chantier.Notes,
		chantier.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeleteChautre(db *sqlx.DB, id int) error {
	query := "delete from chautre where id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
