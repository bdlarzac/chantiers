/*
@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"strings"
	"time"
	"net/http"
"fmt"	
)

type detailsSearchForm struct {
	Periods     [][2]time.Time    // pour choix-date
	EssencesMap map[string]string // pour choix-essence
	PropriosMap map[int]string    // pour choix-proprio
	UrlAction   string
}

type detailsSearchResults struct{
}

/*
    Affiche / process le formulaire de recherche
*/
func Search(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	switch r.Method {
	case "POST":
		//
		// Process form et affiche page de résultats
		//
        if err = r.ParseForm(); err != nil {
            return err
        }
fmt.Println("ici")
fmt.Printf("%+v\n",r.PostForm)
        //
//        filtre_essence := computeFiltreEssence(r)
//fmt.Printf("filtre_essence = %+v\n",filtre_essence)
        filtre_proprio := computeFiltreProprio(r)
fmt.Printf("filtre_proprio = %+v\n",filtre_proprio)
        //
        ctx.TemplateName = "search-result.html"
        ctx.Page = &ctxt.Page{
            Header: ctxt.Header{
                Title: "Activités",
            },
            Menu: "accueil",
            Details: detailsSearchResults{
            },
        }
        return nil
	default:
		//
		// Affiche form
		//
        //periods, hasChantier, err := model.ComputeLimitesSaisons(ctx.DB, ctx.Config.DebutSaison)
        periods, _, err := model.ComputeLimitesSaisons(ctx.DB, ctx.Config.DebutSaison)
        if err != nil {
            return err
        }
        //
        essencesMap, err := model.GetEssencesMap(ctx.DB)
        if err != nil {
            return err
        }
        //
        propriosMap, err := model.GetProprietaires(ctx.DB)
        if err != nil {
            return err
        }
        //
        ctx.TemplateName = "search-form.html"
        ctx.Page = &ctxt.Page{
            Header: ctxt.Header{
                Title: "Recherche d'activité",
                CSSFiles: []string{
                    "/static/css/form.css"},
            },
            Menu: "accueil",
            Details: detailsSearchForm{
                Periods:   periods,
                EssencesMap: essencesMap,
                PropriosMap: propriosMap,
                UrlAction: "/search",
            },
        }
        return nil
    }
}

// ************************** Auxiliaires *******************************

/* 
    Filtre essence :
        - Soit contient un seul élément "ALL" => pas de filtre.
        - Soit contient une liste de codes essence.
*/
func computeFiltreEssence(r *http.Request) (result []string) {
    if r.PostFormValue("choix-ALL-essence") == "true"{
        return []string{"ALL"}
    }
    result = []string{}
    for key, _ := range r.PostForm {
        if strings.Index(key, "choix-essence-") != 0 {
            continue
        }
        if r.PostFormValue(key) != "on" {
            continue
        }
        code := key[14:]
        result = append(result, code)
    }
    return result
}

/* 
    Filtre propriétaire :
        - Soit contient un seul élément "ALL" => pas de filtre.
        - Soit contient une liste d'id propriétaires (dans la table acteur)
          (attention, cette liste contient des strings, pas des ints).
*/
func computeFiltreProprio(r *http.Request) (result []string) {
    if r.PostFormValue("choix-ALL-proprio") == "true"{
        return []string{"ALL"}
    }
    result = []string{}
    for key, _ := range r.PostForm {
        if strings.Index(key, "choix-proprio-") != 0 {
            continue
        }
        if r.PostFormValue(key) != "on" {
            continue
        }
        id := key[14:]
        result = append(result, id)
    }
    return result
}