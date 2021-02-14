/******************************************************************************
    Unités de gestion

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-11-14 23:36:13+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
//"fmt"
)

type UG struct {
	Id                int
	Code              string
	TypeCoupe         string `db:"type_coupe"`
	PrevisionnelCoupe string `db:"previsionnel_coupe"`
	TypePeuplement    string `db:"type_peuplement"`
	// pas stocké en base
	Parcelles        []*Parcelle
	Recaps           map[string]RecapUG
	SortedRecapYears []string // années contenant de l'activité prise en compte dans Recaps
}

// Sert à afficher la liste des activités sur une UG.
// Contient les infos utilisées pour l'affichage, pas les activités.
type UGActivite struct {
	Date        time.Time
	URL         string // URL de la page de l'activité concernée
	NomActivite string
}

type RecapUG struct {
	Annee                  string // YYYY
	Plaquettes             LigneRecapUG
	Chauffage              LigneRecapUG
	PateAPapier            LigneRecapUG
	Palette                LigneRecapUG
	BoisOeuvre             LigneRecapUG
	PlaquettesEtBoisOeuvre LigneRecapUG
}

type LigneRecapUG struct {
	Quantite         float64
	Superficie       float64
	CoutExploitation float64
	Benefice         float64
}

// ************************ Nom *********************************

func (ug *UG) String() string {
	return ug.Code + " -- " + ug.TypePeuplement
}

// ************************ Get *********************************

// Renvoie une UG à partir de son id.
// Ne contient que les champs de la table lieudit.
// Les autres champs ne sont pas remplis.
func GetUG(db *sqlx.DB, id int) (*UG, error) {
	ug := &UG{}
	query := "select * from ug where id=$1"
	row := db.QueryRowx(query, id)
	err := row.StructScan(ug)
	if err != nil {
		return ug, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ug, err
}

func GetUGFull(db *sqlx.DB, id int) (*UG, error) {
	ug, err := GetUG(db, id)
	if err != nil {
		return ug, werr.Wrapf(err, "Erreur appel GetUG()")
	}
	err = ug.ComputeParcelles(db)
	if err != nil {
		return ug, werr.Wrapf(err, "Erreur appel UG.ComputeParcelles()")
	}
	for i, _ := range ug.Parcelles {
		err = ug.Parcelles[i].ComputeLieudits(db)
		if err != nil {
			return ug, werr.Wrapf(err, "Erreur appel Parcelle.ComputeLieudits()")
		}
	}
	return ug, nil
}

// Renvoie une UG à partir de son code, ou nil si aucune UG ne correspond au code
// Ne contient que les champs de la table ug.
// Les autres champs ne sont pas remplis.
// Utilisé par ajax
func GetUGFromCode(db *sqlx.DB, code string) (*UG, error) {
	ug := UG{}
	query := "select * from ug where code=$1"
	err := db.Get(&ug, query, code)
	if err != nil {
		return nil, nil
	}
	return &ug, nil
}

// ************************ Get many *********************************

// Renvoie des UGs à partir d'un lieu-dit.
// Utilise les parcelles pour faire le lien
// Ne contient que les champs de la table ug.
// Les autres champs ne sont pas remplis.
// Utilisé par ajax
func GetUGsFromLieudit(db *sqlx.DB, idLieudit int) ([]*UG, error) {
	ugs := []*UG{}
	// parcelles
	idsParcelles := []int{}
	query := "select id_parcelle from parcelle_lieudit where id_lieudit=$1"
	err := db.Select(&idsParcelles, query, idLieudit)
	if err != nil {
		return ugs, werr.Wrapf(err, "Erreur query : "+query)
	}
	if len(idsParcelles) == 0 {
		return ugs, nil // empty res
	}
	// ids ugs
	strIdsParcelles := tiglib.JoinInt(idsParcelles, ",")
	idsUGs := []int{}
	query = "select distinct id_ug from parcelle_ug where id_parcelle in(" + strIdsParcelles + ")"
	err = db.Select(&idsUGs, query)
	if err != nil {
		return ugs, werr.Wrapf(err, "Erreur query : "+query)
	}
	if len(idsUGs) == 0 {
		return ugs, nil // empty res
	}
	// ugs
	strIdsUGs := tiglib.JoinInt(idsUGs, ",")
	query = "select * from ug where id in(" + strIdsUGs + ") order by code"
	err = db.Select(&ugs, query)
	if err != nil {
		return ugs, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ugs, nil
}

// Renvoie des UGs à partir d'un fermier.
// Utilise les parcelles pour faire le lien
// Ne contient que les champs de la table ug.
// Les autres champs ne sont pas remplis.
// Utilisé par ajax
func GetUGsFromFermier(db *sqlx.DB, idFermier int) ([]*UG, error) {
	ugs := []*UG{}
	query := `
        select * from ug where id in(
            select id_ug from parcelle_ug where id_parcelle in(
                select id_parcelle from parcelle_fermier where id_fermier in(
                    select id from fermier where id=$1
                )
            )
        ) order by code`
	err := db.Select(&ugs, query, idFermier)
	if err != nil {
		return ugs, werr.Wrapf(db.Select(&ugs, query, idFermier), "Erreur query : "+query)
	}
	return ugs, nil
}

// Renvoie les ugs triées par code (nombre romain) et par numéro au sein d'un code (nombres arabes)
// en respectant l'ordre des chiffres romains et arabes.
func GetUGsSortedByCode(db *sqlx.DB) ([]*UG, error) {
    romans := []string {
        "I",
        "II",
        "III",
        "IV",
        "V",
        "VI",
        "VII",
        "VIII",
        "IX",
        "X",
        "XI",
        "XII",
        "XIII",
        "XIV",
        "XV",
        "XVI",
        "XVII",
        "XVIII",
        "XIX",
    }
    res := []*UG{}
	query := `select * from ug`
	err := db.Select(&res, query)
    sort.Slice(res, func(i, j int) bool {
        ug1 := res[i]
        ug2 := res[j]
        code1 := strings.Replace(ug1.Code, ".", "-", -1) // fix typo dans un code (XIX.5)
        tmp1 := strings.Split(code1, "-")
        code2 := strings.Replace(ug2.Code, ".", "-", -1) // fix typo dans un code (XIX.5)
        tmp2 := strings.Split(code2, "-")
        // teste chiffres romains
        idx1 := tiglib.ArraySearchString(romans, tmp1[0])
        idx2 := tiglib.ArraySearchString(romans, tmp2[0])
        if idx1 < idx2 {
            return true
        }
        if idx1 > idx2 {
            return false
        }
        // idx1 = idx2 - chiffres romains identiques
        n1, _ := strconv.Atoi(tmp1[1])
        n2, _ := strconv.Atoi(tmp2[1])
        return n1 < n2
    })
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query : "+query)
	}
    return res, err
}

// Renvoie les ugs triées par code (nombre romain) et par numéro au sein d'un code (nombres arabes)
// en respectant l'ordre des chiffres romains et arabes.
// Renvoie un tableau de tableaux d'UGs dont le code commence par le même nombre romain.
// res[0] : ugs avec code commençant par I-
// res[1] : ugs avec code commençant par II-
// etc.
func GetUGsSortedByCodeAndSeparated(db *sqlx.DB) ([][]*UG, error) {
    res := [][]*UG{}
    ugs, err := GetUGsSortedByCode(db)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur appel GetUGsSortedByCode()")
	}
    curRoman := "I"
    cur := []*UG{}
    for _, ug := range(ugs){
        code := strings.Replace(ug.Code, ".", "-", -1) // fix typo dans un code (XIX.5)
	    roman := ug.Code[: strings.Index(code, "-")]
	    if roman == curRoman {
	        cur = append(cur, ug)
	    } else {
	        // nombre romain différent
            curRoman = roman
            res = append(res, cur)
            cur = []*UG{}
	    }
	}
    res = append(res, cur)
    return res, err
}

// ************************** Compute *******************************

// Remplit le champ Parcelles d'une UG
func (ug *UG) ComputeParcelles(db *sqlx.DB) error {
	if len(ug.Parcelles) != 0 {
		return nil // déjà calculé
	}
	query := `
	    select * from parcelle where id in(
            select id_parcelle from parcelle_ug where id_ug=$1
        ) order by code`
	return db.Select(&ug.Parcelles, query, ug.Id)
}

// Pas inclus dans GetUGFull()
func (ug *UG) ComputeRecap(db *sqlx.DB) error {
	var query string
	var err error
	ug.Recaps = make(map[string]RecapUG)
	//
	// chantiers plaquettes
	//
	ids := []int{}
	query = `select id from plaq where id in(
	    select id_chantier from chantier_ug where type_chantier='plaq' and id_ug =$1
    )`
	err = db.Select(&ids, query, ug.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, idChantier := range ids {
		chantier, err := GetPlaqFull(db, idChantier)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetPlaqFull()")
		}
		y := strconv.Itoa(chantier.DateDebut.Year())
		myrecap := ug.Recaps[y] // à cause de pb "cannot assign"
		myrecap.Annee = y       // au cas où on l'utilise pour la 1e fois
		myrecap.Plaquettes.Quantite += chantier.Volume
		myrecap.Plaquettes.Superficie += chantier.Surface
		myrecap.PlaquettesEtBoisOeuvre.Quantite += chantier.Volume
		myrecap.PlaquettesEtBoisOeuvre.Superficie += chantier.Surface
		// TODO myrecap.Plaquettes.CoutExploitation
		// TODO myrecap.Plaquettes.Benefice
		ug.Recaps[y] = myrecap
	}
	//
	// Chantier autres valorisations
	//
	ids = []int{}
	query = `select id from chautre where id in(
	    select id_chantier from chantier_ug where type_chantier='chautre' and id_ug =$1
    )`
	err = db.Select(&ids, query, ug.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, idChantier := range ids {
		chantier, err := GetChautreFull(db, idChantier)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetChautreFull()")
		}
		y := strconv.Itoa(chantier.DateContrat.Year())
		myrecap := ug.Recaps[y] // à cause de pb "cannot assign"
		myrecap.Annee = y       // au cas où on l'utilise pour la 1e fois
		switch chantier.TypeValo {
		case "BO":
			myrecap.BoisOeuvre.Quantite += chantier.Volume
			myrecap.PlaquettesEtBoisOeuvre.Quantite += chantier.Volume
		case "CH":
			myrecap.Chauffage.Quantite += chantier.Volume
		case "PL":
			myrecap.Palette.Quantite += chantier.Volume
		case "PP":
			myrecap.PateAPapier.Quantite += chantier.Volume
		}
		ug.Recaps[y] = myrecap
	}
	ug.SortedRecapYears = make([]string, 0, len(ug.Recaps))
	for k, _ := range ug.Recaps {
		ug.SortedRecapYears = append(ug.SortedRecapYears, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(ug.SortedRecapYears)))
	//
	return nil
}

// ************************** Activité *******************************

// Renvoie les activités ayant eu lieu sur une UG.
// Ordre chronologique inverse
// Ne renvoie que des infos pour afficher la liste, pas les activités réelles.
func (u *UG) GetActivitesByDate(db *sqlx.DB) ([]*UGActivite, error) {
	res := []*UGActivite{}
	var err error
	var query string
	//
	// Chantiers plaquettes
	//
	list1 := []Plaq{}
	query = `select * from plaq where id in(
	    select id_chantier from chantier_ug where type_chantier='plaq' and id_ug =$1
    )`
	err = db.Select(&list1, query, u.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list1 {
		err = elt.ComputeLieudits(db) // pour le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
		}
		new := &UGActivite{
			Date:        elt.DateDebut,
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.Id),
			NomActivite: "Chantier plaquettes " + elt.String()}
		res = append(res, new)
	}
	//
	// Chantiers bois sur pied
	//
	list2 := []BSPied{}
	query = `select * from bspied where id in(
	    select id_chantier from chantier_ug where type_chantier='bspied' and id_ug =$1
    )`
	err = db.Select(&list2, query, u.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list2 {
		err = elt.ComputeLieudits(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel BSPied.ComputeLieudits()")
		}
		err = elt.ComputeAcheteur(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel BSPied.ComputeAcheteur()")
		}
		new := &UGActivite{
			Date:        elt.DateContrat,
			URL:         "/chantier/plaquette/" + strconv.Itoa(elt.Id),
			NomActivite: "Chantier bois sur pied " + elt.String()}
		res = append(res, new)
	}
	//
	// Chantiers Autres valorisations
	//
	list3 := []Chautre{}
	query = `select * from chautre where id in(
	    select id_chantier from chantier_ug where type_chantier='chautre' and id_ug =$1
    )`
	err = db.Select(&list3, query, u.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list3 {
		err = elt.ComputeClient(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Chautre.ComputeClient()")
		}
		new := &UGActivite{
			Date:        elt.DateContrat,
			URL:         "/chantier/autre/liste/" + strconv.Itoa(elt.DateContrat.Year()),
			NomActivite: "Chantier " + elt.String()}
		res = append(res, new)
	}
	//
	// Chantiers Chauffage fermier
	//
	list4 := []Chaufer{}
	query = `select * from chaufer where id_ug=$1`
	err = db.Select(&list4, query, u.Id)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, elt := range list4 {
		err = elt.ComputeFermier(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Chaufer.ComputeFermier()")
		}
		new := &UGActivite{
			Date:        elt.DateChantier,
			URL:         "/chantier/chauffage-fermier/liste/" + strconv.Itoa(elt.DateChantier.Year()),
			NomActivite: "Chauffage fermier " + elt.String()}
		res = append(res, new)
	}
	// tri par date
	sortedRes := make(ugActiviteSlice, 0, len(res))
	for _, elt := range res {
		sortedRes = append(sortedRes, elt)
	}
	sort.Sort(sortedRes)
	//
	return sortedRes, nil
}

// Auxiliaires de GetActivitesByDate() pour trier par date
type ugActiviteSlice []*UGActivite

func (p ugActiviteSlice) Len() int {
	return len(p)
}
func (p ugActiviteSlice) Less(i, j int) bool {
	return p[i].Date.After(p[j].Date)
}
func (p ugActiviteSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
