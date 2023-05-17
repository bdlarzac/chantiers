/*
*****************************************************************************

	Valorisations
	cf type postgres typevalo

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-05-17 08:35:58+02:00, Thierry Graff : Creation à partir de code existant, refactoring #19

*******************************************************************************
*/
package model

var valoMap = map[string]string{
	"BO": "Bois d'oeuvre",
	"CF": "Chauffage fermier",
	"CH": "Chauffage client",
	"PI": "Piquets",
	"PL": "Palette",
	"PP": "Pâte à papier",
	"PQ": "Plaquettes",
}


// ************************** Labels *******************************

// Labels du type de valorisation (pour Chautre)
// @param abbrev : valorisation telle que stockée en base
func LabelValorisation(abbrev string) string {
	switch abbrev {
	case "PP":
		return "Pâte à papier"
	case "CH":
		return "Chauffage"
	case "PL":
		return "Palette"
	case "PI":
		return "Piquets"
	case "BO":
		return "Bois d'oeuvre"
	}
	return "??? BUG LabelValorisation ("+abbrev+") ???"
}

// Labels du type de valorisation
// Utilisé par bilans, où
// - "Chauffage" doit être séparé en 2 : chauffage fermier ( de chaufer) et chauffage client (de chautre)
// - On doit ajouter chantiers plaquettes
func LabelValorisationAvecChauferEtPlaq(abbrev string) string {
	switch abbrev {
	case "CF":
		return "Chauffage fermier"
	case "CH":
		return "Chauffage client"
	case "PQ":
		return "Plaquettes"
	}
	return LabelValorisation(abbrev)
}


// ************************** Codes *******************************

// L'ordre des valorisations correspond à la demande de BDL
// @return  Tous les codes des valorisations de chautre
// Correspond au type typevalo en base
func AllValorisationCodes() []string {
	return []string{"PP", "CH", "PL", "PI", "BO"}
}

// Pareil que AllValorisationCodes(), avec en plus CF (chauffage fermier) et PQ (plaquettes)
// Utilisé pour la recherche de ventes
func AllValorisationCodesAvecChaufer() []string {
    // codes triés par ordre alphabétique des labels correspondant
	return []string{"BO", "CH", "PP", "PL", "PI", "PQ"}
}

// Pareil que AllValorisationCodes(), avec en plus CF (chauffage fermier) et PQ (plaquettes)
// Utilisé pour la recherche d'activités
func AllValorisationCodesAvecChauferEtPlaq() []string {
    // codes triés par ordre alphabétique des labels correspondant
	return []string{"BO", "CH", "CF", "PP", "PL", "PI", "PQ"}
}

