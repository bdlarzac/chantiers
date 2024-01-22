/*
Vente générique
- vente de plaquettes forestières
- vente issue d'un chantier autres valorisations

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2023-05-03 08:14:58+02:00, Thierry Graff : Creation
*/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"sort"
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
	NumFacture  string
	DateFacture time.Time
	Notes       string
	//
	Details interface{}
	// relations n-n - utiles pour l'application de certains filtres
	LiensParcelles []*ChantierParcelle
}

// ************************** Nom *******************************

func (v *Vente) String() string {
	return v.Titre
}

// ************************** Instance methods *******************************

// Rajouté pour issue #24 (implémenter filtre proprio dans bilans vente)
// Un peu bidouille pour utiliser le code de bilan activité
func (v *Vente) ComputeLiensParcelles(db *sqlx.DB) (err error) {
	if len(v.LiensParcelles) != 0 {
		return nil // déjà calculé
	}
	// Cas simple, une vente "autre" est forcément un chantier autre valorisation
	if v.TypeVente == "autre" {
		v.LiensParcelles, err = computeLiensParcellesOfChantier(db, "chautre", v.Id)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel computeLiensParcellesOfChantier()")
		}
		return nil
	}
	// Pour les ventes plaquettes, il faut calculer les liens parcelles de tous les chantiers
	// liés à cette vente (plusieurs possibles si la vente correspond à des chargements sur
	// différents tas venant de différents chantiers).
	vp, err := GetVentePlaq(db, v.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetVentePlaq")
	}
	err = vp.ComputeChantiers(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel VentePlaq.ComputeChantiers()")
	}
	for _, ch := range vp.Chantiers {
		liensParcelles, err := computeLiensParcellesOfChantier(db, "plaq", ch.Id)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel computeLiensParcellesOfChantier()")
		}
		v.LiensParcelles = append(v.LiensParcelles, liensParcelles...)
	}
	return nil
}

// ************************** Get many *******************************

func ComputeVentesFromFiltres(db *sqlx.DB, filtres map[string][]string) (res []*Vente, err error) {
	res = []*Vente{}
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
	//
	// 2 - Filtres periode et client (et valo pour chautre)
	//
	var tmp []*Vente
	if tiglib.InArray("venteplaq", tables) {
		tmp, err = computeVentePlaqVenteFromFiltresPeriodeEtClient(db, filtres["periode"], filtres["client"])
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel computePlaqFromFiltrePeriode()")
		}
		res = append(res, tmp...)
	}
	if tiglib.InArray("chautre", tables) {
		tmp, err = computeChautreVentreFromFiltresPeriodeEtClientEtValo(db, filtres["periode"], filtres["client"], filtres["valo"])
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel computePlaqFromFiltrePeriode()")
		}
		res = append(res, tmp...)
	}
	//
	// Filtres suivants
	//
	if len(filtres["proprio"]) != 0 {
		res, err = filtreVente_proprio(db, res, filtres["proprio"])
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel filtreVente_proprio()")
		}
	}
	// Tri par date
	sortedRes := make(venteSlice, 0, len(res))
	for _, elt := range res {
		sortedRes = append(sortedRes, elt)
	}
	sort.Sort(sortedRes)
	return sortedRes, nil
}

// ****************************************************************************************************
// ************************** Auxiliaires de ComputeVentesFromFiltres() *******************************
// ****************************************************************************************************

// ************************** Auxiliaire pour trier par date *******************************

type venteSlice []*Vente

func (p venteSlice) Len() int           { return len(p) }
func (p venteSlice) Less(i, j int) bool { return p[i].DateVente.Before(p[j].DateVente) }
func (p venteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// ************************** Selection par période et client et valo *******************************

// Fabrique des Ventes à partir de la table venteplaq
func computeVentePlaqVenteFromFiltresPeriodeEtClient(db *sqlx.DB, filtrePeriode, filtreClient []string) (res []*Vente, err error) {
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
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, vp := range ventePlaqs {
		vente, err := ventePlaq2Vente(db, vp)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel ventePlaq2Vente()")
		}
		res = append(res, vente)
	}
	return res, nil
}

// Fabrique des Ventes à partir de la table chautre
func computeChautreVentreFromFiltresPeriodeEtClientEtValo(db *sqlx.DB, filtrePeriode, filtreClient, filtreValo []string) (res []*Vente, err error) {
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
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, ch := range chantiers {
		vente, err := chautre2Vente(db, ch)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel chautre2Vente()")
		}
		res = append(res, vente)
	}
	//
	return res, nil
}

// ************************** Application des filtres *******************************
// En entrée : liste de ventes
// En sortie : liste de ventes qui satisfont au filtre

func filtreVente_proprio(db *sqlx.DB, input []*Vente, filtre []string) (res []*Vente, err error) {
	res = []*Vente{}
	for _, vente := range input {
		err = vente.ComputeLiensParcelles(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel vente.ComputeLiensParcelles()")
		}
		for _, f := range filtre {
			idProprio, _ := strconv.Atoi(f)
			for _, lienParcelle := range vente.LiensParcelles {
				parcelle, err := GetParcelle(db, lienParcelle.IdParcelle)
				if err != nil {
					return res, werr.Wrapf(err, "Erreur appel GetParcelle()")
				}
				if parcelle.IdProprietaire == idProprio {
					res = append(res, vente)
					break
				}
			}
		}
	}
	return res, nil
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
	v.URL = "/chantier/autre/" + strconv.Itoa(ch.Id)
	v.DateVente = ch.DateContrat
	v.TypeValo = ch.TypeValo
	v.Volume = ch.VolumeRealise
	v.Unite = ch.Unite
	v.CodeEssence = ch.Essence
	v.PrixHT = ch.PUHT * ch.VolumeRealise
	v.PUHT = ch.PUHT
	v.NumFacture = ch.NumFacture
	v.DateFacture = ch.DateFacture
	v.Notes = ch.Notes
	return v, nil
}
