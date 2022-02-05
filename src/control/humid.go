/**
    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
**/
package control

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/webo"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	//"fmt"
)

type detailsHumidForm struct {
	Humid        *model.Humid
	TasOptions   template.HTML
	ListeActeurs map[int]string
	UrlAction    string
}

type detailsHumidList struct {
	Humids []*model.Humid
	Annee  string   // année courante
	Annees []string // toutes les années avec chantier
}

// *********************************************************
func ListHumid(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	annee := vars["annee"]
	if annee == "" {
		// annee non spécifiée, on prend l'année courante
		annee = strconv.Itoa(time.Now().Year())
	}
	humids, err := model.GetHumidsOfYear(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	annees, err := model.GetHumidDifferentYears(ctx.DB, annee)
	if err != nil {
		return err
	}
	//
	ctx.TemplateName = "humid-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Mesures d'humidité " + annee,
		},
		Menu: "accueil",
		Details: detailsHumidList{
			Humids: humids,
			Annee:  annee,
			Annees: annees,
		},
	}
	return nil
}

// *********************************************************
// Process ou affiche formulaire new mesure d'humidité
func NewHumid(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		humid, err := humidForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertHumid(ctx.DB, humid)
		if err != nil {
			return err
		}
		ctx.Redirect = "/humidite/liste/" + strconv.Itoa(humid.DateMesure.Year())
		return nil
	default:
		//
		// Affiche form
		//
		humid := &model.Humid{}
		weboTas, err := WeboTas(ctx)
		if err != nil {
			return err
		}
		// Ici, cas particulier, si on arrive par la route "/humidite/new/tas/{id-tas:[0-9]+}"
		// alors le tas est pré-selectionné
		optionTas := "CHOOSE_TAS"
		vars := mux.Vars(r)
		if vars["id-tas"] != "" {
			optionTas = "tas-" + vars["id-tas"]
		}
		listeActeurs, err := model.GetListeActeurs(ctx.DB)
		if err != nil {
			return err
		}
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Créer une mesure d'humidité",
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Menu: "accueil",
			Footer: ctxt.Footer{
				JSFiles: []string{},
			},
			Details: detailsHumidForm{
				Humid:        humid,
				TasOptions:   webo.FmtOptions(weboTas, optionTas),
				ListeActeurs: listeActeurs,
				UrlAction:    "/humidite/new",
			},
		}
		ctx.TemplateName = "humid-form.html"
		return nil
	}
	return nil
}

// *********************************************************
// Process ou affiche formulaire update mesure d'humidité
func UpdateHumid(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		humid, err := humidForm2var(r)
		if err != nil {
			return err
		}
		idMesure, err := strconv.Atoi(r.PostFormValue("id-mesure"))
		if err != nil {
			return err
		}
		humid.Id = idMesure
		err = model.UpdateHumid(ctx.DB, humid)
		if err != nil {
			return err
		}
		ctx.Redirect = "/humidite/liste/" + strconv.Itoa(humid.DateMesure.Year())
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idMesure, err := strconv.Atoi(vars["id"])
		if err != nil {
			return err
		}
		humid, err := model.GetHumidFull(ctx.DB, idMesure)
		if err != nil {
			return err
		}
		humid.Tas, err = model.GetTas(ctx.DB, humid.IdTas)
		if err != nil {
			return err
		}
		for _, idMesureur := range humid.IdsMesureurs {
			mesureur, err := model.GetActeur(ctx.DB, idMesureur)
			if err != nil {
				return err
			}
			humid.Mesureurs = append(humid.Mesureurs, mesureur)
		}
		weboTas, err := WeboTas(ctx)
		if err != nil {
			return err
		}
		listeActeurs, err := model.GetListeActeurs(ctx.DB)
		if err != nil {
			return err
		}
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier une mesure d'humidité",
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Menu: "accueil",
			Footer: ctxt.Footer{
				JSFiles: []string{},
			},
			Details: detailsHumidForm{
				Humid:        humid,
				TasOptions:   webo.FmtOptions(weboTas, "tas-"+strconv.Itoa(humid.IdTas)),
				ListeActeurs: listeActeurs,
				UrlAction:    "/humidite/update/" + vars["id"],
			},
		}
		ctx.TemplateName = "humid-form.html"
		return nil
	}
}

// *********************************************************
func DeleteHumid(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	err = model.DeleteHumid(ctx.DB, id)
	if err != nil {
		return err
	}
	ctx.Redirect = "/humidite/liste"
	return nil
}

// *********************************************************
// Fabrique un Humid à partir des valeurs d'un formulaire.
// Auxiliaire de NewHumid() et UpdateHumid()
// Ne gère pas le champ Id
// Gère les champs IdsMesureurs et Mesureurs
func humidForm2var(r *http.Request) (*model.Humid, error) {
	humid := &model.Humid{}
	var err error
	if err = r.ParseForm(); err != nil {
		return humid, err
	}
	humid.IdTas, _ = strconv.Atoi(strings.TrimLeft(r.PostFormValue("tas"), "tas-"))
	humid.Valeur, _ = strconv.ParseFloat(r.PostFormValue("valeur-mesure"), 32)
	humid.DateMesure, _ = time.Parse("2006-01-02", r.PostFormValue("date-mesure"))
	ids := strings.Split(r.PostFormValue("ids-mesureurs"), "-")
	for _, id := range ids {
		idMesureur, _ := strconv.Atoi(id)
		humid.IdsMesureurs = append(humid.IdsMesureurs, idMesureur)
	}
	humid.Notes = r.PostFormValue("notes")
	return humid, nil
}
