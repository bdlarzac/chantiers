package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"bdl.local/bdl/generic/wilk/webo"
	"html/template"
	"net/http"
	"time"
)

type detailsBilansForm struct {
	Periods   [][2]time.Time
	ClientsPlaquettesOptions        template.HTML
	UrlAction string
}

type detailsBilansPlaquettes struct {
	DateDebut time.Time
	DateFin   time.Time
}

func FormBilans(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		var err error
		if err = r.ParseForm(); err != nil {
			return err
		}
		//
		// Process form
		//
		if r.PostFormValue("type-bilan") == "clients-plaquettes" {
			err = ShowBilanClientslaquettes(ctx, r.PostFormValue("date-debut"), r.PostFormValue("date-fin"))
            if err != nil {
                return err
            }
		}
		return nil
	default:
		//
		// Affiche form
		//
		periods, err := model.ComputeLimitesSaisons(ctx.DB, ctx.Config.DebutSaison)
		if err != nil {
			return err
		}
		ctx.TemplateName = "bilans-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Choix bilan",
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Details: detailsBilansForm{
				Periods:   periods,
				UrlAction: "/bilans",
				ClientsPlaquettesOptions: webo.FmtOptions(WeboClientsPlaquettes(ctx), "CHOOSE_CLIENT"),
			},
			Menu: "bilans",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
				},
			},
		}
		return nil
	}
}

// *********************************************************
func ShowBilanClientslaquettes(ctx *ctxt.Context, debut, fin string) error {
	dateDebut, err := time.Parse("2006-01-02", debut)
	if err != nil {
		return err
	}
	dateFin, err := time.Parse("2006-01-02", fin)
	if err != nil {
		return err
	}
	ctx.TemplateName = "bilan-clients-plaquettes-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Bilan clients plaquettes",
			CSSFiles: []string{},
			JSFiles:  []string{},
		},
		Menu: "bilans",
		Footer: ctxt.Footer{
			JSFiles: []string{},
		},
		Details: detailsBilansPlaquettes{
			DateDebut: dateDebut,
			DateFin:   dateFin,
		},
	}
	if err != nil {
		return err
	}
	return nil
}
