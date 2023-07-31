package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func GetFermiersFromIdsUGs(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	lds, err := model.GetFermiersFromIdsUGs(ctx.DB, vars["ids"])
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
