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

type detailsChauferForm struct {
	Chantier            *model.Chaufer
	UrlAction           string
	EssenceOptions      template.HTML
	UniteOptions        template.HTML
	ExploitationOptions template.HTML
}

type detailsChauferList struct {
	Chantiers       []*model.Chaufer
	Annee           string   // année courante
	Annees          []string // toutes les années avec chantier
	TotalParEssence map[string]float64
}

// *********************************************************
func ListChaufer(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	annee := vars["annee"]
	if annee == "" {
		// annee non spécifiée, on prend l'année courante
		annee = strconv.Itoa(time.Now().Year())
	}
	chantiers, err := model.GetChaufersOfYear(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	annees, err := model.GetChauferDifferentYears(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	totalParEssence := map[string]float64{}
	for _, essence := range model.AllEssenceCodes() {
		totalParEssence[essence] = 0
	}
	for _, ch := range chantiers {
		totalParEssence[ch.Essence] += ch.Volume
	}
	//
	ctx.TemplateName = "chaufer-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Chantiers chauffage fermier " + annee,
			JSFiles: []string{
				"/static/js/round.js",
				"/static/js/prix.js"},
		},
		Menu: "chantiers",
		Details: detailsChauferList{
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
func NewChaufer(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		chantier, err := chantierChauferForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertChaufer(ctx.DB, chantier)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/chauffage-fermier/liste/" + strconv.Itoa(chantier.DateChantier.Year())
		return nil
	default:
		//
		// Affiche form
		//
		chantier := &model.Chaufer{}
		chantier.Fermier = &model.Acteur{}
		chantier.UG = &model.UG{}
		chantier.LiensParcelles = []*model.ChauferParcelle{}
		ctx.TemplateName = "chaufer-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouveau chantier chauffage fermier",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Details: detailsChauferForm{
				Chantier:            chantier,
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "CHOOSE_ESSENCE"),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "CHOOSE_EXPLOITATION"),
				UniteOptions:        webo.FmtOptions(WeboChauferUnite(), "CHOOSE_UNITE"),
				UrlAction:           "/chantier/chauffage-fermier/new",
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
func UpdateChaufer(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		chantier, err := chantierChauferForm2var(r)
		if err != nil {
			return err
		}
		chantier.Id, err = strconv.Atoi(r.PostFormValue("id-chantier"))
		if err != nil {
			return err
		}
		err = model.UpdateChaufer(ctx.DB, chantier)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/chauffage-fermier/liste/" + strconv.Itoa(chantier.DateChantier.Year())
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
		chantier, err := model.GetChauferFull(ctx.DB, idChantier)
		if err != nil {
			return err
		}
		ctx.TemplateName = "chaufer-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier un chantier chauffage fermier",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Details: detailsChauferForm{
				Chantier:            chantier,
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "essence-"+chantier.Essence),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "exploitation-"+chantier.Exploitation),
				UniteOptions:        webo.FmtOptions(WeboChauferUnite(), "unite-"+chantier.Unite),
				UrlAction:           "/chantier/chauffage-fermier/update/" + vars["id"],
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
func DeleteChaufer(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	chantier, err := model.GetChaufer(ctx.DB, id) // on retient l'année pour le redirect
	if err != nil {
		return err
	}
	err = model.DeleteChaufer(ctx.DB, id)
	if err != nil {
		return err
	}
	ctx.Redirect = "/chantier/chauffage-fermier/liste/" + strconv.Itoa(chantier.DateChantier.Year())
	return nil
}

// *********************************************************
// Fabrique un Chaufer à partir des valeurs d'un formulaire.
// Auxiliaire de NewChaufer() et UpdateChaufer()
// Ne gère pas le champ Id
// Attention : construit des ChauferParcelle ;
//     pour form new, IdChaufer = 0 ; pour form update, IdChaufer a la bonne valeur
func chantierChauferForm2var(r *http.Request) (*model.Chaufer, error) {
	chantier := &model.Chaufer{}
	var err error
	if err = r.ParseForm(); err != nil {
		return chantier, err
	}
	//
	chantier.IdFermier, err = strconv.Atoi(r.PostFormValue("id-fermier"))
	if err != nil {
		return chantier, err
	}
	//
	chantier.IdUG, _ = strconv.Atoi(strings.Replace(r.PostFormValue("ug"), "ug-", "", -1))
	// parcelles ; ex de valeurs :
	// radio-parcelle-148:[radio-parcelle-entiere-148]
	// parcelle-surface-148:[]
	// radio-parcelle-337:[radio-parcelle-surface-337]
	// parcelle-surface-337:[3.5]
	//  => parcelle 148 entière et parcelle 337 = surface de 3.5 ha
	idChaufer, _ := strconv.Atoi(r.PostFormValue("id-chantier"))
	for k, v := range r.PostForm {
		if strings.HasPrefix(k, "radio-parcelle-") {
			lien := model.ChauferParcelle{}
			lien.IdChaufer = idChaufer
			idPString := strings.Replace(k, "radio-parcelle-", "", -1)
			idP, _ := strconv.Atoi(idPString)
			lien.IdParcelle = idP
			if v[0] == "radio-parcelle-entiere-"+idPString {
				lien.Entiere = true
			} else {
				lien.Entiere = false
				lien.Surface, _ = strconv.ParseFloat(r.PostFormValue("parcelle-surface-"+idPString), 32)
				lien.Surface = tiglib.Round(lien.Surface, 2)
			}
			chantier.LiensParcelles = append(chantier.LiensParcelles, &lien)
		}
	}
	//
	chantier.DateChantier, err = time.Parse("2006-01-02", r.PostFormValue("datechantier"))
	if err != nil {
		return chantier, err
	}
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
	chantier.Notes = r.PostFormValue("notes")
	//
	return chantier, nil
}
