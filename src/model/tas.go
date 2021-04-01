/******************************************************************************
    Tas de plaquettes dans un hangar

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-03-09 17:18:27+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
	"strconv"
	"sort"
)

type Tas struct {
	Id         int
	IdStockage int `db:"id_stockage"`
	IdChantier int `db:"id_chantier"`
	Stock      float64
	DateVidage time.Time
	Actif      bool
	// pas stocké en base
	Nom             string
	Chantier        *Plaq
	Stockage        *Stockage
	MesuresHumidite []*Humid
	EvolutionStock  []*MouvementStock
}

// MouvementStock = opération qui fait changer le stock du tas :
// transports, chargements, vidage
type MouvementStock struct{
    Date    time.Time
    Label   string
    URL     string
    Delta   float64
}

func NewTas(idStockage, idChantier int, stock float64, actif bool) *Tas {
	return &Tas{
		IdStockage: idStockage,
		IdChantier: idChantier,
		Stock:      stock,
		Actif:      actif,
	}
}

// ************************** Manipulation Stock *******************************

// Si qte > 0, ajoute des plaquettes au tas
// Si qte < 0, retire des plaquettes au tas
// Fait la maj en BDD
// @param   qte en maps
func (t *Tas) ModifierStock(db *sqlx.DB, qte float64) error {
	t.Stock += qte
	return UpdateTas(db, t)
}

// Pour indiquer qu'un tas est vide
func DesactiverTas(db *sqlx.DB, id int) error {
	tas, err := GetTas(db, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetTas()")
	}
	tas.Actif = false
	err = UpdateTas(db, tas)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel UpdateTas()")
	}
	return nil
}

// ************************** Get one *******************************

func GetTas(db *sqlx.DB, idTas int) (*Tas, error) {
	tas := &Tas{}
	query := "select * from tas where id=$1"
	row := db.QueryRowx(query, idTas)
	err := row.StructScan(tas)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur query : "+query)
	}
	return tas, nil
}

func GetTasFull(db *sqlx.DB, idTas int) (*Tas, error) {
	tas, err := GetTas(db, idTas)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur appel GetTas()")
	}
	err = tas.ComputeStockage(db)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur appel Tas.ComputeStockage()")
	}
	err = tas.ComputeChantier(db)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur appel Tas.ComputeChantier()")
	}
	err = tas.ComputeNom(db)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur appel Tas.ComputeNom()")
	}
	return tas, nil
}

// ************************** Get many *******************************

// Utilisé pour select html
// Obligé d'avoir tas full, car besoin du nom du tas, qui a besoin de chantier et stockage
func GetAllTasActifsFull(db *sqlx.DB) ([]*Tas, error) {
	tas := []*Tas{}
	ids := []int{}
	query := "select id from tas where actif"
	err := db.Select(&ids, query)
	if err != nil {
		return tas, werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, id := range ids {
		t, err := GetTasFull(db, id)
		if err != nil {
			return tas, werr.Wrapf(err, "Erreur appel GetTasFull()")
		}
		tas = append(tas, t)
	}
	return tas, nil
}

// ************************** Compute *******************************

func (t *Tas) ComputeNom(db *sqlx.DB) error {
	if t.Chantier == nil || t.Stockage == nil {
		return errors.New("Impossible de calculer le nom du tas - appeler d'abord ComputeStockage() et ComputeChantier()")
	}
	err := t.Chantier.ComputeLieudits(db) // Pour le nom du chantier
	if err != nil {
		return werr.Wrapf(err, "Erreur appel t.Chantier.ComputeLieudits()")
	}
	t.Nom = t.Chantier.String() + " - " + t.Stockage.Nom
	return nil
}

func (t *Tas) ComputeStockage(db *sqlx.DB) error {
	var err error
	t.Stockage, err = GetStockage(db, t.IdStockage)
	return err
}

func (t *Tas) ComputeChantier(db *sqlx.DB) error {
	var err error
	t.Chantier, err = GetPlaq(db, t.IdChantier)
	return err
}

// pas inclus par défaut dans GetTasFull()
func (t *Tas) ComputeMesuresHumidite(db *sqlx.DB) error {
	var err error
	query := "select * from humid where id_tas=$1"
	err = db.Select(&t.MesuresHumidite, query, t.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, h := range t.MesuresHumidite {
		err = h.ComputeMesureurs(db)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel humid.ComputeMesureurs()")
		}
	}
	return nil
}

// pas inclus par défaut dans GetTasFull()
// note : en théorie, l'url des mouvements ne devrait pas être calculée dans le model mais dans le controller
func (t *Tas) ComputeEvolutionStock(db *sqlx.DB) error {
	var err error
	res := []*MouvementStock{}
	// transports
	var tmp1 = []*struct{
		Id int
		Id_chantier int
		Qte float64
		PourcentPerte float64
		DateTrans  time.Time
	}{}
	query := "select id,id_chantier,qte,pourcentperte,datetrans from plaqtrans where id_tas=$1"
	err = db.Select(&tmp1, query, t.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, line := range(tmp1){
	    mvt := MouvementStock{
	        Date: line.DateTrans,
	        Label: "Transport",
	        URL: "/chantier/plaquette/" + strconv.Itoa(line.Id_chantier) + "/chantiers",
	        Delta: line.Qte*(1-line.PourcentPerte/100),
	    }
	    res = append(res, &mvt)
	}
	// chargements (utilise GetVenteCharge() pour récupérer id vente)
    idsCharge := []int{}
	query = "select id from ventecharge where id_tas=$1"
	err = db.Select(&idsCharge, query, t.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, idCharge := range(idsCharge){
	    vc, err := GetVenteCharge(db, idCharge)
        if err != nil {
            return werr.Wrapf(err, "Erreur appel GetVenteCharge()")
        }
        err = vc.ComputeIdVente(db)
        if err != nil {
            return werr.Wrapf(err, "Erreur appel ComputeIdVente()")
        }
	    mvt := MouvementStock{
	        Date: vc.DateCharge,
	        Label: "Chargement",
	        URL: "/vente/plaquette/" + strconv.Itoa(vc.IdVente),
	        Delta: -vc.Qte,
	    }
	    res = append(res, &mvt)
	}
	// vidage
// TODO 
	// tri par date
	sortedRes := make(mouvementStockSlice, 0, len(res))
	for _, elt := range res {
		sortedRes = append(sortedRes, elt)
	}
	sort.Sort(sortedRes)
	t.EvolutionStock = sortedRes
	//
	return nil
}

// Auxiliaires de ComputeEvolutionStock() pour trier par date
type mouvementStockSlice []*MouvementStock

func (m mouvementStockSlice) Len() int {
	return len(m)
}
func (m mouvementStockSlice) Less(i, j int) bool {
	return m[i].Date.After(m[j].Date)
}
func (m mouvementStockSlice) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}



// ************************** CRUD *******************************

func InsertTas(db *sqlx.DB, tas *Tas) (int, error) {
	query := `insert into tas(
        id_stockage,                              
        id_chantier,
        stock,
        actif
        ) values($1,$2,$3,$4) returning id`
	id := int(0)
	err := db.QueryRow(
		query,
		tas.IdStockage,
		tas.IdChantier,
		tas.Stock,
		tas.Actif).Scan(&id)
	if err != nil {
		return id, werr.Wrapf(err, "Erreur query : "+query)
	}
	return id, nil
}

func UpdateTas(db *sqlx.DB, tas *Tas) error {
	query := `update tas set(
        id_stockage,
        id_chantier,
        stock,
        actif
        ) = ($1,$2,$3,$4) where id=$5`
	_, err := db.Exec(
		query,
		tas.IdStockage,
		tas.IdChantier,
		tas.Stock,
		tas.Actif,
		tas.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func DeleteTas(db *sqlx.DB, id int) error {
	var err error
	var query string
	var ids []int
	var deletedId int
	// delete transports associés à ce tas
	query = "select id from plaqtrans where id_tas=$1"
	err = db.Select(&ids, query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, deletedId = range ids {
		err = DeletePlaqTrans(db, deletedId)
		if err != nil {
			return werr.Wrapf(err, "Erreur DeletePlaqTrans()")
		}
	}
	// delete chargements liés à ce tas
	query = "select id from ventecharge where id_tas=$1"
	err = db.Select(&ids, query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, deletedId = range ids {
		err = DeleteVenteCharge(db, deletedId)
		if err != nil {
			return werr.Wrapf(err, "Erreur DeleteVenteCharge()")
		}
	}
	// delete le tas
	query = "delete from tas where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
