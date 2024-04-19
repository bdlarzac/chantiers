package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Pour la modification de la date dans la liste des chantiers autres valorisations
func UpdateChautre_datePaiement(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := "ok"
	msg := ""
	err = model.UpdateChautre_datePaiement(ctx.DB, id, vars["date"])
	if err != nil {
		ok = "nok"
		msg = err.Error()
	}
	json, _ := json.Marshal(map[string]string{"ok": ok, "message": msg})
	w.Write(json)
	return err
}
