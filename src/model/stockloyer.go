/******************************************************************************
    Loyers des hangars pour stocker des plaquettes

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-12-03 16:27:28+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"time"

	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type StockLoyer struct {
	Id         int
	IdStockage int `db:"id_stockage"`
	Montant    float64
	DateDebut  time.Time `db:"datedeb"`
	DateFin    time.Time `db:"datefin"`
}

// *********************************************************
func GetStockLoyer(db *sqlx.DB, id int) (*StockLoyer, error) {
	ls := &StockLoyer{}
	query := "select * from stockloyer where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(ls)
	if err != nil {
		return ls, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ls, nil
}

// *********************************************************
func InsertStockLoyer(db *sqlx.DB, sl *StockLoyer) (int, error) {
	query := `insert into stockloyer(
	    id_stockage,
	    montant,
	    datedeb,
	    datefin
	    ) values($1, $2, $3, $4) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		sl.IdStockage,
		sl.Montant,
		sl.DateDebut,
		sl.DateFin).Scan(&id)
	return id, err
}

// *********************************************************
func UpdateStockLoyer(db *sqlx.DB, sl *StockLoyer) error {
	query := `update stockloyer set(
        montant,
        datedeb,
        datefin
        ) = ($1,$2,$3) where id=$4`
	_, err := db.Exec(
		query,
		sl.Montant,
		sl.DateDebut,
		sl.DateFin,
		sl.Id)
	return err
}

// *********************************************************
func DeleteStockLoyer(db *sqlx.DB, id int) error {
	query := "delete from stockloyer where id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
