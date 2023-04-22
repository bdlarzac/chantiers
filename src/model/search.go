/*
    Fonctions liées à la recherche

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2023-04-22 07:38:21+02:00, Thierry Graff : Creation
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/generic/tiglib"
	"github.com/jmoiron/sqlx"
	"time"
	//"strconv"
"fmt"
)

func ComputeRecapFiltres(db *sqlx.DB, filtres map[string][]string) (result string, err error){
    result = ""
    // Si aucun filtre
    aucun := true
    for k, _ := range(filtres){
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
    for k, filtre := range(filtres){
fmt.Printf("%s : %+v\n", k, filtre)
        if len(filtre) == 0 {
            continue
        }
        switch k {
        case "periode":
            deb, err := time.Parse("2006-01-02", filtre[0])
            if err != nil {
                return result, werr.Wrapf(err, "Erreur appel time.Parse("+filtre[0]+")")
            }
            strDeb := tiglib.DateFr(deb)
            //
            fin, err := time.Parse("2006-01-02", filtre[1])
            if err != nil {
                return result, werr.Wrapf(err, "Erreur appel time.Parse("+filtre[1]+")")
            }
            strFin := tiglib.DateFr(fin)
            result += "<tr><td>Période :</td><td>"+strDeb+" - "+strFin+"</td></tr>\n"
        break
        case "proprio":
            result += "<tr><td>Propriétaire :</td><td>"+""+"</td></tr>\n"
        break
        }
    }
    result += "</table>\n"
    //
    return result, nil
}

