/******************************************************************************
   Entry point of the web application

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2018-04, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"
	"html/template"

	"bdl.local/bdl/control"
	"bdl.local/bdl/control/ajax"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/gorilla/mux"
)

// *********************************************************
func main() {

    defer func(){
        if p := recover(); p != nil{
            err := fmt.Errorf("%w", p)
            ctxt.LogError(err)
        }
    }()

	r := mux.NewRouter()

	// ajax acteurs
	r.HandleFunc("/ajax/autocomplete/acteur/{str}/{type-acteur}", Hajax(ajax.AutocompleteActeur))
	r.HandleFunc("/ajax/autocomplete/acteur/{str}", Hajax(ajax.AutocompleteActeur))
	r.HandleFunc("/ajax/check/acteur/{str}", Hajax(ajax.CheckNomActeur))
	r.HandleFunc("/ajax/check/acteur/", Hajax(ajax.CheckNomActeur))
	r.HandleFunc("/ajax/get/fermiers-from-lieudit/{id}", Hajax(ajax.GetFermiersFromLieudit))
	r.HandleFunc("/ajax/get/ugs-from-fermier/{id}", Hajax(ajax.GetUGsFromFermier))
	// ajax lieux-dits
	r.HandleFunc("/ajax/autocomplete/lieudit/{str}", Hajax(ajax.AutocompleteLieudit))
	r.HandleFunc("/ajax/check/lieudit/{str}", Hajax(ajax.CheckNomLieudit))
	r.HandleFunc("/ajax/get/ugs-from-lieudit/{id}", Hajax(ajax.GetUGsFromLieudit))
	r.HandleFunc("/ajax/get/parcelles-from-lieudit/{id}", Hajax(ajax.GetParcellesFromLieudit))
	// ajax ug
	r.HandleFunc("/ajax/get/parcelles-from-ug/{id}", Hajax(ajax.GetParcellesFromUG))

	r.HandleFunc("/", H(control.Accueil))
	r.HandleFunc("/admin/maj-foncier", H(control.MajFoncier))
	r.HandleFunc("/admin/maj-psg", H(control.MajPSG))

	r.HandleFunc("/facture/vente-plaquette/{id:[0-9]+}", HPDF(control.ShowFactureVentePlaq))
	r.HandleFunc("/facture/bois-sur-pied/{id:[0-9]+}", HPDF(control.ShowFactureBSPied))
	r.HandleFunc("/facture/autre/{id:[0-9]+}", HPDF(control.ShowFactureChautre))

	r.HandleFunc("/acteur/liste", H(control.ListActeur))
	r.HandleFunc("/acteur/new", H(control.NewActeur))
	r.HandleFunc("/acteur/update/{id:[0-9]+}", H(control.UpdateActeur))
	r.HandleFunc("/acteur/{id:[0-9]+}", H(control.ShowActeur))

	r.HandleFunc("/affacture/form/{id:[0-9]+}", H(control.FormAffacture))
	r.HandleFunc("/affacture/show", HPDF(control.ShowAffacture))

	r.HandleFunc("/chantier/autre/liste", H(control.ListChautre))
	r.HandleFunc("/chantier/autre/liste/{annee:[0-9]+}", H(control.ListChautre))
	r.HandleFunc("/chantier/autre/new", H(control.NewChautre))
	r.HandleFunc("/chantier/autre/update/{id:[0-9]+}", H(control.UpdateChautre))
	r.HandleFunc("/chantier/autre/delete/{id:[0-9]+}", H(control.DeleteChautre))

	r.HandleFunc("/chantier/chauffage-fermier/liste/{annee:[0-9]+}", H(control.ListChaufer))
	r.HandleFunc("/chantier/chauffage-fermier/liste", H(control.ListChaufer))
	r.HandleFunc("/chantier/chauffage-fermier/new", H(control.NewChaufer))
	r.HandleFunc("/chantier/chauffage-fermier/update/{id:[0-9]+}", H(control.UpdateChaufer))
	r.HandleFunc("/chantier/chauffage-fermier/delete/{id:[0-9]+}", H(control.DeleteChaufer))

	r.HandleFunc("/chantier/bois-sur-pied/liste/{annee:[0-9]+}", H(control.ListBspied))
	r.HandleFunc("/chantier/bois-sur-pied/liste", H(control.ListBspied))
	r.HandleFunc("/chantier/bois-sur-pied/new", H(control.NewBSPied))
	r.HandleFunc("/chantier/bois-sur-pied/update/{id:[0-9]+}", H(control.UpdateBSPied))
	r.HandleFunc("/chantier/bois-sur-pied/delete/{id:[0-9]+}", H(control.DeleteBSPied))

	r.HandleFunc("/chantier/plaquette/liste", H(control.ListPlaq))
	r.HandleFunc("/chantier/plaquette/liste/{annee:[0-9]+}", H(control.ListPlaq))
	r.HandleFunc("/chantier/plaquette/new", H(control.NewPlaq))
	r.HandleFunc("/chantier/plaquette/update/{id:[0-9]+}", H(control.UpdatePlaq))
	r.HandleFunc("/chantier/plaquette/{id:[0-9]+}", H(control.ShowPlaq))
	r.HandleFunc("/chantier/plaquette/{id:[0-9]+}/{tab}", H(control.ShowPlaq))
	r.HandleFunc("/chantier/plaquette/delete/{id:[0-9]+}", H(control.DeletePlaq))

	r.HandleFunc("/chantier/plaquette/{id-chantier:[0-9]+}/op/new", H(control.NewPlaqOp))
	r.HandleFunc("/chantier/plaquette/{id-chantier:[0-9]+}/op/update/{id-op:[0-9]+}", H(control.UpdatePlaqOp))
	r.HandleFunc("/chantier/plaquette/{id-chantier:[0-9]+}/op/delete/{id-op:[0-9]+}", H(control.DeletePlaqOp))

	r.HandleFunc("/chantier/plaquette/{id-chantier:[0-9]+}/transport/new", H(control.NewPlaqTrans))
	r.HandleFunc("/chantier/plaquette/{id-chantier:[0-9]+}/transport/update/{id-pt:[0-9]+}", H(control.UpdatePlaqTrans))
	r.HandleFunc("/chantier/plaquette/{id-chantier:[0-9]+}/transport/delete/{id-pt:[0-9]+}", H(control.DeletePlaqTrans))

	r.HandleFunc("/chantier/plaquette/{id-chantier:[0-9]+}/range/new", H(control.NewPlaqRange))
	r.HandleFunc("/chantier/plaquette/{id-chantier:[0-9]+}/range/update/{id-pr:[0-9]+}", H(control.UpdatePlaqRange))
	r.HandleFunc("/chantier/plaquette/{id-chantier:[0-9]+}/range/delete/{id-pr:[0-9]+}", H(control.DeletePlaqRange))

	r.HandleFunc("/vente/liste", H(control.ListVentePlaq))
	r.HandleFunc("/vente/liste/{annee:[0-9]+}", H(control.ListVentePlaq))
	r.HandleFunc("/vente/{id-vente:[0-9]+}", H(control.ShowVentePlaq))
	r.HandleFunc("/vente/new", H(control.NewVentePlaq))
	r.HandleFunc("/vente/update/{id-vente:[0-9]+}", H(control.UpdateVentePlaq))
	r.HandleFunc("/vente/delete/{id-vente:[0-9]+}", H(control.DeleteVentePlaq))

	r.HandleFunc("/vente/{id-vente:[0-9]+}/livraison/new", H(control.NewVenteLivre))
	r.HandleFunc("/vente/{id-vente:[0-9]+}/livraison/update/{id-livraison:[0-9]+}", H(control.UpdateVenteLivre))
	r.HandleFunc("/vente/{id-vente:[0-9]+}/livraison/delete/{id-livraison:[0-9]+}", H(control.DeleteVenteLivre))

	r.HandleFunc("/vente/{id-vente:[0-9]+}/livraison/{id-livraison:[0-9]+}/chargement/new", H(control.NewVenteCharge))
	r.HandleFunc("/vente/{id-vente:[0-9]+}/livraison/{id-livraison:[0-9]+}/chargement/update/{id-chargement:[0-9]+}", H(control.UpdateVenteCharge))
	r.HandleFunc("/vente/{id-vente:[0-9]+}/livraison/{id-livraison:[0-9]+}/chargement/delete/{id-chargement:[0-9]+}", H(control.DeleteVenteCharge))

	r.HandleFunc("/stockage/liste", H(control.ListStockages))
	r.HandleFunc("/stockage/new", H(control.NewStockage))
	r.HandleFunc("/stockage/update/{id:[0-9]+}", H(control.UpdateStockage))
	r.HandleFunc("/stockage/delete/{id:[0-9]+}", H(control.DeleteStockage))

	r.HandleFunc("/tas/empty/{id:[0-9]+}", H(control.SignalerTasVide))

	r.HandleFunc("/loyer-stockage/new/{id-stockage:[0-9]+}", H(control.NewStockLoyer))
	r.HandleFunc("/loyer-stockage/update/{id:[0-9]+}", H(control.UpdateStockLoyer))
	r.HandleFunc("/loyer-stockage/delete/{id:[0-9]+}", H(control.DeleteStockLoyer))

	r.HandleFunc("/humidite/liste", H(control.ListHumid))
	r.HandleFunc("/humidite/liste/{annee:[0-9]+}", H(control.ListHumid))
	r.HandleFunc("/humidite/new", H(control.NewHumid))
	r.HandleFunc("/humidite/update/{id:[0-9]+}", H(control.UpdateHumid))
	r.HandleFunc("/humidite/delete/{id:[0-9]+}", H(control.DeleteHumid))

	// *** geo ***
	r.HandleFunc("/unite-gestion/{id:[0-9]+}", H(control.ShowUG))
	r.HandleFunc("/geo/commune/liste", H(control.ListCommunes))
	r.HandleFunc("/geo/lieudit/{id:[0-9]+}", H(control.ShowLieudit))
	r.HandleFunc("/geo/parcelle/{id:[0-9]+}", H(control.ShowParcelle))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.NotFoundHandler = http.HandlerFunc(notFound)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// *********************************************************
// H = Handler
// Returns a function with same signature as http.Handler.ServeHTTP() usable by r.HandleFunc()
// Adapter between ServeHTTP() and controller function
// @param  h Controller function
func H(h func(*ctxt.Context, http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := ctxt.NewContext()
		//
		err = h(ctx, w, r) // Call controller h ; fills ctx.TemplateName
		if err != nil {
			showErrorPage(err, ctx, w, r)
			return
		}
		if ctx.Redirect != "" {
			http.Redirect(w, r, ctx.Redirect, http.StatusSeeOther)
			return
		}
		//
		tmpl := ctx.Template
		//
		err = tmpl.ExecuteTemplate(w, "header.html", ctx.Page)
		if err != nil {
			ctxt.LogError(err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "menu.html", ctx.Page)
		if err != nil {
			ctxt.LogError(err)
			return
		}
		// Execute template computed by h
		err = tmpl.ExecuteTemplate(w, ctx.TemplateName, ctx.Page)
		if err != nil {
			ctxt.LogError(err)
			return
		}
		err = tmpl.ExecuteTemplate(w, "footer.html", ctx.Page)
		if err != nil {
			ctxt.LogError(err)
			return
		}
	}
}

// *********************************************************
// Hajax = Handler ajax
// Same as H, but for ajax (does not execute templates)
// @param  h Controller function
func Hajax(h func(*ctxt.Context, http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := ctxt.NewContext()
		err = h(ctx, w, r) // Calls controller h
		if err != nil {
			ctxt.LogError(err)
		}
	}
}

// *********************************************************
// HPDF = Handler PDF
// Same as H, but for pdf (does not execute templates)
// @param  h Controller function
func HPDF(h func(*ctxt.Context, http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := ctxt.NewContext()
		err = h(ctx, w, r) // Calls controller h
		if err != nil {
			ctxt.LogError(err)
		}
	}
}

// *********************** Gestion d'erreur **********************************
// A mettre ailleurs, mais où ?

func notFound(w http.ResponseWriter, r *http.Request) {
	ctx := ctxt.NewContext()
	err := fmt.Errorf("Page inexistante<br><code>%s</code>", r.URL)
	showErrorPage(err, ctx, w, r)
}

func showErrorPage(theErr error, ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) {
	type detailsErrorPage struct {
		URL   string
		Details template.HTML
	}
	var err error
	
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "ERREUR",
		},
		Menu: "accueil",
		Details: detailsErrorPage{
			URL:   r.URL.String(),
			Details: template.HTML(strings.Replace(template.HTMLEscapeString(werr.Sprint(theErr)), "\n", "<br>", -1)),
		},
	}
	tmpl := ctx.Template
	err = tmpl.ExecuteTemplate(w, "header.html", ctx.Page)
	if err != nil {
		ctxt.LogError(err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "menu.html", ctx.Page)
	if err != nil {
		ctxt.LogError(err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "error.html", ctx.Page)
	if err != nil {
		ctxt.LogError(err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "footer.html", ctx.Page)
	if err != nil {
		ctxt.LogError(err)
		return
	}
}
