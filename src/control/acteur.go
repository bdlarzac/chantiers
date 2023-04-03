/*

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"strconv"
	"strings"
	"net/http"
	"github.com/gorilla/mux"
)

type detailsActeurList struct {
	List  []*model.Acteur
	Count int
}

type detailsActeurForm struct {
	UrlAction string
	Acteur    *model.Acteur
}

type detailsActeurShow struct {
	Acteur    *model.Acteur
	Activites []*model.ActeurActivite
}

// *********************************************************
func ListActeur(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	list, err := model.GetSortedActeurs(ctx.DB, "nom")
	if err != nil {
		return err
	}
	for _, acteur := range list {
		// TODO bizarre de faire ça ici, devrait être dans model/
		acteur.Deletable, err = acteur.IsDeletable(ctx.DB)
		if err != nil {
			return err
		}
	}
	ctx.TemplateName = "acteur-list.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Acteurs",
		},
		Menu: "acteurs",
		Footer: ctxt.Footer{
			JSFiles: []string{},
		},
		Details: detailsActeurList{
			List:  list,
			Count: model.CountActeurs(ctx.DB),
		},
	}
	return nil
}

// *********************************************************
func ShowActeur(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id1, err := strconv.Atoi(vars["id"])
	id := int(id1)
	if err != nil {
		return err
	}
	acteur, err := model.GetActeur(ctx.DB, id)
	if err != nil {
		return err
	}
	acteur.Deletable, err = acteur.IsDeletable(ctx.DB)
	if err != nil {
		return err
	}
	activites, err := acteur.GetActivitesByDate(ctx.DB)
	if err != nil {
		return err
	}
	ctx.TemplateName = "acteur-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: acteur.String(),
			CSSFiles: []string{
				"/static/css/form.css"},
		},
		Menu: "acteurs",
		Details: detailsActeurShow{
			Acteur:    acteur,
			Activites: activites,
		},
	}
	return nil
}

// *********************************************************
// Process ou affiche form new
func NewActeur(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		acteur, err := acteurForm2var(r)
		if err != nil {
			return err
		}
		acteur.Deletable = true // nouvellement créé, pas SCTL, pas d'activité => effaçable
		id, err := model.InsertActeur(ctx.DB, acteur)
		if err != nil {
			return err
		}
		ctx.Redirect = "/acteur/" + strconv.Itoa(id)
		return nil
	default:
		//
		// Affiche form
		//
		acteur := &model.Acteur{}
		acteur.Actif = true
		ctx.TemplateName = "acteur-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Nouvel acteur",
				CSSFiles: []string{"/static/css/form.css"},
			},
			Menu: "acteurs",
			Footer: ctxt.Footer{
				JSFiles: []string{"/static/js/toogle.js"},
			},
			Details: detailsActeurForm{
				Acteur:    acteur,
				UrlAction: "/acteur/new",
			},
		}
		return nil
	}
}

// *********************************************************
// Process ou affiche form update
func UpdateActeur(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		acteur, err := acteurForm2var(r)
		if err != nil {
			return err
		}
		acteur.Id, err = strconv.Atoi(r.PostFormValue("id"))
		if err != nil {
			return err
		}
		// Actif et Deletable sont gérés lors d'un import SCTL
		// ou lors de l'effacement d'activités le concernant
		err = model.UpdateActeur(ctx.DB, acteur)
		if err != nil {
			return err
		}
		ctx.Redirect = "/acteur/" + r.PostFormValue("id")
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			return err
		}
		acteur, err := model.GetActeurFull(ctx.DB, id)
		if err != nil {
			return err
		}
		ctx.TemplateName = "acteur-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Modifier l'acteur " + acteur.String(),
				CSSFiles: []string{"/static/css/form.css"},
			},
			Menu: "acteurs",
			Footer: ctxt.Footer{
				JSFiles: []string{"/static/js/toogle.js"},
			},
			Details: detailsActeurForm{
				Acteur:    acteur,
				UrlAction: "/acteur/update/" + vars["id"],
			},
		}
		return nil
	}
}

// *********************************************************
func DeleteActeur(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	err = model.DeleteActeur(ctx.DB, id)
	if err != nil {
		return err
	}
	ctx.Redirect = "/acteur/liste"
	return nil
}

// *********************************************************
// Fabrique un Acteur à partir des valeurs d'un formulaire.
// Auxiliaire de NewActeur() et UpdateActeur()
// Ne gère pas les champs Id, Actif et Deletable
func acteurForm2var(r *http.Request) (*model.Acteur, error) {
	acteur := &model.Acteur{}
	var err error
	if err = r.ParseForm(); err != nil {
		return acteur, err
	}
	acteur.Nom = r.PostFormValue("nom")
	//
	tmp := strings.Split(r.PostFormValue("codes-roles"), ";")
	for _, str := range tmp {
		acteur.CodesRoles = append(acteur.CodesRoles, str)
	}
	//
	acteur.Prenom = r.PostFormValue("prenom")
	acteur.Adresse1 = r.PostFormValue("adresse1")
	acteur.Adresse2 = r.PostFormValue("adresse2")
	acteur.Cp = r.PostFormValue("cp")
	acteur.Ville = r.PostFormValue("ville")
	acteur.Tel = r.PostFormValue("tel")
	acteur.Mobile = r.PostFormValue("mobile")
	acteur.Email = r.PostFormValue("email")
	acteur.Bic = r.PostFormValue("bic")
	acteur.Iban = r.PostFormValue("iban")
	acteur.Siret = r.PostFormValue("siret")
	//
	acteur.Proprietaire = false
	if r.PostFormValue("proprietaire") == "on" {
		acteur.Proprietaire = true
	}
	//
	acteur.Fournisseur = false
	if r.PostFormValue("fournisseur") == "on" {
		acteur.Fournisseur = true
	}
	//
	acteur.Actif = false
	if r.PostFormValue("actif") == "on" {
		acteur.Actif = true
	}
	//
	acteur.Notes = r.PostFormValue("notes")
	return acteur, nil
}
