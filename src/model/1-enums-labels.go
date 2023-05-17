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
// @param abbrev : type d'opération tel que stocké en base
func LabelActivite(abbrev string) string {
	switch abbrev {
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
	return "??? BUG LabelActivite ("+abbrev+") ???"
}

// ************************** Essences (= espèces d'arbres) *******************************
// cf type postgres typessence

// @return  Tous les codes des essences
func AllEssenceCodes() []string {
	return []string{"PS", "CN", "CT", "GV"}
}

// Labels des essences  manipulées par l'appli
// @param abbrev : essence telle que stockée en base
func LabelEssence(abbrev string) string {
	switch abbrev {
	case "PS":
		return "Pin sylvestre"
	case "CN":
		return "Chêne"
	case "CT":
		return "Châtaigner"
	case "GV":
		return "Genévrier"
	}
	return "??? BUG LabelEssence ("+abbrev+") ???"
}

// Labels du type de vente (pour Chautre)
// @param abbrev : type de vente tel que stockée en base
func LabelTypeVente(abbrev string) string {
	switch abbrev {
	case "NON":
		return "Non spécifié"
	case "BSP":
		return "Bois-sur-pied"
	case "BDR":
		return "Bord-de-route"
	case "LIV":
		return "Livré"
	}
	return "??? BUG LabelTypeVente ("+abbrev+") ???"
}

// ************************** Type d'exploitation *******************************
// cf type postgres typexploitation

// Labels du type d'exploitation
// @param abbrev : type d'exploitation
func LabelExploitation(abbrev string) string {
	switch abbrev {
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
	return "??? BUG LabelExploitation ("+abbrev+") ???"
}

// ************************** Granulométrie *******************************
// cf type postgres typegranulo

// Labels du type de granulométrie
// @param abbrev : type de granulométrie
func LabelGranulo(abbrev string) string {
	return abbrev // le type stocké en base correspond au label
}

// ************************** Frais des lieux de stockage *******************************
// cf type postgres typestockfrais

// Labels du type de frais (pour lieu de stockage)
// @param abbrev : type de frais tel que stockée en base
func LabelStockFrais(abbrev string) string {
	switch abbrev {
	case "AS":
		return "Assurance"
	case "EL":
		return "Electricité"
	case "LO":
		return "Loyer"
	}
	return "??? BUG LabelStockFrais ("+abbrev+") ???"
}
