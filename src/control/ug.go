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
}

func ShowUG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return werr.Wrapf(err, "Erreur conversion id UG")
	}
	ug, err := model.GetUGFull(ctx.DB, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetUGFull()")
	}
	activites, err := ug.GetActivitesByDate(ctx.DB)
	if err != nil {
		return err
	}
	//
	ctx.TemplateName = "ug-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Unité de gestion " + vars["id"],
			CSSFiles: []string{"/static/css/tabstrip.css"},
		},
		Menu: "unites-gestion",
		Details: detailsUGShow{
			UG:        ug,
			Activites: activites,
		},
		Footer: ctxt.Footer{
			JSFiles: []string{"/static/js/tabstrip.js"},
		},
	}
	return nil
}
