package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
"fmt"
)

// Pour la modification de la date dans la liste des ventes
func UpdateVentePlaq(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
//	id, _ := strconv.Atoi(vars["id"])
id, _ := strconv.Atoi("656565")
	ok := "ok"
	msg := ""
	err = model.UpdateVentePlaq_datePaiement(ctx.DB, id, vars["date"])
fmt.Printf("dans controller, err = %+v\n",err)
	if err != nil {
	    ok = "nok"
	    msg = err.Error()
	}
	json, _ := json.Marshal(map[string]string{"ok": ok, "message": msg})
	w.Write(json)
	return err
}
