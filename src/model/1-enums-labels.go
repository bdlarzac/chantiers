/*
*****************************************************************************

	Enums manipulées par le programme
	En théorie, les fonctions All*Codes() devraient
	renvoyer des valeurs tirées de la base avec des queries du type :
	SELECT unnest(enum_range(NULL::my_type));

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

// ************************** Essences (= espèces d'arbres) *******************************
// cf type postgres typessence

// @return  Tous les codes des essences
func AllEssenceCodes() []string {
	return []string{"PS", "CN", "CT", "GV"}
}

// Labels des essences  manipulées par l'appli
// @param code : essence telle que stockée en base
func LabelEssence(code string) string {
	switch code {
	case "PS":
		return "Pin sylvestre"
	case "CN":
		return "Chêne"
	case "CT":
		return "Châtaigner"
	case "GV":
		return "Genévrier"
	}
	return "??? BUG LabelEssence ("+code+") ???"
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
