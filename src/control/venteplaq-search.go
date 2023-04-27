/*
Recherche / bilans de ventes plaquettes

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type detailsVenteSearchForm struct {
	Periods     [][2]time.Time    // pour choix-date
	EssencesMap map[string]string // pour choix-essence
	PropriosMap map[int]string    // pour choix-proprio
	Fermiers    []*model.Fermier  // pour choix-fermier
	AllUGs      []*model.UG       // pour choix-ug - liens-ugs-modal
	UGs         []*model.UG       // pour choix-ug - liens-ugs - toujours vide, utile que pour compatibilité avec liens-ugs.html
	AllCommunes []*model.Commune  // pour choix-parcelle
	UrlAction   string
}

type detailsVenteSearchResults struct {
	Activites                []*model.Activite
	RecapFiltres             string
	ActiviteMap              map[string]string
	BilansActivitesParSaison []*model.BilanActivitesParSaison
	Tab                      string
}

/*
Affiche / process le formulaire de recherche
*/
func SearchVente(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	switch r.Method {
	case "POST":
		//
		// Process form et affiche page de résultats
		//
		if err = r.ParseForm(); err != nil {
			return err
		}
		vars := mux.Vars(r)
		tab := vars["tab"]
		if tab == "" {
			tab = "liste"
		}
		//
		filtres := map[string][]string{}
		filtres["fermier"] = computeFiltreFermier(r)
		filtres["essence"] = computeFiltreEssence(r)
		filtres["proprio"] = computeFiltreProprio(r)
		filtres["periode"] = computeFiltrePeriode(r)
		filtres["ug"] = computeFiltreUG(r)
		filtres["parcelle"] = computeFiltreParcelle(r)
		activites, err := model.ComputeActivitesFromFiltres(ctx.DB, filtres)
		if err != nil {
			return err
		}
		//
		recapFiltres, err := model.ComputeRecapFiltres(ctx.DB, filtres)
		if err != nil {
			return err
		}
		//
		bilansActivitesParSaison, err := model.ComputeBilansActivitesParSaison(ctx.DB, ctx.Config.DebutSaison, activites)
		if err != nil {
			return err
		}
		//
		ctx.TemplateName = "venteplaq-search-show.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Ventes plaquettes",
				CSSFiles: []string{"/static/lib/tabstrip/tabstrip.css"},
				JSFiles: []string{"/static/js/formatNb.js"},
			},
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/lib/tabstrip/tabstrip.js",
					"/static/lib/table-sort/table-sort.js"},
			},
			Menu: "accueil",
			Details: detailsSearchResults{
				Activites:                activites,
				RecapFiltres:             recapFiltres,
				ActiviteMap:              model.GetActivitesMap(),
				BilansActivitesParSaison: bilansActivitesParSaison,
				Tab:                      tab,
			},
		}
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		tab := vars["tab"]
		if tab == "" {
			tab = "liste"
		}
		//
		periods, _, err := model.ComputeLimitesSaisons(ctx.DB, ctx.Config.DebutSaison)
		if err != nil {
			return err
		}
		essencesMap, err := model.GetEssencesMap(ctx.DB)
		if err != nil {
			return err
		}
		propriosMap, err := model.GetProprietaires(ctx.DB)
		if err != nil {
			return err
		}
		fermiers, err := model.GetSortedFermiers(ctx.DB, "nom")
		if err != nil {
			return err
		}
		allUGs, err := model.GetUGsSortedByCode(ctx.DB)
		if err != nil {
			return err
		}
		allCommunes, err := model.GetSortedCommunes(ctx.DB, "nom")
		if err != nil {
			return err
		}
		//
		ctx.TemplateName = "venteplaq-search-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Recherche de ventes plaquettes",
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Menu: "accueil",
			Details: detailsSearchForm{
				Periods:     periods,
				EssencesMap: essencesMap,
				PropriosMap: propriosMap,
				Fermiers:    fermiers,
				AllUGs:      allUGs,
				UGs:         []*model.UG{},
				AllCommunes: allCommunes,
				UrlAction:   "/vente/recherche/" + tab,
			},
		}
		return nil
	}
}
