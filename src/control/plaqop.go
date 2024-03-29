/*
@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/webo"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type detailsPlaqOpForm struct {
	TVAOptions    template.HTML
	UniteOptions  template.HTML
	TypeOpOptions template.HTML
	Op            *model.PlaqOp
	ListeActeurs  map[int]string
	UrlAction     string
}

// Process ou affiche form new
func NewPlaqOp(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		op, err := plaqOpForm2var(r)
		if err != nil {
			return werr.Wrap(err)
		}
		_, err = model.InsertPlaqOp(ctx.DB, op)
		if err != nil {
			return werr.Wrap(err)
		}
		ctx.Redirect = "/chantier/plaquette/" + strconv.Itoa(op.IdChantier) + "/chantiers"
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idChantierStr := vars["id-chantier"]
		idChantier, err := strconv.Atoi(idChantierStr)
		if err != nil {
			return werr.Wrap(err)
		}
		op := &model.PlaqOp{}
		op.Acteur = &model.Acteur{}
		op.IdChantier = idChantier
		op.Chantier, err = model.GetPlaq(ctx.DB, idChantier)
		if err != nil {
			return werr.Wrap(err)
		}
		err = op.Chantier.ComputeLieudits(ctx.DB) // Pour afficher nom chantier
		if err != nil {
			return werr.Wrap(err)
		}
		listeActeurs, err := model.GetListeActeurs(ctx.DB)
		if err != nil {
			return werr.Wrap(err)
		}
		ctx.TemplateName = "plaqop-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Nouvelle opération chantier plaquettes",
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Menu: "production",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js"},
			},
			Details: detailsPlaqOpForm{
				TVAOptions:    webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA", "tva-"), "CHOOSE_TVA"),
				UniteOptions:  webo.FmtOptions(WeboPlaqOpUnite(), "CHOOSE_UNITE"),
				TypeOpOptions: webo.FmtOptions(WeboTypeOp(), "CHOOSE_TYPEOP"),
				Op:            op,
				ListeActeurs:  listeActeurs,
				UrlAction:     "/chantier/plaquette/" + idChantierStr + "/op/new",
			},
		}
		return nil
	}
}

// Process ou affiche form update
func UpdatePlaqOp(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		op, err := plaqOpForm2var(r)
		if err != nil {
			return werr.Wrap(err)
		}
		op.Id, err = strconv.Atoi(r.PostFormValue("id-op"))
		if err != nil {
			return werr.Wrap(err)
		}
		err = model.UpdatePlaqOp(ctx.DB, op)
		if err != nil {
			return werr.Wrap(err)
		}
		ctx.Redirect = "/chantier/plaquette/" + strconv.Itoa(op.IdChantier) + "/chantiers"
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		idOp, err := strconv.Atoi(vars["id-op"])
		if err != nil {
			return werr.Wrap(err)
		}
		op, err := model.GetPlaqOp(ctx.DB, idOp)
		if err != nil {
			return werr.Wrap(err)
		}
		op.Acteur, err = model.GetActeur(ctx.DB, op.IdActeur)
		if err != nil {
			return werr.Wrap(err)
		}
		op.Chantier, err = model.GetPlaq(ctx.DB, op.IdChantier)
		if err != nil {
			return werr.Wrap(err)
		}
		err = op.Chantier.ComputeLieudits(ctx.DB) // Pour afficher nom chantier
		if err != nil {
			return werr.Wrap(err)
		}
		listeActeurs, err := model.GetListeActeurs(ctx.DB)
		if err != nil {
			return werr.Wrap(err)
		}
		ctx.TemplateName = "plaqop-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier l'opération : " + model.LabelActivite(op.TypOp),
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Menu: "production",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js"},
			},
			Details: detailsPlaqOpForm{
				TVAOptions:    webo.FmtOptions(WeboTVAExt(ctx, "CHOOSE_TVA", "tva-"), strconv.FormatFloat(op.TVA, 'f', 1, 64)),
				UniteOptions:  webo.FmtOptions(WeboPlaqOpUnite(), "unite-"+op.Unite),
				TypeOpOptions: webo.FmtOptions(WeboTypeOp(), "typeop-"+op.TypOp),
				Op:            op,
				ListeActeurs:  listeActeurs,
				UrlAction:     "/chantier/plaquette/" + vars["id-chantier"] + "/op/update/" + vars["id-op"],
			},
		}
		return nil
	}
}

func DeletePlaqOp(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idOp, err := strconv.Atoi(vars["id-op"])
	if err != nil {
		return werr.Wrap(err)
	}
	err = model.DeletePlaqOp(ctx.DB, idOp)
	if err != nil {
		return werr.Wrap(err)
	}
	ctx.Redirect = "/chantier/plaquette/" + vars["id-chantier"] + "/chantiers"
	return nil
}

// Fabrique une PlaqOp à partir des valeurs d'un formulaire.
// Auxiliaire de NewPlaqOp() et UpdatePlaqOp()
// Ne gère pas le champ Id
func plaqOpForm2var(r *http.Request) (*model.PlaqOp, error) {
	op := &model.PlaqOp{}
	var err error
	if err = r.ParseForm(); err != nil {
		return op, werr.Wrap(err)
	}
	//
	op.TypOp = strings.Replace(r.PostFormValue("type-op"), "typeop-", "", -1)
	if err != nil {
		return op, werr.Wrap(err)
	}
	//
	op.IdChantier, err = strconv.Atoi(r.PostFormValue("id-chantier"))
	if err != nil {
		return op, werr.Wrap(err)
	}
	//
	op.IdActeur, err = strconv.Atoi(r.PostFormValue("id-acteur"))
	if err != nil {
		return op, werr.Wrap(err)
	}
	//
	op.DateDebut, err = time.Parse("2006-01-02", r.PostFormValue("date-debut"))
	if err != nil {
		return op, werr.Wrap(err)
	}
	op.DateFin, err = time.Parse("2006-01-02", r.PostFormValue("date-fin"))
	if err != nil {
		return op, werr.Wrap(err)
	}
	//
	op.Qte, err = strconv.ParseFloat(r.PostFormValue("qte"), 32)
	if err != nil {
		return op, werr.Wrap(err)
	}
	op.Qte = tiglib.Round(op.Qte, 2)
	//
	op.Unite = strings.Replace(r.PostFormValue("unite"), "unite-", "", -1)
	//
	op.PUHT, err = strconv.ParseFloat(r.PostFormValue("puht"), 32)
	if err != nil {
		return op, werr.Wrap(err)
	}
	op.PUHT = tiglib.Round(op.PUHT, 2)
	//
	op.TVA, err = strconv.ParseFloat(r.PostFormValue("tva"), 32)
	if err != nil {
		return op, werr.Wrap(err)
	}
	//
	if r.PostFormValue("date-pay") != "" {
		op.DatePay, err = time.Parse("2006-01-02", r.PostFormValue("date-pay"))
		if err != nil {
			return op, werr.Wrap(err)
		}
	}
	//
	op.Notes = r.PostFormValue("notes")
	return op, nil
}
