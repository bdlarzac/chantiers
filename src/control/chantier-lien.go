/**
    Code commun à plaq.go et chautre.go
    (chantiers utilisant view/comon/chantier-lien.html)
    Pour gérer les liens chanter - UG, Lieudit, Fermier

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
**/
package control

import (
	"bdl.local/bdl/model"
	"strconv"
	"strings"
	"net/http"
)


/** 
    Calcule les ids ug à partir des champs d'un formulaire chantier
    Utilisé par
        NewPlaq()    UpdatePlaq()
        NewChautre() UpdateChautre()
        NewChaufer() UpdateChaufer()
**/
func form2IdsUG(r *http.Request)(ids []int) {
	var tmp []string
	var str string
	var id int
	//
	tmp = strings.Split(r.PostFormValue("ids-ugs"), ";")
	for _, str = range tmp {
		id, _ = strconv.Atoi(str)
		ids = append(ids, id)
	}
	return ids
}

/** 
    Calcule les ids lieudit à partir des champs d'un formulaire chantier
    Utilisé par
        NewPlaq()    UpdatePlaq()
        NewChautre() UpdateChautre()
        NewChaufer() UpdateChaufer()
**/
func form2IdsLieudit(r *http.Request)(ids []int) {
	var tmp []string
	var str string
	var id int
	//
	tmp = strings.Split(r.PostFormValue("ids-lieudits"), ";")
	for _, str = range tmp {
		id, _ = strconv.Atoi(str)
		ids = append(ids, id)
	}
	return ids
}

/** 
    Calcule les ids fermier à partir des champs d'un formulaire chantier
    Utilisé par
        NewPlaq()    UpdatePlaq()
        NewChautre() UpdateChautre()
**/
func form2IdsFermier(r *http.Request)(ids []int) {
	var tmp []string
	var str string
	var id int
	//
	tmp = strings.Split(r.PostFormValue("ids-fermiers"), ";")
	for _, str = range tmp {
		id, _ = strconv.Atoi(str)
		ids = append(ids, id)
	}
	return ids
}

/** 
    Utilise la variable liens-parcelles pour calculer les model.ChantierParcelle
    ex de liens-parcelles : [1025:entiere;1239:surface-0.10]
**/
func form2LienParcelles(r *http.Request) (result []*model.ChantierParcelle) {
	result = []*model.ChantierParcelle{}
 	idChaufer, _ := strconv.Atoi(r.PostFormValue("id-chantier"))
	strLiens := r.PostFormValue("liens-parcelles")
	if strLiens == "" {
	    return result // ne se produit pas si le choix des parcelles est obligatoire dans le form
	}
	liens := strings.Split(strLiens, ";")
	for _, lien := range(liens) {
	    newChantierParcelle := model.ChantierParcelle{}
        newChantierParcelle.IdChantier = idChaufer
	    tmp := strings.Split(lien, ":")
	    idParcelle, _ := strconv.Atoi(tmp[0])
        newChantierParcelle.IdParcelle = idParcelle
	    what := tmp[1]
	    newChantierParcelle.Entiere = what == "entiere"
	    if ! newChantierParcelle.Entiere {
            newChantierParcelle.Surface, _ = strconv.ParseFloat(strings.Replace(what, "surface-", "", -1), 32)
	    }
	    result = append(result, &newChantierParcelle);
	}
	return result
}

/** 
    ================= A SUPPRIMER lorsque les controlers utilisent form2*() =================
    Utilisé par
        NewPlaq()    UpdatePlaq()
        NewChautre() UpdateChautre()
**/
func calculeIdsLiensChantier(r *http.Request) (idsUGs, idsParcelles, idsLieudits, idsFermiers []int, err error) {
	rien := []int{}
	var tmp []string
	var str string
	var id int
	//
	tmp = strings.Split(r.PostFormValue("ids-ugs"), ",")
	for _, str = range tmp {
		id, err = strconv.Atoi(str)
		if err != nil {
			return rien, rien, rien, rien, err
		}
		idsUGs = append(idsUGs, id)
	}
	//
	tmp = strings.Split(r.PostFormValue("liens-parcelles"), ",")
/* 
	for _, str = range tmp {
		id, err = strconv.Atoi(str)
		if err != nil {
			return rien, rien, rien, rien, err
		}
		idsParcelles = append(idsParcelles, id)
	}
*/
	//
	tmp = strings.Split(r.PostFormValue("ids-lieudits"), ",")
	for _, str = range tmp {
		id, err = strconv.Atoi(str)
		if err != nil {
			return rien, rien, rien, rien, err
		}
		idsLieudits = append(idsLieudits, id)
	}
	//
	tmp = strings.Split(r.PostFormValue("ids-fermiers"), ",")
	for _, str = range tmp {
		id, err = strconv.Atoi(str)
		if err != nil {
			return rien, rien, rien, rien, err
		}
		idsFermiers = append(idsFermiers, id)
	}
	return idsUGs, idsParcelles, idsLieudits, idsFermiers, nil
}
