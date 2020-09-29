package control

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bdl.local/bdl/generic/wilk/webo"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
)

type detailsStockFraisForm struct {
	UrlAction string
	TypeFraisOptions     template.HTML
	Stockage  *model.Stockage
	Frais     *model.StockFrais
}

// *********************************************************
// Process ou affiche form new
func NewStockFrais(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		frais, err := stockFraisForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertStockFrais(ctx.DB, frais)
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
		frais := &model.StockFrais{}
		if err != nil {
			return err
		}
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Créer un frais",
				CSSFiles: []string{"/static/css/form.css"},
			},
			Menu: "accueil",
			Details: detailsStockFraisForm{
				UrlAction: "/frais-stockage/new/" + vars["id-stockage"],
				TypeFraisOptions:      webo.FmtOptions(WeboStockFrais(), "CHOOSE_TYPEFRAIS"),
				Stockage:  stockage,
				Frais:     frais,
			},
		}
		ctx.TemplateName = "stockfrais-form.html"
		return nil
	}
	return nil
}

// *********************************************************
// Process ou affiche form update
func UpdateStockFrais(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		frais, err := stockFraisForm2var(r)
		frais.Id, err = strconv.Atoi(r.PostFormValue("id-frais"))
		if err != nil {
			return err
		}
		err = model.UpdateStockFrais(ctx.DB, frais)
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
		idFrais, err := strconv.Atoi(vars["id"])
		if err != nil {
			return err
		}
		frais, err := model.GetStockFrais(ctx.DB, idFrais)
		if err != nil {
			return err
		}
		// Le stockage sert à afficher son nom dans le form
		stockage, err := model.GetStockage(ctx.DB, frais.IdStockage)
		if err != nil {
			return err
		}
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title:    "Modifier un frais",
				CSSFiles: []string{"/static/css/form.css"},
			},
			Menu: "accueil",
			Details: detailsStockFraisForm{
				UrlAction: "/frais-stockage/update/" + strconv.Itoa(stockage.Id),
				TypeFraisOptions:      webo.FmtOptions(WeboStockFrais(), "stockfrais-"+frais.TypeFrais),
				Stockage:  stockage,
				Frais:     frais,
			},
		}
		ctx.TemplateName = "stockfrais-form.html"
		return nil
	}
}

// *********************************************************
func DeleteStockFrais(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idFrais, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	err = model.DeleteStockFrais(ctx.DB, idFrais)
	if err != nil {
		return err
	}
	ctx.Redirect = "/stockage/liste"
	return nil
}

// *********************************************************
// Fabrique un StockFrais à partir des valeurs d'un formulaire.
// Auxiliaire de NewStockFrais() et UpdateStockFrais()
// Ne gère pas le champ Id
func stockFraisForm2var(r *http.Request) (*model.StockFrais, error) {
	frais := &model.StockFrais{}
	var err error
	if err = r.ParseForm(); err != nil {
		return frais, err
	}
	frais.IdStockage, err = strconv.Atoi(r.PostFormValue("id-stockage"))
	if err != nil {
		return frais, err
	}
	//
	frais.TypeFrais = strings.Replace(r.PostFormValue("typefrais"), "stockfrais-", "", -1)
	//
	frais.Montant, err = strconv.ParseFloat(r.PostFormValue("montant"), 32)
	if err != nil {
		return frais, err
	}
	frais.DateDebut, err = time.Parse("2006-01-02", r.PostFormValue("date-debut"))
	if err != nil {
		return frais, err
	}
	frais.DateFin, err = time.Parse("2006-01-02", r.PostFormValue("date-fin"))
	if err != nil {
		return frais, err
	}
	return frais, nil
}
