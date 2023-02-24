/*
*****************************************************************************

	La table facture permet la génération automatique du numéro de facture.

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-02-23 17:55:53+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"fmt"
	"github.com/jmoiron/sqlx"
)

/*
*

	Si l'année demandée est déjà présente en base :
	    - récupère lastnum,
	    - incrémente de 1,
	    - enregistre la nouvelle valeur de lastnum
	Si l'année demandée n'est pas déjà présente en base :
	    - crée une nouvelle ligne avec l'année et lastnum à 1
	Dans les deux cas, renvoie la string avec le numéro d'affacture.

	@param  annee   Format AAAA, ex 2023
	@return         String du genre "2023054"

*
*/
func NouveauNumeroFacture(db *sqlx.DB, annee string) (result string, err error) {
	var lastnum int
	query := "select lastnum from facture where annee=$1"
	_ = db.Get(&lastnum, query, annee) // empty => lastnum reste = 0
	// lastnum = 0 si nouvelle année
	lastnum++
	if lastnum == 1 {
		// année pas présente en base
		query = "insert into facture(annee,lastnum) values($1, $2)"
		_, err = db.Exec(query, annee, lastnum)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur query DB : "+query)
		}
	} else {
		// année déjà présente en base
		query = "update facture set lastnum = $1 where annee=$2"
		_, err = db.Exec(query, lastnum, annee)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur query DB : "+query)
		}
	}
	return fmt.Sprintf("%s%03d", annee, lastnum), nil
}
