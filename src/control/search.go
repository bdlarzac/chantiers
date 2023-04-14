/*
@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"fmt"
	"net/http"
	"strings"
	"time"
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
		fmt.Printf("%+v\n", r.PostForm)
		//
		//        filtre_essence := computeFiltreEssence(r)
		//fmt.Printf("filtre_essence = %+v\n",filtre_essence)
		//        filtre_proprio := computeFiltreProprio(r)
		//fmt.Printf("filtre_proprio = %+v\n",filtre_proprio)
		//        filtre_periode := computeFiltrePeriode(r)
		//fmt.Printf("filtre_periode = %+v\n",filtre_periode)
		filtre_ug := computeFiltreUG(r)
		fmt.Printf("filtre_ug = %+v\n", filtre_ug)
		//
		ctx.TemplateName = "search-result.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Activités",
			},
			Menu:    "accueil",
			Details: detailsSearchResults{},
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
		fermiers, err := model.GetSortedFermiers(ctx.DB, "nom")
		if err != nil {
			return err
		}
		//
		allUGs, err := model.GetUGsSortedByCode(ctx.DB)
		if err != nil {
			return err
		}
		//
		allCommunes, err := model.GetSortedCommunes(ctx.DB, "nom")
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
				Periods:     periods,
				EssencesMap: essencesMap,
				PropriosMap: propriosMap,
				Fermiers:    fermiers,
				AllUGs:      allUGs,
				UGs:         []*model.UG{},
				AllCommunes: allCommunes,
				UrlAction:   "/search",
			},
		}
		return nil
	}
}

// ************************** Auxiliaires *******************************

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
	    - Sinon contient 2 strings, dates de début et de fin au format AAAA-MM-JJ.
*/
func computeFiltreUG(r *http.Request) (result []string) {
	if r.PostFormValue("ids-ugs") == "" {
		return []string{}
	}
	result = strings.Split(r.PostFormValue("ids-ugs"), ";")
	return result
}
