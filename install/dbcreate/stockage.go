/******************************************************************************
    Initialise le hangar des liquisses
    Code pas utilisé en fonctionnement normal.

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-12-05 16:54:53+01:00, Thierry Graff : Creation
    @history    2020-09-29 14:46:44+02:00, Thierry Graff : Change de fixture à initialize
********************************************************************************/
package dbcreate

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"time"
	"fmt"
)

// *********************************************************
func FillHangarsInitiaux() {
	ctx := ctxt.NewContext()
	db := ctx.DB
	//
	fmt.Println("Crée le hangar des Liquisses")
	//
	var err error
	var idStock int
	var typefrais string
	var nom, notes string
	var montant float64
	var datedeb, datefin time.Time
	//
	nom = "Hangar des Liquisses"
	idStock, err = model.InsertStockage(db, &model.Stockage{Nom: nom})
	if err != nil {
		panic(err)
	}
	
	// Loyer 6000 E / an
	typefrais = "LO"
	montant = 6000 * 3
	notes = "6000 E / an"
	datedeb, _ = time.Parse("2006-01-02", "2018-01-01")
	datefin, _ = time.Parse("2006-01-02", "2020-12-31")
	_, err = model.InsertStockFrais(db, &model.StockFrais{
	    IdStockage: idStock,
	    TypeFrais: typefrais,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	    Notes: notes,
	})
	if err != nil {
		panic(err)
	}
	
	// Assurance 700 E tous les 6 mois
	typefrais = "AS"
	montant = 700 * 6                                                                       
	notes = "700 E tous les 6 mois"
	datedeb, _ = time.Parse("2006-01-02", "2018-01-01")
	datefin, _ = time.Parse("2006-01-02", "2020-12-31")
	_, err = model.InsertStockFrais(db, &model.StockFrais{
	    IdStockage: idStock,
	    TypeFrais: typefrais,
	    Montant: montant,
	    DateDebut: datedeb,
	    DateFin: datefin,
	    Notes: notes,
	})
	if err != nil {
		panic(err)
	}
	//
	fmt.Println("Crée le hangar des Baumes")
	//
	nom = "Hangar des Baumes"
	idStock, err = model.InsertStockage(db, &model.Stockage{Nom: nom})
	if err != nil {
		panic(err)
	}
}
