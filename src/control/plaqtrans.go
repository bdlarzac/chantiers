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
	//"fmt"
)

type detailsPlaqTransForm struct {
	UrlAction    string
	GlTVAOptions template.HTML
	TrTVAOptions template.HTML
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
		_, err = model.InsertPlaqTrans(ctx.DB, pt)
		if err != nil {
			return err
		}
		// Mise à jour du stock du tas
		err = pt.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		pt.Tas.ModifierStock(ctx.DB, pt.Qte) // Ajoute plaquettes au tas
		err = model.UpdateTas(ctx.DB, pt.Tas)
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
		idChantierStr := vars["id-chantier"]
		idChantier, err := strconv.Atoi(idChantierStr)
		if err != nil {
			return err
		}
		pt := &model.PlaqTrans{}
		pt.TypeCout = "G"
		pt.Transporteur = &model.Acteur{}
		pt.IdChantier = idChantier
		pt.Chantier, err = model.GetPlaq(ctx.DB, idChantier)
		if err != nil {
			return err
		}
		err = pt.Chantier.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		err = pt.Chantier.ComputeLieudit(ctx.DB) // pour le nom du chantier
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
				JSFiles: []string{"/static/autocomplete/autocomplete.js"},
			},
			Details: detailsPlaqTransForm{
				Transport:    pt,
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL"), "CHOOSE_TVA_GL"),
				TrTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_TR"), "CHOOSE_TVA_TR"),
				CaTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_CA"), "CHOOSE_TVA_CA"),
				TbTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_TB"), "CHOOSE_TVA_TB"),
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
		// Mise à jour du stock du tas
		// Enlève la qté du transport avant update transport
		// puis ajoute qté après update transport
		// Attention, le tas avant update n'est pas forcément le même que le tas après update
		// (cas où plusieurs tas pour un chantier plaquette et changement de tas lors de update transport)
		ptAvant, err := model.GetPlaqTrans(ctx.DB, pt.Id)
		if err != nil {
			return err
		}
		err = ptAvant.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		ptAvant.Tas.ModifierStock(ctx.DB, -ptAvant.Qte) // Retire des plaquettes au tas
		err = model.UpdateTas(ctx.DB, ptAvant.Tas)
		if err != nil {
			return err
		}
		//
		err = pt.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		pt.Tas.ModifierStock(ctx.DB, pt.Qte) // Ajoute des plaquettes au tas
		err = model.UpdateTas(ctx.DB, pt.Tas)
		if err != nil {
			return err
		}
		//
		err = model.UpdatePlaqTrans(ctx.DB, pt)
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
		pt.Chantier, err = model.GetPlaq(ctx.DB, pt.IdChantier)
		if err != nil {
			return err
		}
		err = pt.Chantier.ComputeTas(ctx.DB)
		if err != nil {
			return err
		}
		err = pt.Chantier.ComputeLieudit(ctx.DB) // pour le nom du chantier
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
				JSFiles: []string{"/static/autocomplete/autocomplete.js"},
			},
			Details: detailsPlaqTransForm{
				Transport:    pt,
				GlTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_GL"), strconv.FormatFloat(pt.GlTVA, 'f', 1, 64)),
				TrTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_TR"), strconv.FormatFloat(pt.TrTVA, 'f', 1, 64)),
				CaTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_CA"), strconv.FormatFloat(pt.CaTVA, 'f', 1, 64)),
				TbTVAOptions: webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA_TB"), strconv.FormatFloat(pt.TbTVA, 'f', 1, 64)),
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
	// Enlève le stock du tas concerné par le transport
	// avant de supprimer le transport
	pt, err := model.GetPlaqTrans(ctx.DB, idPt)
	if err != nil {
		return err
	}
	err = pt.ComputeTas(ctx.DB)
	if err != nil {
		return err
	}
	pt.Tas.ModifierStock(ctx.DB, -pt.Qte) // Retire des plaquettes au tas
	err = model.UpdateTas(ctx.DB, pt.Tas)
	if err != nil {
		return err
	}
	//
	err = model.DeletePlaqTrans(ctx.DB, idPt)
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
	pt.IdTransporteur, err = strconv.Atoi(r.PostFormValue("id-transporteur"))
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
		// concerne le transporteur
		//
		pt.TrNheure, err = strconv.ParseFloat(r.PostFormValue("trnheure"), 32)
		if err != nil {
			return pt, err
		}
		pt.TrNheure = tiglib.Round(pt.TrNheure, 2)
		//
		pt.TrPrixH, err = strconv.ParseFloat(r.PostFormValue("trprixh"), 32)
		if err != nil {
			return pt, err
		}
		pt.TrPrixH = tiglib.Round(pt.TrPrixH, 2)
		//
		pt.TrTVA, err = strconv.ParseFloat(r.PostFormValue("trtva"), 32)
		if err != nil {
			return pt, err
		}
		pt.TrTVA = tiglib.Round(pt.TrTVA, 2)
		//
		if r.PostFormValue("trdatepay") != "" {
			pt.TrDatePay, err = time.Parse("2006-01-02", r.PostFormValue("trdatepay"))
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
