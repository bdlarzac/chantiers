/******************************************************************************
    Chantier plaquettes - contient infos générales d'un chantier

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"strconv"
	"time"

	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	//"fmt"
)

type Plaq struct {
	Id              int
	IdLieudit       int       `db:"id_lieudit"`
	IdFermier       int       `db:"id_fermier"`
	IdUG            int       `db:"id_ug"`
	DateDebut       time.Time `db:"datedeb"`
	DateFin         time.Time
	Surface         float64
	Exploitation    string
	Essence         string
	FraisRepas      float64
	FraisReparation float64
	// pas stocké en base
	Volume     float64
	Lieudit    *Lieudit
	Fermier    *Acteur
	UG         *UG
	Tas        []*Tas
	Operations []*PlaqOp
	Transports []*PlaqTrans
	Rangements []*PlaqRange
	Ventes     []*VentePlaq
	Cout       *CoutPlaq
}

// Coût exploitation
type CoutPlaq struct {
	// poste - coût / map sèche
	Abattage     float64
	Debardage    float64
	Dechiquetage float64
	Broyage      float64
	Transport    float64
	Livraison    float64
	Stockage     float64
	// récap compta
	/*
	   Travaux bois déchiqueté
	   Transport plateforme
	   Total coût d'exploit sans broyage
	   Total coût d'exploit avec broyage
	   Transport client PF
	*/
	// récap
	/*
	   Transport
	   Rangement
	   Repas
	   Réparation
	   TOTAL
	*/
	PrixParMap float64
}

// Coût exploitation
// Doit être effectué sur un chantier obtenu par GetPlaqFull() - pas de vérification d'erreur
func (ch *Plaq) ComputeCout(db *sqlx.DB, config *Config) error {
	ch.Cout = &CoutPlaq{}
	nMapSec := ch.Volume * (1 - config.PourcentagePerte/100)
	var cout float64
	//
	for _, op := range ch.Operations {
		cout = op.PUHT * op.Qte / nMapSec
		switch op.TypOp {
		case "AB":
			ch.Cout.Abattage += cout
		case "DB":
			ch.Cout.Debardage += cout
		case "BR":
			ch.Cout.Broyage += cout
		case "DC":
			ch.Cout.Dechiquetage += cout
		}
	}
	//
	cout = 0
	for _, tr := range ch.Transports {
		if tr.TypeCout == "G" {
			cout += tr.GlPrix
		} else if tr.TypeCout == "C" {
			cout += tr.CaNkm * tr.CaPrixKm
		} else if tr.TypeCout == "T" {
			cout += float64(tr.TbNbenne) * tr.TbDuree * tr.TbPrixH
		}
	}
	ch.Cout.Transport = cout / nMapSec
	//
	// A finir - spec pas claire
	return nil
}

// ************************** Manipulation Volume *******************************

// @param   vol en maps
func (ch *Plaq) ModifierVolume(db *sqlx.DB, vol float64) {
	ch.Volume += vol
}

// ************************** Nom *******************************

func (ch *Plaq) String() string {
	if ch.Lieudit == nil {
		panic("Erreur dans le code - Le lieu-dit d'un chantier plaquettes doit être calculé avant d'appeler String()")
	}
	return ch.Lieudit.Nom + " " + tiglib.DateFr(ch.DateDebut)
}

// ************************** Get *******************************

// Renvoie un chantier plaquette
// contenant uniquement les données stockées en base
func GetPlaq(db *sqlx.DB, idChantier int) (*Plaq, error) {
	chantier := &Plaq{}
	query := "select * from plaq where id=$1"
	row := db.QueryRowx(query, idChantier)
	err := row.StructScan(chantier)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur query : "+query)
	}
	return chantier, nil
}

// Renvoie un chantier plaquette contenant
// - les données stockées dans la table
// - Lieudit
// - UG
// - Fermier
// - Tas
// - les opérations simples (abattage...)
// - les transports vers le stockage
// - les opérations de rangement
// Toutes les activités liées à ce chantier sont triées par ordre chronologique inverse
func GetPlaqFull(db *sqlx.DB, idChantier int) (*Plaq, error) {
	chantier, err := GetPlaq(db, idChantier)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel GetPlaq()")
	}
	err = chantier.ComputeVolume(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Plaq.ComputeVolume()")
	}
	err = chantier.ComputeLieudit(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudit()")
	}
	err = chantier.ComputeUG(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Plaq.ComputeUG()")
	}
	err = chantier.ComputeFermier(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Plaq.ComputeFermier()")
	}
	err = chantier.ComputeOperations(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Plaq.ComputeOperations()")
	}
	err = chantier.ComputeTransports(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Plaq.ComputeTransports()")
	}
	err = chantier.ComputeRangements(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Plaq.ComputeRangements()")
	}
	err = chantier.ComputeTas(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Plaq.ComputeTas()")
	}
	err = chantier.ComputeVentes(db)
	if err != nil {
		return chantier, werr.Wrapf(err, "Erreur appel Plaq.ComputeVentes()")
	}
	// inclure CoutPlaq ?
	//
	return chantier, nil
}

// Renvoie la liste des années ayant des chantiers bois sur pied,
// triées par ordre chronologique inverse.
// @param exclude   Année à exclure du résultat
func GetPlaqDifferentYears(db *sqlx.DB, exclude string) ([]string, error) {
	res := []string{}
	list := []time.Time{}
	query := "select datedeb from plaq order by datedeb desc"
	err := db.Select(&list, query)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, d := range list {
		y := strconv.Itoa(d.Year())
		if !tiglib.InArray(y, res) && y != exclude {
			res = append(res, y)
		}
	}
	return res, nil
}

// Renvoie la liste des chantiers plaquettes pour une année donnée,
// triés par ordre chronologique inverse.
// Chaque chantier contient les mêmes champs que ceux renvoyés par GetPlaqFull()
func GetPlaqsOfYear(db *sqlx.DB, annee string) ([]*Plaq, error) {
	res := []*Plaq{}
	type ligne struct {
		Id      int
		DateDeb time.Time
	}
	tmp1 := []*ligne{}
	// select aussi datedeb au lieu de seulement id pour pouvoir faire le order by
	query := "select id,datedeb from plaq where extract(year from datedeb)=$1 order by datedeb desc"
	err := db.Select(&tmp1, query, annee)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, tmp2 := range tmp1 {
		chantier, err := GetPlaqFull(db, tmp2.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaqFull()")
		}
		res = append(res, chantier)
	}
	return res, nil
}

// ************************** Compute *******************************

func (ch *Plaq) ComputeVolume(db *sqlx.DB) error {
	var volumes []float64
	query := "select qte from plaqop where id_chantier=$1 and typop='DC'" // déchiquetage
	err := db.Select(&volumes, query, ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	ch.Volume = 0
	for _, volume := range volumes {
		ch.Volume += volume
	}
	return nil
}

func (chantier *Plaq) ComputeLieudit(db *sqlx.DB) error {
	if chantier.Lieudit != nil {
		return nil
	}
	var err error
	chantier.Lieudit, err = GetLieudit(db, chantier.IdLieudit)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetLieudit()")
	}
	return nil
}

func (chantier *Plaq) ComputeUG(db *sqlx.DB) error {
	if chantier.UG != nil {
		return nil
	}
	var err error
	chantier.UG, err = GetUG(db, chantier.IdUG)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetUG()")
	}
	return nil
}

func (chantier *Plaq) ComputeFermier(db *sqlx.DB) error {
	if chantier.Fermier != nil {
		return nil
	}
	var err error
	chantier.Fermier, err = GetActeur(db, chantier.IdFermier)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetActeur()")
	}
	return nil
}

func (chantier *Plaq) ComputeOperations(db *sqlx.DB) error {
	if len(chantier.Operations) != 0 {
		return nil
	}
	query := "select * from plaqop where id_chantier=$1 order by dateop desc"
	err := db.Select(&chantier.Operations, query, &chantier.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for i, _ := range chantier.Operations {
		chantier.Operations[i].ComputeActeur(db)
	}
	return nil
}

func (chantier *Plaq) ComputeTransports(db *sqlx.DB) error {
	if len(chantier.Transports) != 0 {
		return nil
	}
	query := "select * from plaqtrans where id_chantier=$1 order by datetrans desc"
	err := db.Select(&chantier.Transports, query, &chantier.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for i, _ := range chantier.Transports {
		chantier.Transports[i].ComputeTas(db)
		chantier.Transports[i].ComputeTransporteur(db)
	}
	return nil
}

func (chantier *Plaq) ComputeRangements(db *sqlx.DB) error {
	if len(chantier.Rangements) != 0 {
		return nil
	}
	query := "select * from plaqrange where id_chantier=$1 order by daterange desc"
	err := db.Select(&chantier.Rangements, query, &chantier.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for i, _ := range chantier.Rangements {
		chantier.Rangements[i].ComputeTas(db)
		chantier.Rangements[i].ComputeConducteur(db)
	}
	return nil
}

func (chantier *Plaq) ComputeTas(db *sqlx.DB) error {
	query := "select * from tas where id_chantier=$1"
	err := db.Select(&chantier.Tas, query, &chantier.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for i, _ := range chantier.Tas {
		chantier.Tas[i].ComputeStockage(db)
		chantier.Tas[i].Chantier = chantier
		err = chantier.Tas[i].ComputeNom(db)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel Tas.ComputeNom()")
		}
	}
	return nil
}

func (chantier *Plaq) ComputeVentes(db *sqlx.DB) error {
	ids := []int{}
	query := `select id_vente from ventelivre where id in (
                  select id_livraison from ventecharge where id_tas in(
                      select id from tas where id_chantier=$1
                  )
              )`
	err := db.Select(&ids, query, &chantier.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, idVente := range ids {
		vp, err := GetVentePlaq(db, idVente)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetVentePlaq()")
		}
		// Ajoute acteur pour avoir le nom de la vente
		vp.Client, err = GetActeur(db, vp.IdClient)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetActeur() pour client")
		}
		chantier.Ventes = append(chantier.Ventes, vp)
	}
	return nil
}

// ************************** CRUD *******************************

func InsertPlaq(db *sqlx.DB, chantier *Plaq) (int, error) {
	query := `insert into plaq(
        id_lieudit,
        id_fermier,
        id_ug,
        datedeb,
        datefin,
        surface,
        exploitation,
        essence,
        fraisrepas,
        fraisreparation
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		chantier.IdLieudit,
		chantier.IdFermier,
		chantier.IdUG,
		chantier.DateDebut,
		chantier.DateFin,
		chantier.Surface,
		chantier.Exploitation,
		chantier.Essence,
		chantier.FraisRepas,
		chantier.FraisReparation).Scan(&id)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

func UpdatePlaq(db *sqlx.DB, chantier *Plaq) error {
	query := `update plaq set(
        id_lieudit,
        id_fermier, 
        id_ug,
        datedeb,
        datefin,
        surface,
        exploitation,
        essence,
        fraisrepas, 
        fraisreparation
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) where id=$11`
	_, err := db.Exec(
		query,
		chantier.IdLieudit,
		chantier.IdFermier,
		chantier.IdUG,
		chantier.DateDebut,
		chantier.DateFin,
		chantier.Surface,
		chantier.Exploitation,
		chantier.Essence,
		chantier.FraisRepas,
		chantier.FraisReparation,
		chantier.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeletePlaq(db *sqlx.DB, id int) error {
	var query string
	var err error
	query = "delete from plaqop where id_chantier=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from plaqtrans where id_chantier=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from plaqrange where id_chantier=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from plaq where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
