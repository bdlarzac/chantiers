package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// *********************************************************
func ShowParcelle(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// TODO show 404
	}
	parcelle, err := model.GetParcelle(ctx.DB, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetParcelle()")
	}

	err = parcelle.ComputeLieudits(ctx.DB)
	if err != nil {
		return err
	}

	for i, _ := range parcelle.Lieudits {
		err = parcelle.Lieudits[i].ComputeCommunes(ctx.DB)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel ComputeCommunes()")
		}
	}

	err = parcelle.ComputeProprietaire(ctx.DB)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ComputeProprietaire()")
	}

	err = parcelle.ComputeFermiers(ctx.DB)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ComputeFermiers()")
	}

	err = parcelle.ComputeUGs(ctx.DB)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ComputeUGs()")
	}

	ctx.TemplateName = "parcelle-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Parcelle " + parcelle.Code,
		},
		Menu:    "accueil",
		Details: parcelle,
	}
	return nil
}
