/*
*****************************************************************************

	Enums manipulées par le programme

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2020-02-11 17:27:04+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

// ************************** Opérations simples *******************************
// cf type postgres typeop pour opérations simples des chantiers plaquettes

// Labels du type d'opération simple pour PlaqOp
// @param code : type d'opération tel que stocké en base
func LabelActivite(code string) string {
	switch code {
	//
	// opérations simples des chantiers plaquettes
	//
	case "AB":
		return "Abattage"
	case "DB":
		return "Débardage"
	case "DC":
		return "Déchiquetage"
	case "BR":
		return "Broyage"
	//
	// Autres activités
	//
	case "TR":
		return "Transport"
	case "RG":
		return "Rangement"
	case "LV":
		return "Livraison"
	case "CG":
		return "Chargement"
	}
	return "??? BUG LabelActivite ("+code+") ???"
}

// ************************** Type d'exploitation *******************************
// cf type postgres typexploitation

// Labels du type d'exploitation
// @param code : type d'exploitation
func LabelExploitation(code string) string {
	switch code {
	case "1":
		return "1ère"
	case "2":
		return "2ème"
	case "3":
		return "3ème"
	case "4":
		return "4ème"
	case "5":
		return "5ème"
	}
	return "??? BUG LabelExploitation ("+code+") ???"
}

// ************************** Granulométrie *******************************
// cf type postgres typegranulo

// Labels du type de granulométrie
// @param code : type de granulométrie
func LabelGranulo(code string) string {
	return code // le type stocké en base correspond au label
}
