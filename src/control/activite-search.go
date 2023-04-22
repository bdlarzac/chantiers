/*
Recheche d'activité

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
	
	//"fmt"
)

type detailsSearchForm struct {
	Periods     [][2]time.Time    // pour choix-date
	EssencesMap map[string]string // pour choix-essence
	PropriosMap map[int]string    // pour choix-proprio
	Fermiers    []*model.Fermier  // pour choix-fermier
	AllUGs      []*model.UG       // pour choix-ug - liens-ugs-modal
	UGs         []*model.UG       // pour choix-ug - liens-ugs - toujours vide, utile que pour compatibilité avec liens-ugs.html
	AllCommunes []*model.Commune  // pour choix-parcelle
	UrlAction   string
}

type detailsSearchResults struct {
	Activites    []*model.Activite
	RecapFiltres string
	ActiviteMap  map[string]string
	Tab          string
}

/*
Affiche / process le formulaire de recherche
*/
func SearchActivite(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
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
//fmt.Printf("%+v\n", r.PostForm)
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
//fmt.Printf("filtres = %+v\n",filtres)
		//
		recapFiltres, err := model.ComputeRecapFiltres(ctx.DB, filtres)
		if err != nil {
			return err
		}
		ctx.TemplateName = "activite-show.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Activités",
				CSSFiles: []string{"/static/lib/tabstrip/tabstrip.css"},
			},
			Footer: ctxt.Footer{
				JSFiles: []string{
				    "/static/lib/tabstrip/tabstrip.js",
					"/static/lib/table-sort/table-sort.js"},
			},
			Menu: "accueil",
			Details: detailsSearchResults{
				Activites:       activites,
				RecapFiltres:    recapFiltres,
				ActiviteMap:     model.GetActivitesMap(),
				Tab:             tab,
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
		ctx.TemplateName = "activite-search.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Recherche d'activité",
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
				UrlAction:   "/activite",
			},
		}
		return nil
	}
}

// ************************** Auxiliaires *******************************

/*
Filtre fermier : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient une liste avec un seul élément, l'id du fermier sélectionné.
*/
func computeFiltreFermier(r *http.Request) (result []string) {
	choix := r.PostFormValue("select-choix-fermier")
	if choix == "choix-fermier-no-limit" {
		return []string{}
	}
	return []string{choix[14:]}
}

/*
Filtre essence : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient une liste de codes essence.
*/
func computeFiltreEssence(r *http.Request) (result []string) {
	if r.PostFormValue("choix-ALL-essence") == "true" {
		return []string{}
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
Filtre propriétaire : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient une liste d'id propriétaires (dans la table acteur)
    (attention, cette liste contient des strings, pas des ints).
*/
func computeFiltreProprio(r *http.Request) (result []string) {
	if r.PostFormValue("choix-ALL-proprio") == "true" {
		return []string{}
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

/*
Filtre période : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient 2 strings, dates de début et de fin au format AAAA-MM-JJ.
*/
func computeFiltrePeriode(r *http.Request) (result []string) {
	if r.PostFormValue("choix-periode-periodes") == "choix-periode-no-limit" {
		return []string{}
	}
	result = append(result, r.PostFormValue("choix-periode-debut"))
	result = append(result, r.PostFormValue("choix-periode-fin"))
	return result
}

/*
Filtre UG : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient les ids UG
*/
func computeFiltreUG(r *http.Request) (result []string) {
	if r.PostFormValue("ids-ugs") == "" {
		return []string{}
	}
	result = strings.Split(r.PostFormValue("ids-ugs"), ";")
	return result
}

/*
Filtre Parcelles : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient les ids parcelle.
*/
func computeFiltreParcelle(r *http.Request) (result []string) {
	if r.PostFormValue("ids-parcelles") == "" {
		return []string{}
	}
	result = strings.Split(r.PostFormValue("ids-parcelles"), ";")
	return result
}
