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
)

type detailsChautreForm struct {
	Chantier            *model.Chautre
	TypeChantier        string
	UrlAction           string
	EssenceOptions      template.HTML
	ExploitationOptions template.HTML
	ValorisationOptions template.HTML
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
	titrePage := "Chantiers autres valorisations " + annee
	ctx.TemplateName = "chautre-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: titrePage,
			JSFiles: []string{
				"/static/js/round.js",
				"/view/common/prix.js"},
		},
		Menu: "chantiers",
		Details: detailsChautreList{
			Chantiers: chantiers,
			Annee:     annee,
			Annees:    annees,
		},
	}
	err = model.AddRecent(ctx.DB, ctx.Config, &model.Recent{URL: r.URL.String(), Label: titrePage})
	if err != nil {
		return err
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
		if err != nil {
			return err
		}
		chantier.TVA = ctx.Config.TVABDL.AutreValorisation
		// calcul des ids UG, Lieudit et Fermier, pour transmettre à InsertChautre()
		var idsUG, idsLieudit, idsFermier []int
		var id int
		for key, val := range r.PostForm {
			if strings.Index(key, "ug-") == 0 {
				// ex : ug-0:[6] (6 est l'id UG)
				id, err = strconv.Atoi(val[0])
				if err != nil {
					return err
				}
				idsUG = append(idsUG, id)
			}
			if strings.Index(key, "lieudit-") == 0 {
				// ex : lieudit-164:[on] (164 est l'id lieudit)
				id, err = strconv.Atoi(key[8:])
				if err != nil {
					return err
				}
				idsLieudit = append(idsLieudit, id)
			}
			if strings.Index(key, "fermier-") == 0 {
				// ex : fermier-25:[on] (25 est l'id fermier)
				id, err = strconv.Atoi(key[8:])
				if err != nil {
					return err
				}
				idsFermier = append(idsFermier, id)
			}
		}
		//
		_, err = model.InsertChautre(ctx.DB, chantier, idsUG, idsLieudit, idsFermier)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/autre/liste/" + strconv.Itoa(chantier.DateContrat.Year())
		// model.AddRecent() inutile puisqu'on est redirigé vers la liste, où AddRecent() est exécuté
		return nil
	default:
		//
		// Affiche form
		//
		chantier := &model.Chautre{}
		chantier.Acheteur = &model.Acteur{}
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
				UrlAction:           "/chantier/autre/new",
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/autocomplete/autocomplete.js",
					"/view/common/checkActeur.js",
					"/view/common/getActeurPossibles.js"},
			},
		}
		// model.AddRecent() inutile puisqu'on est redirigé vers la liste, où AddRecent() est exécuté
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
		chantier.TVA = ctx.Config.TVABDL.AutreValorisation
		if err != nil {
			return err
		}
		// calcul des ids UG, Lieudit et Fermier, pour transmettre à UpdateChautre()
		var idsUG, idsLieudit, idsFermier []int
		var id int
		for key, val := range r.PostForm {
			if strings.Index(key, "ug-") == 0 {
				// ex : ug-0:[6] (6 est l'id UG)
				id, err = strconv.Atoi(val[0])
				if err != nil {
					return err
				}
				idsUG = append(idsUG, id)
			}
			if strings.Index(key, "lieudit-") == 0 {
				// ex : lieudit-164:[on] (164 est l'id lieudit)
				id, err = strconv.Atoi(key[8:])
				if err != nil {
					return err
				}
				idsLieudit = append(idsLieudit, id)
			}
			if strings.Index(key, "fermier-") == 0 {
				// ex : fermier-25:[on] (25 est l'id fermier)
				id, err = strconv.Atoi(key[8:])
				if err != nil {
					return err
				}
				idsFermier = append(idsFermier, id)
			}
		}
		//
		err = model.UpdateChautre(ctx.DB, chantier, idsUG, idsLieudit, idsFermier)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/autre/liste/" + strconv.Itoa(chantier.DateContrat.Year())
		// model.AddRecent() inutile puisqu'on est redirigé vers la liste, où AddRecent() est exécuté
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
					"/static/autocomplete/autocomplete.js",
					"/view/common/checkActeur.js",
					"/view/common/getActeurPossibles.js"},
			},
			Details: detailsChautreForm{
				Chantier:            chantier,
				TypeChantier:        "chautre",
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "essence-"+chantier.Essence),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "exploitation-"+chantier.Exploitation),
				ValorisationOptions: webo.FmtOptions(WeboChautreValo(), "valorisation-"+chantier.TypeValo),
				UrlAction:           "/chantier/autre/update/" + vars["id"],
			},
		}
		// model.AddRecent() inutile puisqu'on est redirigé vers la liste, où AddRecent() est exécuté
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
// Fabrique un Chautre à partir des valeurs d'un formulaire.
// Auxiliaire de NewChautre() et UpdateChautre()
// Ne gère pas le champ Id
// Ne gère pas le champ TVA (tiré de la config)
// Ne gère pas liens vers UGs, lieux-dits, fermiers
func chautreForm2var(r *http.Request) (*model.Chautre, error) {
	ch := &model.Chautre{}
	var err error
	if err = r.ParseForm(); err != nil {
		return ch, err
	}
	//
	ch.IdAcheteur, err = strconv.Atoi(r.PostFormValue("id-acheteur"))
	if err != nil {
		return ch, err
	}
	//
	ch.TypeVente = strings.Replace(r.PostFormValue("typevente"), "typevente-", "", -1)
	//
	ch.DateContrat, err = time.Parse("2006-01-02", r.PostFormValue("datecontrat"))
	if err != nil {
		return ch, err
	}
	//
	ch.TypeValo = strings.Replace(r.PostFormValue("typevalo"), "valorisation-", "", -1)
	//
	ch.Exploitation = strings.ReplaceAll(r.PostFormValue("exploitation"), "exploitation-", "")
	//
	ch.Essence = strings.ReplaceAll(r.PostFormValue("essence"), "essence-", "")
	//
	if r.PostFormValue("volume-contrat") == "" {
	    ch.VolumeContrat = 0 // car optionnel
	} else {
        ch.VolumeContrat, err = strconv.ParseFloat(r.PostFormValue("volume-contrat"), 32)
        if err != nil {
            return ch, err
        }
    }
	ch.VolumeContrat = tiglib.Round(ch.VolumeContrat, 2)
	//
    ch.VolumeRealise, err = strconv.ParseFloat(r.PostFormValue("volume-realise"), 32)
    if err != nil {
        return ch, err
    }
	ch.VolumeRealise = tiglib.Round(ch.VolumeRealise, 2)
	
	//
	ch.Unite = model.Valorisation2unite(ch.TypeValo)
	//
	ch.PUHT, err = strconv.ParseFloat(r.PostFormValue("puht"), 32)
	if err != nil {
		return ch, err
	}
	ch.PUHT = tiglib.Round(ch.PUHT, 2)
	//
	if r.PostFormValue("datefacture") != "" {
		ch.DateFacture, err = time.Parse("2006-01-02", r.PostFormValue("datefacture"))
		if err != nil {
			return ch, err
		}
	}
	//
	ch.NumFacture = r.PostFormValue("numfacture")
	//
	ch.Notes = r.PostFormValue("notes")
	//
	return ch, nil
}

// *********************************************************
func ShowFactureChautre(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	//
	ch, err := model.GetChautreFull(ctx.DB, id)
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
	pdf.MultiCell(100, 7, tr(StringActeurFacture(ch.Acheteur)), "1", "C", false)
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
	pdf.MultiCell(wi, he, tiglib.DateFr(ch.DateFacture), "LRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	pdf.MultiCell(wi, he, ch.NumFacture, "RB", "C", false)
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
	str = "Vente " + tr(model.LabelValorisation(ch.TypeValo)) + " " + tr(model.LabelEssence(ch.Essence))
	pdf.MultiCell(wi, he, str, "LRB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w2
	pdf.MultiCell(wi, he, strconv.FormatFloat(ch.VolumeRealise, 'f', 2, 64), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w3
	pdf.MultiCell(wi, he, tr(model.LabelUnite(ch.Unite)), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w4
	pdf.MultiCell(wi, he, strconv.FormatFloat(ch.PUHT, 'f', 2, 64), "RB", "C", false)
	x += wi
	pdf.SetXY(x, y)
	wi = w5
	prixHT := ch.VolumeRealise * ch.PUHT
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
	pdf.MultiCell(wi, he, strconv.FormatFloat(ch.TVA, 'f', 2, 64)+" %", "RB", "C", false)
	x += wi
	wi = w5
	pdf.SetXY(x, y)
	prixTVA := prixHT * ch.TVA / 100
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
