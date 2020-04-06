package control

import (
	"fmt"
	"net/http"

	"bdl.local/bdl/ctxt"
)

// *********************************************************
func ShowError(ctx *ctxt.Context, w http.ResponseWriter, err error) error {
	ctx.TemplateName = "error"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "ERREUR",
		},
		Menu:    "accueil",
		Details: fmt.Sprintf("%v", err),
	}
	return nil
}

// *********************************************************
func Show404(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	msg := r.URL.String()
	ctx.TemplateName = "404"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "PAGE INEXISTANTE",
		},
		Menu:    "accueil",
		Details: msg,
	}
	return nil
}
