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

type detailsVenteLivreForm struct {
	VenteLivre   *model.VenteLivre
	GlTVAOptions template.HTML
	MoTVAOptions template.HTML
	OuTVAOptions template.HTML
	UrlAction    string
}
type detailsVenteLivreList struct {
	Ventes []*model.VenteLivre
	Annee  string   // année courante
	Annees []string // toutes les années avec chantier
}

// *********************************************************
// Process ou affiche form new
func NewVenteLivre(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		vl, err := venteLivreForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertVenteLivre(ctx.DB, vl)
		if err != nil {
			return err
		}
		ctx.Redirect = "/vente/" + strconv.Itoa(vl.IdVente)
		return nil
	default:
		//
		// Affiche form
		//
		var err error
		vars := mux.Vars(r)
		idVente, _ := strconv.Atoi(vars["id-vente"])
		vl := &model.VenteLivre{}
		vl.IdVente = idVente
		vl.Livreur = &model.Acteur{}
		vl.Conducteur = &model.Acteur{}
		vl.Proprioutil = &model.Acteur{}
		// pour afficher le nom de la vente => besoin du nom client => besoin de GetVentePlaqFull
		vl.Vente, err = model.GetVentePlaqFull(ctx.DB, idVente)
		if err != nil {
			return err
		}
		ctx.TemplateName = "ventelivre-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouvelle Livraison plaquettes",
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
			Details: detailsVenteLivreForm{
				VenteLivre:   vl,
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL", "gl-"), "CHOOSE_TVA_GL"),
				MoTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_MO", "mo-"), "CHOOSE_TVA_MO"),
				OuTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_OU", "ou-"), "CHOOSE_TVA_OU"),
				UrlAction:    "/vente/" + vars["id-vente"] + "/livraison/new",
			},
		}
		return nil
	}
}

// *********************************************************
// Process ou affiche form update
func UpdateVenteLivre(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		vl, err := venteLivreForm2var(r)
		if err != nil {
			return err
		}
		vl.Id, err = strconv.Atoi(r.PostFormValue("id-ventelivre"))
		if err != nil {
			return err
		}
		err = model.UpdateVenteLivre(ctx.DB, vl)
		if err != nil {
			return err
		}
		ctx.Redirect = "/vente/" + strconv.Itoa(vl.IdVente)
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idLivraison, err := strconv.Atoi(vars["id-livraison"])
		if err != nil {
			return err
		}
		vl, err := model.GetVenteLivreFull(ctx.DB, idLivraison)
		if err != nil {
			return err
		}
		// pour afficher le nom de la vente => besoin du nom client => besoin de GetVentePlaqFull
		vl.Vente, err = model.GetVentePlaqFull(ctx.DB, vl.IdVente)
		if err != nil {
			return err
		}
		ctx.TemplateName = "ventelivre-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier la livraison : " + vl.String(),
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
			Details: detailsVenteLivreForm{
				VenteLivre:   vl,
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL", "gl-"), strconv.FormatFloat(vl.GlTVA, 'f', 1, 64)),
				MoTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_MO", "mo-"), strconv.FormatFloat(vl.MoTVA, 'f', 1, 64)),
				OuTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_OU", "ou-"), strconv.FormatFloat(vl.OuTVA, 'f', 1, 64)),
				UrlAction:    "/vente/" + vars["id-vente"] + "/livraison/update/" + vars["id-livraison"],
			},
		}
		return nil
	}
}

// *********************************************************
func DeleteVenteLivre(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id-livraison"])
	if err != nil {
		return err
	}
	err = model.DeleteVenteLivre(ctx.DB, id)
	if err != nil {
		return err
	}
	ctx.Redirect = "/vente/" + vars["id-vente"] + "/livraison/update/" + vars["id-livraison"]
	return nil
}

// *********************************************************
// Fabrique une VenteLivre à partir des valeurs d'un formulaire.
// Auxiliaire de NewVenteLivre() et UpdateVenteLivre()
// Ne gère pas le champ Id
func venteLivreForm2var(r *http.Request) (*model.VenteLivre, error) {
	vl := &model.VenteLivre{}
	var err error
	if err = r.ParseForm(); err != nil {
		return vl, err
	}
	//
	if r.PostFormValue("type-cout") == "cout-global" {
		vl.TypeCout = "G"
	} else {
		vl.TypeCout = "D"
	}
	//
	vl.IdVente, err = strconv.Atoi(r.PostFormValue("id-vente"))
	if err != nil {
		return vl, err
	}
	//
	vl.DateLivre, err = time.Parse("2006-01-02", r.PostFormValue("datelivre"))
	if err != nil {
		return vl, err
	}
	if vl.TypeCout == "G" {
		//
		// coût global
		//
		vl.IdLivreur, err = strconv.Atoi(r.PostFormValue("id-livreur"))
		if err != nil {
			return vl, err
		}
		//
		vl.GlPrix, err = strconv.ParseFloat(r.PostFormValue("glprix"), 32)
		if err != nil {
			return vl, err
		}
		vl.GlPrix = tiglib.Round(vl.GlPrix, 2)
		//
		vl.GlTVA, err = strconv.ParseFloat(r.PostFormValue("gltva"), 32)
		if err != nil {
			return vl, err
		}
		vl.GlTVA = tiglib.Round(vl.GlTVA, 2)
		//
		if r.PostFormValue("gldatepay") != "" {
			vl.GlDatePay, err = time.Parse("2006-01-02", r.PostFormValue("gldatepay"))
			if err != nil {
				return vl, err
			}
		}
	} else {
		//
		// coût détaillé, conducteur
		//
		vl.IdConducteur, err = strconv.Atoi(r.PostFormValue("id-conducteur"))
		if err != nil {
			return vl, err
		}
		//
		vl.MoNHeure, err = strconv.ParseFloat(r.PostFormValue("monheure"), 32)
		if err != nil {
			return vl, err
		}
		vl.MoNHeure = tiglib.Round(vl.MoNHeure, 2)
		//
		vl.MoPrixH, err = strconv.ParseFloat(r.PostFormValue("moprixh"), 32)
		if err != nil {
			return vl, err
		}
		vl.MoPrixH = tiglib.Round(vl.MoPrixH, 2)
		//
		vl.MoTVA, err = strconv.ParseFloat(r.PostFormValue("motva"), 32)
		if err != nil {
			return vl, err
		}
		vl.MoTVA = tiglib.Round(vl.MoTVA, 2)
		//
		if r.PostFormValue("modatepay") != "" {
			vl.MoDatePay, err = time.Parse("2006-01-02", r.PostFormValue("modatepay"))
			if err != nil {
				return vl, err
			}
		}
		//
		// coût détaillé, outil
		//
		vl.IdProprioutil, err = strconv.Atoi(r.PostFormValue("id-proprioutil"))
		if err != nil {
			return vl, err
		}
		//
		vl.OuPrix, err = strconv.ParseFloat(r.PostFormValue("ouprix"), 32)
		if err != nil {
			return vl, err
		}
		vl.OuPrix = tiglib.Round(vl.OuPrix, 2)
		//
		vl.OuTVA, err = strconv.ParseFloat(r.PostFormValue("outva"), 32)
		if err != nil {
			return vl, err
		}
		vl.OuTVA = tiglib.Round(vl.OuTVA, 2)
		//
		if r.PostFormValue("oudatepay") != "" {
			vl.OuDatePay, err = time.Parse("2006-01-02", r.PostFormValue("oudatepay"))
			if err != nil {
				return vl, err
			}
		}
	}
	//
	vl.Notes = r.PostFormValue("notes")
	//
	return vl, nil
}
