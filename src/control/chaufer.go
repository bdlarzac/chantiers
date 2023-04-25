/*
*

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.

*
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/webo"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type detailsChauferForm struct {
	Chantier            *model.Chaufer
	UrlAction           string
	FermierOptions      template.HTML
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
	titrePage := "Chantiers chauffage fermier " + annee
	ctx.TemplateName = "chaufer-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: titrePage,
			JSFiles: []string{
			    "/static/lib/table-sort/table-sort.js",
				"/static/js/round.js",
				"/view/common/prix.js"},
		},
		Menu: "chantiers",
		Details: detailsChauferList{
			Chantiers:       chantiers,
			Annee:           annee,
			Annees:          annees,
			TotalParEssence: totalParEssence,
		},
		Footer: ctxt.Footer{
			JSFiles: []string{
			    "/static/lib/table-sort/table-sort.js"},
		},
	}
	// Add recent - modifie URL pour éviter des doublons :
	// Année non spécifiée dans URL = Année courante
	url := r.URL.String()
	if strings.HasSuffix(url, "/liste") {
		url += "/" + annee
	}
	err = model.AddRecent(ctx.DB, ctx.Config, &model.Recent{URL: url, Label: titrePage})
	if err != nil {
		return err
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
		chantier, idsUG, err := chantierChauferForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertChaufer(ctx.DB, chantier, idsUG)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/chauffage-fermier/liste/" + strconv.Itoa(chantier.DateChantier.Year())
		// model.AddRecent() inutile puisqu'on est redirigé vers la liste, où AddRecent() est exécuté
		return nil
	default:
		//
		// Affiche form
		//
		chantier := &model.Chaufer{}
		chantier.Fermier = &model.Fermier{}
		chantier.UGs = []*model.UG{}
		chantier.LiensParcelles = []*model.ChantierParcelle{}
		weboFermier, err := WeboFermier(ctx)
		if err != nil {
			return err
		}
		ctx.TemplateName = "chaufer-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouveau chantier chauffage fermier",
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Details: detailsChauferForm{
				Chantier:            chantier,
				FermierOptions:      webo.FmtOptions(weboFermier, "CHOOSE_FERMIER"),
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "CHOOSE_ESSENCE"),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "CHOOSE_EXPLOITATION"),
				UniteOptions:        webo.FmtOptions(WeboChauferUnite(), "CHOOSE_UNITE"),
				UrlAction:           "/chantier/chauffage-fermier/new",
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/js/round.js"},
			},
		}
		// model.AddRecent() inutile puisqu'on est redirigé vers la liste, où AddRecent() est exécuté
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
		chantier, idsUG, err := chantierChauferForm2var(r)
		if err != nil {
			return err
		}
		chantier.Id, err = strconv.Atoi(r.PostFormValue("id-chantier"))
		if err != nil {
			return err
		}
		err = model.UpdateChaufer(ctx.DB, chantier, idsUG)
		if err != nil {
			return err
		}
		ctx.Redirect = "/chantier/chauffage-fermier/liste/" + strconv.Itoa(chantier.DateChantier.Year())
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
		chantier, err := model.GetChauferFull(ctx.DB, idChantier)
		if err != nil {
			return err
		}
		weboFermier, err := WeboFermier(ctx)
		ctx.TemplateName = "chaufer-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier un chantier chauffage fermier",
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Details: detailsChauferForm{
				Chantier:            chantier,
				FermierOptions:      webo.FmtOptions(weboFermier, "fermier-"+strconv.Itoa(chantier.IdFermier)),
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "essence-"+chantier.Essence),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "exploitation-"+chantier.Exploitation),
				UniteOptions:        webo.FmtOptions(WeboChauferUnite(), "unite-"+chantier.Unite),
				UrlAction:           "/chantier/chauffage-fermier/update/" + vars["id"],
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/js/round.js"},
			},
		}
		// model.AddRecent() inutile puisqu'on est redirigé vers la liste, où AddRecent() est exécuté
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
	chantier, err := model.GetChaufer(ctx.DB, id) // chantier sert pour l'année du redirect
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
/**
    Fabrique un Chaufer à partir des valeurs d'un formulaire.
    Auxiliaire de NewChaufer() et UpdateChaufer()
    Ne gère pas le champ Id
    Pour form new, IdChantier = 0 ; pour form update, IdChantier a la bonne valeur
    Renvoie idsUG car ils ne sont pas stockés dans model.chaufer
    Mais les liens avec les parcelles sont stockés dans ch.ChantierParcelle
**/
func chantierChauferForm2var(r *http.Request) (ch *model.Chaufer, idsUG []int, err error) {
	ch = &model.Chaufer{}
	vide := []int{}
	if err = r.ParseForm(); err != nil {
		return ch, vide, err
	}
	//
	ch.Titre = r.PostFormValue("titre")
	//
	ch.IdFermier, err = strconv.Atoi(r.PostFormValue("id-fermier"))
	if err != nil {
		return ch, vide, err
	}
	//
	idsUG = form2IdsUG(r)
	//
	ch.LiensParcelles = form2LienParcelles(r)
	//
	ch.DateChantier, err = time.Parse("2006-01-02", r.PostFormValue("datechantier"))
	if err != nil {
		return ch, vide, err
	}
	//
	ch.Exploitation = strings.ReplaceAll(r.PostFormValue("exploitation"), "exploitation-", "")
	//
	ch.Essence = strings.ReplaceAll(r.PostFormValue("essence"), "essence-", "")
	//
	ch.Volume, err = strconv.ParseFloat(r.PostFormValue("volume"), 32)
	if err != nil {
		return ch, vide, err
	}
	ch.Volume = tiglib.Round(ch.Volume, 2)
	//
	ch.Unite = strings.Replace(r.PostFormValue("unite"), "unite-", "", -1)
	//
	ch.Notes = r.PostFormValue("notes")
	//
	return ch, idsUG, nil
}
