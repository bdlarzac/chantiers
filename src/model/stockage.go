/*
Hangars pour stocker des plaquettes

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2019-12-03 16:27:28+01:00, Thierry Graff : Creation
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type Stockage struct {
	Id  int
	Nom string
	// Lg float64
	// Lat float64
	Archived bool
	// pas stocké en base
	Stock      float64
	Deletable  bool
	Archivable bool
	TasActifs  []*Tas
	Frais      []*StockFrais
}

// ************************** Get one *******************************

// Renvoie un lieu de stockage contenant les données stockées en base.
// Les autres champs ne sont pas remplis.
func GetStockage(db *sqlx.DB, id int) (s *Stockage, err error) {
	s = &Stockage{}
	query := "select * from stockage where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(s)
	if err != nil {
		return s, werr.Wrapf(err, "Erreur query : "+query)
	}
	return s, nil
}

// Renvoie un lieu de stockage contenant
// - les données stockées en base.
// - les frais
// - les champs Deletable et Archivable
// - les tas non vides
func GetStockageFull(db *sqlx.DB, id int) (s *Stockage, err error) {
	s, err = GetStockage(db, id)
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
	err = s.ComputeDeletableAndArchivable(db)
	if err != nil {
		return s, werr.Wrapf(err, "Erreur appel Stockage.ComputeDeletableAndArchivable()")
	}
	//
	return s, nil
}

// ************************** Get many *******************************

// Renvoie la liste de tous les lieux de stockage
// avec uniquement les champs stockés en base
// @param actifs
//
//	true => ne renvoie que les stockages actifs (pas archivés)
//	false => ne renvoie que les stockages archivés
func GetStockages(db *sqlx.DB, actifs bool) (stockages []*Stockage, err error) {
	stockages = []*Stockage{}
	query := "select * from stockage where archived="
	if actifs {
		query += "FALSE"
	} else {
		query += "TRUE"
	}
	err = db.Select(&stockages, query)
	if err != nil {
		return stockages, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	return stockages, nil
}

// Renvoie la liste de tous les lieux de stockage contenant
// les mêmes données que celles renvoyées par GetStockageFull()
// @param actifs
//
//	true => ne renvoie que les stockages actifs (pas archivés)
//	false => ne renvoie que les stockages archivés
func GetStockagesFull(db *sqlx.DB, actifs bool) (stockages []*Stockage, err error) {
	res := []*Stockage{}
	stockages, err = GetStockages(db, actifs)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur appel GetStockages()")
	}
	for _, stockage := range stockages {
		s, err := GetStockageFull(db, stockage.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetStockageFull()")
		}
		res = append(res, s)
	}
	return res, nil
}

// ************************** Compute *******************************

func (s *Stockage) ComputeFrais(db *sqlx.DB) (err error) {
	query := "select * from stockfrais where id_stockage=$1 order by datedeb"
	frais := []*StockFrais{}
	err = db.Select(&frais, query, s.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	s.Frais = frais
	return nil
}

func (s *Stockage) ComputeTasActifs(db *sqlx.DB) (err error) {
	query := "select * from tas where actif and id_stockage=$1"
	err = db.Select(&s.TasActifs, query, &s.Id)
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

func (s *Stockage) ComputeStock(db *sqlx.DB) (err error) {
	var stocks []float64
	query := "select stock from tas where actif and id_stockage=$1"
	err = db.Select(&stocks, query, s.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	s.Stock = 0
	for _, stock := range stocks {
		s.Stock += stock
	}
	return nil
}

// Calcule les champs Deletable et Archivable
// Un stockage est Deletable s'il n'est associé à aucune activité
// Un stockage est Archivable s'il est associé à des activités mais ne contient pas de tas actif
func (s *Stockage) ComputeDeletableAndArchivable(db *sqlx.DB) (err error) {
	var count int
	// Deletable
	// il suffit de compter les tas associés au lieu de stockage
	// pour savoir si des chantiers plaquettes y sont associés
	// car DeletePlaq() efface tous les tas reliés au chantier
	query := `select count(*) from tas where id_stockage=$1`
	err = db.QueryRow(query, s.Id).Scan(&count)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	s.Deletable = (count == 0)
	// Archivable
	query = "select count(*) from tas where actif and id_stockage=$1"
	err = db.QueryRow(query, s.Id).Scan(&count)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	s.Archivable = (count == 0)
	return nil
}

// Calcule le coût du stockage pour une période donnée.
// Prend en compte tous les frais du hangar (loyer, elec, assurance)
// Le coût est ramené à la période considérée.
// Ex : pour un loyer de 6000 E / an, si j2 - j1 = 6 mois, va compter 3000
// @param   jourD, jour2 jours de début / fin de la période au format YYYY-MM-DD
// @return  Tableau contenant les coûts pour chaque jour de la période [jourD, jourF]
//
//	res[0] = frais pour jourD, res[1] = frais pour jourD + 1, etc.
func (s *Stockage) ComputeCout(db *sqlx.DB, jourD, jourF string) (res []float64, err error) {
	res = []float64{}
	jD, err := time.Parse("2006-01-02", jourD)
	if err != nil {
		return res, werr.Wrapf(err, "Format de date incorrect : "+jourD)
	}
	jF, err := time.Parse("2006-01-02", jourF)
	if err != nil {
		return res, werr.Wrapf(err, "Format de date incorrect : "+jourF)
	}
	if jF.Before(jD) {
		return res, errors.New("ComputeCout() a besoin de jourD < jourF")
	}
	//	duree := jF.Sub(jD)
	// Récupère les frais
	// tels que datedeb ou datefin sont dans [jD, jF]
	var frais []StockFrais
	query := `select * from stockfrais where id_stockage=$1
	    and ( (datedeb>=$2 and datedeb<=$3) or (datefin>=$2 and datefin<=$3) )`
	err = db.Select(&frais, query, s.Id, &jD, &jF)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query : "+query)
	}
	/*
	   	var debFrais, finFrais time.Time
	   	var dureeFraisTotale time.Duration  // durée totale du frais, peut dépasser [jD, jF]
	   	var dureeFraisPeriode time.Duration // durée du frais dans [jD, jF]
	   	//
	   	for _, f := range frais {
	   		dureeFraisTotale = f.DateFin.Sub(f.DateDebut)
	   		if jD.After(f.DateDebut) {
	   			debFrais = jD
	   		} else {
	   			debFrais = f.DateDebut
	   		}
	   		if jF.Before(f.DateFin) {
	   			finFrais = jF
	   		} else {
	   			finFrais = f.DateFin
	   		}
	   		dureeFraisPeriode = finFrais.Sub(debFrais)
	   		// contribution du frais sur la période [jD, jF]
	   //		total += (dureeFraisPeriode.Hours() * f.Montant / dureeFraisTotale.Hours())
	   	}
	*/
	return res, nil
}

/*
// Calcule le coût du stockage pour une période donnée.
// Prend en compte tous les frais du hangar (loyer, elec, assurance)
// Le coût est ramené à la période considérée.
// Ex : pour un loyer de 6000 E / an, si j2 - j1 = 6 mois, va compter 3000
// @param j1, j2 jours de début / fin de la période au format YYYY-MM-DD
func (s *Stockage) ComputeCout(db *sqlx.DB, jour1, jour2 string) (total float64, err error) {
	j1, err := time.Parse("2006-01-02", jour1)
	if err != nil {
		return 0, werr.Wrapf(err, "Format de date incorrect : "+jour1)
	}
	j2, err := time.Parse("2006-01-02", jour2)
	if err != nil {
		return 0, werr.Wrapf(err, "Format de date incorrect : "+jour2)
	}
	if j2.Before(j1) {
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
	var dureeFraisTotale time.Duration  // durée totale du frais, peut dépasser [j1, j2]
	var dureeFraisPeriode time.Duration // durée du frais dans [j1, j2]
	//
	for _, f := range frais {
		dureeFraisTotale = f.DateFin.Sub(f.DateDebut)
		if j1.After(f.DateDebut) {
			debFrais = j1
		} else {
			debFrais = f.DateDebut
		}
		if j2.Before(f.DateFin) {
			finFrais = j2
		} else {
			finFrais = f.DateFin
		}
		dureeFraisPeriode = finFrais.Sub(debFrais)
		// contribution du frais sur la période [j1, j2]
		total += (dureeFraisPeriode.Hours() * f.Montant / dureeFraisTotale.Hours())
	}
	return total, nil
}
*/

// ************************** CRUD *******************************

func InsertStockage(db *sqlx.DB, s *Stockage) (id int, err error) {
	query := `insert into stockage(nom) values($1) returning id`
	err = db.QueryRow(
		query,
		s.Nom).Scan(&id)
	if err != nil {
		return 0, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

func UpdateStockage(db *sqlx.DB, s *Stockage) (err error) {
	query := `update stockage set(
	    nom,
	    archived
	    ) = ($1,$2) where id=$3`
	_, err = db.Exec(
		query,
		s.Nom,
		s.Archived,
		s.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeleteStockage(db *sqlx.DB, id int) (err error) {
	query := "delete from stockfrais where id_stockage=$1"
	_, err = db.Exec(query, id)
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
