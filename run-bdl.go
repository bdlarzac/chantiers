/******************************************************************************
   Entry point of the web application

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2018-04, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"bdl.local/bdl/control"
	"bdl.local/bdl/control/ajax"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	"bdl.local/bdl/static"
	"bdl.local/bdl/view"
	"github.com/gorilla/mux"
)

// *********************************************************
func main() {

	defer func() {
		if p := recover(); p != nil {
			err := fmt.Errorf("%w", p)
			ctxt.LogError(err)
		}
	}()

	model.MustLoadEnv()
	ctxt.MustLoadConfig()
	ctxt.MustInitDB()
	ctxt.MustInitTemplates()

	r := mux.NewRouter()

	r.HandleFunc("/ajax/get/lieudits-from-code-ug/{code}", Hajax(ajax.GetLieuditsFromCodeUG))
	r.HandleFunc("/ajax/get/fermiers-from-code-ug/{code}", Hajax(ajax.GetFermiersFromCodeUG))
	r.HandleFunc("/ajax/get/ugs-from-fermier/{id}", Hajax(ajax.GetUGsFromFermier))
	r.HandleFunc("/ajax/get/ug-from-code/{code}", Hajax(ajax.GetUGFromCode))
	r.HandleFunc("/ajax/get/parcelles-from-ug/{id}", Hajax(ajax.GetParcellesFromUG))

	r.HandleFunc("/", H(control.Accueil))
	r.HandleFunc("/doc", H(control.ShowDoc))
	r.HandleFunc("/backup", H(control.BackupDB))
	r.HandleFunc("/sctl-update", H(control.UpdateSCTL))

	r.HandleFunc("/facture/vente-plaquette/{id:[0-9]+}", HPDF(control.ShowFactureVentePlaq))
	r.HandleFunc("/facture/autre/{id:[0-9]+}", HPDF(control.ShowFactureChautre))

	r.HandleFunc("/affacture/form/{id:[0-9]+}", H(control.FormAffacture))
	r.HandleFunc("/affacture/show", HPDF(control.ShowAffacture))

	r.HandleFunc("/acteur/liste", H(control.ListActeur))
	r.HandleFunc("/acteur/new", H(control.NewActeur))
	r.HandleFunc("/acteur/update/{id:[0-9]+}", H(control.UpdateActeur))
	r.HandleFunc("/acteur/delete/{id:[0-9]+}", H(control.DeleteActeur))
	r.HandleFunc("/acteur/{id:[0-9]+}", H(control.ShowActeur))

	r.HandleFunc("/fermier/liste", H(control.ListFermier))
	r.HandleFunc("/fermier/{id:[0-9]+}", H(control.ShowFermier))

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
	r.HandleFunc("/stockage/delete/{id:[0-9]+}", H(control.DeleteOrArchiveStockage))

	r.HandleFunc("/tas-vides", H(control.ShowTasVides))
	r.HandleFunc("/tas/vider/{id:[0-9]+}/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", H(control.SignalerTasVide))

	r.HandleFunc("/frais-stockage/new/{id-stockage:[0-9]+}", H(control.NewStockFrais))
	r.HandleFunc("/frais-stockage/update/{id:[0-9]+}", H(control.UpdateStockFrais))
	r.HandleFunc("/frais-stockage/delete/{id:[0-9]+}", H(control.DeleteStockFrais))

	r.HandleFunc("/humidite/liste", H(control.ListHumid))
	r.HandleFunc("/humidite/liste/{annee:[0-9]+}", H(control.ListHumid))
	r.HandleFunc("/humidite/new", H(control.NewHumid))
	r.HandleFunc("/humidite/new/tas/{id-tas:[0-9]+}", H(control.NewHumid))
	r.HandleFunc("/humidite/update/{id:[0-9]+}", H(control.UpdateHumid))
	r.HandleFunc("/humidite/delete/{id:[0-9]+}", H(control.DeleteHumid))

	r.HandleFunc("/bilans", H(control.FormBilans))

	r.HandleFunc("/ug/search", H(control.SearchUG))
	r.HandleFunc("/ug/{id:[0-9]+}", H(control.ShowUG))
	r.HandleFunc("/ug/{id:[0-9]+}/{tab}", H(control.ShowUG))
	r.HandleFunc("/commune/liste", H(control.ListCommunes))
	r.HandleFunc("/lieudit/{id:[0-9]+}", H(control.ShowLieudit))
	r.HandleFunc("/parcelle/{id:[0-9]+}", H(control.ShowParcelle))

	r.PathPrefix("/doc/").Handler(http.StripPrefix("/doc/", http.FileServer(http.Dir(filepath.Join("doc")))))
	r.HandleFunc("/dbdump/", notFound) // pour empêcher de lister le rep contenant les db dumps
	r.PathPrefix("/dbdump/").Handler(http.StripPrefix("/dbdump/", http.FileServer(http.Dir(model.SERVER_ENV.BACKUP_DIR))))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.StaticFiles))))
	r.PathPrefix("/view/").Handler(http.StripPrefix("/view/", http.FileServer(http.FS(view.ViewFiles))))

	r.NotFoundHandler = http.HandlerFunc(notFound)

	r.Use(contentTypeMiddleware)

	addr := model.SERVER_ENV.RUN_SERVER_ADDR + ":" + model.SERVER_ENV.PORT
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listen %s", addr)
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
		//
		if ctx.Page != nil {
			// ctx.Page == nil if contentTypeMiddleware was called
			ctx.Page.RunMode = model.SERVER_ENV.RUN_MODE // "dev" or "prod", available in all pages
		}
		//
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
	err := fmt.Errorf("Page inexistante :<br><code><b>%s</b></code>", r.URL)
	showErrorPage(err, ctx, w, r)
}

func showErrorPage(theErr error, ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) {
	type detailsErrorPage struct {
		URL     string
		Details string
	}
	var err error

	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "ERREUR",
		},
		Menu: "accueil",
		Details: detailsErrorPage{
			URL:     r.URL.String(),
			Details: werr.SprintHTML(theErr),
		},
		RunMode: model.SERVER_ENV.RUN_MODE,
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

/**
    Adds a "Content-Type" header to the response
**/
func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(r.URL.String())
		var contentType string
		switch ext {
		case ".htm", ".html":
			contentType = "text/html"
		case ".css":
			contentType = "text/css"
		case ".png":
			contentType = "image/png"
		case ".svg":
			contentType = "image/svg+xml"
		case ".js":
			contentType = "application/javascript"
		default:
			contentType = mime.TypeByExtension(ext)
		}
		w.Header().Add("Content-Type", contentType+";charset=utf-8")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
