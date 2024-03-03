/*
Contrôle l'affichage de la page cartographie.

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"net/http"
)

func ShowCarte(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	ctx.TemplateName = "carte.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Cartographie",
			JSFiles: []string{
				"/static/lib/leaflet/leaflet.js",
				"/static/data-carto/ug.geojson.js",
//				"/static/data-carto/parcelle.geojson.js",
			},
			CSSFiles: []string{
				"/static/lib/leaflet/leaflet.css",
				"/static/css/form.css",
			},
		},
		Menu: "accueil",
	}
	return nil
}

