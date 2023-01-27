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
*/

/** 
    Renvoie les parcelles d'une UG.
    @param  vars["id"]  id numérique d'une UG
    TODO A priori, supprimer après #9
**/
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
		Surface  float32          `json:"surface"`
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
		resp = append(resp, respElement{p.Id, p.Code, p.Communes, p.Surface})
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

/** 
    Renvoie les parcelles correspondant à plusieurs UGs.
    @param  vars["ids"] string contenant les ids numériques des UGs, séparés par des virgules.0
            ex : 12,35,87
**/
func GetParcellesFromUGs(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	// pas terrible, vars["ids"] est directement utilisé pour fabriquer le sql
	// risque d'injection, mais le routing impose [1-9,]
	parcelles, err := model.GetParcellesFromUGs(ctx.DB, vars["ids"])
	if err != nil {
		return err
	}
	type respElement struct {
		Id       int              `json:"id"`
		Name     string           `json:"name"`
		Commune  *model.Commune   `json:"commune"`
		Surface  float32          `json:"surface"`
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
