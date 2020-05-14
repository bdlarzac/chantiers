package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetParcellesFromUG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idUG, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	type respElement struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var resp []respElement
	ug, err := model.GetUG(ctx.DB, idUG)
	if err != nil {
		return err
	}
	err = ug.ComputeParcelles(ctx.DB)
	for _, p := range ug.Parcelles {
		resp = append(resp, respElement{p.Id, p.Code})
	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

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

func GetUGFromCode1(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	type respElement struct {
		Id   int    `json:"id"`
	}
	var resp []respElement
	ug, err := model.GetUGFromCode(ctx.DB, vars["code"])
	if err != nil {
		return err
	}
	if ug != nil {
        resp = append(resp, respElement{ug.Id})
    }
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

