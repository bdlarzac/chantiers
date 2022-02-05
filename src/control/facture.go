/******************************************************************************
    Fonctions communes à plusieurs factures

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    
    @history    2020-02-27 13:52:20+01:00, Thierry Graff : Creation
********************************************************************************/
package control

import (
	"bdl.local/bdl/model"
	"github.com/jung-kurt/gofpdf"
)

// Initialisations communes à toutes les factures émises par BDL
// Ne concerne pas les affactures (factures à l'envers)
func InitializeFacture(pdf *gofpdf.Fpdf) {
	pdf.SetAutoPageBreak(true, 1) // met la marge du bas à 1cm
}

// Métadonnées du PDF (pour factures et affactures)
func MetaDataPDF(pdf *gofpdf.Fpdf, tr func(string) string, conf *model.Config, titre string) {
	pdf.SetTitle(tr(titre), true)
	pdf.SetAuthor(conf.Facture.Auteur, true)
	pdf.SetCreator(conf.Facture.Createur, true)
}

// Header commun à toutes les factures
func HeaderFacture(pdf *gofpdf.Fpdf, tr func(string) string, conf *model.Config) {
	//
	var opt gofpdf.ImageOptions
	opt.ImageType = "jpg"
	pdf.ImageOptions("static/logo-bdl-facture.jpg", 10, 10, 70, 0, false, opt, 0, "")
	//
	pdf.SetXY(150, 20)
	pdf.SetFont("Arial", "B", 24)
	pdf.Cell(100, 15, "FACTURE")
	//
	pdf.SetXY(10, 30)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(100, 15, tr(conf.Facture.Adresse))
	//
	pdf.SetXY(10, 33)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(100, 20, tr(conf.Facture.Tel))
	//
	pdf.SetXY(40, 33)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(100, 20, tr(conf.Facture.Email))
}

// Footer commun à toutes les factures
// Attention, ne marche que si pdf.InitializeFacture() a été appelé avant (pour réduire la marge du bas)
func FooterFacture(pdf *gofpdf.Fpdf, tr func(string) string, conf *model.Config) {
	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(80, 284)
	pdf.MultiCell(50, 5, "Site internet : "+tr(conf.Facture.SiteWeb), "", "C", false)
	pdf.SetXY(50, 290)
	pdf.MultiCell(100, 3, tr("N° SIRET : "+tr(conf.Facture.Siret)+" - N° TVA : "+tr(conf.Facture.TVA)), "", "C", false)
}

// Renvoie une string permettant d'afficher un acteur avec son adresse dans une facture
func StringActeurFacture(acteur *model.Acteur) string {
	str := acteur.String()
	if acteur.Adresse1 != "" {
		str += "\n" + acteur.Adresse1
	}
	if acteur.Adresse2 != "" {
		str += "\n" + acteur.Adresse2
	}
	if acteur.Cp != "" && acteur.Ville != "" {
		str += "\n" + acteur.Cp + " " + acteur.Ville
	} else if acteur.Cp != "" {
		str += "\n" + acteur.Cp
	} else if acteur.Ville != "" {
		str += "\n" + acteur.Ville
	}
	return str
}
