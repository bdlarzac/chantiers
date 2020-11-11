/******************************************************************************
    Modifie les acteurs de la base BDL.

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-01-09 11:38:50+01:00, Thierry Graff : Creation
********************************************************************************/
package fixture

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"fmt"
	"math/rand"
)

func AnonymizeActeurs() {
	ctx := ctxt.NewContext()
	table := "acteur"
	fmt.Println("MAJ " + table + " avec des données de test")
	// noms de famille les plus répandus en France
	noms := [...]string{
		"Martin",
		"Bernard",
		"Thomas",
		"Petit",
		"Robert",
		"Richard",
		"Durand",
		"Dubois",
		"Moreau",
		"Laurent",
	}
	// prénoms les plus répandus en France + qq sigles pour personnes morales
	prenoms := [...]string{
		"SARL",
		"GAEC",
		"S.A.",
		"Jean",
		"Pierre",
		"Michel",
		"André",
		"Philippe",
		"René",
		"Louis",
		"Alain",
		"Jacques",
		"Bernard",
		"Marie",
		"Jeanne",
		"Françoise",
		"Monique",
		"Catherine",
		"Nathalie",
		"Isabelle",
		"Jacqueline",
		"Anne",
		"Sylvie",
	}
	autres := model.Acteur{
		Adresse1:    "Le bourg",
		Adresse2:    "",
		Tel:         "01 02 03 04 05",
		Mobile: "06 07 08 09 10",
		Email:       "test@mail.org",
		Bic:         "",
		Iban:        "",
		Siret:       "",
		Notes:       "Données de test - ne correspond à aucune personne réelle",
	}

	acteurs, _ := model.SortedActeurs(ctx.DB, "id")

	for _, a := range acteurs {
		if a.Nom == "BDL" || a.Nom == "SCTL"  || a.Nom == "GFA" {
			continue
		}
		idxNom := rand.Intn(len(noms) - 1)
		idxPrenom := rand.Intn(len(prenoms) - 1)
		a.Nom = noms[idxNom]
		a.Prenom = prenoms[idxPrenom]
		a.Adresse1 = autres.Adresse1
		a.Tel = autres.Tel
		a.Mobile = autres.Mobile
		a.Email = autres.Email
		a.Notes = autres.Notes
		fmt.Printf("Anonymise %+v\n", a)
		err := model.UpdateActeur(ctx.DB, a)
		if err != nil {
			panic(err)
		}
	}
}
