/*
*****************************************************************************

		Vente générique
	    - vente de plaquettes forestières
	    - vente issue d'un chantier autres valorisations

		@copyright  BDL, Bois du Larzac.
		@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
		@history    2023-05-03 08:14:58+02:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

type Vente struct {
	Id          int
	TypeVente   string // "plaq" ou "autre"
	Titre       string
	URL         string // Chaîne vide ou URL du détail de l'entité, ex "/plaq/32"
	DateVente   time.Time
	TypeValo    string
	Volume      float64
	Unite       string // pour le volume
	CodeEssence string
	PrixHT      float64
	PUHT        float64
	TVA         float64
	NumFacture  string
	DateFacture time.Time
	Notes       string
	//
	Details interface{}
}

// ************************** Nom *******************************

func (v *Vente) String() string {
	return v.Titre
}

// ************************** Get many *******************************

func ComputeVentesFromFiltres(db *sqlx.DB, filtres map[string][]string) (result []*Vente, err error) {
	fmt.Printf("=== model.ComputeVentesFromFiltres() - filtres = %+v\n", filtres)
	result = []*Vente{}
	//
	// 1 - détermine dans quelles tables rechercher à partir du filtre valorisations
	// on pourrait déterminer plus précisément à partir du filtre clients (avec les rôles)
	// mais compliqué sans gagner en efficacité
	//
	tables := []string{}
	if len(filtres["valo"]) == 0 {
		tables = []string{"venteplaq", "chautre"}
	} else {
		if tiglib.InArray("PQ", filtres["valo"]) {
			tables = append(tables, "venteplaq")
			if len(filtres["valo"]) > 1 {
				tables = append(tables, "chautre")
			}
		} else {
			tables = append(tables, "chautre") // car len(filtres["valo"]) > 0
		}
	}
	fmt.Printf("tables = %v\n", tables)
	//
	// 2 - Filtres periode et client (et valo pour chautre)
	//
	var tmp []*Vente
	if tiglib.InArray("venteplaq", tables) {
		tmp, err = computeVentePlaqVenteFromFiltresPeriodeEtClient(db, filtres["periode"], filtres["client"])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel computePlaqFromFiltrePeriode()")
		}
		result = append(result, tmp...)
	}
	if tiglib.InArray("chautre", tables) {
		tmp, err = computeChautreVentreFromFiltresPeriodeEtClientEtValo(db, filtres["periode"], filtres["client"], filtres["valo"])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel computePlaqFromFiltrePeriode()")
		}
		result = append(result, tmp...)
	}
	//
	// Filtres suivants
	//
	if len(filtres["proprio"]) != 0 {
		result, err = filtreVente_proprio(db, result, filtres["proprio"])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel filtreVente_proprio()")
		}
	}
	//fmt.Printf("result = %+v\n",result)
	// TODO éventuellement, trier par date
	return result, nil
}

// ****************************************************************************************************
// ************************** Auxiliaires de ComputeVentesFromFiltres() *******************************
// ****************************************************************************************************

// ************************** Selection par période et client et valo *******************************

/*  Fabrique des Ventes à partir de la table venteplaq */
func computeVentePlaqVenteFromFiltresPeriodeEtClient(db *sqlx.DB, filtrePeriode, filtreClient []string) (result []*Vente, err error) {
	ventePlaqs := []*VentePlaq{}
	and := []string{}
	var args []interface{}
	if len(filtreClient) == 1 {
		and = append(and, "id_client=?")
		args = append(args, filtreClient[0])
	}
	if len(filtrePeriode) == 2 {
		and = append(and, "datevente >= ? and datevente <= ?")
		args = append(args, filtrePeriode[0])
		args = append(args, filtrePeriode[1])
	}
	//
	query := "select * from venteplaq"
	if len(and) > 0 {
		query += " where " + strings.Join(and, " and ")
		query = db.Rebind(query) // transforme les ? en $1, $2 etc.
	}
	err = db.Select(&ventePlaqs, query, args...)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, vp := range ventePlaqs {
		vente, err := ventePlaq2Vente(db, vp)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel ventePlaq2Vente()")
		}
		result = append(result, vente)
	}
	return result, nil
}

/*  Fabrique des Ventes à partir de la table chautre */
func computeChautreVentreFromFiltresPeriodeEtClientEtValo(db *sqlx.DB, filtrePeriode, filtreClient, filtreValo []string) (result []*Vente, err error) {
	chantiers := []*Chautre{}
	and := []string{}
	var args []interface{}
	if len(filtreClient) == 1 {
		and = append(and, "id_acheteur=?")
		args = append(args, filtreClient[0])
	}
	if len(filtrePeriode) == 2 {
		and = append(and, "datecontrat >= ? and datecontrat <= ?")
		args = append(args, filtrePeriode[0])
		args = append(args, filtrePeriode[1])
	}
	if len(filtreValo) > 0 {
		inClause := "typevalo in('" + strings.Join(filtreValo, "','") + "')"
		and = append(and, inClause)
	}
	//
	query := "select * from chautre"
	if len(and) > 0 {
		query += " where " + strings.Join(and, " and ")
		query = db.Rebind(query) // transforme les ? en $1, $2 etc.
	}
	err = db.Select(&chantiers, query, args...)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, ch := range chantiers {
		vente, err := chautre2Vente(db, ch)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel chautre2Vente()")
		}
		result = append(result, vente)
	}
	//
	return result, nil
}

// ************************** Application des filtres *******************************
// En entrée : liste de ventes
// En sortie : liste de ventes qui satisfont au filtre

// TODO implement
func filtreVente_proprio(db *sqlx.DB, input []*Vente, filtre []string) (result []*Vente, err error) {
	result = []*Vente{}
	/*
		for _, a := range input {
			for _, f := range filtre {
				id, _ := strconv.Atoi(f)
				for _, lienParcelle := range a.LiensParcelles {
					parcelle, err := GetParcelle(db, lienParcelle.IdParcelle)
					if err != nil {
						return result, werr.Wrapf(err, "Erreur appel GetParcelle()")
					}
					if parcelle.IdProprietaire == id {
						result = append(result, a)
						break
					}
				}
			}
		}
	*/
	return result, nil
}

// ************************** Conversion de struct vers une Vente *******************************

func ventePlaq2Vente(db *sqlx.DB, vp *VentePlaq) (v *Vente, err error) {
	v = &Vente{}
	err = vp.ComputeClient(db) // Obligatoire pour pouvoir utiliser String()
	if err != nil {
		return v, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeClient()")
	}
	v.Id = vp.Id
	v.TypeVente = "plaq"
	v.Titre = vp.StringSansDate()
	v.URL = "/vente/" + strconv.Itoa(vp.Id)
	v.DateVente = vp.DateVente
	v.TypeValo = "PQ"
	err = vp.ComputeQte(db)
	if err != nil {
		return v, werr.Wrapf(err, "Erreur appel VentePlaq.ComputeQte()")
	}
	v.Volume = vp.Qte
	v.Unite = "MA"
	v.CodeEssence = "PS"
	v.PrixHT = vp.PUHT * vp.Qte
	v.PUHT = vp.PUHT
	v.TVA = vp.TVA
	v.NumFacture = vp.NumFacture
	v.DateFacture = vp.DateFacture
	v.Notes = vp.Notes
	return v, nil
}

func chautre2Vente(db *sqlx.DB, ch *Chautre) (v *Vente, err error) {
	v = &Vente{}
	v.Id = ch.Id
	v.TypeVente = "autre"
	v.Titre = ch.Titre
	//	v.URL = "/chantier/autre/" + strconv.Itoa(idChautre)
	v.DateVente = ch.DateContrat
	v.TypeValo = ch.TypeValo
	v.Volume = ch.VolumeRealise
	v.Unite = ch.Unite
	v.CodeEssence = ch.Essence
	v.PrixHT = ch.PUHT * ch.VolumeRealise
	v.PUHT = ch.PUHT
	v.TVA = ch.TVA
	v.NumFacture = ch.NumFacture
	v.DateFacture = ch.DateFacture
	v.Notes = ch.Notes
	return v, nil
}
