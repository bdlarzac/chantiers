/*
@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type detailsUGShow struct {
	UG  *model.UG
	Tab string
}

type detailsUGSearch struct {
	AllUgs [][]*model.UG
}

func ShowUG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	tab := vars["tab"]
	if tab == "" {
		tab = "general"
	}
	//
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return werr.Wrapf(err, "Erreur conversion id UG")
	}
	//
	ug, err := model.GetUGFull(ctx.DB, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetUGFull()")
	}
	//
	err = ug.ComputeRecap(ctx.DB)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel UG.ComputeRecap()")
	}
	//
	err = ug.ComputeActivites(ctx.DB)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel UG.ComputeActivites()")
	}
	//
	ctx.TemplateName = "ug-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Unité de gestion " + ug.Code,
			CSSFiles: []string{"/static/lib/tabstrip/tabstrip.css"},
			JSFiles:  []string{"/static/js/formatNb.js"},
		},
		Menu: "accueil",
		Details: detailsUGShow{
			UG:  ug,
			Tab: tab,
		},
		Footer: ctxt.Footer{
			JSFiles: []string{
				"/static/lib/tabstrip/tabstrip.js",
				"/static/lib/table-sort/table-sort.js",
			},
		},
	}
	return nil
}

// Affiche à la fois le formulaire et la liste des UGs
func ListUGs(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	allUgs, err := model.GetUGsSortedByCodeAndSeparated(ctx.DB)
	if err != nil {
		return werr.Wrap(err)
	}
	ctx.TemplateName = "ug-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Recherche UG",
			CSSFiles: []string{"/static/css/form.css"},
		},
		Details: detailsUGSearch{
			AllUgs: allUgs,
		},
		Menu: "accueil",
	}
	return nil
}
