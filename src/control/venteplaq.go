package control

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/webo"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	//"fmt"
)

type detailsVentePlaqForm struct {
	Vente              *model.VentePlaq
	FournisseurOptions template.HTML
	UrlAction          string
}
type detailsVentePlaqList struct {
	Ventes []*model.VentePlaq
	Annee  string   // année courante
	Annees []string // toutes les années avec chantier
}
type detailsVentePlaqShow struct {
	Vente *model.VentePlaq
}

// *********************************************************
func ListVentePlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	annee := vars["annee"]
	if annee == "" {
		// annee non spécifiée, on prend l'année courante
		annee = strconv.Itoa(time.Now().Year())
	}
	ventes, err := model.GetVentePlaqsOfYear(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	annees, err := model.GetVentePlaqDifferentYears(ctx.DB, annee)
	if err != nil {
		return err
	}
	ctx.TemplateName = "venteplaq-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Ventes plaquettes " + annee,
			JSFiles: []string{
				"/static/js/venteplaq.js",
				"/static/js/round.js",
				"/static/js/prix.js"},
		},
		Menu:   "ventes",
		Footer: ctxt.Footer{},
		Details: detailsVentePlaqList{
			Ventes: ventes,
			Annee:  annee,
			Annees: annees,
		},
	}
	return nil
}

// *********************************************************
func ShowVentePlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idVente, _ := strconv.Atoi(vars["id-vente"])
	vente, err := model.GetVentePlaqFull(ctx.DB, idVente)
	if err != nil {
		return err
	}
	ctx.TemplateName = "venteplaq-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Vente " + vente.String(),
			JSFiles: []string{
				"/static/js/venteplaq.js",
				"/static/js/round.js",
				"/static/js/prix.js"},
		},
		Menu: "ventes",
		Footer: ctxt.Footer{
			JSFiles: []string{},
		},
		Details: detailsVentePlaqShow{
			Vente: vente,
		},
	}
	return nil
}

// *********************************************************
// Process ou affiche form new
func NewVentePlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		vente, err := ventePlaqForm2var(r)
		if err != nil {
			return err
		}
		vente.TVA = ctx.Config.TVABDL.VentePlaquettes
		vente.FactureLivraisonTVA = ctx.Config.TVABDL.Livraison
		idVente, err := model.InsertVentePlaq(ctx.DB, vente)
		if err != nil {
			return err
		}
		ctx.Redirect = "/vente/" + strconv.Itoa(idVente)
		return nil
	default:
		//
		// Affiche form
		//
		vente := &model.VentePlaq{}
		vente.Client = &model.Acteur{}
		vente.Fournisseur = &model.Acteur{}
		vente.TVA = ctx.Config.TVABDL.VentePlaquettes
		ctx.TemplateName = "venteplaq-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouvelle vente de plaquettes",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Menu: "ventes",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/autocomplete/autocomplete.js"},
			},
			Details: detailsVentePlaqForm{
				Vente:              vente,
				FournisseurOptions: webo.FmtOptions(WeboFournisseur(ctx), "CHOOSE_FOURNISSEUR"),
				UrlAction:          "/vente/new",
			},
		}
		return nil
	}
}

// *********************************************************
// Process ou affiche form update
func UpdateVentePlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		vente, err := ventePlaqForm2var(r)
		if err != nil {
			return err
		}
		vente.Id, err = strconv.Atoi(r.PostFormValue("id-vente"))
		if err != nil {
			return err
		}
		vente.TVA = ctx.Config.TVABDL.VentePlaquettes
		vente.FactureLivraisonTVA = ctx.Config.TVABDL.Livraison
		err = model.UpdateVentePlaq(ctx.DB, vente)
		if err != nil {
			return err
		}
		ctx.Redirect = "/vente/" + r.PostFormValue("id-vente")
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idVente, err := strconv.Atoi(vars["id-vente"])
		if err != nil {
			return err
		}
		vente, err := model.GetVentePlaqFull(ctx.DB, idVente)
		if err != nil {
			return err
		}
		ctx.TemplateName = "venteplaq-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier la vente : " + vente.String(),
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Menu: "ventes",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/autocomplete/autocomplete.js"},
			},
			Details: detailsVentePlaqForm{
				FournisseurOptions: webo.FmtOptions(WeboFournisseur(ctx), strconv.Itoa(vente.IdFournisseur)),
				Vente:              vente,
				UrlAction:          "/vente/update/" + vars["id-vente"],
			},
		}
		return nil
	}
}

// *********************************************************
func DeleteVentePlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id-vente"])
	if err != nil {
		return err
	}
	vente, err := model.GetVentePlaq(ctx.DB, id)
	err = model.DeleteVentePlaq(ctx.DB, id)
	if err != nil {
		return err
	}
	ctx.Redirect = "/vente/liste/" + strconv.Itoa(vente.DateVente.Year())
	return nil
}

// *********************************************************
// Fabrique une VentePlaq à partir des valeurs d'un formulaire.
// Auxiliaire de NewVentePlaq() et UpdateVentePlaq()
// Ne gère pas le champ Id
// Ne gère pas les champs TVA et FactureLivraisonTVA (car viennent de la config)
func ventePlaqForm2var(r *http.Request) (*model.VentePlaq, error) {
	vente := &model.VentePlaq{}
	var err error
	if err = r.ParseForm(); err != nil {
		return vente, err
	}
	//
	vente.IdClient, err = strconv.Atoi(r.PostFormValue("id-client"))
	if err != nil {
		return vente, err
	}
	//
	vente.IdFournisseur, err = strconv.Atoi(r.PostFormValue("fournisseur"))
	if err != nil {
		return vente, err
	}
	//
	vente.PUHT, err = strconv.ParseFloat(r.PostFormValue("puht"), 32)
	if err != nil {
		return vente, err
	}
	vente.PUHT = tiglib.Round(vente.PUHT, 2)
	//
	vente.DateVente, err = time.Parse("2006-01-02", r.PostFormValue("datevente"))
	if err != nil {
		return vente, err
	}
	//
	// Facture
	//
	vente.NumFacture = r.PostFormValue("numfacture")
	//
	if r.PostFormValue("datefacture") != "" {
		vente.DateFacture, err = time.Parse("2006-01-02", r.PostFormValue("datefacture"))
		if err != nil {
			return vente, err
		}
	}
	//
	vente.FactureLivraison = false
	if r.PostFormValue("facturelivraison") == "on" {
		vente.FactureLivraison = true
		//
		vente.FactureLivraisonPUHT, err = strconv.ParseFloat(r.PostFormValue("facturelivraisonpuht"), 32)
		if err != nil {
			return vente, err
		}
		vente.FactureLivraisonPUHT = tiglib.Round(vente.FactureLivraisonPUHT, 2)
	} // sinon FactureLivraisonPUHT reste à 0 => ok
	//
	vente.FactureNotes = false
	if r.PostFormValue("facturenotes") == "on" {
		vente.FactureNotes = true
	}
	//
	vente.Notes = r.PostFormValue("notes")
	//
	return vente, nil
}

// *********************************************************
func ShowFactureVentePlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	//
	vente, err := model.GetVentePlaqFull(ctx.DB, id)
	if err != nil {
		return err
	}
	//
	pdf := gofpdf.New("P", "mm", "A4", "")
	InitializeFacture(pdf)
	tr := pdf.UnicodeTranslatorFromDescriptor("") // "" defaults to "cp1252"
	pdf.AddPage()
	MetaDataFacture(pdf, "Facture vente plaquettes", tr)
	HeaderFacture(pdf, tr)
	FooterFacture(pdf, tr)
	//
	// Client
	//
	str := tr(vente.Client.String())
	if vente.Client.Adresse1 != "" {
		str += "\n" + tr(vente.Client.Adresse1)
	}
	if vente.Client.Adresse2 != "" {
		str += "\n" + tr(vente.Client.Adresse2)
	}
	if vente.Client.Cp != "" && vente.Client.Ville != "" {
		str += "\n" + vente.Client.Cp + " " + tr(vente.Client.Ville)
	} else if vente.Client.Cp != "" {
		str += "\n" + vente.Client.Cp
	} else if vente.Client.Ville != "" {
		str += "\n" + tr(vente.Client.Ville)
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
	pdf.MultiCell(wi, he, tiglib.DateFr(vente.DateFacture), "LRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	pdf.MultiCell(wi, he, "???", "RB", "C", false)
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
	pdf.SetFont("Arial", "B", 10)
	x = x0
	y += he
	pdf.SetXY(x, y)
	wi = w1
	pdf.MultiCell(wi, he, tr("Vente de plaquettes forestières"), "LRB", "C", false)
	pdf.SetFont("Arial", "", 10)
	x += wi
	pdf.SetXY(x, y)
	wi = w2
	pdf.MultiCell(wi, he, strconv.FormatFloat(vente.Qte, 'f', 2, 64), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w3
	pdf.MultiCell(wi, he, "MAP", "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w4
	pdf.MultiCell(wi, he, strconv.FormatFloat(vente.PUHT, 'f', 2, 64), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w5
	prixHT := vente.Qte * vente.PUHT
	pdf.MultiCell(wi, he, strconv.FormatFloat(prixHT, 'f', 2, 64), "RB", "C", false)
	//
	pdf.SetFont("Arial", "B", 10)
	x = x0 + w1
	y += he
	pdf.SetXY(x, y)
	wi = w2 + w3 + w4 + w5
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
	TVA := 20.0 ////// @todo recup de la conf ou à stocker en base
	pdf.MultiCell(wi, he, strconv.FormatFloat(TVA, 'f', 2, 64)+" %", "RB", "C", false)
	x += wi
	wi = w5
	pdf.SetXY(x, y)
	//prixTVA := prixHT * vente.TVA / 100
	prixTVA := prixHT * TVA / 100
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
