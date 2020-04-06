package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// *********************************************************
func SignalerTasVide(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	err = model.DesactiverTas(ctx.DB, id)
	if err != nil {
		return err
	}
	ctx.Redirect = "/stockage/liste"
	return nil
}
