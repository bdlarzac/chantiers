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
	"bdl.local/bdl/generic/tiglib"
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
	tmp = strings.Split(r.PostFormValue("ids-lieudits"), ",")
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
	tmp = strings.Split(r.PostFormValue("ids-fermiers"), ",")
	for _, str = range tmp {
		id, _ = strconv.Atoi(str)
		ids = append(ids, id)
	}
	return ids
}

func form2LienParcelles(r *http.Request) (result []*model.ChantierParcelle) {
	//
	// parcelles ; ex de valeurs :
	// radio-parcelle-148:[radio-parcelle-entiere-148]
	// parcelle-surface-148:[]
	// radio-parcelle-337:[radio-parcelle-surface-337]
	// parcelle-surface-337:[3.5]
	//  => parcelle 148 entière et parcelle 337 = surface de 3.5 ha
	result = []*model.ChantierParcelle{}
	idChaufer, _ := strconv.Atoi(r.PostFormValue("id-chantier"))
	for k, v := range r.PostForm {
		if strings.HasPrefix(k, "radio-parcelle-") {
			lien := model.ChantierParcelle{}
			lien.IdChantier = idChaufer
			idPString := strings.Replace(k, "radio-parcelle-", "", -1)
			idP, _ := strconv.Atoi(idPString)
			lien.IdParcelle = idP
			if v[0] == "radio-parcelle-entiere-"+idPString {
				lien.Entiere = true
			} else if v[0] == "radio-parcelle-surface-"+idPString {
				lien.Entiere = false
				lien.Surface, _ = strconv.ParseFloat(r.PostFormValue("parcelle-surface-"+idPString), 32)
				lien.Surface = tiglib.Round(lien.Surface, 2)
			} else {
				continue
			}
			result = append(result, &lien)
		}
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
	tmp = strings.Split(r.PostFormValue("ids-parcelles"), ",")
	for _, str = range tmp {
		id, err = strconv.Atoi(str)
		if err != nil {
			return rien, rien, rien, rien, err
		}
		idsParcelles = append(idsParcelles, id)
	}
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
