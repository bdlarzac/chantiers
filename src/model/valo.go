/*
*****************************************************************************

	Valorisations
	cf type postgres typevalorisation

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-05-17 08:35:58+02:00, Thierry Graff : Creation à partir de code existant, refactoring #19

*******************************************************************************
*/
package model

// Association code valorisation => label
// Les codes correspondent aux valeurs stockées en base dans chautre.typevalo
var ValoMap = map[string]string{
	"BO": "Bois d'oeuvre",
	"CF": "Chauffage fermier",
	"CH": "Chauffage client",
	"PI": "Piquets",
	"PL": "Palette",
	"PP": "Pâte à papier",
	"PQ": "Plaquettes",
}

// ************************** Codes *******************************

// Codes utilisés pour chautre
func AllValoCodes() []string {
    // L'ordre des valorisations correspond à la demande de BDL
	return []string{"PP", "CH", "PL", "PI", "BO"}
}

// Pareil que AllValoCodes(), avec en plus CF (chauffage fermier)
// Utilisé pour la recherche de ventes
func AllValoCodesAvecChaufer() []string {
    // codes triés par ordre alphabétique des labels correspondant
	return []string{"BO", "CH", "PP", "PL", "PI", "PQ"}
}

// Pareil que AllValoCodes(), avec en plus CF (chauffage fermier) et PQ (plaquettes)
// Utilisé pour la recherche d'activités
func AllValoCodesAvecChauferEtPlaq() []string {
    // codes triés par ordre alphabétique des labels correspondant
	return []string{"BO", "CH", "CF", "PP", "PL", "PI", "PQ"}
}

// ************************** Unités *******************************

// Renvoie le code de l'unité correspondant à une valorisation, tel que stocké en base
// Comprend les valos pour chautre + CF pour chauffage fermier + PQ pour plaquettes
func CodeValo2CodeUnite(codeValo string) string {
	switch codeValo {
	case "PP":
		return "TO"
	case "CH":
		return "ST"
	case "CF":
		return "ST"
	case "PL":
		return "ST"
	case "PI":
		// 2023-05-17 je ne comprends pas cette histoire
		// l'unité "nb de piquets" n'est plus utilisée
		// suite à demande BDL de pouvoir choisir stères ou nb de piquets
		return "ST"
	case "BO":
		return "M3"
	case "PQ":
		return "MA"
	}
	return "??? Code inconnu dans CodeValo2CodeUnite ("+codeValo+")  ???"
}

