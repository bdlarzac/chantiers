package control

import (
	"io/ioutil"
	"net/http"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
)

type detailsConfigShow struct {
	Data string
}

func ShowConfig(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	tmp, err := ioutil.ReadFile("../config.yml")
	if err != nil {
		return werr.Wrapf(err, "Erreur lecture config.yml")
	}
	data := string(tmp)
	ctx.TemplateName = "config-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Configuration",
			CSSFiles: []string{"/static/css/form.css"},
		},
		Menu: "accueil",
		Details: detailsConfigShow{
			Data: data,
		},
	}
	return nil
}
