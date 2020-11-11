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
	"github.com/jung-kurt/gofpdf"
	//"fmt"
)

type detailsChautreForm struct {
	UrlAction           string
	EssenceOptions      template.HTML
	ExploitationOptions template.HTML
	ValorisationOptions template.HTML
	UniteOptions        template.HTML
	Chantier            *model.Chautre
}

type detailsChautreList struct {
	Chantiers []*model.Chautre
	Annee     string   // année courante
	Annees    []string // toutes les années avec chantier
}

// *********************************************************
func ListChautre(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	annee := vars["annee"]
	if annee == "" {
		// annee non spécifiée, on prend l'année courante
		annee = strconv.Itoa(time.Now().Year())
	}
	chantiers, err := model.GetChautresOfYear(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	annees, err := model.GetChautreDifferentYears(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	ctx.TemplateName = "chautre-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Chantiers autres valorisations " + annee,
			JSFiles: []string{
				"/static/js/round.js",
				"/static/js/prix.js"},
		},
		Menu: "chantiers",
		Details: detailsChautreList{
			Chantiers: chantiers,
			Annee:     annee,
			Annees:    annees,
		},
	}
	return nil
}

// *********************************************************
// Process ou affiche form new
func NewChautre(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		chantier, err := chautreForm2var(r)
		chantier.TVA = ctx.Config.TVABDL.AutreValorisation
		if err != nil {
			return err
		}
		_, err = model.InsertChautre(ctx.DB, chantier)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/autre/liste/" + strconv.Itoa(chantier.DateContrat.Year())
		return nil
	default:
		//
		// Affiche form
		//
		chantier := &model.Chautre{}
		chantier.Client = &model.Acteur{}
		chantier.Lieudit = &model.Lieudit{}
		chantier.TVA = ctx.Config.TVABDL.AutreValorisation
		ctx.TemplateName = "chautre-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouveau chantier autres valorisations",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Details: detailsChautreForm{
				Chantier:            chantier,
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "CHOOSE_ESSENCE"),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "CHOOSE_EXPLOITATION"),
				ValorisationOptions: webo.FmtOptions(WeboChautreValo(), "CHOOSE_VALORISATION"),
				UniteOptions:        webo.FmtOptions(WeboChautreUnite(), "CHOOSE_UNITE"),
				UrlAction:           "/chantier/autre/new",
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/autocomplete/autocomplete.js"},
			},
		}
		return nil
	}
}

// *********************************************************
// Process ou affiche form update
func UpdateChautre(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		chantier, err := chautreForm2var(r)
		if err != nil {
			return err
		}
		chantier.Id, err = strconv.Atoi(r.PostFormValue("id-chantier"))
		if err != nil {
			return err
		}
		err = model.UpdateChautre(ctx.DB, chantier)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/autre/liste/" + strconv.Itoa(chantier.DateContrat.Year())
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
		chantier, err := model.GetChautreFull(ctx.DB, idChantier)
		if err != nil {
			return err
		}
		ctx.TemplateName = "chautre-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier un chantier autres valorisations",
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
			Details: detailsChautreForm{
				Chantier:            chantier,
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "essence-"+chantier.Essence),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "exploitation-"+chantier.Exploitation),
				ValorisationOptions: webo.FmtOptions(WeboChautreValo(), "valorisation-"+chantier.TypeValo),
				UniteOptions:        webo.FmtOptions(WeboChautreUnite(), "unite-"+chantier.Unite),
				UrlAction:           "/chantier/autre/update/" + vars["id"],
			},
		}
		return nil
	}
}

// *********************************************************
func DeleteChautre(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	chantier, err := model.GetChautre(ctx.DB, id) // on retient l'année pour le redirect
	if err != nil {
		return err
	}
	err = model.DeleteChautre(ctx.DB, id)
	if err != nil {
		return err
	}
	ctx.Redirect = "/chantier/autre/liste/" + strconv.Itoa(chantier.DateContrat.Year())
	return nil
}

// *********************************************************
// Fabrique un BSPied à partir des valeurs d'un formulaire.
// Auxiliaire de NewBSPied() et UpdateBSPied()
// Ne gère pas le champ Id
// Ne gère pas le champ TVA (tiré de la config)
func chautreForm2var(r *http.Request) (*model.Chautre, error) {
	chantier := &model.Chautre{}
	var err error
	if err = r.ParseForm(); err != nil {
		return chantier, err
	}
	//
	chantier.IdClient, err = strconv.Atoi(r.PostFormValue("id-client"))
	if err != nil {
		return chantier, err
	}
	//
	chantier.DateContrat, err = time.Parse("2006-01-02", r.PostFormValue("datecontrat"))
	if err != nil {
		return chantier, err
	}
	//
	chantier.IdLieudit, err = strconv.Atoi(r.PostFormValue("id-lieudit"))
	if err != nil {
		return chantier, err
	}
	//
	chantier.IdUG, _ = strconv.Atoi(strings.Replace(r.PostFormValue("ug"), "ug-", "", -1))
	//
	chantier.TypeValo = strings.Replace(r.PostFormValue("typevalo"), "valorisation-", "", -1)
	//
	chantier.Exploitation = strings.ReplaceAll(r.PostFormValue("exploitation"), "exploitation-", "")
	//
	chantier.Essence = strings.ReplaceAll(r.PostFormValue("essence"), "essence-", "")
	//
	chantier.Volume, err = strconv.ParseFloat(r.PostFormValue("volume"), 32)
	if err != nil {
		return chantier, err
	}
	chantier.Volume = tiglib.Round(chantier.Volume, 2)
	//
	chantier.Unite = strings.Replace(r.PostFormValue("unite"), "unite-", "", -1)
	//
	chantier.PUHT, err = strconv.ParseFloat(r.PostFormValue("puht"), 32)
	if err != nil {
		return chantier, err
	}
	chantier.PUHT = tiglib.Round(chantier.PUHT, 2)
	//
	if r.PostFormValue("datefacture") != "" {
		chantier.DateFacture, err = time.Parse("2006-01-02", r.PostFormValue("datefacture"))
		if err != nil {
			return chantier, err
		}
	}
	//
	chantier.NumFacture = r.PostFormValue("numfacture")
	//
	chantier.Notes = r.PostFormValue("notes")
	//
	return chantier, nil
}

// *********************************************************
func ShowFactureChautre(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	//
	chantier, err := model.GetChautreFull(ctx.DB, id)
	if err != nil {
		return err
	}
	//
	pdf := gofpdf.New("P", "mm", "A4", "")
	InitializeFacture(pdf)
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.AddPage()
	//
	MetaDataPDF(pdf, tr, ctx.Config, "Facture bois sur pied")
	HeaderFacture(pdf, tr, ctx.Config)
	FooterFacture(pdf, tr, ctx.Config)
	//
	var str string
	//
	// Acheteur
	//
	pdf.SetXY(60, 70)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(100, 7, tr(StringActeurFacture(chantier.Client)), "1", "C", false)
	//
	// Date  + n° facture
	//
	var x0, x, y, wi, he float64
	//
	x0 = 10
	x = x0
	y = 110
	wi = 50
	he = 6
	//
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(x, y)
	pdf.MultiCell(wi, he, "Date", "1", "C", false)
	//
	x += wi
	pdf.SetXY(x, y)
	pdf.MultiCell(wi, he, tr("Facture n°"), "TRB", "C", false)
	//
	pdf.SetFont("Arial", "", 10)
	x = 10
	y += he
	//
	pdf.SetXY(x, y)
	pdf.MultiCell(wi, he, tiglib.DateFr(chantier.DateFacture), "LRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	pdf.MultiCell(wi, he, chantier.NumFacture, "RB", "C", false)
	//
	// Tableau principal
	//
	var w1, w2, w3, w4, w5 = 70.0, 20.0, 20.0, 30.0, 30.0
	x = x0
	y = 140
	pdf.SetXY(x, y)
	wi = w1
	pdf.MultiCell(wi, he, tr("Désignation"), "1", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w2
	pdf.MultiCell(wi, he, tr("Quantité"), "TRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w3
	pdf.MultiCell(wi, he, tr("Unité"), "TRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w4
	pdf.MultiCell(wi, he, tr("P.U. € H.T"), "TRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w5
	pdf.MultiCell(wi, he, tr("Montant € H.T"), "TRB", "C", false)
	//
	x = x0
	y += he
	pdf.SetXY(x, y)
	wi = w1
	str = "Vente " + tr(model.LabelValorisation(chantier.TypeValo)) + " " + tr(model.LabelEssence(chantier.Essence))
	pdf.MultiCell(wi, he, str, "LRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w2
	pdf.MultiCell(wi, he, strconv.FormatFloat(chantier.Volume, 'f', 2, 64), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w3
	pdf.MultiCell(wi, he, tr(model.LabelUnite(chantier.Unite)), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w4
	pdf.MultiCell(wi, he, strconv.FormatFloat(chantier.PUHT, 'f', 2, 64), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w5
	prixHT := chantier.Volume * chantier.PUHT
	pdf.MultiCell(wi, he, strconv.FormatFloat(prixHT, 'f', 2, 64), "RB", "C", false)
	//
	pdf.SetFont("Arial", "B", 10)
	x = x0 + w1
	y += he
	pdf.SetXY(x, y)
	wi = w2 + w3 + w4 + w5
	// @todo arriver à dire euro : € \u20AC
	pdf.MultiCell(wi, he, tr("Montant total € HT"), "RBL", "C", false)
	//
	pdf.SetFont("Arial", "", 10)
	x = x0 + w1
	y += he
	pdf.SetXY(x, y)
	wi = w2 + w3
	pdf.MultiCell(wi, he, "Montant TVA", "RBL", "C", false)
	x += wi
	wi = w4
	pdf.SetXY(x, y)
	pdf.MultiCell(wi, he, strconv.FormatFloat(chantier.TVA, 'f', 2, 64)+" %", "RB", "C", false)
	x += wi
	wi = w5
	pdf.SetXY(x, y)
	prixTVA := prixHT * chantier.TVA / 100
	pdf.MultiCell(wi, he, strconv.FormatFloat(prixTVA, 'f', 2, 64), "RB", "C", false)
	//
	pdf.SetFont("Arial", "B", 10)
	x = x0 + w1
	y += 2 * he
	pdf.SetXY(x, y)
	wi = w2 + w3 + w4 + w5
	pdf.MultiCell(wi, he, "Montant total TTC", "1", "C", false)
	pdf.SetFont("Arial", "", 10)
	x = x0 + w1
	y += he
	pdf.SetXY(x, y)
	wi = w2 + w3 + w4
	pdf.MultiCell(wi, he, tr("Net à payer en euros"), "RBL", "C", false)
	pdf.SetFont("Arial", "B", 10)
	x += wi
	wi = w5
	pdf.SetXY(x, y)
	prixTTC := prixHT + prixTVA
	pdf.MultiCell(wi, he, strconv.FormatFloat(prixTTC, 'f', 2, 64), "RB", "C", false)
	//
	return pdf.Output(w)
}
