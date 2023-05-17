/*
Calcul d'activités regroupées par UG.

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2023-05-12 12:02:12+02:00, Thierry Graff : Creation
*/
package model

import (
    "sort"
//	"github.com/jmoiron/sqlx"
//"fmt"
)

/* Liste d'activités ayant lieu dans une UG donnée */
type ActivitesParUG struct {
	UG        *UG
	Activites []*Activite
}

func ComputeActivitesParUG(activites []*Activite) (result []*ActivitesParUG) {
	result = []*ActivitesParUG{}
	mapug := map[string]*ActivitesParUG{} // map intermédiaire pour regrouper par UG - clé = code ug
    for _, activite := range(activites){
//fmt.Printf("len(activite.UGs) = %d\n",len(activite.UGs))
        for _, ug := range(activite.UGs){
            code := ug.Code
            if _, ok := mapug[code]; !ok {
                //mapug[code] = *ActivitesParUG{}
                empty := ActivitesParUG{}
                mapug[code] = &empty
                mapug[code].UG = ug
            }
            mapug[code].Activites = append(mapug[code].Activites, activite)
        }
    }
/* for k, v := range(mapug){
    fmt.Printf("=== %s ===\n",k)
    for _, a := range(v.Activites){
        fmt.Printf("%d - %s - %s\n",a.Id, a.TypeActivite, a.Titre)
    }
} */
    keys := []string{}
    for k, _ := range(mapug){
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, key := range(keys){
        result = append(result, mapug[key])
    }
//fmt.Printf("mapug = %+v\n",mapug)
//fmt.Printf("%+v\n",result)
	return result
}

