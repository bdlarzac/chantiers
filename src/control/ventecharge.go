package control

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/webo"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	//"fmt"
)

type detailsVenteChargeForm struct {
	VenteCharge  *model.VenteCharge
	TasOptions   template.HTML
	GlTVAOptions template.HTML
	MoTVAOptions template.HTML
	OuTVAOptions template.HTML
	UrlAction    string
}

// *********************************************************
// Process ou affiche form new
func NewVenteCharge(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		vl, err := venteChargeForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertVenteCharge(ctx.DB, vl)
		if err != nil {
			return err
		}
		ctx.Redirect = "/vente/" + r.PostFormValue("id-vente")
		return nil
	default:
		//
		// Affiche form
		//
		var err error
		vars := mux.Vars(r)
		vc := &model.VenteCharge{}
		vc.IdVente, _ = strconv.Atoi(vars["id-vente"])
		vc.IdLivraison, _ = strconv.Atoi(vars["id-livraison"])
		vc.Livraison, err = model.GetVenteLivreFull(ctx.DB, vc.IdLivraison)
		if err != nil {
			return err
		}
		vc.Livraison.Vente, err = model.GetVentePlaqFull(ctx.DB, vc.IdVente)
		if err != nil {
			return err
		}
		vc.Chargeur = &model.Acteur{}
		weboTas, err := WeboTas(ctx)
		if err != nil {
			return err
		}
		ctx.TemplateName = "ventecharge-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouveau chargement plaquettes",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Menu: "ventes",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/autocomplete/autocomplete.js"},
			},
			Details: detailsVenteChargeForm{
				VenteCharge:  vc,
				TasOptions:   webo.FmtOptions(weboTas, "CHOOSE_TAS"),
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL"), "CHOOSE_TVA_GL"),
				OuTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_OU"), "CHOOSE_TVA_OU"),
				MoTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_MO"), "CHOOSE_TVA_MO"),
				UrlAction:    "/vente/" + vars["id-vente"] + "/livraison/" + vars["id-livraison"] + "/chargement/new",
			},
		}
		return nil
	}
}

// *********************************************************
// Process ou affiche form update
func UpdateVenteCharge(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		vars := mux.Vars(r)
		vc, err := venteChargeForm2var(r)
		if err != nil {
			return err
		}
		err = model.UpdateVenteCharge(ctx.DB, vc)
		if err != nil {
			return err
		}
		ctx.Redirect = "/vente/" + vars["id-vente"]
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idChargement, err := strconv.Atoi(vars["id-chargement"])
		if err != nil {
			return err
		}
		vc, err := model.GetVenteChargeFull(ctx.DB, idChargement) // full pour avoir le nom du chargeur
		if err != nil {
			return err
		}
		vc.Livraison, err = model.GetVenteLivreFull(ctx.DB, vc.IdLivraison) // full pour avoir le nom de la livraison
		if err != nil {
			return err
		}
		vc.Livraison.Vente, err = model.GetVentePlaqFull(ctx.DB, vc.Livraison.IdVente) // full pour avoir le nom de la vente
		if err != nil {
			return err
		}
		weboTas, err := WeboTas(ctx)
		if err != nil {
			return err
		}
		ctx.TemplateName = "ventecharge-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier le chargement : " + vc.String(),
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Menu: "ventes",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/autocomplete/autocomplete.js"},
			},
			Details: detailsVenteChargeForm{
				VenteCharge:  vc,
				TasOptions:   webo.FmtOptions(weboTas, strconv.Itoa(vc.IdTas)),
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL"), strconv.FormatFloat(vc.GlTVA, 'f', 1, 64)),
				OuTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_OU"), strconv.FormatFloat(vc.OuTVA, 'f', 1, 64)),
				MoTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_MO"), strconv.FormatFloat(vc.MoTVA, 'f', 1, 64)),
				UrlAction:    "/vente/" + vars["id-vente"] + "/livraison/" + vars["id-livraison"] + "/chargement/update/" + vars["id-chargement"],
			},
		}
		return nil
	}
}

// *********************************************************
func DeleteVenteCharge(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id-chargement"])
	if err != nil {
		return err
	}
	err = model.DeleteVenteCharge(ctx.DB, id)
	if err != nil {
		return err
	}
	ctx.Redirect = "/vente/" + vars["id-vente"]
	return nil
}

// *********************************************************
// Fabrique une VenteCharge à partir des valeurs d'un formulaire.
// Auxiliaire de NewVenteCharge() et UpdateVenteCharge()
// Ne gère pas le champ Id
func venteChargeForm2var(r *http.Request) (*model.VenteCharge, error) {
	vc := &model.VenteCharge{}
	var err error
	if err = r.ParseForm(); err != nil {
		return vc, err
	}
	//
	if r.PostFormValue("type-cout") == "cout-global" {
		vc.TypeCout = "G"
	} else {
		vc.TypeCout = "D"
	}
	//
	vc.IdLivraison, err = strconv.Atoi(r.PostFormValue("id-livraison"))
	if err != nil {
		return vc, err
	}
	//
	vc.IdChargeur, err = strconv.Atoi(r.PostFormValue("id-chargeur"))
	if err != nil {
		return vc, err
	}
	vc.IdVente, err = strconv.Atoi(r.PostFormValue("id-vente"))
	if err != nil {
		return vc, err
	}
	//
	vc.IdTas, err = strconv.Atoi(r.PostFormValue("tas"))
	if err != nil {
		return vc, err
	}
	//
	vc.Qte, err = strconv.ParseFloat(r.PostFormValue("qte"), 32)
	if err != nil {
		return vc, err
	}
	vc.Qte = tiglib.Round(vc.Qte, 2)
	//
	vc.DateCharge, err = time.Parse("2006-01-02", r.PostFormValue("datecharge"))
	if err != nil {
		return vc, err
	}
	if vc.TypeCout == "G" {
		//
		// coût global
		//
		vc.GlPrix, err = strconv.ParseFloat(r.PostFormValue("glprix"), 32)
		if err != nil {
			return vc, err
		}
		vc.GlPrix = tiglib.Round(vc.GlPrix, 2)
		//
		vc.GlTVA, err = strconv.ParseFloat(r.PostFormValue("gltva"), 32)
		if err != nil {
			return vc, err
		}
		vc.GlTVA = tiglib.Round(vc.GlTVA, 2)
		//
		if r.PostFormValue("gldatepay") != "" {
			vc.GlDatePay, err = time.Parse("2006-01-02", r.PostFormValue("gldatepay"))
			if err != nil {
				return vc, err
			}
		}
	} else {
		//
		// coût détaillé, outil
		//
		vc.OuPrix, err = strconv.ParseFloat(r.PostFormValue("ouprix"), 32)
		if err != nil {
			return vc, err
		}
		vc.OuPrix = tiglib.Round(vc.OuPrix, 2)
		//
		vc.OuTVA, err = strconv.ParseFloat(r.PostFormValue("outva"), 32)
		if err != nil {
			return vc, err
		}
		vc.OuTVA = tiglib.Round(vc.OuTVA, 2)
		//
		if r.PostFormValue("oudatepay") != "" {
			vc.OuDatePay, err = time.Parse("2006-01-02", r.PostFormValue("oudatepay"))
			if err != nil {
				return vc, err
			}
		}
		//
		// coût détaillé, main d'oeuvre
		//
		vc.MoNHeure, err = strconv.ParseFloat(r.PostFormValue("monheure"), 32)
		if err != nil {
			return vc, err
		}
		vc.MoNHeure = tiglib.Round(vc.MoNHeure, 2)
		//
		vc.MoPrixH, err = strconv.ParseFloat(r.PostFormValue("moprixh"), 32)
		if err != nil {
			return vc, err
		}
		vc.MoPrixH = tiglib.Round(vc.MoPrixH, 2)
		//
		vc.MoTVA, err = strconv.ParseFloat(r.PostFormValue("motva"), 32)
		if err != nil {
			return vc, err
		}
		vc.MoTVA = tiglib.Round(vc.MoTVA, 2)
		//
		if r.PostFormValue("modatepay") != "" {
			vc.MoDatePay, err = time.Parse("2006-01-02", r.PostFormValue("modatepay"))
			if err != nil {
				return vc, err
			}
		}
	}
	//
	vc.Notes = r.PostFormValue("notes")
	//
	return vc, nil
}
