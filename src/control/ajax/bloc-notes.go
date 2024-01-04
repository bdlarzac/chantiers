package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"encoding/json"
	"net/http"
)

// Renvoie les parcelles correspondant à plusieurs UGs.
// @param  vars["ids"] string contenant les ids numériques des UGs, séparés par des virgules.0
//
//	ex : 12,35,87
func GetBlocnotes(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	contenu, err := model.GetBlocnotes(ctx.DB)
	if err != nil {
		return err
	}
	json, err := json.Marshal(contenu)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}
