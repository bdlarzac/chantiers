/*
@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
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
