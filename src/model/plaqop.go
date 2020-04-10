/******************************************************************************
    Opération simple associée à un chantier plaquette
    Peut être abattage, débardage, déchiquetage, broyage

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-01-15 02:37:25+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"time"

	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type PlaqOp struct {
	Id         int
	TypOp      string // faire un type ?
	IdChantier int    `db:"id_chantier"`
	IdActeur   int    `db:"id_acteur"`
	DateOp     time.Time
	Qte        float64
	Unite      string // j, h ou m
	PUHT       float64
	TVA        float64
	DatePay    time.Time
	Notes      string
	// Pas stocké en base
	Chantier *Plaq
	Acteur   *Acteur
}

// ************************** Instance methods *******************************
// Rôle de l'Acteur d'une opération simple
func (op *PlaqOp) RoleName() string {
	switch op.TypOp {
	case "AB":
		return "abatteur"
	case "BR":
		return "broyeur"
	case "DB":
		return "débardeur"
	case "DC":
		return "déchiqueteur"
	}
	return "??? Rôle inconnu ???"
}

// ************************** Get *******************************

func GetPlaqOp(db *sqlx.DB, id int) (*PlaqOp, error) {
	op := &PlaqOp{}
	query := "select * from plaqop where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(op)
	if err != nil {
		return op, werr.Wrapf(err, "Erreur query : "+query)
	}
	return op, nil
}

// ************************** Compute *******************************

// Remplit le champ Acteur d'une opération
func (op *PlaqOp) ComputeActeur(db *sqlx.DB) error {
	var err error
	op.Acteur, err = GetActeur(db, op.IdActeur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return err
}

// ************************** CRUD *******************************

func InsertPlaqOp(db *sqlx.DB, op *PlaqOp) (int, error) {
	query := `insert into plaqop(
	    typop,
	    id_chantier,
	    id_acteur,
	    dateop,
	    qte,
	    unite,
	    puht,
	    tva,                                                              
	    datepay,
	    notes
	    ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		op.TypOp,
		op.IdChantier,
		op.IdActeur,
		op.DateOp,
		op.Qte,
		op.Unite,
		op.PUHT,
		op.TVA,
		op.DatePay,
		op.Notes).Scan(&id)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

func UpdatePlaqOp(db *sqlx.DB, op *PlaqOp) error {
	query := `update plaqop set(
        typop,
        id_chantier,
        id_acteur,
        dateop,
        qte,
        unite,
        puht,
        tva,
        datepay,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) where id=$11`
	_, err := db.Exec(
		query,
		op.TypOp,
		op.IdChantier,
		op.IdActeur,
		op.DateOp,
		op.Qte,
		op.Unite,
		op.PUHT,
		op.TVA,
		op.DatePay,
		op.Notes,
		op.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeletePlaqOp(db *sqlx.DB, id int) error {
	query := "delete from plaqop where id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}