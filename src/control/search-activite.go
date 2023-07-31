/*
Recherche / bilans d'activité

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/model"
	"golang.org/x/exp/slices"
	"net/http"
	"strconv"
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
	HasPlaquettes            bool // est-ce que les plaquettes sont demandées ? (facilite l'affichage de la template)
	ActivitesParUG           []*model.ActivitesParUG
	LabelProprios            map[int]string
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
			return werr.Wrap(err)
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
		//
		activites, err := model.ComputeActivitesFromFiltres(ctx.DB, filtres)
		if err != nil {
			return werr.Wrap(err)
		}
		//
		recapFiltres, err := model.ComputeRecapFiltres(ctx.DB, filtres) // pour l'affichage
		if err != nil {
			return werr.Wrap(err)
		}
		//
		hasPlaquettes := false
		if len(filtres["valo"]) == 0 {
		    hasPlaquettes = true
		} else {
            if tiglib.InArray("PQ", filtres["valo"]) {
                hasPlaquettes = true
            }
		}
		//
		bilansActivitesParSaison, err := model.ComputeBilansActivitesParSaison(ctx.DB, ctx.Config.DebutSaison, activites)
		if err != nil {
			return werr.Wrap(err)
		}
		//
		labelProprios, err := model.LabelActeurs(ctx.DB, "DIV-PF") // "DIV-PF" = "divers - propriétaire foncier"
		if err != nil {
			return werr.Wrap(err)
		}
		// on ne garde dans labelProprios que les propriétaires choisis,
		// utilisé dans la template pour n'afficher que ces propriétaires.
		if len(filtres["proprio"]) != 0 {
			for idProprio, _ := range labelProprios {
				if !slices.Contains(filtres["proprio"], strconv.Itoa(idProprio)) {
					delete(labelProprios, idProprio)
				}
			}
		}
		//
		ctx.TemplateName = "search-activite-show.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Activités",
				CSSFiles: []string{"/static/lib/tabstrip/tabstrip.css"},
				JSFiles: []string{
					"/static/js/formatNb.js",
					"/static/js/round.js",
				},
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
				HasPlaquettes:            hasPlaquettes,
				ActivitesParUG:           model.ComputeActivitesParUG(activites),
				LabelProprios:            labelProprios,
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
			return werr.Wrap(err)
		}
		propriosMap, err := model.GetProprietaires(ctx.DB)
		if err != nil {
			return werr.Wrap(err)
		}
		fermiers, err := model.GetSortedFermiers(ctx.DB, "nom")
		if err != nil {
			return werr.Wrap(err)
		}
		allUGs, err := model.GetUGsSortedByCode(ctx.DB)
		if err != nil {
			return werr.Wrap(err)
		}
		allCommunes, err := model.GetSortedCommunes(ctx.DB, "nom")
		if err != nil {
			return werr.Wrap(err)
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
