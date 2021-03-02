package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"bdl.local/bdl/generic/wilk/webo"
	"time"
	"strconv"
	"html/template"
	"net/http"
	"net/url"
)

type detailsBilansForm struct {
	Periods   [][2]time.Time
	ClientsPlaquettesOptions        template.HTML
	UrlAction string
}

type detailsBilanClientPlaquettes struct {
	DateDebut time.Time
	DateFin   time.Time
    Client    *model.Acteur
    Ventes    []*model.VentePlaq
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
		if r.PostFormValue("type-bilan") == "client-plaquettes" {
            err = showBilanClientlaquettes(ctx, r.PostForm)
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
func showBilanClientlaquettes(ctx *ctxt.Context, formValues url.Values) error {
	dateDebut, err := time.Parse("2006-01-02", formValues.Get("date-debut"))
	if err != nil {
		return err
	}
	dateFin, err := time.Parse("2006-01-02", formValues.Get("date-fin"))
	if err != nil {
		return err
	}
    idClient, err := strconv.Atoi(formValues.Get("liste-client-plaquettes"))
	if err != nil {
		return err
	}
	client, err := model.GetActeur(ctx.DB, idClient)
	if err != nil {
		return err
	}
	ventes, err := model.GetVentePlaqsOfClient(ctx.DB, idClient, dateDebut, dateFin)
	//
	ctx.TemplateName = "bilan-client-plaquettes-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Bilan vente plaquettes",
			CSSFiles: []string{},
			JSFiles:  []string{
				"/static/js/round.js"},
		},
		Menu: "bilans",
		Footer: ctxt.Footer{
			JSFiles: []string{},
		},
		Details: detailsBilanClientPlaquettes{
			DateDebut: dateDebut,
			DateFin:   dateFin,
		    Client:    client,
		    Ventes:    ventes,
		},
	}
	if err != nil {
		return err
	}
	return nil
}
