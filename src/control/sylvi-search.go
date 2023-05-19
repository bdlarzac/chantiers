/*
Recherche / bilans liés à la sylviculture

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"net/http"
)

type detailsSylviForm struct {
	EssenceCodes []string         // pour choix-essence
	Fermiers     []*model.Fermier // pour choix-fermier
	AllCommunes  []*model.Commune // pour choix-parcelle
	UrlAction    string
}

type detailsSylviResults struct {
	UGs          []*model.UG
	RecapFiltres string
////// supprimer si finalement pas de tab
	Tab          string
}

/*
Affiche / process le formulaire de recherche
*/
func SearchSylvi(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	switch r.Method {
	case "POST":
		//
		// Process form et affiche page de résultats
		//
		if err = r.ParseForm(); err != nil {
			return err
		}
////// supprimer si finalement pas de tab
		vars := mux.Vars(r)
		tab := vars["tab"]
		if tab == "" {
            tab = "liste"
		}
		//
		filtres := map[string][]string{}
		filtres["fermier"] = computeFiltreFermier(r)
		filtres["essence"] = computeFiltreEssence(r)
		filtres["commune"] = computeFiltreCommune(r)
		ugs, err := model.ComputeUGsFromFiltres(ctx.DB, filtres)
		if err != nil {
			return err
		}
		//
		recapFiltres, err := model.ComputeRecapFiltres(ctx.DB, filtres)
		if err != nil {
			return err
		}
		//
		ctx.TemplateName = "sylvi-search-show.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Sylviculture",
				CSSFiles: []string{"/static/lib/tabstrip/tabstrip.css"},
				// JSFiles: []string{"/static/js/formatNb.js"},
			},
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/lib/tabstrip/tabstrip.js",
					"/static/lib/table-sort/table-sort.js"},
			},
			Menu: "production",
			Details: detailsSylviResults{
				UGs:          ugs,
				RecapFiltres: recapFiltres,
////// supprimer si finalement pas de tab
				Tab:                      tab,
			},
		}
		return nil
	default:
		//
		// Affiche form
		//
////// supprimer si finalement pas de tab
		vars := mux.Vars(r)
		tab := vars["tab"]
		if tab == "" {
            tab = "liste"
		}
		//
		fermiers, err := model.GetSortedFermiers(ctx.DB, "nom")
		if err != nil {
			return err
		}
		allCommunes, err := model.GetSortedCommunes(ctx.DB, "nom")
		if err != nil {
			return err
		}
		//
		ctx.TemplateName = "sylvi-search-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Recherche sylviculture",
				CSSFiles: []string{
					"/static/css/form.css",
				},
			},
			Menu: "production",
			Details: detailsSylviForm{
				EssenceCodes: model.EssenceCodes,
				Fermiers:     fermiers,
				AllCommunes:  allCommunes,
////// supprimer si finalement pas de tab
				UrlAction:   "/sylviculture/recherche/" + tab,
				//UrlAction: "/sylviculture/recherche",
			},
		}
		return nil
	}
}
