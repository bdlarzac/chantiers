/*
*****************************************************************************

	Code lié aux unites
	cf type postgres typeunite

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-05-17 09:46:52+02:00, Thierry Graff : Creation à partir de code existant, refactoring #19

*******************************************************************************
*/
package model

// Labels des unités manipulées par l'appli - format adapté à select html (à cause de m3)
// @param abbrev : unité telle que stockée en base
func LabelUnite(abbrev string) string {
	switch abbrev {
	case "HE":
		return "heures"
	case "JO":
		return "jours"
	case "M3":
		return "m<sup>3</sup>"
	case "MA":
		return "maps"
	case "NP": // nb de piquets
		return "piquets"
	case "ST":
		return "stères"
	case "TO":
		return "tonnes"
	}
	return "??? BUG LabelUnite ("+abbrev+") ???"
}
