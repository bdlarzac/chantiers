/******************************************************************************
    Initialise un hangar de test
    Code pas utilisé en fonctionnement normal.

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-12-05 16:54:53+01:00, Thierry Graff : Creation
********************************************************************************/
package fixture

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"time"
	"fmt"
)

// *********************************************************
func FillStockage() {
	ctx := ctxt.NewContext()
	db := ctx.DB
	fmt.Println("Crée le hangar de test")
	//
	var err error
	var idStock int
	var typefrais string
	var nom string
	var montant float64
	var datedeb, datefin time.Time
	//
	nom = "Hangar de test"
	idStock, err = model.InsertStockage(db, &model.Stockage{Nom: nom})
	if err != nil {
		panic(err)
	}
	
	// Loyer 1000 E / an
	typefrais = "LO"
	montant = 1000 * 3
	datedeb, _ = time.Parse("2006-01-02", "2018-01-01")
	datefin, _ = time.Parse("2006-01-02", "2020-12-31")
	_, err = model.InsertStockFrais(db, &model.StockFrais{
	    IdStockage: idStock,
	    TypeFrais: typefrais,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
	
	// Assurance 100 E / mois
	typefrais = "AS"
	montant = 100 * 36
	datedeb, _ = time.Parse("2006-01-02", "2018-01-01")
	datefin, _ = time.Parse("2006-01-02", "2020-12-31")
	_, err = model.InsertStockFrais(db, &model.StockFrais{
	    IdStockage: idStock,
	    TypeFrais: typefrais,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
	
	// Elec 52.5 E tous les 2 mois
	typefrais = "AS"
	montant = 52.5 * 18
	datedeb, _ = time.Parse("2006-01-02", "2018-01-01")
	datefin, _ = time.Parse("2006-01-02", "2020-12-31")
	_, err = model.InsertStockFrais(db, &model.StockFrais{
	    IdStockage: idStock,
	    TypeFrais: typefrais,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
}
