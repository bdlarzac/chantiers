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
	Chantier            *model.Plaq
	TypeChantier        string
	UrlAction           string
	EssenceOptions      template.HTML
	ExploitationOptions template.HTML
	GranuloOptions      template.HTML
	AllStockages        []*model.Stockage
	AllUGs              []*model.UG
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
			JSFiles: []string{"/view/common/plaq.js"},
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
	err = chantier.ComputeCouts(ctx.DB, ctx.Config)
	if err != nil {
		return err
	}
	ctx.TemplateName = "plaq-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title:    chantier.FullString(),
			CSSFiles: []string{"/static/css/tabstrip.css"},
			JSFiles: []string{
				"/static/js/round.js",
				"/view/common/prix.js"},
		},
		Menu: "chantiers",
		Footer: ctxt.Footer{
			JSFiles: []string{
				"/static/js/tabstrip.js",
				"/view/common/plaq.js"},
		},
		Details: detailsPlaqShow{
			Chantier:         chantier,
			PourcentagePerte: ctx.Config.PourcentagePerte,
			Tab:              tab,
		},
	}
	// Add recent - modifie URL pour éviter des doublons (même chantier mais urls différentes)
	url := r.URL.String()
	url = strings.TrimSuffix(url, "/general")
	url = strings.TrimSuffix(url, "/chantiers")
	url = strings.TrimSuffix(url, "/cout")
	err = model.AddRecent(ctx.DB, ctx.Config, &model.Recent{URL: url, Label: chantier.FullString()})
	if err != nil {
		return err
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
		allStockages, err := model.GetStockages(ctx.DB, true)
		if err != nil {
			return err
		}
		idsStockages := []int{}
		for _, stockage := range allStockages {
			if r.PostFormValue("stockage-"+strconv.Itoa(stockage.Id)) == "on" {
				idsStockages = append(idsStockages, stockage.Id)
			}
		}
		// calcul des ids UG, Lieudit et Fermier, pour transmettre à InsertPlaq()
		idsUGs, idsLieudits, idsFermiers, err := calculeIdsLiensChantier(r)
        if err != nil {
            return err
        }
		//
		id, err := model.InsertPlaq(ctx.DB, chantier, idsStockages, idsUGs, idsLieudits, idsFermiers)
		if err != nil {
			return err
		}                                                                                                  
		//
		redirect := "/chantier/plaquette/" + strconv.Itoa(id)
		err = chantier.ComputeLieudits(ctx.DB) // nécessaire pour appeler chantier.FullString()
		if err != nil {
			return err
		}
		// model.AddRecent() inutile puisqu'on est redirigé vers ShowPlaq(), où AddRecent() est exécuté
		ctx.Redirect = redirect
		return nil
	default:
		//
		// Affiche form
		//
		chantier := &model.Plaq{}
		allStockages, err := model.GetStockages(ctx.DB, true)
		if err != nil {
			return err
		}
		allUGs, err := model.GetUGsSortedByCode(ctx.DB)
		if err != nil {
			return err
		}
		ctx.TemplateName = "plaq-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouveau chantier plaquettes",
				CSSFiles: []string{
					"/static/css/form.css",
				    "/static/css/modal.css"},
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js"},
			},
			Details: detailsPlaqForm{
				Chantier:            chantier,
				TypeChantier:        "plaq",
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "CHOOSE_ESSENCE"),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "CHOOSE_EXPLOITATION"),
				GranuloOptions:      webo.FmtOptions(WeboGranulo(), "CHOOSE_GRANULO"),
				AllStockages:        allStockages,
				AllUGs:              allUGs,
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
		chantier.Id, err = strconv.Atoi(r.PostFormValue("id-chantier"))
		if err != nil {
			return err
		}
		// calcul des ids stockage, pour transmettre à model.UpdatePlaq(),
		// qui va créer ou supprimer ou ne pas changer le(s) tas
		allStockages, err := model.GetStockages(ctx.DB, true)
		if err != nil {
			return err
		}
		idsStockages := []int{}
		for _, stockage := range allStockages {
			if r.PostFormValue("stockage-"+strconv.Itoa(stockage.Id)) == "on" {
				idsStockages = append(idsStockages, stockage.Id)
			}
		}
		// calcul des ids UG, Lieudit et Fermier, pour transmettre à InsertPlaq()
		idsUGs, idsLieudits, idsFermiers, err := calculeIdsLiensChantier(r)
        if err != nil {
            return err
        }
		//
		err = model.UpdatePlaq(ctx.DB, chantier, idsStockages, idsUGs, idsLieudits, idsFermiers)
		if err != nil {
			return err
		}
		//
		redirect := "/chantier/plaquette/" + r.PostFormValue("id-chantier")
		err = chantier.ComputeLieudits(ctx.DB) // nécessaire pour appeler chantier.FullString()
		if err != nil {
			return err
		}
		// model.AddRecent() inutile puisqu'on est redirigé vers ShowPlaq(), où AddRecent() est exécuté
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
		allStockages, err := model.GetStockages(ctx.DB, true)
		if err != nil {
			return err
		}
		allUGs, err := model.GetUGsSortedByCode(ctx.DB)
		if err != nil {
			return err
		}
		ctx.TemplateName = "plaq-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier " + chantier.FullString(),
				CSSFiles: []string{
					"/static/css/form.css",
				    "/static/css/modal.css"},
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js"},
			},
			Details: detailsPlaqForm{
				Chantier:            chantier,
				TypeChantier:        "plaq",
				EssenceOptions:      webo.FmtOptions(WeboEssence(), "essence-"+chantier.Essence),
				ExploitationOptions: webo.FmtOptions(WeboExploitation(), "exploitation-"+chantier.Exploitation),
				GranuloOptions:      webo.FmtOptions(WeboGranulo(), "granulo-"+chantier.Granulo),
				AllStockages:        allStockages,
				AllUGs:              allUGs,
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
// Code commun à NewPlaq() et UpdatePlaq()
func calculeIdsLiensChantier(r *http.Request) (idsUGs, idsLieudits, idsFermiers []int, err error) {
    rien := []int{}
    var tmp []string
    var str string
    var id int
    //
    tmp = strings.Split(r.PostFormValue("ids-ugs"), ",")
    for _, str = range(tmp){
        id, err = strconv.Atoi(str)
        if err != nil {
            return rien, rien, rien, err
        }
        idsUGs = append(idsUGs, id)
    }
    //
    tmp = strings.Split(r.PostFormValue("ids-lieudits"), ",")
    for _, str = range(tmp){
        id, err = strconv.Atoi(str)
        if err != nil {
            return rien, rien, rien, err
        }
        idsLieudits = append(idsLieudits, id)
    }
    //
    tmp = strings.Split(r.PostFormValue("ids-fermiers"), ",")
    for _, str = range(tmp){
        id, err = strconv.Atoi(str)
        if err != nil {
            return rien, rien, rien, err
        }
        idsFermiers = append(idsFermiers, id)
    }
    return idsUGs, idsLieudits, idsFermiers, nil
}

// *********************************************************
// Fabrique un Plaq à partir des valeurs d'un formulaire.
// Auxiliaire de NewPlaq() et UpdatePlaq()
// Ne gère pas le champ Id
// Ne gère pas les stockages (tas) (parce que besoinde DB pour le calculer)
// Ne gère pas liens vers UGs, lieux-dits, fermiers (parce que model.Plaq ne possède pas ces champs)
func chantierPlaquetteForm2var(r *http.Request) (*model.Plaq, error) {
	ch := &model.Plaq{}
	var err error
	if err = r.ParseForm(); err != nil {
		return ch, err
	}
	//
	ch.DateDebut, err = time.Parse("2006-01-02", r.PostFormValue("date-debut"))
	if err != nil {
		return ch, err
	}
	//
	ch.DateFin, err = time.Parse("2006-01-02", r.PostFormValue("date-fin"))
	if err != nil {
		return ch, err
	}
	//
	if r.PostFormValue("surface") != "" {
		ch.Surface, err = strconv.ParseFloat(r.PostFormValue("surface"), 32)
		if err != nil {
			return ch, err
		}
		ch.Surface = tiglib.Round(ch.Surface, 2)
	}
	//
	ch.Granulo = strings.ReplaceAll(r.PostFormValue("granulo"), "granulo-", "")
	//
	ch.Exploitation = strings.ReplaceAll(r.PostFormValue("exploitation"), "exploitation-", "")
	//
	ch.Essence = strings.ReplaceAll(r.PostFormValue("essence"), "essence-", "")
	//
	if r.PostFormValue("frais-repas") != "" {
		ch.FraisRepas, err = strconv.ParseFloat(r.PostFormValue("frais-repas"), 32)
		if err != nil {
			return ch, err
		}
		ch.FraisRepas = tiglib.Round(ch.FraisRepas, 2)
	}
	//
	if r.PostFormValue("frais-reparation") != "" {
		ch.FraisReparation, err = strconv.ParseFloat(r.PostFormValue("frais-reparation"), 32)
		if err != nil {
			return ch, err
		}
		ch.FraisReparation = tiglib.Round(ch.FraisReparation, 2)
	}
	//
	return ch, nil
}
