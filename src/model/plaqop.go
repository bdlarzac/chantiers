/*
*****************************************************************************

	Opération simple associée à un chantier plaquette
	Peut être abattage, débardage, déchiquetage, broyage

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2020-01-15 02:37:25+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"time"

	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type PlaqOp struct {
	Id         int
	TypOp      string
	IdChantier int       `db:"id_chantier"`
	IdActeur   int       `db:"id_acteur"`
	DateDebut  time.Time `db:"datedeb"`
	DateFin    time.Time
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

func GetPlaqOp(db *sqlx.DB, id int) (op *PlaqOp, err error) {
	op = &PlaqOp{}
	query := "select * from plaqop where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(op)
	if err != nil {
		return op, werr.Wrapf(err, "Erreur query : "+query)
	}
	return op, nil
}

// ************************** Compute *******************************

// Remplit le champ Acteur d'une opération
func (op *PlaqOp) ComputeActeur(db *sqlx.DB) (err error) {
	if op.Acteur != nil {
		return nil // déjà calculé
	}
	op.Acteur, err = GetActeur(db, op.IdActeur)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

// ************************** CRUD *******************************

func InsertPlaqOp(db *sqlx.DB, op *PlaqOp) (id int, err error) {
	query := `insert into plaqop(
	    typop,
	    id_chantier,
	    id_acteur,
	    datedeb,
	    datefin,
	    qte,
	    unite,
	    puht,
	    tva,                                                              
	    datepay,
	    notes
	    ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning id`
	err = db.QueryRow(
		query,
		op.TypOp,
		op.IdChantier,
		op.IdActeur,
		op.DateDebut,
		op.DateFin,
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

func UpdatePlaqOp(db *sqlx.DB, op *PlaqOp) (err error) {
	query := `update plaqop set(
        typop,
        id_chantier,
        id_acteur,
        datedeb,
        datefin,
        qte,
        unite,
        puht,
        tva,
        datepay,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) where id=$12`
	_, err = db.Exec(
		query,
		op.TypOp,
		op.IdChantier,
		op.IdActeur,
		op.DateDebut,
		op.DateFin,
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

func DeletePlaqOp(db *sqlx.DB, id int) (err error) {
	query := "delete from plaqop where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
