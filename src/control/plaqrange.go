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

type detailsPlaqRangeForm struct {
	GlTVAOptions template.HTML
	CoTVAOptions template.HTML
	OuTVAOptions template.HTML
	Rangement    *model.PlaqRange
	UrlAction    string
}

// *********************************************************
// Process ou affiche form new
func NewPlaqRange(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		pr, err := plaqRangeForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertPlaqRange(ctx.DB, pr)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/plaquette/" + strconv.Itoa(pr.IdChantier) + "/chantiers"
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idChantierStr := vars["id-chantier"]
		idChantier, err := strconv.Atoi(idChantierStr)
		if err != nil {
			return err
		}
		pr := &model.PlaqRange{}
		pr.Conducteur = &model.Acteur{}
		pr.IdChantier = idChantier
		pr.Chantier, err = model.GetPlaq(ctx.DB, idChantier)
		if err != nil {
			return err
		}
		err = pr.Chantier.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		err = pr.Chantier.ComputeLieudit(ctx.DB) // Pour le nom du chantier
		if err != nil {
			return err
		}
		ctx.TemplateName = "plaqrange-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouveau rangement plaquettes",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{"/static/autocomplete/autocomplete.js"},
			},
			Details: detailsPlaqRangeForm{
				Rangement:    pr,
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL"), "CHOOSE_TVA_GL"),
				CoTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_CO"), "CHOOSE_TVA_CO"),
				OuTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_OU"), "CHOOSE_TVA_OU"),
				UrlAction:    "/chantier/plaquette/" + idChantierStr + "/range/new",
			},
		}
		return nil
	}
}

// *********************************************************
// Process ou affiche form update
func UpdatePlaqRange(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		pr, err := plaqRangeForm2var(r)
		if err != nil {
			return err
		}
		pr.Id, err = strconv.Atoi(r.PostFormValue("id-pr"))
		if err != nil {
			return err
		}
		err = model.UpdatePlaqRange(ctx.DB, pr)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/plaquette/" + strconv.Itoa(pr.IdChantier) + "/chantiers"
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idPr, err := strconv.Atoi(vars["id-pr"])
		if err != nil {
			return err
		}
		pr, err := model.GetPlaqRange(ctx.DB, idPr)
		if err != nil {
			return err
		}
		pr.Conducteur, err = model.GetActeur(ctx.DB, pr.IdConducteur)
		if err != nil {
			return err
		}
		pr.Chantier, err = model.GetPlaq(ctx.DB, pr.IdChantier)
		if err != nil {
			return err
		}
		err = pr.Chantier.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		err = pr.Chantier.ComputeLieudit(ctx.DB) // pour le nom du chantier
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		ctx.TemplateName = "plaqrange-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier un rangement plaquette ",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{"/static/autocomplete/autocomplete.js"},
			},
			Details: detailsPlaqRangeForm{
				Rangement:    pr,
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL"), strconv.FormatFloat(pr.GlTVA, 'f', 1, 64)),
				CoTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_CO"), strconv.FormatFloat(pr.CoTVA, 'f', 1, 64)),
				OuTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_OU"), strconv.FormatFloat(pr.OuTVA, 'f', 1, 64)),
				UrlAction:    "/chantier/plaquette/" + vars["id-chantier"] + "/range/update/" + vars["id-pr"],
			},
		}
		return nil
	}
	return nil
}

// *********************************************************
func DeletePlaqRange(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idPr, err := strconv.Atoi(vars["id-pr"])
	if err != nil {
		return err
	}
	err = model.DeletePlaqRange(ctx.DB, idPr)
	if err != nil {
		return err
	}
	ctx.Redirect = "/chantier/plaquette/" + vars["id-chantier"] + "/chantiers"
	return nil
}

// *********************************************************
// Fabrique une PlaqRange à partir des valeurs d'un formulaire.
// Auxiliaire de NewPlaqRange() et UpdatePlaqRange()
// Ne gère pas le champ Id
func plaqRangeForm2var(r *http.Request) (*model.PlaqRange, error) {
	pr := &model.PlaqRange{}
	var err error
	if err = r.ParseForm(); err != nil {
		return pr, err
	}
	pr.IdChantier, err = strconv.Atoi(r.PostFormValue("id-chantier"))
	if err != nil {
		return pr, err
	}
	pr.IdTas, err = strconv.Atoi(strings.ReplaceAll(r.PostFormValue("tas"), "tas-", ""))
	if err != nil {
		return pr, err
	}
	pr.IdConducteur, err = strconv.Atoi(r.PostFormValue("id-conducteur"))
	if err != nil {
		return pr, err
	}
	pr.DateRange, err = time.Parse("2006-01-02", r.PostFormValue("daterange"))
	if err != nil {
		return pr, err
	}
	//
	if r.PostFormValue("type-cout") == "cout-global" {
		pr.TypeCout = "G"
	} else {
		pr.TypeCout = "D"
	}
	//
	if pr.TypeCout == "G" {
		//
		// coût global
		//
		pr.GlPrix, err = strconv.ParseFloat(r.PostFormValue("glprix"), 32)
		if err != nil {
			return pr, err
		}
		pr.GlPrix = tiglib.Round(pr.GlPrix, 2)
		//
		pr.GlTVA, err = strconv.ParseFloat(r.PostFormValue("gltva"), 32)
		if err != nil {
			return pr, err
		}
		pr.GlTVA = tiglib.Round(pr.GlTVA, 2)
		//
		if r.PostFormValue("gldatepay") != "" {
			pr.GlDatePay, err = time.Parse("2006-01-02", r.PostFormValue("gldatepay"))
			if err != nil {
				return pr, err
			}
		}
	} else {
		//
		// conducteur
		//
		pr.CoNheure, err = strconv.ParseFloat(r.PostFormValue("conheure"), 32)
		if err != nil {
			return pr, err
		}
		pr.CoNheure = tiglib.Round(pr.CoNheure, 2)
		//
		pr.CoPrixH, err = strconv.ParseFloat(r.PostFormValue("coprixh"), 32)
		if err != nil {
			return pr, err
		}
		pr.CoPrixH = tiglib.Round(pr.CoPrixH, 2)
		//
		pr.CoTVA, err = strconv.ParseFloat(r.PostFormValue("cotva"), 32)
		if err != nil {
			return pr, err
		}
		//
		if r.PostFormValue("codatepay") != "" {
			pr.CoDatePay, err = time.Parse("2006-01-02", r.PostFormValue("codatepay"))
			if err != nil {
				return pr, err
			}
		}
		//
		// outil
		//
		pr.OuPrix, err = strconv.ParseFloat(r.PostFormValue("ouprix"), 32)
		if err != nil {
			return pr, err
		}
		pr.OuPrix = tiglib.Round(pr.OuPrix, 2)
		//
		pr.OuTVA, err = strconv.ParseFloat(r.PostFormValue("outva"), 32)
		if err != nil {
			return pr, err
		}
		pr.OuTVA = tiglib.Round(pr.OuTVA, 2)
		//
		if r.PostFormValue("oudatepay") != "" {
			pr.OuDatePay, err = time.Parse("2006-01-02", r.PostFormValue("oudatepay"))
			if err != nil {
				return pr, err
			}
		}
	}
	//
	pr.Notes = r.PostFormValue("notes")
	//
	return pr, nil
}
