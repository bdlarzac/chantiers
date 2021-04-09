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
)

type detailsVenteChargeForm struct {
	VenteCharge  *model.VenteCharge
	TasOptions   template.HTML
	GlTVAOptions template.HTML
	MoTVAOptions template.HTML
	OuTVAOptions template.HTML
    ListeActeurs map[int]string
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
		vc, err := venteChargeForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertVenteCharge(ctx.DB, vc)
		if err != nil {
			return err
		}
		// Mise à jour du stock du tas
		err = vc.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		err = vc.Tas.ModifierStock(ctx.DB, -vc.Qte) // Retire plaquettes au tas
		if err != nil {
			return err
		}
		//
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
		vc.Conducteur = &model.Acteur{}
		vc.Proprioutil = &model.Acteur{}
		weboTas, err := WeboTas(ctx)
		if err != nil {
			return err
		}
		listeActeurs, err := model.GetListeActeurs(ctx.DB)
		if err != nil {
			return err
		}
		ctx.TemplateName = "ventecharge-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouveau chargement plaquettes",
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Menu: "ventes",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js"},
			},
			Details: detailsVenteChargeForm{
				VenteCharge:  vc,
				TasOptions:   webo.FmtOptions(weboTas, "CHOOSE_TAS"),
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL", "gl-"), "CHOOSE_TVA_GL"),
				MoTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_MO", "mo-"), "CHOOSE_TVA_MO"),
				OuTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_OU", "ou-"), "CHOOSE_TVA_OU"),
			    ListeActeurs:  listeActeurs,
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
		vc.Id, err = strconv.Atoi(r.PostFormValue("id-chargement"))
		if err != nil {
			return err
		}
		// Mise à jour du stock du tas
		// Ajoute la qté du chargement avant update chargement
		// puis enlève la qté après update chargement
		// Attention, le tas avant update n'est pas forcément le même que le tas après update
		// (cas où changement de tas lors de update chargement)
		vcAvant, err := model.GetVenteCharge(ctx.DB, vc.Id)
		if err != nil {
			return err
		}
		err = vcAvant.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		err = vcAvant.Tas.ModifierStock(ctx.DB, vcAvant.Qte) // Ajoute des plaquettes au tas
		if err != nil {
			return err
		}
		//
		err = vc.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		err = vc.Tas.ModifierStock(ctx.DB, -vc.Qte) // Retire des plaquettes au tas
		if err != nil {
			return err
		}
		//
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
		// full pour avoir les noms du chargeur, conducteur et proprioutil
		vc, err := model.GetVenteChargeFull(ctx.DB, idChargement)
		if err != nil {
			return err
		}
		// full pour avoir la vente et le nom de la livraison
		vc.Livraison, err = model.GetVenteLivreFull(ctx.DB, vc.IdLivraison)
		if err != nil {
			return err
		}
		// full pour avoir le nom de la vente
		vc.Livraison.Vente, err = model.GetVentePlaqFull(ctx.DB, vc.Livraison.IdVente)
		if err != nil {
			return err
		}
		weboTas, err := WeboTas(ctx)
		if err != nil {
			return err
		}
		listeActeurs, err := model.GetListeActeurs(ctx.DB)
		if err != nil {
			return err
		}
		ctx.TemplateName = "ventecharge-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier le chargement : " + vc.String(),
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Menu: "ventes",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js"},
			},
			Details: detailsVenteChargeForm{
				VenteCharge:  vc,
				TasOptions:   webo.FmtOptions(weboTas, "tas-" + strconv.Itoa(vc.IdTas)),
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL", "gl-"), strconv.FormatFloat(vc.GlTVA, 'f', 1, 64)),
				MoTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_MO", "mo-"), strconv.FormatFloat(vc.MoTVA, 'f', 1, 64)),
				OuTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_OU", "ou-"), strconv.FormatFloat(vc.OuTVA, 'f', 1, 64)),
			    ListeActeurs:  listeActeurs,
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
	vc.IdVente, err = strconv.Atoi(r.PostFormValue("id-vente"))
	if err != nil {
		return vc, err
	}
	//
	vc.IdTas, err = strconv.Atoi(strings.TrimLeft(r.PostFormValue("tas"), "tas-"))
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
		vc.IdChargeur, err = strconv.Atoi(r.PostFormValue("id-chargeur"))
		if err != nil {
			return vc, err
		}
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
		// coût détaillé, conducteur
		//
		vc.IdConducteur, err = strconv.Atoi(r.PostFormValue("id-conducteur"))
		if err != nil {
			return vc, err
		}
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
		//
		// coût détaillé, outil
		//
		vc.IdProprioutil, err = strconv.Atoi(r.PostFormValue("id-proprioutil"))
		if err != nil {
			return vc, err
		}
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
	}
	//
	vc.Notes = r.PostFormValue("notes")
	//
	return vc, nil
}
