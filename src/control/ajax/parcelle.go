package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

/*
	    Renvoie les parcelles correspondant à plusieurs UGs.
	    @param  vars["ids"] string contenant les ids numériques des UGs, séparés par des virgules.0
		        ex : 12,35,87
*/
func GetParcellesFromIdsUGs(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	parcelles, err := model.GetParcellesFromIdsUGs(ctx.DB, vars["ids"])
	if err != nil {
		return err
	}
	type respElement struct {
		Id      int            `json:"id"`
		Name    string         `json:"name"`
		Commune *model.Commune `json:"commune"`
		Surface float64        `json:"surface"`
	}
	var resp []respElement
	for _, p := range parcelles {
		err = p.ComputeCommune(ctx.DB)
		if err != nil {
			return err
		}
		resp = append(resp, respElement{p.Id, p.Code, p.Commune, p.Surface})
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

func GetParcelleFromCodeAndCommuneId(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	idCommune, err := strconv.Atoi(vars["id-commune"])
	if err != nil {
		return err
	}
	parcelle, err := model.GetParcelleFromCodeAndCommuneId(ctx.DB, vars["code-parcelle"], idCommune)
	if err != nil {
		return err
	}
	json, err := json.Marshal(parcelle)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}
