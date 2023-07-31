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
	"time"
)

type detailsTasVidesShow struct {
	Chantiers []*model.Plaq
}

func ShowTasVides(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	chantiers, err := model.GetAllPlaqsVides(ctx.DB)
	if err != nil {
		return werr.Wrap(err)
	}
	ctx.TemplateName = "stockage-tas-vides.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:   "Tas vides",
			JSFiles: []string{"/static/js/round.js"},
		},
		Menu: "accueil",
		Footer: ctxt.Footer{
			JSFiles: []string{},
		},
		Details: detailsTasVidesShow{
			Chantiers: chantiers,
		},
	}
	return nil
}

func SignalerTasVide(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return werr.Wrap(err)
	}
	date, err := time.Parse("2006-01-02", vars["date"])
	if err != nil {
		return werr.Wrap(err)
	}
	err = model.DesactiverTas(ctx.DB, id, date)
	if err != nil {
		return werr.Wrap(err)
	}
	ctx.Redirect = "/stockage/liste"
	return nil
}
