package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	//"fmt"
)

func GetParcellesFromLieudit(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idLieudit, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	type respElement struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var resp []respElement
	parcelles, err := model.GetParcellesFromLieudit(ctx.DB, idLieudit)
	if err != nil {
		return err
	}
	for _, parcelle := range parcelles {
		resp = append(resp, respElement{parcelle.Id, parcelle.Code})
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

func GetParcellesFromUG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idUG, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	type respElement struct {
		Id       int              `json:"id"`
		Name     string           `json:"name"`
		Communes []*model.Commune `json:"communes"`
	}
	var resp []respElement
	ug, err := model.GetUG(ctx.DB, idUG)
	if err != nil {
		return err
	}
	err = ug.ComputeParcelles(ctx.DB)
	if err != nil {
		return err
	}
	for _, p := range ug.Parcelles {
		err = p.ComputeCommunes(ctx.DB)
		if err != nil {
			return err
		}
		resp = append(resp, respElement{p.Id, p.Code, p.Communes})
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}
