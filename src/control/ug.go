package control

import (
	"net/http"
	"strconv"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
)

type detailsUGShow struct {
	UG        *model.UG
	Activites []*model.UGActivite
	Tab       string
}

func ShowUG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	tab := vars["tab"]
	if tab == "" {
		tab = "general"
	}
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return werr.Wrapf(err, "Erreur conversion id UG")
	}
	ug, err := model.GetUGFull(ctx.DB, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetUGFull()")
	}
	err = ug.ComputeRecap(ctx.DB)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel UG.ComputeRecap()")
	}
	activites, err := ug.GetActivitesByDate(ctx.DB)
	if err != nil {
		return err
	}
	//
	ctx.TemplateName = "ug-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Unité de gestion " + ug.Code,
			CSSFiles: []string{"/static/css/tabstrip.css"},
		},
		Menu: "accueil",
		Details: detailsUGShow{
			UG:        ug,
			Activites: activites,
			Tab:       tab,
		},
		Footer: ctxt.Footer{
			JSFiles: []string{"/static/js/tabstrip.js"},
		},
	}
	return nil
}

func SearchUG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	ctx.TemplateName = "ug-search.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Recherche UG",
			CSSFiles: []string{"/static/css/form.css"},
		},
		Menu: "accueil",
	}
	return nil
}
