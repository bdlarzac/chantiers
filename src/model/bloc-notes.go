/*
Bloc-notes = zone de texte qui sert de pense-bête
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

func UpdateBlocnotes(db *sqlx.DB, contenu string) (err error) {
	query := "update blocnotes set contenu=$1"
	_, err = db.Exec(query, contenu)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	return nil
}

func GetBlocnotes(db *sqlx.DB) (contenu string, err error) {
	query := "select contenu from blocnotes"
    err = db.QueryRow(query).Scan(&contenu)
	if err != nil {
		return contenu, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	return contenu, nil
}
