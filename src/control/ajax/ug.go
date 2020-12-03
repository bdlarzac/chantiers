package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Renvoie l'id de l'UG correspondant au code,
// ou 0 si aucune UG ne correspond
func GetUGFromCode(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var resp int
	ug, err := model.GetUGFromCode(ctx.DB, vars["code"])
	if err != nil {
		return err
	}
	if ug != nil {
		resp = ug.Id
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

func GetUGsFromFermier(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idActeur, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	type respElement struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var resp []respElement
	ugs, err := model.GetUGsFromFermier(ctx.DB, idActeur)
	if err != nil {
		return err
	}
	for _, ug := range ugs {
		resp = append(resp, respElement{ug.Id, ug.String()})
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

func GetUGsFromLieudit(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
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
	ugs, err := model.GetUGsFromLieudit(ctx.DB, idLieudit)
	if err != nil {
		return err
	}
	for _, ug := range ugs {
		resp = append(resp, respElement{ug.Id, ug.String()})
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

