package control

import (
	"net/http"
	"strconv"
	"time"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	//"fmt"
)

type detailsStockLoyerForm struct {
	UrlAction string
	Stockage  *model.Stockage
	Loyer     *model.StockLoyer
}

// *********************************************************
// Process ou affiche form new
func NewStockLoyer(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		loyer, err := stockLoyerForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertStockLoyer(ctx.DB, loyer)
		if err != nil {
			return err
		}
		ctx.Redirect = "/stockage/liste"
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idStockage, err := strconv.Atoi(vars["id-stockage"])
		if err != nil {
			return err
		}
		stockage, err := model.GetStockage(ctx.DB, idStockage)
		loyer := &model.StockLoyer{}
		if err != nil {
			return err
		}
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Créer un loyer",
				CSSFiles: []string{"/static/css/form.css"},
			},
			Menu: "accueil",
			Details: detailsStockLoyerForm{
				UrlAction: "/loyer-stockage/new/" + vars["id-stockage"],
				Stockage:  stockage,
				Loyer:     loyer,
			},
		}
		ctx.TemplateName = "stockloyer-form.html"
		return nil
	}
	return nil
}

// *********************************************************
// Process ou affiche form update
func UpdateStockLoyer(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		loyer, err := stockLoyerForm2var(r)
		loyer.Id, err = strconv.Atoi(r.PostFormValue("id-loyer"))
		if err != nil {
			return err
		}
		err = model.UpdateStockLoyer(ctx.DB, loyer)
		if err != nil {
			return err
		}
		ctx.Redirect = "/stockage/liste"
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idLoyer, err := strconv.Atoi(vars["id"])
		if err != nil {
			return err
		}
		loyer, err := model.GetStockLoyer(ctx.DB, idLoyer)
		if err != nil {
			return err
		}
		// Le stockage sert à afficher son nom dans le form
		stockage, err := model.GetStockage(ctx.DB, loyer.IdStockage)
		if err != nil {
			return err
		}
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Modifier un loyer",
				CSSFiles: []string{"/static/css/form.css"},
			},
			Menu: "accueil",
			Details: detailsStockLoyerForm{
				UrlAction: "/loyer-stockage/update/" + strconv.Itoa(stockage.Id),
				Stockage:  stockage,
				Loyer:     loyer,
			},
		}
		ctx.TemplateName = "stockloyer-form.html"
		return nil
	}
}

// *********************************************************
func DeleteStockLoyer(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idLoyer, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	err = model.DeleteStockLoyer(ctx.DB, idLoyer)
	if err != nil {
		return err
	}
	ctx.Redirect = "/stockage/liste"
	return nil
}

// *********************************************************
// Fabrique un StockLoyer à partir des valeurs d'un formulaire.
// Auxiliaire de NewStockLoyer() et UpdateStockLoyer()
// Ne gère pas le champ Id
func stockLoyerForm2var(r *http.Request) (*model.StockLoyer, error) {
	loyer := &model.StockLoyer{}
	var err error
	if err = r.ParseForm(); err != nil {
		return loyer, err
	}
	loyer.IdStockage, err = strconv.Atoi(r.PostFormValue("id-stockage"))
	if err != nil {
		return loyer, err
	}
	loyer.Montant, err = strconv.ParseFloat(r.PostFormValue("montant"), 32)
	if err != nil {
		return loyer, err
	}
	loyer.DateDebut, err = time.Parse("2006-01-02", r.PostFormValue("date-debut"))
	if err != nil {
		return loyer, err
	}
	loyer.DateFin, err = time.Parse("2006-01-02", r.PostFormValue("date-fin"))
	if err != nil {
		return loyer, err
	}
	return loyer, nil
}
