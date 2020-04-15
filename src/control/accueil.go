/*
Contient des fonctions en vrac - à dispatcher dans d'autres contrôleurs

*/

package control

import (
	"net/http"

	"bdl.local/bdl/ctxt"
)

func Accueil(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	ctx.TemplateName = "accueil.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Accueil",
		},
		Menu: "accueil",
	}
	return nil
}

func MajFoncier(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	ctx.TemplateName = "maj-foncier.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Mise à jour données foncières",
		},
		Menu: "accueil",
	}
	return nil
}

func MajPSG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	ctx.TemplateName = "maj-psg.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Mise à jour PSG",
		},
		Menu: "accueil",
	}
	return nil
}

func ShowDoc(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	ctx.TemplateName = "doc.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Documentation",
		},
		Menu: "accueil",
	}
	return nil
}
