/*
*****************************************************************************

	Fonctions liées aux rôles des acteurs

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-04-04 16:49:32+02:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

func GetRolesMap(db *sqlx.DB) (rolesMap map[string]string, err error) {
	type role struct {
		Code string
		Nom  string
	}
	roles := []role{}
	query := "select * from role"
	err = db.Select(&roles, query)
	if err != nil {
		return rolesMap, werr.Wrapf(err, "Erreur query : "+query)
	}
	rolesMap = make(map[string]string, len(roles))
	for _, role := range roles {
		rolesMap[role.Code] = role.Nom
	}
	return rolesMap, nil
}
