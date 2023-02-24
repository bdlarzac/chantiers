/*
*****************************************************************************

	Affacture = "facture à l'envers", que BDL doit payer à un acteur
	Les affacures ne sont pas stockées en base.

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2020-03-03 11:16:58+01:00, Thierry Graff : Creation

*******************************************************************************
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

// Contient les données nécessaires pour afficher le PDF
type Affacture struct {
	// Renseigné via le formulaire
	IdActeur       int
	DateDebut      time.Time
	DateFin        time.Time
	TypesActivites []string
	// Calculé à partir de la BDD
	Items    []*AffactureItem
	TotalHT  float64
	TotalTTC float64
}

// Les données hétérogènes du model (PlaqOp, PlaqTrans...) sont traduites dans un vocabulaire commun pour l'affichage
// 1 AffactureItem <-> 1 activité (1 ligne ds la BDD)
type AffactureItem struct {
	Titre    string
	Date     time.Time
	Lignes   []AffactureLigne
	TotalHT  float64
	TotalTTC float64
}

// La plupart des AffactureItem sont constitués d'une seule AffactureLigne.
// Plusieurs AffactureLigne dans le cas de détails :
// par ex un transport peut être constitué de 2 lignes : une pour la m.o. et une pour le camion
type AffactureLigne struct {
	Titre    string
	Colonnes []AffactureColonne
}

// Chaque AffactureLigne est constituée de plusieurs AffactureColonne
type AffactureColonne struct {
	Titre  string
	Valeur string
}

func (aff *Affacture) ComputeItems(db *sqlx.DB) (err error) {
	for _, typeActivite := range aff.TypesActivites {
		switch typeActivite {
		// Opérations simples
		case "AB", "DB", "DC", "BR":
			err = aff.computeItemsOperationSimple(db, typeActivite)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsOperationSimple()")
			}
			// Transport
		case "TR":
			err = aff.computeItemsTransportGlobal(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsTransportGlobal()")
			}
		case "TR-CO":
			err = aff.computeItemsTransportConducteur(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsTransportConducteur()")
			}
		case "TR-OU":
			err = aff.computeItemsTransportProprioutil(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsTransportProprioutil()")
			}
		// Rangement
		case "RG":
			err = aff.computeItemsRangementGlobal(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsRangementGlobal()")
			}
		case "RG-CO":
			err = aff.computeItemsRangementConducteur(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsRangementConducteur()")
			}
		case "RG-OU":
			err = aff.computeItemsRangementProprioutil(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsRangementProprioutil()")
			}
			// Chargement
		case "CG":
			err = aff.computeItemsChargementGlobal(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsChargementGlobal()")
			}
		case "CG-CO":
			err = aff.computeItemsChargementConducteur(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsChargementConducteur()")
			}
		case "CG-OU":
			err = aff.computeItemsChargementProprioutil(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsChargementOutil()")
			}
		// Livraison
		case "LV":
			err = aff.computeItemsLivraisonGlobal(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsLivraisonGlobal()")
			}
		case "LV-CO":
			err = aff.computeItemsLivraisonConducteur(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsLivraisonConducteur()")
			}
		case "LV-OU":
			err = aff.computeItemsLivraisonOutil(db)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Affacture.computeItemsLivraisonOutil()")
			}
		}
	}
	// tri par date
	sortedItems := make(affactureItemSlice, 0, len(aff.Items))
	for _, elt := range aff.Items {
		sortedItems = append(sortedItems, elt)
	}
	sort.Sort(sortedItems)
	aff.Items = sortedItems
	//
	return nil
}

// Auxiliaires de Affacture.ComputeItems() pour trier par date
type affactureItemSlice []*AffactureItem

func (p affactureItemSlice) Len() int           { return len(p) }
func (p affactureItemSlice) Less(i, j int) bool { return p[i].Date.Before(p[j].Date) }
func (p affactureItemSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (aff *Affacture) computeItemsOperationSimple(db *sqlx.DB, typeActivite string) (err error) {
	list := []PlaqOp{}
	query := "select * from plaqop where id_acteur=$1 and typop=$2 and datedeb>=$3 and datedeb<=$4"
	err = db.Select(&list, query, aff.IdActeur, typeActivite, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		montantHT = elt.Qte * elt.PUHT
		montantTVA = montantHT * elt.TVA / 100
		montantTTC = montantHT + montantTVA
		item := AffactureItem{
			Titre: LabelActivite(typeActivite),
			Date:  elt.DateDebut,
		}
		ligne = AffactureLigne{
			Titre: "M.O. " + LabelActivite(typeActivite),
			Colonnes: []AffactureColonne{
				{
					Titre:  "Nb " + LabelUnite(elt.Unite),
					Valeur: strconv.FormatFloat(elt.Qte, 'f', 2, 64),
				},
				{
					Titre:  "Prix / " + strings.TrimSuffix(LabelUnite(elt.Unite), "s"),
					Valeur: strconv.FormatFloat(elt.PUHT, 'f', 2, 64),
				},
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.TVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

//
// Transport
//

func (aff *Affacture) computeItemsTransportGlobal(db *sqlx.DB) (err error) {
	list := []PlaqTrans{}
	query := "select * from plaqtrans where id_transporteur=$1 and datetrans>=$2 and datetrans<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("TR"),
			Date:  elt.DateTrans,
		}
		montantHT = elt.GlPrix
		montantTVA = montantHT * elt.GlTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Transport",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.GlTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

func (aff *Affacture) computeItemsTransportConducteur(db *sqlx.DB) (err error) {
	list := []PlaqTrans{}
	query := "select * from plaqtrans where id_conducteur=$1 and datetrans>=$2 and datetrans<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("TR"),
			Date:  elt.DateTrans,
		}
		montantHT = elt.CoNheure * elt.CoPrixH
		montantTVA = montantHT * elt.CoTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Conducteur",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Nb heures",
					Valeur: strconv.FormatFloat(elt.CoNheure, 'f', 2, 64),
				},
				{
					Titre:  "Prix / h",
					Valeur: strconv.FormatFloat(elt.CoPrixH, 'f', 2, 64),
				},
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.CoTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

func (aff *Affacture) computeItemsTransportProprioutil(db *sqlx.DB) (err error) {
	list := []PlaqTrans{}
	query := "select * from plaqtrans where id_proprioutil=$1 and datetrans>=$2 and datetrans<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("TR"),
			Date:  elt.DateTrans,
		}
		if elt.TypeCout == "C" {
			//
			// Camion
			//
			montantHT = elt.CaNkm * elt.CaPrixKm
			montantTVA = montantHT * elt.CaTVA / 100
			montantTTC = montantHT + montantTVA
			ligne = AffactureLigne{
				Titre: "Camion",
				Colonnes: []AffactureColonne{
					{
						Titre:  "Nb km",
						Valeur: strconv.FormatFloat(elt.CaNkm, 'f', 2, 64),
					},
					{
						Titre:  "Prix / km",
						Valeur: strconv.FormatFloat(elt.CaPrixKm, 'f', 2, 64),
					},
					{
						Titre:  "Montant HT",
						Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
					},
					{
						Titre:  "TVA " + strconv.FormatFloat(elt.CaTVA, 'f', -1, 64) + "%",
						Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
					},
					{
						Titre:  "Montant TTC",
						Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
					},
				},
			}
			item.Lignes = append(item.Lignes, ligne)
			aff.TotalHT += montantHT
			aff.TotalTTC += montantTTC
			item.TotalHT += montantHT
			item.TotalTTC += montantTTC
		} else {
			//
			// Tracteur + Benne
			//
			montantHT = float64(elt.TbNbenne) * elt.TbDuree * elt.TbPrixH
			montantTVA = montantHT * elt.TbTVA / 100
			montantTTC = montantHT + montantTVA
			ligne = AffactureLigne{
				Titre: "Tract. + benne",
				Colonnes: []AffactureColonne{
					{
						Titre:  "Nb de bennes",
						Valeur: strconv.Itoa(elt.TbNbenne),
					},
					{
						Titre:  "Durée / benne",
						Valeur: strconv.FormatFloat(elt.TbDuree, 'f', 2, 64),
					},
					{
						Titre:  "Prix HT / heure",
						Valeur: strconv.FormatFloat(elt.TbPrixH, 'f', 2, 64),
					},
					{
						Titre:  "Montant HT",
						Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
					},
					{
						Titre:  "TVA " + strconv.FormatFloat(elt.TbTVA, 'f', -1, 64) + "%",
						Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
					},
					{
						Titre:  "Montant TTC",
						Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
					},
				},
			}
			item.Lignes = append(item.Lignes, ligne)
			aff.TotalHT += montantHT
			aff.TotalTTC += montantTTC
			item.TotalHT += montantHT
			item.TotalTTC += montantTTC
		}
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

//
// Rangement
//

func (aff *Affacture) computeItemsRangementGlobal(db *sqlx.DB) (err error) {
	list := []PlaqRange{}
	query := "select * from plaqrange where id_rangeur=$1 and daterange>=$2 and daterange<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("RG"),
			Date:  elt.DateRange,
		}
		montantHT = elt.GlPrix
		montantTVA = montantHT * elt.GlTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Rangement",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.GlTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		//
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

func (aff *Affacture) computeItemsRangementConducteur(db *sqlx.DB) (err error) {
	list := []PlaqRange{}
	query := "select * from plaqrange where id_conducteur=$1 and daterange>=$2 and daterange<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("RG"),
			Date:  elt.DateRange,
		}
		montantHT = elt.CoNheure * elt.CoPrixH
		montantTVA = montantHT * elt.CoTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Conducteur",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Nb heures",
					Valeur: strconv.FormatFloat(elt.CoNheure, 'f', 2, 64),
				},
				{
					Titre:  "Prix / h",
					Valeur: strconv.FormatFloat(elt.CoPrixH, 'f', 2, 64),
				},
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.CoTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		//
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

func (aff *Affacture) computeItemsRangementProprioutil(db *sqlx.DB) (err error) {
	list := []PlaqRange{}
	query := "select * from plaqrange where id_proprioutil=$1 and daterange>=$2 and daterange<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("RG"),
			Date:  elt.DateRange,
		}
		montantHT = elt.OuPrix
		montantTVA = montantHT * elt.OuTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Outil",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.OuTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		//
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

//
// Chargement
//

func (aff *Affacture) computeItemsChargementGlobal(db *sqlx.DB) (err error) {
	list := []VenteCharge{}
	query := "select * from ventecharge where id_chargeur=$1 and datecharge>=$2 and datecharge<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("CG"),
			Date:  elt.DateCharge,
		}
		montantHT = elt.GlPrix
		montantTVA = montantHT * elt.GlTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Chargement",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.GlTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		//
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

func (aff *Affacture) computeItemsChargementConducteur(db *sqlx.DB) (err error) {
	list := []VenteCharge{}
	query := "select * from ventecharge where id_conducteur=$1 and datecharge>=$2 and datecharge<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("CG"),
			Date:  elt.DateCharge,
		}
		montantHT = elt.MoNHeure * elt.MoPrixH
		montantTVA = montantHT * elt.MoTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Conducteur",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Nb heures",
					Valeur: strconv.FormatFloat(elt.MoNHeure, 'f', 2, 64),
				},
				{
					Titre:  "Prix / h",
					Valeur: strconv.FormatFloat(elt.MoPrixH, 'f', 2, 64),
				},
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.MoTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		//
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

func (aff *Affacture) computeItemsChargementProprioutil(db *sqlx.DB) (err error) {
	list := []VenteCharge{}
	query := "select * from ventecharge where id_proprioutil=$1 and datecharge>=$2 and datecharge<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("CG"),
			Date:  elt.DateCharge,
		}
		montantHT = elt.OuPrix
		montantTVA = montantHT * elt.OuTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Outil",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.OuTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		//
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

//
// Livraison
//

func (aff *Affacture) computeItemsLivraisonGlobal(db *sqlx.DB) (err error) {
	list := []VenteLivre{}
	query := "select * from ventelivre where id_livreur=$1 and datelivre>=$2 and datelivre<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("LV"),
			Date:  elt.DateLivre,
		}
		montantHT = elt.GlPrix
		montantTVA = montantHT * elt.GlTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Livraison",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.GlTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		//
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

func (aff *Affacture) computeItemsLivraisonConducteur(db *sqlx.DB) (err error) {
	list := []VenteLivre{}
	query := "select * from ventelivre where id_conducteur=$1 and datelivre>=$2 and datelivre<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("LV"),
			Date:  elt.DateLivre,
		}
		//
		// Detail - Main d'oeuvre (livreur)
		//
		montantHT = elt.MoNHeure * elt.MoPrixH
		montantTVA = montantHT * elt.MoTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Conducteur",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Nb heures",
					Valeur: strconv.FormatFloat(elt.MoNHeure, 'f', 2, 64),
				},
				{
					Titre:  "Prix / h",
					Valeur: strconv.FormatFloat(elt.MoPrixH, 'f', 2, 64),
				},
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.MoTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		//
		aff.Items = append(aff.Items, &item)
	}
	return nil
}

func (aff *Affacture) computeItemsLivraisonOutil(db *sqlx.DB) (err error) {
	list := []VenteLivre{}
	query := "select * from ventelivre where id_proprioutil=$1 and datelivre>=$2 and datelivre<=$3"
	err = db.Select(&list, query, aff.IdActeur, tiglib.DateIso(aff.DateDebut), tiglib.DateIso(aff.DateFin))
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	var montantHT, montantTVA, montantTTC float64
	var ligne AffactureLigne
	for _, elt := range list {
		item := AffactureItem{
			Titre: LabelActivite("RG"),
			Date:  elt.DateLivre,
		}
		montantHT = elt.OuPrix
		montantTVA = montantHT * elt.OuTVA / 100
		montantTTC = montantHT + montantTVA
		ligne = AffactureLigne{
			Titre: "Outil",
			Colonnes: []AffactureColonne{
				{
					Titre:  "Montant HT",
					Valeur: strconv.FormatFloat(montantHT, 'f', 2, 64),
				},
				{
					Titre:  "TVA " + strconv.FormatFloat(elt.OuTVA, 'f', -1, 64) + "%",
					Valeur: strconv.FormatFloat(montantTVA, 'f', 2, 64),
				},
				{
					Titre:  "Montant TTC",
					Valeur: strconv.FormatFloat(montantTTC, 'f', 2, 64),
				},
			},
		}
		item.Lignes = append(item.Lignes, ligne)
		aff.TotalHT += montantHT
		aff.TotalTTC += montantTTC
		item.TotalHT += montantHT
		item.TotalTTC += montantTTC
		//
		aff.Items = append(aff.Items, &item)
	}
	return nil
}
