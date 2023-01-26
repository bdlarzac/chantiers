/**
    Code commun à plaq.go et chautre.go
    (chantiers utilisant view/comon/chantier-lien.html)
    Pour gérer les liens chanter - UG, Lieudit, Fermier

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
**/
package control

import (
	"net/http"
	"strconv"
	"strings"
)

// *********************************************************
// Utilisé par
// NewPlaq()
// UpdatePlaq()
// NewChautre()
// UpdateChautre()
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
