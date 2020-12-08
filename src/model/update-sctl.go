/******************************************************************************
    Mise à jour des données foncières :
    Prise en compte des mises à jour effectuées sur la base SCTL
    Compare l'état de la base BDL (qui est le résultat d'un import d'une version antérieure de la base SCTL)
    avec une nouvelle version de la base SCTL

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-04-18 14:20:54+02:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

// Chaque item indique une entité de la base à mettre à jour
type UpdatedItem struct {
	Type   string      // "acteur"
	Action string      // "insert" "update" "delete"
	Value  interface{} // contient les nouvelles valeurs de l'entité
}

// Calcule les entités qui seront mises à jour si l'utilisateur effectue la maj.
func ComputeUpdateSCTL(db *sqlx.DB, conf *Config) ([]*UpdatedItem, error) {
	res := []*UpdatedItem{}
	newActeurs, err := loadNewActeurs(conf)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur appel loadNewActeurs()")
	}
	oldActeurs, err := loadOldActeurs(db)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur appel loadNewActeurs()")
	}
	/*
	   fmt.Println("========= old ==========")
	       for idOld, a := range(oldActeurs){
	   fmt.Println(idOld, a.Nom, a.Prenom)
	       }
	   fmt.Println("========= new ==========")
	       for idNew, a := range(newActeurs){
	   fmt.Println(idNew, a.Nom, a.Prenom)
	       }
	*/
	// ids sctl des acteurs présents dans new et old,
	// donc possiblement modifiés
	toCheckUpdate := []int{}
	//
	// Acteurs ayant été supprimés
	//
	for idOld, a := range oldActeurs {
		if _, ok := newActeurs[idOld]; !ok {
			res = append(res, &UpdatedItem{
				Type:   "acteur",
				Action: "delete",
				Value:  a,
			})
		} else {
			toCheckUpdate = append(toCheckUpdate, idOld)
		}
	}
	//
	// Acteurs ayant été ajoutés
	//
	for idNew, a := range newActeurs {
		if _, ok := oldActeurs[idNew]; !ok {
			res = append(res, &UpdatedItem{
				Type:   "acteur",
				Action: "add",
				Value:  a,
			})
		} else {
			toCheckUpdate = append(toCheckUpdate, idNew)
		}
	}
fmt.Println(toCheckUpdate)
	toCheckUpdate = tiglib.ArrayUniqueInt(toCheckUpdate)
fmt.Println(toCheckUpdate)
	//
	// Acteurs ayant été modifiés
	//
	// @todo Acteurs ayant été modifiés
	return res, nil
}

// Charge les acteurs de la base BDL
// Clé dans la map renvoyée = id sctl
func loadOldActeurs(db *sqlx.DB) (map[int]Acteur, error) {
	res := map[int]Acteur{}
	acteurs := []*Acteur{}
	query := "select * from acteur where id_sctl<>0 order by id_sctl"
	err := db.Select(&acteurs, query)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, a := range acteurs {
		res[a.IdSctl] = *a
	}
	return res, nil
}

// Charge les acteurs de la base SCTL
// Clé dans la map renvoyée = id sctl
//
// @todo Changer le code pour lire dans le .mdb, au lieu du .csv
//
func loadNewActeurs(conf *Config) (map[int]Acteur, error) {
	res := map[int]Acteur{}
	records, err := tiglib.CsvMap(conf.Paths.LogicielFoncier, ';')
	if err != nil {
		return res, werr.Wrapf(err, "Erreur appel tiglib.CsvMap()")
	}
	for _, record := range records {
		idExploitant, err := strconv.Atoi(record["IdExploitant"])
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel strconv.Atoi(%s) pour IdExploitant", record["IdExploitant"])
		}
		cp := record["CPExp"]
		if len(cp) > 5 {
			cp = cp[:5] // fix une typo dans la base SCTL
		}
		res[idExploitant] = Acteur{
			IdSctl:   idExploitant,
			Nom:      record["NOMEXP"],
			Prenom:   record["Prenom"],
			Adresse1: record["AdresseExp"],
			Cp:       cp,
			Ville:    record["VilleExp"],
			Tel:      record["Telephone"],
			Email:    record["Mail"],
		}
	}
	return res, nil
}
