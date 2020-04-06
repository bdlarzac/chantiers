/******************************************************************************
    Livraison de plaquettes à un client - inclut le chargement

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-01-22 02:56:23+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"time"

	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type PlaqLivre struct {
	Id         int
	IdVente    int     `db:"id_vente"`
	IdChargeur int     `db:"id_chargeur"`
	IdLivreur  int     `db:"id_livreur"`
	Qte        float64 // m3
	// chargement
	ChDate      time.Time
	ChOutil     string
	ChPrixOutil float64 // prix ht
	ChPrixMOH   float64 // prix M.O. ht / heure
	ChDatePay   time.Time
	ChNotes     string
	// livraison
	LiDate    time.Time
	LiPrixMOH float64 // prix M.O. ht / heure
	LiNheure  float64
	LiDatePay time.Time
	LiNotes   string
	// Pas stocké en base
	Vente    *VentePlaq
	Chargeur *Acteur
	Livreur  *Acteur
}

// ************************** Get *******************************

func GetPlaqLivre(db *sqlx.DB, id int) (*PlaqLivre, error) {
	pl := &PlaqLivre{}
	query := "select * from plaqlivre where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(pl)
	if err != nil {
		return pl, werr.Wrapf(err, "Erreur query : "+query)
	}
	return pl, nil
}

// ************************** CRUD *******************************

func InsertPlaqLivre(db *sqlx.DB, pl *PlaqLivre) (int, error) {
	query := `insert into plaqlivre(
        id_vente,
        id_chargeur,
        id_livreur,
        qte,
        chdate,
        choutil,
        chprixoutil,
        chprixmoh,
        chdatepay,
        chnotes,
        lidate,
        liprixmoh,
        linheure,
        lidatepay,
        linotes)
        values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		pl.IdVente,
		pl.IdChargeur,
		pl.IdLivreur,
		pl.Qte,
		pl.ChDate,
		pl.ChOutil,
		pl.ChPrixOutil,
		pl.ChPrixMOH,
		pl.ChDatePay,
		pl.ChNotes,
		pl.LiDate,
		pl.LiPrixMOH,
		pl.LiNheure,
		pl.LiDatePay,
		pl.LiNotes).Scan(&id)
	return id, err
}

func UpdatePlaqLivre(db *sqlx.DB, pl *PlaqLivre) error {
	query := `update plaqlivre set(
        id_vente,
        id_chargeur,
        id_livreur,
        qte,
        chdate,
        choutil,
        chprixoutil,
        chprixmoh,
        chdatepay,
        chnotes,
        lidate,
        liprixmoh,
        linheure,
        lidatepay,
        linotes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) where id=$16`
	_, err := db.Exec(
		query,
		pl.IdVente,
		pl.IdChargeur,
		pl.IdLivreur,
		pl.Qte,
		pl.ChDate,
		pl.ChOutil,
		pl.ChPrixOutil,
		pl.ChPrixMOH,
		pl.ChDatePay,
		pl.ChNotes,
		pl.LiDate,
		pl.LiPrixMOH,
		pl.LiNheure,
		pl.LiDatePay,
		pl.LiNotes,
		pl.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeletePlaqLivre(db *sqlx.DB, id int) error {
	query := "delete from plaqlivre where id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
