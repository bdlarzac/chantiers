/*
*****************************************************************************

	Code lié aux unites.

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-05-17 09:46:52+02:00, Thierry Graff : Creation à partir de code existant, refactoring #19

*******************************************************************************
*/
package model

// Association code unité => label
// Les codes correspondent aux valeurs stockées en base dans les tables plaqop, chaufer, chautre
// = valeurs du type postgres typeunite
var UniteMap = map[string]string{
	"HE": "heures",
	"JO": "jours",
	"M3": "m<sup>3</sup>",
	"MA": "maps",
	"NP": "piquets", // nb de piquets
	"ST": "stères",
	"TO": "tonnes",
}
