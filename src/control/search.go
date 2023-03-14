/*
@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"net/http"
	"time"
)

type detailsSearch struct {
	Periods   [][2]time.Time // pour choix-date
	UrlAction string
}

/*
Affiche / process le formulaire de recherche
*/
func Search(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	ctx.TemplateName = "search.html"

	//periods, hasChantier, err := model.ComputeLimitesSaisons(ctx.DB, ctx.Config.DebutSaison)
	periods, _, err := model.ComputeLimitesSaisons(ctx.DB, ctx.Config.DebutSaison)
	if err != nil {
		return err
	}
	//fmt.Printf("periods = %+v\n",periods)

	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Accueil",
			CSSFiles: []string{
				"/static/css/form.css"},
		},
		Menu: "accueil",
		Details: detailsSearch{
			Periods:   periods,
			UrlAction: "/search",
		},
	}
	return nil
}
