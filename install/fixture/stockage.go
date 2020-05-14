/******************************************************************************
    Remplit la base avec des lieux de stockage de test
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
	fmt.Println("Remplit table stockage avec des données de test")
	//
	var err error
	var idStock, idLoyer int
	var nom string
	var montant float64
	var datedeb, datefin time.Time
	//
	nom = "Hangar des liquisses"
	idStock, err = model.InsertStockage(db, &model.Stockage{Nom: nom})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion stockage %d : %s\n", idStock, nom)
	//
	montant = 6000
	datedeb, _ = time.Parse("2006-01-02", "2018-01-01")
	datefin, _ = time.Parse("2006-01-02", "2050-01-01")
	idLoyer, err = model.InsertStockLoyer(db, &model.StockLoyer{
	    IdStockage: idStock,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion loyer %d : %d %.2f %s %s\n", idLoyer, idStock, montant, datedeb.Format("2006-01-02"), datefin.Format("2006-01-02"))
	//
	//
	//
	nom = "Hangar des Baumes"
	idStock, err = model.InsertStockage(db, &model.Stockage{Nom: nom})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion stockage %d : %s\n", idStock, nom)
	//
	montant = 0
	datedeb, _ = time.Parse("2006-01-02", "2018-01-01")
	datefin, _ = time.Parse("2006-01-02", "2050-01-01")
	idLoyer, err = model.InsertStockLoyer(db, &model.StockLoyer{
	    IdStockage: idStock,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion loyer %d : %d %.2f %s %s\n", idLoyer, idStock, montant, datedeb.Format("2006-01-02"), datefin.Format("2006-01-02"))
	/* 
	//
	montant = 3000
	datedeb, _ = time.Parse("2006-01-02", "2020-05-02")
	datefin, _ = time.Parse("2006-01-02", "2020-12-05")
	idLoyer, err = model.InsertStockLoyer(db, &model.StockLoyer{
	    IdStockage: idStock,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion loyer %d : %d %.2f %s %s\n", idLoyer, idStock, montant, datedeb.Format("2006-01-02"), datefin.Format("2006-01-02"))
	//
	montant = 4000
	datedeb, _ = time.Parse("2006-01-02", "2020-12-05")
	datefin, _ = time.Parse("2006-01-02", "2050-01-01")
	idLoyer, err = model.InsertStockLoyer(db, &model.StockLoyer{
	    IdStockage: idStock,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion loyer %d : %d %.2f %s %s\n", idLoyer, idStock, montant, datedeb.Format("2006-01-02"), datefin.Format("2006-01-02"))
	*/	
}
