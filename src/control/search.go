/*
Fonctions auxiliaires liées à la recherche au niveau controller.
Fonctions communes à la recherche d'activités et de à la recherche de ventes.
Voir activite-search.go, vente-search.go pour le contrôle des formulaires de recherche.

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"net/http"
	"strings"
)

// ************************** computeFiltre *******************************

/*
Filtre fermier : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient une liste avec un seul élément, l'id du fermier sélectionné.
*/
func computeFiltreFermier(r *http.Request) (result []string) {
	choix := r.PostFormValue("select-choix-fermier")
	if choix == "choix-fermier-no-limit" {
		return []string{}
	}
	return []string{choix[14:]}
}

/*
Filtre commune : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient une liste avec un seul élément, l'id de la commune sélectionnée.
*/
func computeFiltreCommune(r *http.Request) (result []string) {
	choix := r.PostFormValue("select-choix-commune")
	if choix == "choix-commune-no-limit" {
		return []string{}
	}
	return []string{choix[14:]}
}

/*
Filtre client : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient une liste avec un seul élément, l'id du client sélectionné (id dans la table acteur).
*/
func computeFiltreClient(r *http.Request) (result []string) {
	choix := r.PostFormValue("select-choix-client")
	if choix == "choix-client-no-limit" {
		return []string{}
	}
	return []string{choix[13:]}
}

/*
Filtre essence : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient une liste de codes essence.
*/
func computeFiltreEssence(r *http.Request) (result []string) {
	if r.PostFormValue("choix-ALL-essence") == "true" {
		return []string{}
	}
	result = []string{}
	for key, _ := range r.PostForm {
		if strings.Index(key, "choix-essence-") != 0 {
			continue
		}
		if r.PostFormValue(key) != "on" {
			continue
		}
		code := key[14:]
		result = append(result, code)
	}
	return result
}

/*
Filtre propriétaire : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient une liste d'id propriétaires (dans la table acteur)
    (attention, cette liste contient des strings, pas des ints).
*/
func computeFiltreProprio(r *http.Request) (result []string) {
	if r.PostFormValue("choix-ALL-proprio") == "true" {
		return []string{}
	}
	result = []string{}
	for key, _ := range r.PostForm {
		if strings.Index(key, "choix-proprio-") != 0 {
			continue
		}
		if r.PostFormValue(key) != "on" {
			continue
		}
		id := key[14:]
		result = append(result, id)
	}
	return result
}

/*
Filtre période : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient 2 strings, dates de début et de fin au format AAAA-MM-JJ.
*/
func computeFiltrePeriode(r *http.Request) (result []string) {
	if r.PostFormValue("choix-periode-periodes") == "choix-periode-no-limit" {
		return []string{}
	}
	result = append(result, r.PostFormValue("choix-periode-debut"))
	result = append(result, r.PostFormValue("choix-periode-fin"))
	return result
}

/*
Filtre UG : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient les ids UG
*/
func computeFiltreUG(r *http.Request) (result []string) {
	if r.PostFormValue("ids-ugs") == "" {
		return []string{}
	}
	result = strings.Split(r.PostFormValue("ids-ugs"), ";")
	return result
}

/*
Filtre Parcelles : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient les ids parcelle.
*/
func computeFiltreParcelle(r *http.Request) (result []string) {
	if r.PostFormValue("ids-parcelles") == "" {
		return []string{}
	}
	result = strings.Split(r.PostFormValue("ids-parcelles"), ";")
	return result
}

/*
Filtre valorisations : renvoie un tableau de strings.
  - Si pas de filtre, contient un tableau vide.
  - Sinon contient une liste de codes essence.
*/
func computeFiltreValo(r *http.Request) (result []string) {
	if r.PostFormValue("choix-ALL-valo") == "true" {
		return []string{}
	}
	result = []string{}
	for key, _ := range r.PostForm {
		if strings.Index(key, "choix-valo-") != 0 {
			continue
		}
		if r.PostFormValue(key) != "on" {
			continue
		}
		code := key[11:]
		result = append(result, code)
	}
	return result
}
