/*
Typo = couche typologique venant du PSG

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2023-03-03 12:11:29+01:00, Thierry Graff : Creation
*/
package model

// Association code typo => label
// Les codes correspondent aux valeurs stockées en base dans :
// - ug.code_typo
// Nom raisonnablement court pour affichage
var TypoMap = map[string]string{
	"A": "Conifère - Petit Bois (PB)",
	"B": "Conifère - Bois Moyen (BM)",
	"C": "Cèdre",
	"D": "Douglas",
	"E": "Chênaie",
	"F": "Hêtraie",
	"G": "Châtaignier",
	"H": "Mélange",
	"I": "Boisement complet, mélange feuillus/résineux possible",
	"J": "Peuplement clair évoluant vers la forêt",
	"K": "Majoritairement conifère",
	"L": "Majoritairement feuillus",
	"M": "Majoritairement conifère",
	"N": "Majoritairement feuillus",
	"O": "Majoritairement conifère",
	"P": "Majoritairement conifère",
	"Q": "Zone improductive",
}

// Nom = nom complet venant du PSG
var TypoMap_long = map[string]string{
	"A": "Conifère - Petit Bois (PB)",
	"B": "Conifère - Bois Moyen (BM)",
	"C": "Cèdre",
	"D": "Douglas",
	"E": "Chênaie",
	"F": "Hêtraie",
	"G": "Châtaignier",
	"H": "Mélange",
	"I": "Boisement complet avec classes d'ages différenciées et mélange feuillus/résineux possible",
	"J": "Peuplement (feuillus ou résineux) clair évoluant vers la forêt mais avec un potentiel retour au sylvopastoralisme possible",
	"K": "Majoritairement conifère",
	"L": "Majoritairement feuillus",
	"M": "Majoritairement conifère",
	"N": "Majoritairement feuillus",
	"O": "Majoritairement conifère",
	"P": "Majoritairement conifère",
	"Q": "Pas de potentiel de production forestière ou pastoral (zones improductives : rochers, ...)",
}
