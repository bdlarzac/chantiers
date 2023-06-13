/*
Recherche / bilans d'activité

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

type detailsActiviteSearchForm struct {
	Periods      [][2]time.Time   // pour choix-date
	EssenceCodes []string         // pour choix-essence
	ValoCodes    []string         // pour choix-valo
	PropriosMap  map[int]string   // pour choix-proprio
	Fermiers     []*model.Fermier // pour choix-fermier
	AllUGs       []*model.UG      // pour choix-ug - liens-ugs-modal
	UGs          []*model.UG      // pour choix-ug - liens-ugs - toujours vide, utile que pour compatibilité avec liens-ugs.html
	AllCommunes  []*model.Commune // pour choix-parcelle
	UrlAction    string
}

type detailsActiviteSearchResults struct {
	Activites                []*model.Activite
	RecapFiltres             string
	ActiviteMap              map[string]string
	BilansActivitesParSaison []*model.BilanActivitesParSaison
	ActivitesParUG           []*model.ActivitesParUG
	Tab                      string
}

// Affiche / process le formulaire de recherche
func SearchActivite(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	switch r.Method {
	case "POST":
		//
		// Process form et affiche page de résultats
		//
		if err = r.ParseForm(); err != nil {
			return err
		}
		//
		filtres := map[string][]string{}
		filtres["fermier"] = computeFiltreFermier(r)
		filtres["essence"] = computeFiltreEssence(r)
		filtres["valo"] = computeFiltreValo(r)
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
		ctx.TemplateName = "search-activite-show.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Activités",
				CSSFiles: []string{"/static/lib/tabstrip/tabstrip.css"},
				JSFiles:  []string{"/static/js/formatNb.js"},
			},
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/lib/tabstrip/tabstrip.js",
					"/static/lib/table-sort/table-sort.js",
				},
			},
			Menu: "accueil",
			Details: detailsActiviteSearchResults{
				Activites:                activites,
				RecapFiltres:             recapFiltres,
				ActiviteMap:              model.GetActivitesMap(),
				BilansActivitesParSaison: bilansActivitesParSaison,
				ActivitesParUG:           model.ComputeActivitesParUG(activites),
				Tab:                      r.PostFormValue("type-resultat"),
			},
		}
		return nil
	default:
		//
		// Affiche form
		//
		periods, _, err := model.ComputeLimitesSaisons(ctx.DB, ctx.Config.DebutSaison)
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
		ctx.TemplateName = "search-activite-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Recherche d'activité",
				CSSFiles: []string{
					"/static/css/form.css",
				},
			},
			Menu: "accueil",
			Details: detailsActiviteSearchForm{
				Periods:      periods,
				EssenceCodes: model.EssenceCodes,
				ValoCodes:    model.AllValoCodesAvecChauferEtPlaq(),
				PropriosMap:  propriosMap,
				Fermiers:     fermiers,
				AllUGs:       allUGs,
				UGs:          []*model.UG{},
				AllCommunes:  allCommunes,
				UrlAction:    "/activite/recherche",
			},
		}
		return nil
	}
}
