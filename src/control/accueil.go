package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"net/http"
)

type detailsAccueil struct {
	Recents []*model.Recent
}

func Accueil(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	ctx.TemplateName = "accueil.html"
	recents, err := model.GetRecents(ctx.DB)
	if err != nil {
		return err
	}
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Accueil",
		},
		Menu: "accueil",
		Details: detailsAccueil{
			Recents: recents,
		},
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
