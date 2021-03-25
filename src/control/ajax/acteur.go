package ajax

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func AutocompleteActeur(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	str := vars["str"]

	type respElement struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var resp []respElement

	var err error
	acteurs, err := model.GetActeursAutocomplete(ctx.DB, str)
	if err != nil {
		return err
	}
	var nom string
	for _, a := range acteurs {
		if a.Prenom == "" {
			nom = a.Nom
		} else {
			nom = a.Nom + " " + a.Prenom
		}
		resp = append(resp, respElement{a.Id, nom})

	}
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}

/**
    @return  Json contenant id de l'acteur correspondant à str,
            ou 0 si lieu-dit pas trouvé.
**/
func CheckNomActeur(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	str := vars["str"]
	ld, err := model.GetActeurByNomAutocomplete(ctx.DB, str)
	resp := 0
	if err == nil {
		resp = ld.Id
	}
	// else error => acteur inexistant => id reste à 0
	json, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(json)
	return nil
}
