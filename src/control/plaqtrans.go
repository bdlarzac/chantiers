package control

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/generic/wilk/webo"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
)

type detailsPlaqTransForm struct {
	UrlAction    string
	GlTVAOptions template.HTML
	CoTVAOptions template.HTML
	CaTVAOptions template.HTML
	TbTVAOptions template.HTML
	Transport    *model.PlaqTrans
}

// *********************************************************
// Process ou affiche form new
func NewPlaqTrans(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		pt, err := plaqTransForm2var(r)
		if err != nil {
			return err
		}
		_, err = model.InsertPlaqTrans(ctx.DB, pt) // gère la modif du stock du tas
		if err != nil {
			//return err
			return werr.Wrapf(err, "Erreur appel model.InsertPlaqTrans()")
		}
		ctx.Redirect = "/chantier/plaquette/" + strconv.Itoa(pt.IdChantier) + "/chantiers"
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idChantierStr := vars["id-chantier"]
		idChantier, err := strconv.Atoi(idChantierStr)
		if err != nil {
			return err
		}
		pt := &model.PlaqTrans{}
		pt.TypeCout = "G"
		pt.Transporteur = &model.Acteur{}
		pt.Conducteur = &model.Acteur{}
		pt.Proprioutil = &model.Acteur{}
		pt.IdChantier = idChantier
		pt.Chantier, err = model.GetPlaq(ctx.DB, idChantier)
		if err != nil {
			return err
		}
		err = pt.Chantier.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		err = pt.Chantier.ComputeLieudits(ctx.DB) // pour le nom du chantier
		if err != nil {
			return err
		}
		ctx.TemplateName = "plaqtrans-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Créer un transport plaquette",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/autocomplete/autocomplete.js",
					"/view/common/checkActeur.js",
					"/view/common/getActeurPossibles.js"},
			},
			Details: detailsPlaqTransForm{
				Transport:    pt,
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL", "gl-"), "CHOOSE_TVA_GL"),
				CoTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_CO", "co-"), "CHOOSE_TVA_CO"),
				CaTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_CA", "ca-"), "CHOOSE_TVA_CA"),
				TbTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_TB", "tb-"), "CHOOSE_TVA_TB"),
				UrlAction:    "/chantier/plaquette/" + idChantierStr + "/transport/new",
			},
		}
		return nil
	}
}

// *********************************************************
// Process ou affiche form update
func UpdatePlaqTrans(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		pt, err := plaqTransForm2var(r)
		if err != nil {
			return err
		}
		pt.Id, err = strconv.Atoi(r.PostFormValue("id-plaqtrans"))
		if err != nil {
			return err
		}
		err = model.UpdatePlaqTrans(ctx.DB, pt) // gère la modif du stock du tas
		if err != nil {
			return err
		}
		//
		ctx.Redirect = "/chantier/plaquette/" + strconv.Itoa(pt.IdChantier) + "/chantiers"
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idPt, err := strconv.Atoi(vars["id-pt"])
		if err != nil {
			return err
		}
		pt, err := model.GetPlaqTrans(ctx.DB, idPt)
		if err != nil {
			return err
		}
		pt.Transporteur, err = model.GetActeur(ctx.DB, pt.IdTransporteur)
		if err != nil {
			return err
		}
		pt.Conducteur, err = model.GetActeur(ctx.DB, pt.IdConducteur)
		if err != nil {
			return err
		}
		pt.Proprioutil, err = model.GetActeur(ctx.DB, pt.IdProprioutil)
		if err != nil {
			return err
		}
		pt.Chantier, err = model.GetPlaq(ctx.DB, pt.IdChantier)
		if err != nil {
			return err
		}
		err = pt.Chantier.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		err = pt.Chantier.ComputeLieudits(ctx.DB) // pour le nom du chantier
		if err != nil {
			return err
		}
		ctx.TemplateName = "plaqtrans-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier un transport plaquette",
				CSSFiles: []string{
					"/static/css/form.css",
					"/static/autocomplete/autocomplete.css"},
			},
			Menu: "chantiers",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
					"/static/autocomplete/autocomplete.js",
					"/view/common/checkActeur.js",
					"/view/common/getActeurPossibles.js"},
			},
			Details: detailsPlaqTransForm{
				Transport:    pt,
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL", "gl-"), strconv.FormatFloat(pt.GlTVA, 'f', 1, 64)),
				CoTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_CO", "co-"), strconv.FormatFloat(pt.CoTVA, 'f', 1, 64)),
				CaTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_CA", "ca-"), strconv.FormatFloat(pt.CaTVA, 'f', 1, 64)),
				TbTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_TB", "tb-"), strconv.FormatFloat(pt.TbTVA, 'f', 1, 64)),
				UrlAction:    "/chantier/plaquette/" + vars["id-chantier"] + "/transport/update/" + vars["id-pt"],
			},
		}
		return nil
	}
}

// *********************************************************
func DeletePlaqTrans(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idPt, err := strconv.Atoi(vars["id-pt"])
	if err != nil {
		return err
	}
	err = model.DeletePlaqTrans(ctx.DB, idPt) // gère la modif du stock du tas
	if err != nil {
		return err
	}
	ctx.Redirect = "/chantier/plaquette/" + vars["id-chantier"] + "/chantiers"
	return nil
}

// *********************************************************
// Fabrique un PlaqTrans à partir des valeurs d'un formulaire.
// Auxiliaire de NewPlaqTrans() et UpdatePlaqTrans()
// Ne gère pas le champ Id
func plaqTransForm2var(r *http.Request) (*model.PlaqTrans, error) {
	pt := &model.PlaqTrans{}
	var err error
	var tmp string
	if err = r.ParseForm(); err != nil {
		return pt, err
	}
	//
	// champs généraux
	//
	pt.IdChantier, err = strconv.Atoi(r.PostFormValue("id-chantier"))
	if err != nil {
		return pt, err
	}
	//
	pt.IdTas, err = strconv.Atoi(strings.ReplaceAll(r.PostFormValue("tas"), "tas-", ""))
	if err != nil {
		return pt, err
	}
	//
	pt.DateTrans, err = time.Parse("2006-01-02", r.PostFormValue("datetrans"))
	if err != nil {
		return pt, err
	}
	//
	pt.Qte, err = strconv.ParseFloat(r.PostFormValue("qte"), 32)
	if err != nil {
		return pt, err
	}
	pt.Qte = tiglib.Round(pt.Qte, 2)
	//
	tmp = r.PostFormValue("type-cout")
	if tmp == "cout-global" {
		pt.TypeCout = "G"
	} else {
		tmp = r.PostFormValue("type-transport")
		if tmp == "camion" {
			pt.TypeCout = "C"
		} else {
			pt.TypeCout = "T"
		}
	}
	//
	// concerne le coût global
	//
	if pt.TypeCout == "G" {
		//
		pt.IdTransporteur, err = strconv.Atoi(r.PostFormValue("id-transporteur"))
		if err != nil {
			return pt, err
		}
		pt.GlPrix, err = strconv.ParseFloat(r.PostFormValue("glprix"), 32)
		if err != nil {
			return pt, err
		}
		pt.GlPrix = tiglib.Round(pt.GlPrix, 2)
		//
		pt.GlTVA, err = strconv.ParseFloat(r.PostFormValue("gltva"), 32)
		if err != nil {
			return pt, err
		}
		pt.GlTVA = tiglib.Round(pt.GlTVA, 2)
		//
		if r.PostFormValue("gldatepay") != "" {
			pt.GlDatePay, err = time.Parse("2006-01-02", r.PostFormValue("gldatepay"))
			if err != nil {
				return pt, err
			}
		}
	} else {
		//
		// concerne le propriétaire outil
		//
		pt.IdProprioutil, err = strconv.Atoi(r.PostFormValue("id-proprioutil"))
		if err != nil {
			return pt, err
		}
		//
		// concerne le conducteur
		//
		pt.IdConducteur, err = strconv.Atoi(r.PostFormValue("id-conducteur"))
		if err != nil {
			return pt, err
		}
		//
		pt.CoNheure, err = strconv.ParseFloat(r.PostFormValue("conheure"), 32)
		if err != nil {
			return pt, err
		}
		pt.CoNheure = tiglib.Round(pt.CoNheure, 2)
		//
		pt.CoPrixH, err = strconv.ParseFloat(r.PostFormValue("coprixh"), 32)
		if err != nil {
			return pt, err
		}
		pt.CoPrixH = tiglib.Round(pt.CoPrixH, 2)
		//
		pt.CoTVA, err = strconv.ParseFloat(r.PostFormValue("cotva"), 32)
		if err != nil {
			return pt, err
		}
		pt.CoTVA = tiglib.Round(pt.CoTVA, 2)
		//
		if r.PostFormValue("codatepay") != "" {
			pt.CoDatePay, err = time.Parse("2006-01-02", r.PostFormValue("codatepay"))
			if err != nil {
				return pt, err
			}
		}
		if pt.TypeCout == "C" {
			//
			// Transport camion
			//
			pt.CaNkm, err = strconv.ParseFloat(r.PostFormValue("cankm"), 32)
			if err != nil {
				return pt, err
			}
			pt.CaNkm = tiglib.Round(pt.CaNkm, 2)
			//
			pt.CaPrixKm, err = strconv.ParseFloat(r.PostFormValue("caprixkm"), 32)
			if err != nil {
				return pt, err
			}
			pt.CaPrixKm = tiglib.Round(pt.CaPrixKm, 2)
			//
			pt.CaTVA, err = strconv.ParseFloat(r.PostFormValue("catva"), 32)
			if err != nil {
				return pt, err
			}
			pt.CaTVA = tiglib.Round(pt.CaTVA, 2)
			//
			if r.PostFormValue("cadatepay") != "" {
				pt.CaDatePay, err = time.Parse("2006-01-02", r.PostFormValue("cadatepay"))
				if err != nil {
					return pt, err
				}
			}
		} else {
			//
			// Transport tracteur + benne
			//
			pt.TbNbenne, err = strconv.Atoi(r.PostFormValue("tbnbenne"))
			if err != nil {
				return pt, err
			}
			//
			pt.TbDuree, err = strconv.ParseFloat(r.PostFormValue("tbduree"), 32)
			if err != nil {
				return pt, err
			}
			pt.TbDuree = tiglib.Round(pt.TbDuree, 2)
			//
			pt.TbPrixH, err = strconv.ParseFloat(r.PostFormValue("tbprixh"), 32)
			if err != nil {
				return pt, err
			}
			pt.TbPrixH = tiglib.Round(pt.TbPrixH, 2)
			//
			pt.TbTVA, err = strconv.ParseFloat(r.PostFormValue("tbtva"), 32)
			if err != nil {
				return pt, err
			}
			pt.TbTVA = tiglib.Round(pt.TbTVA, 2)
			//
			if r.PostFormValue("tbdatepay") != "" {
				pt.TbDatePay, err = time.Parse("2006-01-02", r.PostFormValue("tbdatepay"))
				if err != nil {
					return pt, err
				}
			}
		} // end tracteur + benne
	} // end prix détaillé
	pt.Notes = r.PostFormValue("notes")
	return pt, nil
}
