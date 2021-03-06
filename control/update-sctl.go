/**
    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
**/
package control

import (
	"net/http"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
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
		if true {
			// TODO supprimer
			ctx.TemplateName = "update-sctl-form--tmp.html"
			ctx.Page = &ctxt.Page{
				Header: ctxt.Header{
					Title:    "Mise à jour données foncières (SCTL)",
					CSSFiles: []string{"/static/css/form.css"},
				},
				Menu: "accueil",
				Details: detailsUpdateSCTLForm{
					UrlAction: "",
					Items:     []*model.UpdatedItem{},
				},
			}
		} else {
			// TODO finir et garder
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
				Menu: "accueil",
				Details: detailsUpdateSCTLForm{
					UrlAction: "/update-sctl",
					Items:     items,
				},
			}
		}
		return nil
	}
}
