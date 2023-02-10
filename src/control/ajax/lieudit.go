package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"strconv"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

func AutocompleteLieudit(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	str := vars["str"]
	type respElement struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var resp []respElement
	lds, err := model.GetLieuditsAutocomplete(ctx.DB, str)
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

///// remove apres #9 
/* func GetLieuditsFromCodeUG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	code := vars["code"]
	lds, err := model.GetLieuditsFromCodeUG(ctx.DB, code)
	if err != nil {
		return err
	}
	json, err := json.Marshal(lds)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
} */

func GetLieuditsFromIdUG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		return err
	}
	lds, err := model.GetLieuditsFromIdUG(ctx.DB, id)
	if err != nil {
		return err
	}
	json, err := json.Marshal(lds)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}
