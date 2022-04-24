package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	//	"strconv"
)

// Renvoie une liste de fermiers associés à une UG
func GetFermiersFromCodeUG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	codeUG := vars["code"]
	type respElement struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var resp []respElement
	fermiers, err := model.GetFermiersFromCodeUG(ctx.DB, codeUG)
	if err != nil {
		return err
	}
	for _, f := range fermiers {
		// on met bien id, pas id_fermier
		resp = append(resp, respElement{f.Id, f.String()})
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

/*
// TODO supprimer si toujours inutile
func GetFermiersFromLieudit(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
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
	fermiers, err := model.GetFermiersFromLieudit(ctx.DB, idLieudit)
	if err != nil {
		return err
	}
	for _, f := range fermiers {
		resp = append(resp, respElement{f.Id, f.String()})
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}
*/
