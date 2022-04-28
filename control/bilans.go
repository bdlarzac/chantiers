/**
    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
**/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/webo"
	"bdl.local/bdl/model"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type detailsBilansForm struct {
	Periods                  [][2]time.Time
	HasChantier              bool
	ClientsPlaquettesOptions template.HTML
	UrlAction                string
	Proprietaires            map[int]string
}

// ventes d'un seul client
type detailsBilanClientPlaquettes struct {
	DateDebut time.Time
	DateFin   time.Time
	Client    *model.Acteur
	Ventes    []*model.VentePlaq
}

// ventes de tous les clients
type detailsBilanVentesPlaquettes struct {
	DateDebut time.Time
	DateFin   time.Time
	Ventes    []*model.VentePlaq
}

type detailsBilanValoEssences struct {
	DateDebut     time.Time
	DateFin       time.Time
	Valorisations model.Valorisations
	EssenceCodes  []string
	ValoCodes     []string
}

func FormBilans(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		var err error
		if err = r.ParseForm(); err != nil {
			return err
		}
		if r.PostFormValue("type-bilan") == "client-plaquettes" {
			err = showBilanClientPlaquettes(ctx, r.PostForm)
			if err != nil {
				return err
			}
		} else if r.PostFormValue("type-bilan") == "ventes-plaquettes" {
			err = showBilanVentesPlaquettes(ctx, r.PostForm)
			if err != nil {
				return err
			}
		} else if r.PostFormValue("type-bilan") == "valorisations" {
			err = showBilanValoEssences(ctx, r.PostForm, "valorisations")
			if err != nil {
				return err
			}
		} else if r.PostFormValue("type-bilan") == "essences" {
			err = showBilanValoEssences(ctx, r.PostForm, "essences")
			if err != nil {
				return err
			}
		}
		return nil
	default:
		//
		// Affiche form
		//
		periods, hasChantier, err := model.ComputeLimitesSaisons(ctx.DB, ctx.Config.DebutSaison)
		if err != nil {
			return err
		}
		ctx.TemplateName = "bilans-form.html"
		
		proprietaires, err := model.GetProprietaires(ctx.DB)
		
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Choix bilan",
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Details: detailsBilansForm{
				Periods:                  periods,
				HasChantier:              hasChantier,
				UrlAction:                "/bilans",
				ClientsPlaquettesOptions: webo.FmtOptions(WeboClientsPlaquettes(ctx), "CHOOSE_CLIENT"),
				Proprietaires:            proprietaires,
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
// Bilan ventes pour un client donné
func showBilanClientPlaquettes(ctx *ctxt.Context, formValues url.Values) error {
	dateDebut, err := time.Parse("2006-01-02", formValues.Get("date-debut"))
	if err != nil {
		return err
	}
	dateFin, err := time.Parse("2006-01-02", formValues.Get("date-fin"))
	if err != nil {
		return err
	}
	idClient, err := strconv.Atoi(strings.TrimLeft(formValues.Get("liste-client-plaquettes"), "acteur-"))
	if err != nil {
		return err
	}
	client, err := model.GetActeur(ctx.DB, idClient)
	if err != nil {
		return err
	}
	ventes, err := model.GetVentePlaqsOfClient(ctx.DB, idClient, dateDebut, dateFin)
	if err != nil {
		return err
	}
	//
	ctx.TemplateName = "bilan-client-plaquettes-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Bilan client plaquettes",
			CSSFiles: []string{},
			JSFiles: []string{
				"/static/js/round.js",
				"/static/js/formatNb.js",
				"/view/common/prix.js"},
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

// *********************************************************
// Bilan ventes pour tous clients confondus
func showBilanVentesPlaquettes(ctx *ctxt.Context, formValues url.Values) error {
	dateDebut, err := time.Parse("2006-01-02", formValues.Get("date-debut"))
	if err != nil {
		return err
	}
	dateFin, err := time.Parse("2006-01-02", formValues.Get("date-fin"))
	if err != nil {
		return err
	}
	ventes, err := model.GetVentePlaqsOfPeriod(ctx.DB, dateDebut, dateFin)
	if err != nil {
		return err
	}
	//
	ctx.TemplateName = "bilan-ventes-plaquettes-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Bilan ventes plaquettes",
			CSSFiles: []string{},
			JSFiles: []string{
				"/static/js/formatNb.js",
				"/static/js/round.js"},
		},
		Menu: "bilans",
		Footer: ctxt.Footer{
			JSFiles: []string{},
		},
		Details: detailsBilanVentesPlaquettes{
			DateDebut: dateDebut,
			DateFin:   dateFin,
			Ventes:    ventes,
		},
	}
	if err != nil {
		return err
	}
	return nil
}

// *********************************************************
// Bilan par valorisation et par essences
// (les deux sont presque identiques)
// @param  what "valorisations" ou "essences"
func showBilanValoEssences(ctx *ctxt.Context, formValues url.Values, what string) error {
	dateDebut, err := time.Parse("2006-01-02", formValues.Get("date-debut"))
	if err != nil {
		return err
	}
	dateFin, err := time.Parse("2006-01-02", formValues.Get("date-fin"))
	if err != nil {
		return err
	}
	valos, err := model.ComputeBilanValoEssences(ctx.DB, dateDebut, dateFin)
	if err != nil {
		return err
	}
	var titre, templateName string
	if what == "valorisations" {
		titre = "Bilan valorisations"
		templateName = "bilan-valorisations-show.html"
	} else {
		titre = "Bilan essences"
		templateName = "bilan-essences-show.html"
	}
	//
	ctx.TemplateName = templateName
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    titre,
			CSSFiles: []string{},
			JSFiles: []string{
				"/static/js/formatNb.js"},
		},
		Menu: "bilans",
		Footer: ctxt.Footer{
			JSFiles: []string{},
		},
		Details: detailsBilanValoEssences{
			DateDebut:     dateDebut,
			DateFin:       dateFin,
			Valorisations: valos,
			EssenceCodes:  model.AllEssenceCodes(),
			ValoCodes:     model.AllValorisationCodesAvecChaufer(),
		},
	}
	if err != nil {
		return err
	}
	return nil
}
