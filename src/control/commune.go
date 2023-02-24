/*
*

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.

*
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	"net/http"
)

func ListCommunes(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	communes, err := model.ListCommunesEtLieudits(ctx.DB)
	if err != nil {
		return werr.Wrapf(err, "Erreur ListCommunes")
	}
	ctx.TemplateName = "commune-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Liste des communes",
		},
		Menu:    "accueil",
		Details: communes,
	}
	return nil
}
