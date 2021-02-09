package control

import (
	"net/http"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
)

type detailsFermierList struct {
	List  []*model.Fermier
	Count int
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
