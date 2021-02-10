package control

import (
	"net/http"
	"strconv"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
)

type detailsFermierList struct {
	List  []*model.Fermier
	Count int
}

type detailsFermierShow struct {
	Fermier    *model.Fermier
}

// *********************************************************
func ListFermier(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	list, err := model.GetSortedFermiers(ctx.DB, "nom")
	if err != nil {
		return err
	}
	for _, fermier := range list {
		fermier.ComputeParcelles(ctx.DB)
		if err != nil {
			return err
		}
	}
	ctx.TemplateName = "fermier-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Fermiers SCTL",
		},
		Menu: "acteurs",
		Footer: ctxt.Footer{
			JSFiles: []string{"/static/js/toogleTR.js"},
		},
		Details: detailsFermierList{
			List:  list,
			Count: model.CountFermiers(ctx.DB),
		},
	}
	return nil
}

// *********************************************************
func ShowFermier(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id1, err := strconv.Atoi(vars["id"])
	id := int(id1)
	if err != nil {
		return err
	}
	fermier, err := model.GetFermier(ctx.DB, id)
	if err != nil {
		return err
	}
    err = fermier.ComputeParcelles(ctx.DB)
	if err != nil {
		return err
	}
    // UGs et Lieudits des parcelles
	for _, p := range fermier.Parcelles {
	    err = p.ComputeUGs(ctx.DB)
        if err != nil {
            return err
        }
	    err = p.ComputeLieudits(ctx.DB)
        if err != nil {
            return err
        }
	}

	ctx.TemplateName = "fermier-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: fermier.String(),
			CSSFiles: []string{
				"/static/css/form.css"},
		},
		Menu: "acteurs",
		Details: detailsFermierShow{
			Fermier:    fermier,
		},
	}
	return nil
}

