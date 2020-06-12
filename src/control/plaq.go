package control

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/webo"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
//"fmt"
)

type detailsPlaqForm struct {
	UrlAction           string
	EssenceOptions      template.HTML
	ExploitationOptions template.HTML
	GranuloOptions      template.HTML
	AllStockages        []*model.Stockage
	Chantier            *model.Plaq
}

type detailsPlaqList struct {
	Chantiers []*model.Plaq
	Annee     string   // année courante
	Annees    []string // toutes les années avec chantier
}

type detailsPlaqShow struct {
	Chantier         *model.Plaq
	PourcentagePerte float64
	Tab              string
}

// *********************************************************
func ListPlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	annee := vars["annee"]
	if annee == "" {
		// annee non spécifiée, on prend l'année courante
		annee = strconv.Itoa(time.Now().Year())
	}
	chantiers, err := model.GetPlaqsOfYear(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	annees, err := model.GetPlaqDifferentYears(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	ctx.TemplateName = "plaq-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Chantiers plaquettes " + annee,
		},
		Menu: "chantiers",
		Footer: ctxt.Footer{
			JSFiles: []string{"/static/js/plaq.js"},
		},
		Details: detailsPlaqList{
			Chantiers: chantiers,
			Annee:     annee,
			Annees:    annees,
		},
	}
	return nil
}

// *********************************************************
func ShowPlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	tab := vars["tab"]
	if tab == "" {
		tab = "general"
	}
	idChantier, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	chantier, err := model.GetPlaqFull(ctx.DB, idChantier)
	if err != nil {
		return err
	}
	err = chantier.ComputeCout(ctx.DB, ctx.Config)
	if err != nil {
		return err
	}
	ctx.TemplateName = "plaq-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Chantier " + chantier.String(),
			CSSFiles: []string{"/static/css/tabstrip.css"},
			JSFiles: []string{
				"/static/js/round.js",
				"/static/js/prix.js"},
		},
		Menu: "chantiers",
		Footer: ctxt.Footer{
			JSFiles: []string{
				"/static/js/tabstrip.js",
				"/static/js/plaq.js"},
		},
		Details: detailsPlaqShow{
			Chantier:         chantier,
			PourcentagePerte: ctx.Config.PourcentagePerte,
			Tab:              tab,
		},
	}
	return nil
}

// *********************************************************
// Process ou affiche form new
func NewPlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		chantier, err := chantierPlaquetteForm2var(r)
		if err != nil {
			return err
		}
		// calcul des ids stockage, pour transmettre à InsertPlaq(), qui va créer le(s) tas
        allStockages, err := model.GetStockagesActifs(ctx.DB)
		if err != nil {
			return err
		}
		idsStockages := []int{}
        for _, stockage := range(allStockages){
            if r.PostFormValue("stockage-" + strconv.Itoa(stockage.Id)) == "on" {
                idsStockages = append(idsStockages, stockage.Id)
            }
        }
		//
		id, err := model.InsertPlaq(ctx.DB, chantier, idsStockages)
		if err != nil {
			return err
		}
		//
		redirect := "/chantier/plaquette/" + strconv.Itoa(id)
		err = chantier.ComputeLieudit(ctx.DB) // nécessaire pour appeler chantier.FullString()
		if err != nil {
			return err
		}
		err = model.AddRecent(ctx.DB, ctx.Config, &model.Recent{URL:redirect, Label:chantier.FullString()})
		if err != nil {
			return err
		}
		ctx.Redirect = redirect
		return nil
	default:
		//
		// Affiche form
		//
		chantier := &model.Plaq{}
		chantier.Lieudit = &model.Lieudit{}
		chantier.Fermier = &model.Acteur{}
		chantier.UG = &model.UG{}
        allStockages, err := model.GetStockagesActifs(ctx.DB)
		if err != nil {
			return err
		}
		ctx.TemplateName = "plaq-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouveau chantier plaquettes",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/autocomplete/autocomplete.js"},
			},
			Details: detailsPlaqForm{
				Chantier:            chantier,
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "CHOOSE_ESSENCE"),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "CHOOSE_EXPLOITATION"),
				GranuloOptions:      webo.FmtOptions(WeboGranulo(), "CHOOSE_GRANULO"),
				AllStockages:     	 allStockages,
				UrlAction:           "/chantier/plaquette/new",
			},
		}
		return nil
	}
}

// *********************************************************
// Process ou affiche form update
func UpdatePlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		chantier, err := chantierPlaquetteForm2var(r)
		if err != nil {
			return err
		}
		chantier.Id, err = strconv.Atoi(r.PostFormValue("id"))
		if err != nil {
			return err
		}
		// calcul des ids stockage, pour transmettre à UpdatePlaq(),
		// qui va créer ou supprimer ou ne pas changer le(s) tas
        allStockages, err := model.GetStockagesActifs(ctx.DB)
		if err != nil {
			return err
		}
		idsStockages := []int{}
        for _, stockage := range(allStockages){
            if r.PostFormValue("stockage-" + strconv.Itoa(stockage.Id)) == "on" {
                idsStockages = append(idsStockages, stockage.Id)
            }
        }
		//
		err = model.UpdatePlaq(ctx.DB, chantier, idsStockages)
		if err != nil {
			return err
		}
		//
		redirect := "/chantier/plaquette/" + r.PostFormValue("id")
		err = chantier.ComputeLieudit(ctx.DB) // nécessaire pour appeler chantier.FullString()
		if err != nil {
			return err
		}
		err = model.AddRecent(ctx.DB, ctx.Config, &model.Recent{URL:redirect, Label:chantier.FullString()})
		if err != nil {
			return err
		}
		ctx.Redirect = redirect
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idChantier, err := strconv.Atoi(vars["id"])
		if err != nil {
			return err
		}
		chantier, err := model.GetPlaqFull(ctx.DB, idChantier)
		if err != nil {
			return err
		}
        allStockages, err := model.GetStockagesActifs(ctx.DB)
		if err != nil {
			return err
		}
		ctx.TemplateName = "plaq-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier le chantier : " + chantier.String(),
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/autocomplete/autocomplete.js"},
			},
			Details: detailsPlaqForm{
				Chantier:            chantier,
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "essence-"+chantier.Essence),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "exploitation-" + chantier.Exploitation),
				GranuloOptions:      webo.FmtOptions(WeboGranulo(), "granulo-" + chantier.Granulo),
				AllStockages:     	 allStockages,
				UrlAction:           "/chantier/plaquette/update/" + vars["id"],
			},
		}
		return nil
	}
}

// *********************************************************
func DeletePlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	chantier, err := model.GetPlaq(ctx.DB, id) // on retient l'année pour le redirect
	err = model.DeletePlaq(ctx.DB, id)
	if err != nil {
		return err
	}
	ctx.Redirect = "/chantier/plaquette/liste/" + strconv.Itoa(chantier.DateDebut.Year())
	return nil
}

// *********************************************************
// Fabrique un Plaq à partir des valeurs d'un formulaire.
// Auxiliaire de NewPlaq() et UpdatePlaq()
// Ne gère pas le champ Id
// Ne gère pas les stockages (tas)
func chantierPlaquetteForm2var(r *http.Request) (*model.Plaq, error) {
	chantier := &model.Plaq{}
	var err error
	if err = r.ParseForm(); err != nil {
		return chantier, err
	}
	//
	chantier.IdLieudit, err = strconv.Atoi(r.PostFormValue("id-lieudit"))
	if err != nil {
		return chantier, err
	}
	//
	chantier.DateDebut, err = time.Parse("2006-01-02", r.PostFormValue("date-debut"))
	if err != nil {
		return chantier, err
	}
	//
	chantier.DateFin, err = time.Parse("2006-01-02", r.PostFormValue("date-fin"))
	if err != nil {
		return chantier, err
	}
	//
	chantier.IdFermier, err = strconv.Atoi(r.PostFormValue("fermier"))
	if err != nil {
		return chantier, err
	}
	//
	chantier.IdUG, err = strconv.Atoi(r.PostFormValue("ug"))
	if err != nil {
		return chantier, err
	}
	//
	if r.PostFormValue("surface") != "" {
		chantier.Surface, err = strconv.ParseFloat(r.PostFormValue("surface"), 32)
		if err != nil {
			return chantier, err
		}
		chantier.Surface = tiglib.Round(chantier.Surface, 2)
	}
	//
	chantier.Granulo = strings.ReplaceAll(r.PostFormValue("granulo"), "granulo-", "")
	//
	chantier.Exploitation = strings.ReplaceAll(r.PostFormValue("exploitation"), "exploitation-", "")
	//
	chantier.Essence = strings.ReplaceAll(r.PostFormValue("essence"), "essence-", "")
	//
	if r.PostFormValue("frais-repas") != "" {
		chantier.FraisRepas, err = strconv.ParseFloat(r.PostFormValue("frais-repas"), 32)
		if err != nil {
			return chantier, err
		}
		chantier.FraisRepas = tiglib.Round(chantier.FraisRepas, 2)
	}
	//
	if r.PostFormValue("frais-reparation") != "" {
		chantier.FraisReparation, err = strconv.ParseFloat(r.PostFormValue("frais-reparation"), 32)
		if err != nil {
			return chantier, err
		}
		chantier.FraisReparation = tiglib.Round(chantier.FraisReparation, 2)
	}
	//
	return chantier, nil
}
