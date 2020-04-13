/******************************************************************************
    Tas de plaquettes dans un hangar

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-03-09 17:18:27+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Tas struct {
	Id         int
	IdStockage int `db:"id_stockage"`
	IdChantier int `db:"id_chantier"`
	Stock      float64
	Actif      bool
	// pas stocké en base
	Nom             string
	Chantier        *Plaq
	Stockage        *Stockage
	MesuresHumidite []*Humid
}

func NewTas(idStockage, idChantier int, stock float64, actif bool) *Tas {
	return &Tas{
		IdStockage: idStockage,
		IdChantier: idChantier,
		Stock:      stock,
		Actif:      actif,
	}
}

// ************************** Manipulation Stock *******************************

// Si qte > 0, ajoute des plaquettes au tas
// Si qte < 0, retire des plaquettes au tas
// Fait la maj en BDD
// @param   qte en maps
func (t *Tas) ModifierStock(db *sqlx.DB, qte float64) error {
	t.Stock += qte
	return UpdateTas(db, t)
}

// Pour indiquer qu'un tas est vide
func DesactiverTas(db *sqlx.DB, id int) error {
	query := `update tas set actif=false where id=$1`
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

// ************************** Get one *******************************

func GetTas(db *sqlx.DB, idTas int) (*Tas, error) {
	tas := &Tas{}
	query := "select * from tas where id=$1"
	row := db.QueryRowx(query, idTas)
	err := row.StructScan(tas)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur query : "+query)
	}
	return tas, nil
}

func GetTasFull(db *sqlx.DB, idTas int) (*Tas, error) {
	tas, err := GetTas(db, idTas)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur appel GetTas()")
	}
	err = tas.ComputeStockage(db)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur appel Tas.ComputeStockage()")
	}
	err = tas.ComputeChantier(db)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur appel Tas.ComputeChantier()")
	}
	err = tas.ComputeNom(db)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur appel Tas.ComputeNom()")
	}
	return tas, nil
}

// ************************** Get many *******************************

// Utilisé pour faire un select html
// Obligé d'avoir tas full, car besoin du nom du tas, qui a besoin de chantier et stockage
func GetAllTasActifsFull(db *sqlx.DB) ([]*Tas, error) {
	tas := []*Tas{}
	ids := []int{}
	query := "select id from tas where actif"
	err := db.Select(&ids, query)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, id := range ids {
		t, err := GetTasFull(db, id)
		if err != nil {
			return tas, werr.Wrapf(err, "Erreur appel GetTasFull()")
		}
		tas = append(tas, t)
	}
	return tas, nil
}

// ************************** Compute *******************************

func (t *Tas) ComputeNom(db *sqlx.DB) error {
	if t.Chantier == nil || t.Stockage == nil {
		return errors.New("Impossible de calculer le nom du tas - appeler d'abord ComputeStockage() et ComputeChantier()")
	}
	err := t.Chantier.ComputeLieudit(db) // Pour le nom du chantier
	if err != nil {
		return werr.Wrapf(err, "Erreur appel t.Chantier.ComputeLieudit()")
	}
	t.Nom = t.Chantier.String() + " - " + t.Stockage.Nom
	return nil
}

func (t *Tas) ComputeStockage(db *sqlx.DB) error {
	var err error
	t.Stockage, err = GetStockage(db, t.IdStockage)
	return err
}

func (t *Tas) ComputeChantier(db *sqlx.DB) error {
	var err error
	t.Chantier, err = GetPlaq(db, t.IdChantier)
	return err
}

// pas inclus par défaut dans GetTasFull()
func (t *Tas) ComputeMesuresHumidite(db *sqlx.DB) error {
	var err error
	mesures := []*Humid{}
	query := "select * from humid where id_tas=$1"
	err = db.Select(&mesures, query, t.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	t.MesuresHumidite = mesures
	return err
}

// ************************** CRUD *******************************

func InsertTas(db *sqlx.DB, tas *Tas) (int, error) {
	query := `insert into tas(
        id_stockage,                              
        id_chantier,
        stock,
        actif
        ) values($1,$2,$3,$4) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		tas.IdStockage,
		tas.IdChantier,
		tas.Stock,
		tas.Actif).Scan(&id)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

func UpdateTas(db *sqlx.DB, tas *Tas) error {
	query := `update tas set(
        id_stockage,
        id_chantier,
        stock,
        actif
        ) = ($1,$2,$3,$4) where id=$5`
	_, err := db.Exec(
		query,
		tas.IdStockage,
		tas.IdChantier,
		tas.Stock,
		tas.Actif,
		tas.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeleteTas(db *sqlx.DB, id int) error {
	query := "delete from tas where id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
