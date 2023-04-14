/*
*
*****************************************************************************

	Essence (= espèce d'arbre)

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-03-03 08:57:51+01:00, Thierry Graff : Creation

*******************************************************************************
*
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type Essence struct {
	Code    string
	Nom     string
	NomLong string
}

/*
	Renvoie une map code essence => nom long
*/
func GetEssencesMap(db *sqlx.DB) (essencesMap map[string]string, err error) {
	type essence struct {
		Code    string
		NomLong string
	}
	essences := []essence{}
	query := "select code,nomlong from essence"
	err = db.Select(&essences, query)
	if err != nil {
		return essencesMap, werr.Wrapf(err, "Erreur query : "+query)
	}
	essencesMap = make(map[string]string, len(essences))
	for _, essence := range essences {
		essencesMap[essence.Code] = essence.NomLong
	}
	return essencesMap, nil
}

/*
   Renvoie un tableau de codes essence
*/
/*
// pas utilisée, peut être supprimée
func GetEssenceCodes(db *sqlx.DB) (result []string, err error) {
    result = []string{}
	query := "select code from essence"
	err = db.Select(&result, query)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query : "+query)
	}
    return result, nil
}
*/
