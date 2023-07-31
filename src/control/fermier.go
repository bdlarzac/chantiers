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

type detailsFermierList struct {
	List  []*model.Fermier
	Count int
}

type detailsFermierShow struct {
	Fermier *model.Fermier
}

func ListFermier(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	list, err := model.GetSortedFermiers(ctx.DB, "nom")
	if err != nil {
		return werr.Wrap(err)
	}
	for _, fermier := range list {
		fermier.ComputeParcelles(ctx.DB)
		if err != nil {
			return werr.Wrap(err)
		}
	}
	ctx.TemplateName = "fermier-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Fermiers SCTL",
		},
		Menu: "acteurs",
		Footer: ctxt.Footer{
			JSFiles: []string{
				"/static/lib/table-sort/table-sort.js"},
		},
		Details: detailsFermierList{
			List:  list,
			Count: model.CountFermiers(ctx.DB),
		},
	}
	return nil
}

func ShowFermier(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id1, err := strconv.Atoi(vars["id"])
	id := int(id1)
	if err != nil {
		return werr.Wrap(err)
	}
	fermier, err := model.GetFermier(ctx.DB, id)
	if err != nil {
		return werr.Wrap(err)
	}
	err = fermier.ComputeParcelles(ctx.DB)
	if err != nil {
		return werr.Wrap(err)
	}
	// UGs et Lieudits des parcelles
	for _, p := range fermier.Parcelles {
		err = p.ComputeUGs(ctx.DB)
		if err != nil {
			return werr.Wrap(err)
		}
		err = p.ComputeLieudits(ctx.DB)
		if err != nil {
			return werr.Wrap(err)
		}
	}
	ctx.TemplateName = "fermier-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    fermier.String(),
			CSSFiles: []string{},
		},
		Menu: "acteurs",
		Details: detailsFermierShow{
			Fermier: fermier,
		},
	}
	return nil
}
