/*
*****************************************************************************

	Code pour fabriquer des select html

	@copyright  BDL, Bois du Larzac
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2020-02-11 00:54:33+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/webo"
	"bdl.local/bdl/model"
	"strconv"
)

// Renvoie la liste des essences possibles
// dans un format utilisable par webo
func WeboEssence() []webo.OptionString {
	return []webo.OptionString{
		webo.OptionString{OptionValue: "CHOOSE_ESSENCE", OptionLabel: "--- Choisir ---"},
		webo.OptionString{OptionValue: "essence-PS", OptionLabel: model.LabelEssence("PS")},
		webo.OptionString{OptionValue: "essence-CN", OptionLabel: model.LabelEssence("CN")},
		webo.OptionString{OptionValue: "essence-CT", OptionLabel: model.LabelEssence("CT")},
		webo.OptionString{OptionValue: "essence-GV", OptionLabel: model.LabelEssence("GV")},
	}
}

// Renvoie la liste des unités possibles dans un Chautre
// dans un format utilisable par webo
// Utilisé uniquement dans le contrôleur de Chautre
func WeboChautreUnite() []webo.OptionString {
	return []webo.OptionString{
		webo.OptionString{OptionValue: "CHOOSE_UNITE", OptionLabel: "--- Choisir ---"},
		webo.OptionString{OptionValue: "unite-ST", OptionLabel: model.LabelUnite("ST")},
		webo.OptionString{OptionValue: "unite-TO", OptionLabel: model.LabelUnite("TO")},
		webo.OptionString{OptionValue: "unite-M3", OptionLabel: model.LabelUnite("M3")},
	}
}

// Renvoie la liste des unités possibles dans un Chaufer
// dans un format utilisable par webo
// Utilisé uniquement dans le contrôleur de Chaufer
func WeboChauferUnite() []webo.OptionString {
	return []webo.OptionString{
		webo.OptionString{OptionValue: "CHOOSE_UNITE", OptionLabel: "--- Choisir ---"},
		webo.OptionString{OptionValue: "unite-MA", OptionLabel: model.LabelUnite("MA")},
		webo.OptionString{OptionValue: "unite-ST", OptionLabel: model.LabelUnite("ST")},
	}
}

// Renvoie la liste des unités possibles dans un PlaqOp
// dans un format utilisable par webo
func WeboPlaqOpUnite() []webo.OptionString {
	return []webo.OptionString{
		webo.OptionString{OptionValue: "CHOOSE_UNITE", OptionLabel: "--- Choisir ---"},
		webo.OptionString{OptionValue: "unite-JO", OptionLabel: model.LabelUnite("JO")}, // jour
		webo.OptionString{OptionValue: "unite-HE", OptionLabel: model.LabelUnite("HE")}, // heure
		webo.OptionString{OptionValue: "unite-MA", OptionLabel: model.LabelUnite("MA")}, // map
		webo.OptionString{OptionValue: "unite-ST", OptionLabel: model.LabelUnite("ST")}, // stère
	}
}

// Renvoie la liste des unités possibles dans un PlaqOp
// dans un format utilisable par webo
func WeboTypeOp() []webo.OptionString {
	return []webo.OptionString{
		webo.OptionString{OptionValue: "CHOOSE_TYPEOP", OptionLabel: "--- Choisir ---"},
		webo.OptionString{OptionValue: "typeop-AB", OptionLabel: model.LabelActivite("AB")},
		webo.OptionString{OptionValue: "typeop-DB", OptionLabel: model.LabelActivite("DB")},
		webo.OptionString{OptionValue: "typeop-DC", OptionLabel: model.LabelActivite("DC")},
		webo.OptionString{OptionValue: "typeop-BR", OptionLabel: model.LabelActivite("BR")},
	}
}

// Renvoie la liste des types d'exploitation possibles
// dans un format utilisable par webo
func WeboExploitation() []webo.OptionString {
	return []webo.OptionString{
		webo.OptionString{OptionValue: "CHOOSE_EXPLOITATION", OptionLabel: "--- Choisir ---"},
		webo.OptionString{OptionValue: "exploitation-1", OptionLabel: model.LabelExploitation("1")},
		webo.OptionString{OptionValue: "exploitation-2", OptionLabel: model.LabelExploitation("2")},
		webo.OptionString{OptionValue: "exploitation-3", OptionLabel: model.LabelExploitation("3")},
		webo.OptionString{OptionValue: "exploitation-4", OptionLabel: model.LabelExploitation("4")},
		webo.OptionString{OptionValue: "exploitation-5", OptionLabel: model.LabelExploitation("5")},
	}
}

// Renvoie la liste des valorisations possibles
// dans un format utilisable par webo
// Utilisé uniquement dans le contrôleur de Chautre
// L'ordre des valorisations correspond à la demande de BDL
func WeboChautreValo() []webo.OptionString {
	return []webo.OptionString{
		webo.OptionString{OptionValue: "CHOOSE_VALORISATION", OptionLabel: "--- Choisir ---"},
		webo.OptionString{OptionValue: "valorisation-PP", OptionLabel: model.LabelValorisation("PP")},
		webo.OptionString{OptionValue: "valorisation-CH", OptionLabel: model.LabelValorisation("CH")},
		webo.OptionString{OptionValue: "valorisation-PL", OptionLabel: model.LabelValorisation("PL")},
		webo.OptionString{OptionValue: "valorisation-PI", OptionLabel: model.LabelValorisation("PI")},
		webo.OptionString{OptionValue: "valorisation-BO", OptionLabel: model.LabelValorisation("BO")},
	}
}

// Renvoie la liste des taux de TVA utilisés pour facturer un chantier "autres valorisations"
// dans un format utilisable par webo
// @param  chooseId     Chaîne utilisée pour désigner l'id et la value de l'option "---Choisir ---"
//
//	Permet d'avoir plusieurs formulaires de choix de taux de TVA dans un même form
//
// @param  idPrefix     l'attribut "id" de chaque option sera = idPrefix suivi de la valeur de l'option
//
//	Permet que chaque option soit unique dans tous les formulaires de TVA
func WeboChautreTVA(ctx *ctxt.Context, chooseId, idPrefix string) []webo.OptionString {
	res := []webo.OptionString{}
	res = append(res, webo.OptionString{OptionValue: chooseId, OptionId: idPrefix + chooseId, OptionLabel: "--- Choisir ---"})
	for _, taux := range ctx.Config.TVABDL.AutresValorisations {
		tmp := strconv.FormatFloat(taux, 'f', -1, 64)
		res = append(res, webo.OptionString{OptionValue: tmp, OptionId: idPrefix + tmp, OptionLabel: tmp + " %"})
	}
	return res
}

// Renvoie la liste des granulométries possibles
// dans un format utilisable par webo
func WeboGranulo() []webo.OptionString {
	return []webo.OptionString{
		webo.OptionString{OptionValue: "CHOOSE_GRANULO", OptionLabel: "--- Choisir ---"},
		webo.OptionString{OptionValue: "granulo-P16", OptionLabel: model.LabelGranulo("P16")},
		webo.OptionString{OptionValue: "granulo-P45", OptionLabel: model.LabelGranulo("P45")},
	}
}

// Renvoie la liste des types de frais possibles (pour lieux de stockage)
// dans un format utilisable par webo
func WeboStockFrais() []webo.OptionString {
	return []webo.OptionString{
		webo.OptionString{OptionValue: "CHOOSE_STOCKFRAIS", OptionLabel: "--- Choisir ---"},
		webo.OptionString{OptionValue: "stockfrais-AS", OptionLabel: model.LabelStockFrais("AS")},
		webo.OptionString{OptionValue: "stockfrais-EL", OptionLabel: model.LabelStockFrais("EL")},
		webo.OptionString{OptionValue: "stockfrais-LO", OptionLabel: model.LabelStockFrais("LO")},
	}
}

// Renvoie la liste des taux de TVA utilisés pour payer un intervenant extérieur
// dans un format utilisable par webo
// @param  chooseId     Chaîne utilisée pour désigner l'id et la value de l'option "---Choisir ---"
//
//	Permet d'avoir plusieurs formulaires de choix de taux de TVA dans un même form
//
// @param  idPrefix     l'attribut "id" de chaque option sera = idPrefix suivi de la valeur de l'option
//
//	Permet que chaque option soit unique dans tous les formulaires de TVA
func WeboTVAExt(ctx *ctxt.Context, chooseId, idPrefix string) []webo.OptionString {
	res := []webo.OptionString{}
	res = append(res, webo.OptionString{OptionValue: chooseId, OptionId: idPrefix + chooseId, OptionLabel: "--- Choisir ---"})
	for _, taux := range ctx.Config.TVAExt {
		tmp := strconv.FormatFloat(taux, 'f', 1, 64)
		res = append(res, webo.OptionString{OptionValue: tmp, OptionId: idPrefix + tmp, OptionLabel: tmp + " %"})
	}
	return res
}

// Renvoie la liste des fournisseurs de plaquettes
// dans un format utilisable par webo
func WeboFournisseur(ctx *ctxt.Context) []webo.OptionString {
	res := []webo.OptionString{}
	fournisseurs, err := model.GetFournisseurs(ctx.DB)
	if err != nil {
		panic("La base de donnée doit contenir au moins un fournisseur de plaquettes")
	}
	res = append(res, webo.OptionString{OptionValue: "CHOOSE_FOURNISSEUR", OptionLabel: "--- Choisir ---"})
	for _, fournisseur := range fournisseurs {
		res = append(res, webo.OptionString{OptionValue: "fournisseur-" + strconv.Itoa(fournisseur.Id), OptionLabel: fournisseur.Nom})
	}
	return res
}

// Renvoie la liste des clients plaquettes
// dans un format utilisable par webo
func WeboClientsPlaquettes(ctx *ctxt.Context) []webo.OptionString {
	res := []webo.OptionString{}
	clients, err := model.GetClientsPlaquettes(ctx.DB)
	if err != nil {
		panic(err)
	}
	res = append(res, webo.OptionString{OptionValue: "CHOOSE_CLIENT", OptionLabel: "--- Choisir ---"})
	for _, client := range clients {
		res = append(res, webo.OptionString{OptionValue: "acteur-" + strconv.Itoa(client.Id), OptionLabel: client.Nom})
	}
	return res
}

// Renvoie la liste des tas actifs
// dans un format utilisable par webo
// les attributs value et id des options sont de la forme tas-<id du tas> (ex: tas-3)
func WeboTas(ctx *ctxt.Context) ([]webo.OptionString, error) {
	res := []webo.OptionString{}
	tas, err := model.GetAllTasActifsFull(ctx.DB)
	if err != nil {
		return res, err
	}
	res = append(res, webo.OptionString{OptionValue: "CHOOSE_TAS", OptionLabel: "--- Choisir ---"})
	for _, t := range tas {
		res = append(res, webo.OptionString{OptionValue: "tas-" + strconv.Itoa(t.Id), OptionLabel: t.Nom})
	}
	return res, nil
}

// Renvoie la liste des fermiers SCTL
// dans un format utilisable par webo
func WeboFermier(ctx *ctxt.Context) ([]webo.OptionString, error) {
	res := []webo.OptionString{}
	fermiers, err := model.GetSortedFermiers(ctx.DB, "nom")
	if err != nil {
		return res, err
	}
	res = append(res, webo.OptionString{OptionValue: "CHOOSE_FERMIER", OptionLabel: "--- Choisir ---"})
	for _, f := range fermiers {
		res = append(res, webo.OptionString{OptionValue: "fermier-" + strconv.Itoa(f.Id), OptionLabel: f.String()})
	}
	return res, nil
}
