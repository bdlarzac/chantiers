/*
@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type detailsStockageForm struct {
	UrlAction string
	Stockage  *model.Stockage
}

type detailsStockageList struct {
	Actifs   []*model.Stockage
	Archives []*model.Stockage
}

func ListStockages(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	actifs, err := model.GetStockagesFull(ctx.DB, true)
	if err != nil {
		return werr.Wrap(err)
	}
	for _, s := range actifs {
		for _, t := range s.TasActifs {
			err = t.ComputeMesuresHumidite(ctx.DB)
			if err != nil {
				return werr.Wrap(err)
			}
			err = t.ComputeEvolutionStock(ctx.DB)
			if err != nil {
				return werr.Wrap(err)
			}
		}
	}
	archives, err := model.GetStockagesFull(ctx.DB, false)
	if err != nil {
		return werr.Wrap(err)
	}
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Stockage plaquettes",
			CSSFiles: []string{
				"/static/css/form.css",
				"/static/css/modal.css"},
			JSFiles: []string{
				"/static/js/formatNb.js",
				"/static/js/round.js"},
		},
		Menu: "accueil",
		Details: detailsStockageList{
			Actifs:   actifs,
			Archives: archives,
		},
	}
	ctx.TemplateName = "stockage-list.html"
	return nil
}

// Process ou affiche form new
func NewStockage(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		stockage, err := stockageForm2var(r)
		if err != nil {
			return werr.Wrap(err)
		}
		_, err = model.InsertStockage(ctx.DB, stockage)
		if err != nil {
			return werr.Wrap(err)
		}
		ctx.Redirect = "/stockage/liste"
		return nil
	default:
		//
		// Affiche form
		//
		stockage := &model.Stockage{}
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Créer un lieu de stockage",
				CSSFiles: []string{"/static/css/form.css"},
			},
			Menu: "accueil",
			Footer: ctxt.Footer{
				JSFiles: []string{"/static/js/toogle.js"},
			},
			Details: detailsStockageForm{
				UrlAction: "/stockage/new",
				Stockage:  stockage,
			},
		}
		ctx.TemplateName = "stockage-form.html"
		return nil
	}
}

// Process ou affiche form update
func UpdateStockage(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		if err := r.ParseForm(); err != nil {
			return werr.Wrap(err)
		}
		id, err := strconv.Atoi(r.PostFormValue("id"))
		if err != nil {
			return werr.Wrap(err)
		}
		stockage := &model.Stockage{
			Id:  id,
			Nom: r.PostFormValue("nom"),
		}
		err = model.UpdateStockage(ctx.DB, stockage)
		if err != nil {
			return werr.Wrap(err)
		}
		ctx.Redirect = "/stockage/liste"
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			return werr.Wrap(err)
		}
		stockage, err := model.GetStockage(ctx.DB, id)
		if err != nil {
			return werr.Wrap(err)
		}
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Modifier un lieu de stockage",
				CSSFiles: []string{"/static/css/form.css"},
			},
			Menu: "accueil",
			Footer: ctxt.Footer{
				JSFiles: []string{"/static/js/toogle.js"},
			},
			Details: detailsStockageForm{
				UrlAction: "/stockage/update/" + strconv.Itoa(stockage.Id),
				Stockage:  stockage,
			},
		}
		ctx.TemplateName = "stockage-form.html"
		return nil
	}
}

// Delete ou archive
func DeleteOrArchiveStockage(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return werr.Wrap(err)
	}
	stockage, err := model.GetStockage(ctx.DB, id)
	if err != nil {
		return werr.Wrap(err)
	}
	err = stockage.ComputeDeletableAndArchivable(ctx.DB)
	if err != nil {
		return werr.Wrap(err)
	}
	if stockage.Deletable {
		err = model.DeleteStockage(ctx.DB, id)
		if err != nil {
			return werr.Wrap(err)
		}
	} else if stockage.Archivable {
		stockage.Archived = true
		err = model.UpdateStockage(ctx.DB, stockage)
		if err != nil {
			return werr.Wrap(err)
		}
	}
	// Si ni Archivable ni Deletable, on ne fait rien
	ctx.Redirect = "/stockage/liste"
	return nil
}

// Fabrique un Stockage à partir des valeurs d'un formulaire.
// Auxiliaire de NewStockage() et UpdateStockage()
// Ne gère pas le champ Id
func stockageForm2var(r *http.Request) (*model.Stockage, error) {
	stockage := &model.Stockage{}
	var err error
	if err = r.ParseForm(); err != nil {
		return stockage, werr.Wrap(err)
	}
	stockage.Nom = r.PostFormValue("nom")
	return stockage, nil
}
