/******************************************************************************
    Frais des hangars pour stocker des plaquettes

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

type StockFrais struct {
	Id         int
	IdStockage int `db:"id_stockage"`
	TypeFrais  string
	Montant    float64
	DateDebut  time.Time `db:"datedeb"`
	DateFin    time.Time `db:"datefin"`
	Notes string
}

// *********************************************************
func GetStockFrais(db *sqlx.DB, id int) (*StockFrais, error) {
	sf := &StockFrais{}
	query := "select * from stockfrais where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(sf)
	if err != nil {
		return sf, werr.Wrapf(err, "Erreur query : "+query)
	}
	return sf, nil
}

// *********************************************************
func InsertStockFrais(db *sqlx.DB, sf *StockFrais) (int, error) {
	query := `insert into stockfrais(
	    id_stockage,
	    typefrais,
	    montant,
	    datedeb,
	    datefin,
	    notes
	    ) values($1,$2,$3,$4,$5,$6) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		sf.IdStockage,
		sf.TypeFrais,
		sf.Montant,
		sf.DateDebut,
		sf.DateFin,
		sf.Notes).Scan(&id)
	return id, err
}

// *********************************************************
func UpdateStockFrais(db *sqlx.DB, sf *StockFrais) error {
	query := `update stockfrais set(
	    typefrais,
        montant,
        datedeb,
        datefin,
        notes
        ) = ($1,$2,$3,$4,$5) where id=$6`
	_, err := db.Exec(
		query,
		sf.TypeFrais,
		sf.Montant,
		sf.DateDebut,
		sf.DateFin,
		sf.Notes,
		sf.Id)
	return err
}

// *********************************************************
func DeleteStockFrais(db *sqlx.DB, id int) error {
	query := "delete from stockfrais where id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
