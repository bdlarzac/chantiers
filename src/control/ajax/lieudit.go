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

func AutocompleteLieudit(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	str := vars["str"]
	type respElement struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var resp []respElement
	lds, err := model.GetLieuDitAutocomplete(ctx.DB, str)
	if err != nil {
		return err
	}
	for _, a := range lds {
		resp = append(resp, respElement{a.Id, a.Nom})
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

// @return  Json contenant id du lieu-dit correspondant à str,
//          ou 0 si lieu-dit pas trouvé.
func CheckNomLieudit(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	str := vars["str"]
	ld, err := model.GetLieuditByNom(ctx.DB, str)
	resp := 0
	if err == nil {
		resp = ld.Id
	}
	// else error => LD inexistant => id reste à 0
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

// Renvoie une liste d'acteurs associés à un lieu-dit (ces acteurs sont forcément des fermiers).
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
	acteurs, err := model.GetFermiersFromLieudit(ctx.DB, idLieudit)
	if err != nil {
		return err
	}
	for _, a := range acteurs {
		// on met bien id, pas id_sctl
		resp = append(resp, respElement{a.Id, a.String()})
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
