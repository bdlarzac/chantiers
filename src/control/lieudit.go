/*
@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"net/http"
	"strconv"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
)

func ShowLieudit(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return werr.Wrapf(err, "Erreur conversion id lieu-dit")
	}
	//
	lieudit, err := model.GetLieudit(ctx.DB, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetLieudit()")
	}
	//
	err = lieudit.ComputeParcelles(ctx.DB)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ComputeParcelles()")
	}
	//
	for i, _ := range lieudit.Parcelles {
		err = lieudit.Parcelles[i].ComputeUGs(ctx.DB)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel ComputeUGs()")
		}
		err = lieudit.Parcelles[i].ComputeProprietaire(ctx.DB)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel ComputeProprietaire()")
		}
		err = lieudit.Parcelles[i].ComputeFermiers(ctx.DB)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel ComputeFermiers()")
		}
		err = lieudit.Parcelles[i].ComputeCommune(ctx.DB)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel Parcelle.ComputeCommune()")
		}
	}
	//
	err = lieudit.ComputeCommune(ctx.DB)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel Lieudit.ComputeCommunes()")
	}
	//
	ctx.TemplateName = "lieudit-show.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: lieudit.Nom,
		},
		Menu:    "accueil",
		Details: lieudit,
	}
	return nil
}
