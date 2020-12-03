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

func AutocompleteActeur(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	str := vars["str"]
	typeActeur := vars["type-acteur"]

	type respElement struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var resp []respElement

	var err error
	acteurs := []*model.Acteur{}
	switch typeActeur {
	case "":
		acteurs, err = model.GetActeursAutocomplete(ctx.DB, str)
	case "fermiers":
		acteurs, err = model.GetFermiersAutocomplete(ctx.DB, str)
	}
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

// Renvoie une liste d'acteurs associés à une UG (ces acteurs sont forcément des fermiers).
func GetFermiersFromCodeUG(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	codeUG := vars["code"]
	
	type respElement struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var resp []respElement
	acteurs, err := model.GetFermiersFromCodeUG(ctx.DB, codeUG)
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

