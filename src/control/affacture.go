/*
@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"strconv"
	"time"
)

type detailsAffactureForm struct {
	Acteur    *model.Acteur
	UrlAction string
}

// Affiche formulaire pour une "facture à l'envers"
func FormAffacture(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idActeur, err := strconv.Atoi(vars["id"])
	if err != nil {
		return werr.Wrap(err)
	}
	acteur, err := model.GetActeur(ctx.DB, idActeur)
	if err != nil {
		return werr.Wrap(err)
	}
	ctx.TemplateName = "affacture-form.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    "Affacture pour " + acteur.String(),
			CSSFiles: []string{"/static/css/form.css"},
		},
		Footer: ctxt.Footer{
			JSFiles: []string{"/static/js/toogle.js"},
		},
		Details: detailsAffactureForm{
			Acteur:    acteur,
			UrlAction: "/affacture/show",
		},
	}
	return nil
}

// *********************************************************
func ShowAffacture(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	var err error
	//
	// 1 - parse form
	//
	if err = r.ParseForm(); err != nil {
		return werr.Wrap(err)
	}
	var aff = model.Affacture{}
	aff.IdActeur, err = strconv.Atoi(r.PostFormValue("id-acteur"))
	if err != nil {
		return werr.Wrap(err)
	}
	aff.DateDebut, err = time.Parse("2006-01-02", r.PostFormValue("date-debut"))
	if err != nil {
		return werr.Wrap(err)
	}
	aff.DateFin, err = time.Parse("2006-01-02", r.PostFormValue("date-fin"))
	if err != nil {
		return werr.Wrap(err)
	}
	if r.PostFormValue("AB") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "AB")
	}
	if r.PostFormValue("DB") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "DB")
	}
	if r.PostFormValue("DC") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "DC")
	}
	if r.PostFormValue("BR") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "BR")
	}
	if r.PostFormValue("TR") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "TR")
	}
	if r.PostFormValue("TR-CO") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "TR-CO")
	}
	if r.PostFormValue("TR-OU") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "TR-OU")
	}
	if r.PostFormValue("RG") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "RG")
	}
	if r.PostFormValue("RG-CO") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "RG-CO")
	}
	if r.PostFormValue("RG-OU") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "RG-OU")
	}
	if r.PostFormValue("CG") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "CG")
	}
	if r.PostFormValue("CG-CO") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "CG-CO")
	}
	if r.PostFormValue("CG-OU") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "CG-OU")
	}
	if r.PostFormValue("LV") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "LV")
	}
	if r.PostFormValue("LV-CO") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "LV-CO")
	}
	if r.PostFormValue("LV-OU") == "on" {
		aff.TypesActivites = append(aff.TypesActivites, "LV-OU")
	}
	//
	// 2- récup info dans model
	//
	acteur, err := model.GetActeur(ctx.DB, aff.IdActeur)
	if err != nil {
		return werr.Wrap(err)
	}
	err = aff.ComputeItems(ctx.DB)
	if err != nil {
		return werr.Wrap(err)
	}
	//
	// 3 - génère PDF
	//
	var str string
	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("") // "" defaults to "cp1252"
	pdf.AddPage()
	MetaDataPDF(pdf, tr, ctx.Config, "Affacture")
	// Emetteur de la facture
	str = StringActeurFacture(acteur)
	pdf.SetXY(10, 10)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(100, 7, tr(str), "", "L", false)
	// Destinataire de l'affacture (BDL)
	str = tiglib.DateFrText(time.Now())
	pdf.SetXY(145, 15)
	pdf.Cell(50, 50, "Le "+tr(str))
	str = "FACTURE"
	pdf.SetXY(145, 25)
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(50, 70, str)
	str = ctx.Config.Affacture.Adresse
	pdf.SetXY(145, 68)
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(100, 7, tr(str), "", "L", false)
	//
	// items
	//
	var x, x0, colW, colH float64
	x0 = 10 // démarrage de l'affichage à gauche
	_, fontSize := pdf.GetFontSize()
	colW = 28 // largeur de toutes les colonnes - devrait en fait être calculé col par col
	colH = fontSize + 2
	maxY := 274 // mesuré empiriquement
	//_, topMarg,_,bottomMarg := pdf.GetMargins()
	//maxY := 297 - topMarg - bottomMarg
	//
	pdf.Ln(4 * colH)
	for _, item := range aff.Items {
		heightNeeded := float64(colH * float64(4+2*len(item.Lignes)))
		if heightNeeded+pdf.GetY() > float64(maxY) {
			pdf.AddPage()
		}
		// titre item
		pdf.SetX(x0)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(50, colH*2, tr(item.Titre+" "+tiglib.DateFrText(item.Date)))
		pdf.Ln(2 * colH)
		for _, ligne := range item.Lignes {
			// titre colonnes
			x = x0 + colW // Une cell vide (décalage pour les titres des lignes)
			pdf.SetX(x)
			pdf.SetFont("Arial", "B", 10)
			for _, col := range ligne.Colonnes {
				pdf.CellFormat(colW, colH, tr(col.Titre), "1", 0, "CM", false, 0, "")
			}
			pdf.Ln(colH)
			// titre ligne
			pdf.SetX(x0)
			pdf.CellFormat(colW, colH, tr(ligne.Titre), "1", 0, "CM", false, 0, "")
			// valeurs colonnes
			pdf.SetFont("Arial", "", 10)
			for _, col := range ligne.Colonnes {
				pdf.CellFormat(colW, colH, tr(col.Valeur), "1", 0, "CM", false, 0, "")
			}
			pdf.Ln(colH * 1.5)
		}
		pdf.Ln(2 * colH)
	}
	pdf.SetX(130)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(50, colH, tr("TOTAL GENERAL : "+strconv.FormatFloat(aff.TotalTTC, 'f', 2, 64)+" € TTC"))
	//
	return pdf.Output(w)
}
