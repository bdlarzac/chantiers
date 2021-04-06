/******************************************************************************
    Enums manipulées par le programme
    En théorie, les fonctions All*Codes() devraient
    renvoyer des valeurs tirées de la base avec des queries du type :
    SELECT unnest(enum_range(NULL::my_type));

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-02-11 17:27:04+01:00, Thierry Graff : Creation
********************************************************************************/
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
	return "??? BUG LabelActivite ???"
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
	return "??? BUG LabelEssence ???"
}

// ************************** Unités *******************************
// cf type postgres typeunite

// Renvoie le code de l'unité correspondant à une valorisation, tel que stocké en base
// Comprend les valos pour chautre + CF pour chauffage fermier
func Valorisation2unite(codeValo string) string {
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
		return "ST"
	case "BO":
		return "M3"
	}
	return "??? Code inconnu dans Valorisation2unite() ???"
}

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
	case "ST":
		return "stères"
	case "TO":
		return "tonnes"
	}
	return "??? BUG LabelUnite ???"
}

// Labels des unités manipulées par l'appli
// @param abbrev : unité telle que stockée en base
func LabelUniteHTML(abbrev string) string {
	if abbrev == "M3" {
		return "m<sup>3</sup>"
	}
	return LabelUnite(abbrev)
}

// ************************** Chantiers autres valorisations *******************************
// cf type postgres typevalo

// L'ordre des valorisations correspond à la demande de BDL
// @return  Tous les codes des valorisations de chautre
// Correspond au type typevalo en base
func AllValorisationCodes() []string {
	return []string{"PP", "CH", "PL", "PI", "BO"}
}

// Pareil que AllValorisationCodes(), avec CF en plus (chauffage fermier)
func AllValorisationCodesAvecChaufer() []string {
	return []string{"PP", "CH", "CF", "PL", "PI", "BO"}
}

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
	return "??? BUG LabelValorisation ???"
}

// Labels du type de valorisation
// Utilisé par bilans, où "Chauffage" doit être séparé en 2 :
// chauffage fermier ( de chaufer) et chauffage client (de chautre)
func LabelValorisationAvecChaufer(abbrev string) string {
	switch abbrev {
	case "CF":
		return "Chauffage fermier"
	case "CH":
		return "Chauffage client"
	}
	return LabelValorisation(abbrev)
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
	return "??? BUG LabelTypeVente ???"
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
	return "??? BUG LabelExploitation ???"
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
	return "??? BUG LabelStockFrais ???"
}
