/******************************************************************************
    Hangars pour stocker des plaquettes

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-12-03 16:27:28+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"time"
	"errors"
//	"fmt"
)

type Stockage struct {
	Id  int
	Nom string
	// Lg float64
	// Lat float64
	Archived bool
	// pas stocké en base
	Stock     float64
	Deletable bool
	TasActifs []*Tas
	Frais     []*StockFrais
}

// ************************** Get one *******************************

// Renvoie un lieu de stockage contenant les données stockées en base.
// Les autres champs ne sont pas remplis.
func GetStockage(db *sqlx.DB, id int) (*Stockage, error) {
	s := &Stockage{}
	query := "select * from stockage where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(s)
	if err != nil {
		return s, werr.Wrapf(err, "Erreur query : "+query)
	}
	return s, nil
}

// Renvoie un lieu de stockage contenant
// - les données stockées en base.
// - les frais
// - le champ Deletable
// - les tas non vides
func GetStockageFull(db *sqlx.DB, id int) (*Stockage, error) {
	s, err := GetStockage(db, id)
	if err != nil {
		return s, werr.Wrapf(err, "Erreur appel GetStockage()")
	}
	err = s.ComputeFrais(db)
	if err != nil {
		return s, werr.Wrapf(err, "Erreur appel Stockage.ComputeFrais()")
	}
	err = s.ComputeTasActifs(db)
	if err != nil {
		return s, werr.Wrapf(err, "Erreur appel Stockage.ComputeTasActifs()")
	}
	err = s.ComputeStock(db)
	if err != nil {
		return s, werr.Wrapf(err, "Erreur appel Stockage.ComputeStock()")
	}
	err = s.ComputeDeletable(db)
	if err != nil {
		return s, werr.Wrapf(err, "Erreur appel Stockage.ComputeDeletable()")
	}
	//
	return s, nil
}

// ************************** Get many *******************************

// Renvoie la liste de tous les lieux de stockage
// avec uniquement les champs stockés en base
func GetStockages(db *sqlx.DB) ([]*Stockage, error) {
	stockages := []*Stockage{}
	query := "select * from stockage order by nom"
	err := db.Select(&stockages, query)
	if err != nil {
		return stockages, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	return stockages, nil
}

// Renvoie la liste de tous les lieux de stockage actifs (= pas archivés)
// avec uniquement les champs stockés en base
func GetStockagesActifs(db *sqlx.DB) ([]*Stockage, error) {
	stockages := []*Stockage{}
	query := "select * from stockage where archived=FALSE order by nom"
	err := db.Select(&stockages, query)
	if err != nil {
		return stockages, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	return stockages, nil
}

// Renvoie la liste de tous les lieux de stockage contenant
// les mêmes données que celles renvoyées par GetStockageFull()
func GetStockagesFull(db *sqlx.DB) ([]*Stockage, error) {
	stockages := []*Stockage{}
	query := "select id from stockage order by nom"
	ids := []int{}
	err := db.Select(&ids, query)
	if err != nil {
		return stockages, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, id := range ids {
		s, err := GetStockageFull(db, id)
		if err != nil {
			return stockages, werr.Wrapf(err, "Erreur appel GetStockageFull()")
		}
		stockages = append(stockages, s)
	}
	return stockages, nil
}

// ************************** Compute *******************************

func (s *Stockage) ComputeFrais(db *sqlx.DB) error {
	query := "select * from stockfrais where id_stockage=$1 order by datedeb desc"
	frais := []*StockFrais{}
	err := db.Select(&frais, query, s.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	s.Frais = frais
	return nil
}

func (s *Stockage) ComputeTasActifs(db *sqlx.DB) error {
	query := "select * from tas where actif and id_stockage=$1"
	err := db.Select(&s.TasActifs, query, &s.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for i, _ := range s.TasActifs {
		s.TasActifs[i].Stockage = s
		err = s.TasActifs[i].ComputeChantier(db)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel Tas.ComputeChantier()")
		}
		err = s.TasActifs[i].ComputeNom(db)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel Tas.ComputeNom()")
		}
	}
	return nil
}

func (s *Stockage) ComputeStock(db *sqlx.DB) error {
	var stocks []float64
	query := "select stock from tas where actif and id_stockage=$1"
	err := db.Select(&stocks, query, s.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	s.Stock = 0
	for _, stock := range stocks {
		s.Stock += stock
	}
	return nil
}

func (s *Stockage) ComputeDeletable(db *sqlx.DB) error {
	var count int
	query := "select count(*) from tas where actif and id_stockage=$1"
	err := db.QueryRow(query, s.Id).Scan(&count)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	s.Deletable = (count == 0)
	return nil
}

// Calcule le coût du stockage pour une période donnée.
// Prend en compte tous les frais du hangar (loyer, elec, assurance)
// Le coût est ramené à la période considérée.
// Ex : pour un loyer de 6000 E / an, si j2 - j1 = 6 mois, va compter 3000
// @param j1, j2 jours de début / fin de la période au format YYYY-MM-DD
func (s *Stockage) ComputeCout(db *sqlx.DB, jour1, jour2 string) (float64, error) {
    var err error
    j1, err := time.Parse("2006-01-02", jour1)
	if err != nil {
		return 0, werr.Wrapf(err, "Format de date incorrect : "+jour1)
	}
    j2, err := time.Parse("2006-01-02", jour2)
	if err != nil {
		return 0, werr.Wrapf(err, "Format de date incorrect : "+jour2)
	}
    if j2.Before(j1){
        return 0, errors.New("ComputeCout() a besoin de j1 < j2")
    }
    // Récupère les frais
    // tels que datedeb ou datefin sont dans [j1, j2]
    var frais []StockFrais
	query := `select * from stockfrais where id_stockage=$1
	    and ( (datedeb>=$2 and datedeb<=$3) or (datefin>=$2 and datefin<=$3) )`
	err = db.Select(&frais, query, s.Id, &jour1, &jour2)
	if err != nil {
		return 0, werr.Wrapf(err, "Erreur query : "+query)
	}
	var debFrais, finFrais time.Time
	var dureeFraisTotale time.Duration // durée totale du frais, peut dépasser [j1, j2]
	var dureeFraisPeriode time.Duration // durée du frais dans [j1, j2]
	total := float64(0)
	for _, f := range frais {
	    dureeFraisTotale = f.DateFin.Sub(f.DateDebut)
	    if j1.After(f.DateDebut){
	        debFrais = j1
	    } else{
	        debFrais = f.DateDebut
	    }
	    if j2.Before(f.DateFin){
	        finFrais = j2
	    } else{
	        finFrais = f.DateFin
	    }
	    dureeFraisPeriode = finFrais.Sub(debFrais)
	    // contribution du frais sur la période [j1, j2]
	    total += (dureeFraisPeriode.Hours() * f.Montant / dureeFraisTotale.Hours())
	}
    return total, nil
}

// ************************** CRUD *******************************

func InsertStockage(db *sqlx.DB, s *Stockage) (int, error) {
	query := `insert into stockage(nom) values($1) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		s.Nom).Scan(&id)
	return id, err
}

func UpdateStockage(db *sqlx.DB, s *Stockage) error {
	query := `update stockage set(
	    nom,
	    archived
	    ) = ($1,$2) where id=$3`
	_, err := db.Exec(
		query,
		s.Nom,
		s.Archived,
		s.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeleteStockage(db *sqlx.DB, id int) error {
	query := "delete from stockfrais where id_stockage=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from stockage where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
