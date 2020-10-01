package control

import (
	"net/http"
	//	"strconv"

	"bdl.local/bdl/ctxt"
	//	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	//	"github.com/gorilla/mux"
)

type detailsUpdateSCTLForm struct {
	Items     []*model.UpdatedItem
	UrlAction string
}

func UpdateSCTL(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		//		ctx.Redirect = "/acteur/" + strconv.Itoa(id)
		return nil
	default:
		//
		// Affiche form
		//
		items, err := model.ComputeUpdateSCTL(ctx.DB, ctx.Config)
		if err != nil {
			return err
		}
		ctx.TemplateName = "update-sctl-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Mise à jour données foncières (SCTL)",
				CSSFiles: []string{"/static/css/form.css"},
			},
			Menu: "acteurs",
			Details: detailsUpdateSCTLForm{
				UrlAction: "/update-sctl",
				Items:     items,
			},
		}
		return nil
	}
}
