/*
Fonctions liées à la recherche, communes à plusieurs recherches.

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2023-04-22 07:38:21+02:00, Thierry Graff : Creation
*/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
	// "fmt"
)

func filtreValoContientAutreValo(filtreValo []string) bool {
    codes := AllValorisationCodes() // codes utilisés par chautre - ne contient pas PQ et CF
    for _, code := range(codes){
        if tiglib.InArrayString(code, filtreValo) {
            return true
        }
    }
    return false
}

func ComputeRecapFiltres(db *sqlx.DB, filtres map[string][]string) (result string, err error) {
	result = ""
	// Si aucun filtre
	aucun := true
	for k, _ := range filtres {
		if len(filtres[k]) != 0 {
			aucun = false
			break
		}
	}
	if aucun {
		return "Aucun filtre, tout est affiché", nil
	}
	//
	result += "<table>\n"
	//
	if len(filtres["periode"]) != 0 {
		deb, err := time.Parse("2006-01-02", filtres["periode"][0])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel time.Parse("+filtres["periode"][0]+")")
		}
		strDeb := tiglib.DateFr(deb)
		//
		fin, err := time.Parse("2006-01-02", filtres["periode"][1])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel time.Parse("+filtres["periode"][1]+")")
		}
		strFin := tiglib.DateFr(fin)
		result += "<tr><td>Période :</td><td>" + strDeb + " - " + strFin + "</td></tr>\n"
	}
	if len(filtres["proprio"]) != 0 {
		tmp := []string{} // Comme il n'y a que 2 propriétaires, tmp ne contient qu'un élément - mais code écrit pour un cas plus général
		for _, value := range filtres["proprio"] {
			id, _ := strconv.Atoi(value)
			proprio, err := GetActeur(db, id)
			if err != nil {
				return result, werr.Wrapf(err, "Erreur appel GetActeur()")
			}
			tmp = append(tmp, "<a href=\"/acteur/"+strconv.Itoa(proprio.Id)+"\">"+proprio.String()+"</a>")
		}
		result += "<tr><td>Propriétaire :</td><td>" + strings.Join(tmp, ", ") + "</td></tr>\n"
	}
	//
	if len(filtres["fermier"]) != 0 {
		id, _ := strconv.Atoi(filtres["fermier"][0])
		fermier, err := GetFermier(db, id)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel GetFermier()")
		}
		result += "<tr><td>Fermier :</td><td><a href=\"/fermier/" + strconv.Itoa(fermier.Id) + "\">" + fermier.String() + "</a>" + "</td></tr>\n"
	}
	//
	if len(filtres["essence"]) != 0 {
		tmp := []string{}
		essenceMap, err := GetEssencesMap(db)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel GetEssencesMap()")
		}
		for _, code := range filtres["essence"] {
			essenceLabel := essenceMap[code]
			tmp = append(tmp, essenceLabel)
		}
		result += "<tr><td>Essences :</td><td>" + strings.Join(tmp, ", ") + "</td></tr>\n"
	}
	//
	if len(filtres["ug"]) != 0 {
		tmp := []string{}
		for _, value := range filtres["ug"] {
			id, _ := strconv.Atoi(value)
			ug, err := GetUG(db, id)
			if err != nil {
				return result, werr.Wrapf(err, "Erreur appel GetUG()")
			}
			tmp = append(tmp, "<a href=\"/ug/"+strconv.Itoa(ug.Id)+"\">"+ug.String()+"</a>")
		}
		result += "<tr><td>UGs :</td><td>" + strings.Join(tmp, ", ") + "</td></tr>\n"
	}
	//
	if len(filtres["parcelle"]) != 0 {
		tmp := []string{}
		for _, value := range filtres["parcelle"] {
			id, _ := strconv.Atoi(value)
			parcelle, err := GetParcelle(db, id)
			if err != nil {
				return result, werr.Wrapf(err, "Erreur appel GetParcelle()")
			}
			tmp = append(tmp, "<a href=\"/parcelle/"+strconv.Itoa(parcelle.Id)+"\">"+parcelle.String()+"</a>")
		}
		result += "<tr><td>Parcelles :</td><td>" + strings.Join(tmp, ", ") + "</td></tr>\n"
	}
	//
	result += "</table>\n"
	//
	return result, nil
}
