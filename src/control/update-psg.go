package control

import (
	"net/http"
	//	"strconv"

	"bdl.local/bdl/ctxt"
	//	"bdl.local/bdl/generic/wilk/werr"
	//	"bdl.local/bdl/model"
	//	"github.com/gorilla/mux"
)

type detailsUpdatePSGForm struct {
	UrlAction string
}

func UpdatePSG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
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
		ctx.TemplateName = "update-psg-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Mise à jour données PSG",
				CSSFiles: []string{"/static/css/form.css"},
			},
			Menu: "acuueil",
			Details: detailsUpdatePSGForm{
				UrlAction: "/update-psg",
			},
		}
		return nil
	}
}
