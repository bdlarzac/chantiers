/*

	Essence (= espèce d'arbre)

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2023-03-03 08:57:51+01:00, Thierry Graff : Creation

*/
package model

// Association code essence => label
// Les codes correspondent aux valeurs stockées en base dans :
// - chautre.essence
// - chaufer.essence
// - plaq.essence
var EssenceMap = map[string]string{
    "AL": "Alisier",
    "CN": "Chêne",
    "CT": "Châtaigner",
    "CD": "Cèdre",
    "DG": "Douglas",
    "EM": "Erable de Montpellier",
    "ER": "Erable",
    "FL": "Mélange feuillus",
    "GV": "Genévrier",
    "HT": "Hêtre",
    "PN": "Pin noir",
    "PS": "Pin sylvestre",
}

// Pour avoir les codes dans le bon ordre
var EssenceCodes = []string{
    "AL",
    "CN",
    "CT",
    "CD",
    "DG",
    "EM",
    "ER",
    "FL",
    "GV",
    "HT",
    "PN",
    "PS",
}
