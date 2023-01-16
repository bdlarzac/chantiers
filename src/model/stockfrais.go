/******************************************************************************
    Frais des hangars pour stocker des plaquettes

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2019-12-03 16:27:28+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"time"

	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type StockFrais struct {
	Id         int
	IdStockage int `db:"id_stockage"`
	TypeFrais  string
	Montant    float64
	DateDebut  time.Time `db:"datedeb"`
	DateFin    time.Time `db:"datefin"`
	Notes      string
}

// *********************************************************
func GetStockFrais(db *sqlx.DB, id int) (sf *StockFrais, err error) {
	sf = &StockFrais{}
	query := "select * from stockfrais where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(sf)
	if err != nil {
		return sf, werr.Wrapf(err, "Erreur query : "+query)
	}
	return sf, nil
}

// *********************************************************
func InsertStockFrais(db *sqlx.DB, sf *StockFrais) (id int, err error) {
	query := `insert into stockfrais(
	    id_stockage,
	    typefrais,
	    montant,
	    datedeb,
	    datefin,
	    notes
	    ) values($1,$2,$3,$4,$5,$6) returning id`
	err = db.QueryRow(
		query,
		sf.IdStockage,
		sf.TypeFrais,
		sf.Montant,
		sf.DateDebut,
		sf.DateFin,
		sf.Notes).Scan(&id)
	if err != nil {
		return 0, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

// *********************************************************
func UpdateStockFrais(db *sqlx.DB, sf *StockFrais) (err error) {
	query := `update stockfrais set(
	    typefrais,
        montant,
        datedeb,
        datefin,
        notes
        ) = ($1,$2,$3,$4,$5) where id=$6`
	_, err = db.Exec(
		query,
		sf.TypeFrais,
		sf.Montant,
		sf.DateDebut,
		sf.DateFin,
		sf.Notes,
		sf.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

// *********************************************************
func DeleteStockFrais(db *sqlx.DB, id int) (err error) {
	query := "delete from stockfrais where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
