/******************************************************************************
    Initialise le hangar des liquisses
    Code pas utilisé en fonctionnement normal.

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-12-05 16:54:53+01:00, Thierry Graff : Creation
    @history    2020-09-29 14:46:44+02:00, Thierry Graff : Change de fixture à initialize
********************************************************************************/
package initialize

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"time"
	"fmt"
)

// *********************************************************
func FillHangarLiquisses() {
	ctx := ctxt.NewContext()
	db := ctx.DB
	fmt.Println("Crée le hangar des liquisses")
	//
	var err error
	var idStock, idFrais int
	var typefrais string
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
	typefrais = "LO"
	montant = 6000
	datedeb, _ = time.Parse("2006-01-02", "2018-01-01")
	datefin, _ = time.Parse("2006-01-02", "2050-01-01")
	idFrais, err = model.InsertStockFrais(db, &model.StockFrais{
	    IdStockage: idStock,
	    TypeFrais: typefrais,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion loyer %d : %d %.2f %s %s\n", idFrais, idStock, montant, datedeb.Format("2006-01-02"), datefin.Format("2006-01-02"))
	//
	//
	//
	/* 
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
	idFrais, err = model.InsertStockLoyer(db, &model.StockLoyer{
	    IdStockage: idStock,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion loyer %d : %d %.2f %s %s\n", idFrais, idStock, montant, datedeb.Format("2006-01-02"), datefin.Format("2006-01-02"))
	
	//
	montant = 3000
	datedeb, _ = time.Parse("2006-01-02", "2020-05-02")
	datefin, _ = time.Parse("2006-01-02", "2020-12-05")
	idFrais, err = model.InsertStockLoyer(db, &model.StockLoyer{
	    IdStockage: idStock,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion loyer %d : %d %.2f %s %s\n", idFrais, idStock, montant, datedeb.Format("2006-01-02"), datefin.Format("2006-01-02"))
	//
	montant = 4000
	datedeb, _ = time.Parse("2006-01-02", "2020-12-05")
	datefin, _ = time.Parse("2006-01-02", "2050-01-01")
	idFrais, err = model.InsertStockLoyer(db, &model.StockLoyer{
	    IdStockage: idStock,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion loyer %d : %d %.2f %s %s\n", idFrais, idStock, montant, datedeb.Format("2006-01-02"), datefin.Format("2006-01-02"))
	*/	
}
