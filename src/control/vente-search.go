/*
Recherche / bilans de ventes génériques (plaquettes et chantiers autres valorisations)

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
	Periods     [][2]time.Time  // pour choix-date
	ValoCodes   []string        // pour choix-valo
	PropriosMap map[int]string  // pour choix-proprio
	Clients     []*model.Acteur // pour choix-client
	UrlAction   string
}

// ventes de tous les clients
type detailsVenteSearchResults struct {
	Ventes                []*model.Vente
	RecapFiltres          string
	DateDebut             time.Time
	DateFin               time.Time
	BilansVentesParSaison []*model.BilanVentesParSaison
	Tab                   string
}

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
		filtres["periode"] = computeFiltrePeriode(r)
		filtres["valo"] = computeFiltreValo(r)
		filtres["client"] = computeFiltreClient(r)
		filtres["proprio"] = computeFiltreProprio(r)
		ventes, err := model.ComputeVentesFromFiltres(ctx.DB, filtres)
		if err != nil {
			return err
		}
		//
		recapFiltres, err := model.ComputeRecapFiltres(ctx.DB, filtres)
		if err != nil {
			return err
		}
		//
		bilansVentesParSaison, err := model.ComputeBilansVentesParSaison(ctx.DB, ctx.Config.DebutSaison, ventes)
		if err != nil {
			return err
		}
		//
		ctx.TemplateName = "vente-search-show.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Ventes",
				CSSFiles: []string{
					"/static/lib/tabstrip/tabstrip.css",
				},
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
			Details: detailsVenteSearchResults{
				RecapFiltres:          recapFiltres,
				Ventes:                ventes,
				BilansVentesParSaison: bilansVentesParSaison,
				Tab:                   tab,
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
		clients, err := model.GetClients(ctx.DB)
		if err != nil {
			return err
		}
		//
		propriosMap, err := model.GetProprietaires(ctx.DB)
		if err != nil {
			return err
		}
		//
		ctx.TemplateName = "vente-search-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Recherche de ventes",
				CSSFiles: []string{
					"/static/css/form.css",
				},
			},
			Menu: "ventes",
			Details: detailsVenteSearchForm{
				Periods:     periods,
				PropriosMap: propriosMap,
				Clients:     clients,
				ValoCodes:   model.AllValoCodesAvecChaufer(),
				UrlAction:   "/vente/recherche",
			},
		}
		return nil
	}
}

func SearchVenteParClient(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	return nil
}
