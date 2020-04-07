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
	//"golang.org/x/text/encoding.Encoder"
	//"golang.org/x/text/encoding/charmap.Charmap"
	//"fmt"
)

type detailsBSPiedForm struct {
	Chantier            *model.BSPied
	UrlAction           string
	EssenceOptions      template.HTML
	ExploitationOptions template.HTML
}

type detailsBSPiedList struct {
	Chantiers       []*model.BSPied
	Annee           string   // année courante
	Annees          []string // toutes les années avec chantier
	TotalParEssence map[string]float64
}

// *********************************************************
func ListBspied(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	annee := vars["annee"]
	if annee == "" {
		// annee non spécifiée, on prend l'année courante
		annee = strconv.Itoa(time.Now().Year())
	}
	chantiers, err := model.GetBSPiedsOfYear(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	annees, err := model.GetBSPiedDifferentYears(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	totalParEssence := map[string]float64{}
	for _, essence := range model.AllEssenceCodes() {
		totalParEssence[essence] = 0
	}
	for _, ch := range chantiers {
		totalParEssence[ch.Essence] += ch.NStereCoupees
	}
	//
	ctx.TemplateName = "bspied-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Chantiers bois sur pied " + annee,
			JSFiles: []string{
				"/static/js/round.js",
				"/static/js/prix.js"},
		},
		Menu: "chantiers",
		Details: detailsBSPiedList{
			Chantiers:       chantiers,
			Annee:           annee,
			Annees:          annees,
			TotalParEssence: totalParEssence,
		},
	}
	return nil
}

// *********************************************************
// Process ou affiche form new
func NewBSPied(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		chantier, err := chantierBSPiedForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertBSPied(ctx.DB, chantier)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/bois-sur-pied/liste/" + strconv.Itoa(chantier.DateContrat.Year())
		return nil
	default:
		//
		// Affiche form
		//
		chantier := &model.BSPied{}
		chantier.Acheteur = &model.Acteur{}
		chantier.Lieudit = &model.Lieudit{}
		chantier.TVA = ctx.Config.TVABDL.BoisSurPied
		ctx.TemplateName = "bspied-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouveau chantier bois sur pied",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Details: detailsBSPiedForm{
				Chantier:            chantier,
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "CHOOSE_ESSENCE"),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "CHOOSE_EXPLOITATION"),
				UrlAction:           "/chantier/bois-sur-pied/new",
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
func UpdateBSPied(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		chantier, err := chantierBSPiedForm2var(r)
		chantier.TVA = ctx.Config.TVABDL.BoisSurPied
		if err != nil {
			return err
		}
		chantier.Id, err = strconv.Atoi(r.PostFormValue("id-chantier"))
		if err != nil {
			return err
		}
		err = model.UpdateBSPied(ctx.DB, chantier)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/bois-sur-pied/liste/" + strconv.Itoa(chantier.DateContrat.Year())
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
		chantier, err := model.GetBSPiedFull(ctx.DB, idChantier)
		if err != nil {
			return err
		}
		ctx.TemplateName = "bspied-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier un chantier bois sur pied",
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
			Details: detailsBSPiedForm{
				Chantier:            chantier,
				EssenceOptions:      webo.FmtOptions(WeboEssence(), chantier.Essence),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "exploitation-"+chantier.Exploitation),
				UrlAction:           "/chantier/bois-sur-pied/update/" + vars["id"],
			},
		}
		return nil
	}
}

// *********************************************************
func DeleteBSPied(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	chantier, err := model.GetBSPied(ctx.DB, id) // on retient l'année pour le redirect
	if err != nil {
		return err
	}
	err = model.DeleteBSPied(ctx.DB, id)
	if err != nil {
		return err
	}
	ctx.Redirect = "/chantier/bois-sur-pied/liste/" + strconv.Itoa(chantier.DateContrat.Year())
	return nil
}

// *********************************************************
// Fabrique un BSPied à partir des valeurs d'un formulaire.
// Auxiliaire de NewBSPied() et UpdateBSPied()
// Ne gère pas le champ Id
// Ne gère pas le champ TVA
func chantierBSPiedForm2var(r *http.Request) (*model.BSPied, error) {
	bsp := &model.BSPied{}
	var err error
	if err = r.ParseForm(); err != nil {
		return bsp, err
	}
	//
	bsp.IdAcheteur, err = strconv.Atoi(r.PostFormValue("id-acheteur"))
	if err != nil {
		return bsp, err
	}
	//
	bsp.DateContrat, err = time.Parse("2006-01-02", r.PostFormValue("datecontrat"))
	if err != nil {
		return bsp, err
	}
	//
	bsp.IdLieudit, err = strconv.Atoi(r.PostFormValue("id-lieudit"))
	if err != nil {
		return bsp, err
	}
	//
	for k, _ := range r.PostForm {
		if strings.HasPrefix(k, "parcelle-") {
			idP, _ := strconv.Atoi(strings.Replace(k, "parcelle-", "", -1))
			bsp.IdsParcelles = append(bsp.IdsParcelles, idP)
		}
	}
	//
	bsp.IdUG, _ = strconv.Atoi(strings.Replace(r.PostFormValue("ug"), "ug-", "", -1))
	//
	bsp.Exploitation = strings.ReplaceAll(r.PostFormValue("exploitation"), "exploitation-", "")
	//
	bsp.Essence = strings.ReplaceAll(r.PostFormValue("essence"), "essence-", "")
	//
	bsp.NStereContrat, err = strconv.ParseFloat(r.PostFormValue("nsterecontrat"), 32)
	if err != nil {
		return bsp, err
	}
	bsp.NStereContrat = tiglib.Round(bsp.NStereContrat, 2)
	//
	if r.PostFormValue("nsterecoupees") != "" {
		bsp.NStereCoupees, err = strconv.ParseFloat(r.PostFormValue("nsterecoupees"), 32)
		if err != nil {
			return bsp, err
		}
		bsp.NStereCoupees = tiglib.Round(bsp.NStereCoupees, 2)
	}
	//
	bsp.PrixStere, err = strconv.ParseFloat(r.PostFormValue("prixstere"), 32)
	if err != nil {
		return bsp, err
	}
	bsp.PrixStere = tiglib.Round(bsp.PrixStere, 2)
	//
	if r.PostFormValue("datefacture") != "" {
		bsp.DateFacture, err = time.Parse("2006-01-02", r.PostFormValue("datefacture"))
		if err != nil {
			return bsp, err
		}
	}
	//
	bsp.NumFacture = r.PostFormValue("numfacture")
	//
	bsp.Notes = r.PostFormValue("notes")
	//
	return bsp, nil
}

// *********************************************************
func ShowFactureBSPied(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	//
	bsp, err := model.GetBSPiedFull(ctx.DB, id)
	if err != nil {
		return err
	}
	//
	pdf := gofpdf.New("P", "mm", "A4", "")
	InitializeFacture(pdf)
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.AddPage()
	//
	MetaDataFacture(pdf, tr, ctx.Config, "Facture bois sur pied")
	HeaderFacture(pdf, tr, ctx.Config)
	FooterFacture(pdf, tr, ctx.Config)
	//
	// Acheteur
	//
	str := tr(bsp.Acheteur.String())
	if bsp.Acheteur.Adresse1 != "" {
		str += "\n" + tr(bsp.Acheteur.Adresse1)
	}
	if bsp.Acheteur.Adresse2 != "" {
		str += "\n" + tr(bsp.Acheteur.Adresse2)
	}
	if bsp.Acheteur.Cp != "" && bsp.Acheteur.Ville != "" {
		str += "\n" + bsp.Acheteur.Cp + " " + tr(bsp.Acheteur.Ville)
	} else if bsp.Acheteur.Cp != "" {
		str += "\n" + tr(bsp.Acheteur.Cp)
	} else if bsp.Acheteur.Ville != "" {
		str += "\n" + tr(bsp.Acheteur.Ville)
	}
	pdf.SetXY(60, 70)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(100, 7, str, "1", "C", false)
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
	pdf.MultiCell(wi, he, tiglib.DateFr(bsp.DateFacture), "LRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	pdf.MultiCell(wi, he, bsp.NumFacture, "RB", "C", false)
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
	pdf.MultiCell(wi, he, "P.U. H.T", "TRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w5
	pdf.MultiCell(wi, he, "Montant H.T", "TRB", "C", false)
	//
	x = x0
	y += he
	pdf.SetXY(x, y)
	wi = w1
	str = "Vente de bois de " + tr(model.LabelEssence(bsp.Essence)) + " sur pied"
	pdf.MultiCell(wi, he, str, "LRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w2
	pdf.MultiCell(wi, he, strconv.FormatFloat(bsp.NStereCoupees, 'f', 2, 64), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w3
	pdf.MultiCell(wi, he, tr("Stère"), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w4
	pdf.MultiCell(wi, he, strconv.FormatFloat(bsp.PrixStere, 'f', 2, 64), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w5
	prixHT := bsp.NStereCoupees * bsp.PrixStere
	pdf.MultiCell(wi, he, strconv.FormatFloat(prixHT, 'f', 2, 64), "RB", "C", false)
	//
	pdf.SetFont("Arial", "B", 10)
	x = x0 + w1
	y += he
	pdf.SetXY(x, y)
	wi = w2 + w3 + w4 + w5
	// @todo arriver à dire euro : € \u20AC
	pdf.MultiCell(wi, he, "Montant total E HT", "RBL", "C", false)
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
	pdf.MultiCell(wi, he, strconv.FormatFloat(bsp.TVA, 'f', 2, 64)+" %", "RB", "C", false)
	x += wi
	wi = w5
	pdf.SetXY(x, y)
	prixTVA := prixHT * bsp.TVA / 100
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
