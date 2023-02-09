/******************************************************************************
    Chantier plaquettes - contient infos générales d'un chantier

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2019, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"strconv"
	"strings"
	"time"

	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

// Représente un chantier plaquettes
type Plaq struct {
	Id              int
	DateDebut       time.Time `db:"datedeb"`
	DateFin         time.Time
	Surface         float64
	Granulo         string
	Exploitation    string
	Essence         string
	FraisRepas      float64
	FraisReparation float64
	Notes           string
	// pas stocké en base
	UGs        []*UG
	LiensParcelles []*ChantierParcelle
	Lieudits   []*Lieudit
	Fermiers   []*Fermier
	Volume     float64
	Tas        []*Tas
	TasVides   []*Tas
	TasActifs  []*Tas
	Operations []*PlaqOp
	Transports []*PlaqTrans
	Rangements []*PlaqRange
	Ventes     []*VentePlaq
	CoutTotal  *CoutPlaq
	CoutParMap *CoutPlaq
}

// Coût exploitation
type CoutPlaq struct {
	Abattage     float64
	Debardage    float64
	Dechiquetage float64
	Broyage      float64
	FauxFrais    float64 // repas et réparation
	Transport    float64
	Rangement    float64
	Stockage     float64
	Chargement   float64
	Livraison    float64
	//
	Total float64
}

// ************************** Manipulation Volume *******************************

// @param   vol en maps
func (ch *Plaq) ModifierVolume(db *sqlx.DB, vol float64) {
	ch.Volume += vol
}

// ************************** Nom *******************************

func (ch *Plaq) String() string {
	if len(ch.Lieudits) == 0 {
		panic("Erreur dans le code - Les lieux-dits d'un chantier plaquettes doivent être calculés avant d'appeler String()")
	}
	res := ""
	var noms []string
	for _, ld := range ch.Lieudits {
		noms = append(noms, ld.Nom)
	}
	res += strings.Join(noms, " - ")
	res += " " + tiglib.DateFr(ch.DateDebut)
	return res
}

func (ch *Plaq) FullString() string {
	return "Chantier plaquettes " + ch.String()
}

// ************************** Get one *******************************

// Renvoie un chantier plaquette
// contenant uniquement les données stockées en base
func GetPlaq(db *sqlx.DB, idChantier int) (*Plaq, error) {
	ch := &Plaq{}
	query := "select * from plaq where id=$1"
	row := db.QueryRowx(query, idChantier)
	err := row.StructScan(ch)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur query : "+query)
	}
	return ch, nil
}

// Renvoie un chantier plaquette contenant
//      - les données stockées dans la table
//      - les lieux-dits
//      - les UGs
//      - les parcelles
//      - les fermiers
//      - les Tas
//      - les opérations simples (abattage...)
//      - les transports vers le stockage
//      - les opérations de rangement
func GetPlaqFull(db *sqlx.DB, idChantier int) (*Plaq, error) {
	ch, err := GetPlaq(db, idChantier)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel GetPlaq()")
	}
	err = ch.ComputeVolume(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Plaq.ComputeVolume()")
	}
	err = ch.ComputeUGs(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Plaq.ComputeUGs()")
	}
	err = ch.ComputeLiensParcelles(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Plaq.ComputeLiensParcelles()")
	}
	err = ch.ComputeLieudits(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
	}
	err = ch.ComputeFermiers(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Plaq.ComputeFermiers()")
	}
	err = ch.ComputeOperations(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Plaq.ComputeOperations()")
	}
	err = ch.ComputeTransports(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Plaq.ComputeTransports()")
	}
	err = ch.ComputeRangements(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Plaq.ComputeRangements()")
	}
	err = ch.ComputeTas(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Plaq.ComputeTas()")
	}
	err = ch.ComputeVentes(db)
	if err != nil {
		return ch, werr.Wrapf(err, "Erreur appel Plaq.ComputeVentes()")
	}
	// ComputeCouts pas inclus
	//
	return ch, nil
}

// ************************** Get many *******************************

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
		if !tiglib.InArrayString(y, res) && y != exclude {
			res = append(res, y)
		}
	}
	return res, nil
}

// Renvoie la liste des chantiers plaquettes pour une année donnée,
// Chaque chantier contient les mêmes champs que ceux renvoyés par GetPlaqFull()
func GetPlaqsOfYear(db *sqlx.DB, annee string) ([]*Plaq, error) {
	res := []*Plaq{}
	type ligne struct {
		Id      int
		DateDeb time.Time
	}
	tmp1 := []*ligne{}
	// select aussi datedeb au lieu de seulement id pour pouvoir faire le order by
	query := "select id,datedeb from plaq where extract(year from datedeb)=$1 order by datedeb"
	err := db.Select(&tmp1, query, annee)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, tmp2 := range tmp1 {
		ch, err := GetPlaqFull(db, tmp2.Id)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaqFull()")
		}
		res = append(res, ch)
	}
	return res, nil
}

// Renvoie la liste des chantiers plaquettes ayant des tas (stockage) vides
func GetAllPlaqsVides(db *sqlx.DB) ([]*Plaq, error) {
	res := []*Plaq{}
	var err error
	var idsChantier []int
	query := "select id_chantier from tas where actif=false order by datevidage desc"
	err = db.Select(&idsChantier, query)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, idChantier := range idsChantier {
		ch, err := GetPlaq(db, idChantier)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel GetPlaq()")
		}
		err = ch.ComputeTas(db)
		for _, tas := range ch.TasVides {
			err = tas.ComputeEvolutionStock(db)
			if err != nil {
				return res, werr.Wrapf(err, "Erreur appel Tas.ComputeEvolutionStock()")
			}
		}
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeTas()")
		}
		err = ch.ComputeVolume(db)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeVolume()")
		}
		err = ch.ComputeLieudits(db) // nécessaire pour afficher le nom du chantier
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel Plaq.ComputeLieudits()")
		}
		res = append(res, ch)
	}
	return res, err
}

// ************************** Compute *******************************

func (ch *Plaq) ComputeUGs(db *sqlx.DB) (err error) {
	if len(ch.UGs) != 0 {
		return nil // déjà calculé
	}
	ch.UGs, err = computeUGsOfChantier(db, "plaq", ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeUGsOfChantier")
	}
	return nil
}

func (ch *Plaq) ComputeLiensParcelles(db *sqlx.DB) (err error) {
	if len(ch.LiensParcelles) != 0 {
		return nil
	}
	ch.LiensParcelles, err = computeLiensParcellesOfChantier(db, "plaq", ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeLiensParcellesOfChantier")
	}
	return nil
}

func (ch *Plaq) ComputeLieudits(db *sqlx.DB) (err error) {
	if len(ch.Lieudits) != 0 {
		return nil // déjà calculé
	}
	ch.Lieudits, err = computeLieuditsOfChantier(db, "plaq", ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeLieuditsOfChantier")
	}
	return nil
}

func (ch *Plaq) ComputeFermiers(db *sqlx.DB) (err error) {
	if len(ch.Fermiers) != 0 {
		return nil // déjà calculé
	}
	ch.Fermiers, err = computeFermiersOfChantier(db, "plaq", ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeFermiersOfChantier")
	}
	return nil
}

// Attention, deux calculs possibles du volume vert :
// - Somme des quantités déchiquetées (utilisée ici)
// - Somme des quantités transportées dans les tas
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

func (ch *Plaq) ComputeOperations(db *sqlx.DB) error {
	if len(ch.Operations) != 0 {
		return nil
	}
	query := "select * from plaqop where id_chantier=$1 order by datedeb"
	err := db.Select(&ch.Operations, query, &ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for i, _ := range ch.Operations {
		ch.Operations[i].ComputeActeur(db)
	}
	return nil
}

func (ch *Plaq) ComputeTransports(db *sqlx.DB) error {
	if len(ch.Transports) != 0 {
		return nil
	}
	query := "select * from plaqtrans where id_chantier=$1 order by datetrans"
	err := db.Select(&ch.Transports, query, &ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for i, _ := range ch.Transports {
		ch.Transports[i].ComputeTas(db)
		ch.Transports[i].ComputeTransporteur(db)
		ch.Transports[i].ComputeConducteur(db)
		ch.Transports[i].ComputeProprioutil(db)
	}
	return nil
}

func (ch *Plaq) ComputeRangements(db *sqlx.DB) error {
	if len(ch.Rangements) != 0 {
		return nil
	}
	query := "select * from plaqrange where id_chantier=$1 order by daterange"
	err := db.Select(&ch.Rangements, query, &ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for i, _ := range ch.Rangements {
		ch.Rangements[i].ComputeTas(db)
		ch.Rangements[i].ComputeRangeur(db)
		ch.Rangements[i].ComputeConducteur(db)
		ch.Rangements[i].ComputeProprioutil(db)
	}
	return nil
}

func (ch *Plaq) ComputeTas(db *sqlx.DB) error {
	query := "select * from tas where id_chantier=$1"
	err := db.Select(&ch.Tas, query, &ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for i, _ := range ch.Tas {
		ch.Tas[i].ComputeStockage(db)
		ch.Tas[i].Chantier = ch
		err = ch.Tas[i].ComputeNom(db)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel Tas.ComputeNom()")
		}
		if ch.Tas[i].Actif {
			ch.TasActifs = append(ch.TasActifs, ch.Tas[i])
		} else {
			ch.TasVides = append(ch.TasVides, ch.Tas[i])
		}
	}
	return nil
}

func (ch *Plaq) ComputeVentes(db *sqlx.DB) error {
	ids := []int{}
	query := `select id_vente from ventelivre where id in (
                  select id_livraison from ventecharge where id_tas in(
                      select id from tas where id_chantier=$1
                  )
              )`
	err := db.Select(&ids, query, &ch.Id)
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
		ch.Ventes = append(ch.Ventes, vp)
	}
	return nil
}

// Calcule les différents coûts d'exploitation
// Doit être effectué sur un chantier obtenu par GetPlaqFull() - pas de vérification d'erreur
func (ch *Plaq) ComputeCouts(db *sqlx.DB, config *Config) (err error) {
	if ch.Volume == 0 {
		// valeurs par défaut, tous les coûts restent à 0
		return nil
	}
	ch.CoutParMap = &CoutPlaq{}
	ch.CoutTotal = &CoutPlaq{}
	nMapSec := ch.Volume * (1 - config.PourcentagePerte/100)
	var cout float64
	//
	// Opérations simples
	//
	for _, op := range ch.Operations {
		cout = op.PUHT * op.Qte
		ch.CoutParMap.Total += cout / nMapSec
		ch.CoutTotal.Total += cout
		switch op.TypOp {
		case "AB":
			ch.CoutParMap.Abattage += cout / nMapSec
			ch.CoutTotal.Abattage += cout
		case "DB":
			ch.CoutParMap.Debardage += cout / nMapSec
			ch.CoutTotal.Debardage += cout
		case "BR":
			ch.CoutParMap.Broyage += cout / nMapSec
			ch.CoutTotal.Broyage += cout
		case "DC":
			ch.CoutParMap.Dechiquetage += cout / nMapSec
			ch.CoutTotal.Dechiquetage += cout
		}
	}
	//
	// Faux frais
	//
	ch.CoutParMap.FauxFrais = (ch.FraisReparation + ch.FraisRepas) / nMapSec
	ch.CoutParMap.Total += ch.CoutParMap.FauxFrais
	ch.CoutTotal.FauxFrais = ch.FraisReparation + ch.FraisRepas
	ch.CoutTotal.Total += ch.CoutTotal.FauxFrais
	//
	// Transport
	//
	cout = 0
	for _, t := range ch.Transports {
		if t.TypeCout == "G" {
			cout += t.GlPrix
		} else if t.TypeCout == "C" {
			cout += t.CoNheure * t.CoPrixH // conducteur
			cout += t.CaNkm * t.CaPrixKm   // outil
		} else if t.TypeCout == "T" {
			cout += t.CoNheure * t.CoPrixH                      // conducteur
			cout += float64(t.TbNbenne) * t.TbDuree * t.TbPrixH // outil
		}
	}
	ch.CoutParMap.Transport = cout / nMapSec
	ch.CoutParMap.Total += ch.CoutParMap.Transport
	ch.CoutTotal.Transport = cout
	ch.CoutTotal.Total += ch.CoutTotal.Transport
	//
	// Rangement
	//
	cout = 0
	for _, r := range ch.Rangements {
		if r.TypeCout == "G" {
			cout += r.GlPrix
		} else {
			cout += r.CoPrixH * r.CoNheure // conducteur
			cout += r.OuPrix               // outil
		}
	}
	ch.CoutParMap.Rangement = cout / nMapSec
	ch.CoutParMap.Total += ch.CoutParMap.Rangement
	ch.CoutTotal.Rangement = cout
	ch.CoutTotal.Total += ch.CoutTotal.Rangement
	//
	// Chargement et livraisons
	//
	var coutC, coutL float64
	for _, v := range ch.Ventes {
		// ch.Ventes ne contient que les champs de la base
		// donc appel de GetVentePlaqFull() pour avoir une vente et ses livraisons
		vf, err := GetVentePlaqFull(db, v.Id)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetVentePlaqFull()")
		}
		for _, l := range vf.Livraisons {
			if l.TypeCout == "G" {
				coutL += l.GlPrix
			} else {
				coutL += l.OuPrix               // outil
				coutL += l.MoNHeure * l.MoPrixH // main d'oeuvre
			}
			for _, c := range l.Chargements {
				if c.TypeCout == "G" {
					coutC += c.GlPrix
				} else {
					coutC += c.OuPrix               // outil
					coutC += c.MoNHeure * c.MoPrixH // main d'oeuvre
				}
			}
		}
	}
	//
	// Stockage
	//
	err = ch.computeCoutStockage(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeCoutStockage()")
	}
	ch.CoutTotal.Chargement = coutC
	ch.CoutTotal.Livraison = coutL
	ch.CoutTotal.Total += ch.CoutTotal.Chargement
	ch.CoutTotal.Total += ch.CoutTotal.Livraison
	ch.CoutParMap.Chargement = coutC / nMapSec
	ch.CoutParMap.Livraison = coutL / nMapSec
	ch.CoutParMap.Total += ch.CoutParMap.Chargement
	ch.CoutParMap.Total += ch.CoutParMap.Livraison
	ch.CoutParMap.Total += ch.CoutParMap.Stockage
	//
	return nil
}

// Calcule ch.CoutParMap.Stockage
// Auxiliaire de ComputeCouts(), donc ch est obtenu par GetPlaqFull()
func (ch *Plaq) computeCoutStockage(db *sqlx.DB) (err error) {
	//
	// Calcule tous les hangars (Stockage) contenant des tas liés à ce chantier
	//
	var stockages []*Stockage
	for _, t := range ch.Tas {
		s, err := GetStockage(db, t.IdStockage)
		if err != nil {
			return werr.Wrapf(err, "Erreur appel GetStockage()")
		}
		stockages = append(stockages, s)
	}

	// j1, j2, DATE_MIN, DATE_MAX utilisés pour coût stockage
	// j1 = date du premier transport
	// j2 = date du dernier chargement
	/*
		DATE_MAX, _ := time.Parse("2006-01-02", "2999-12-31")
		DATE_MIN, _ := time.Parse("2006-01-02", "1999-12-31")
		j1 := DATE_MAX
		j2 := DATE_MIN
	*/
	/*
	   // TODO commenté car besoin de trouver le mode de calcul
	   var tas *Tas
	   cout = 0
	   // s'il y a au moins un transport et un chargement
	   if j1 != DATE_MAX && j2 != DATE_MIN {
	       // vérifie que tous les tas du chantier ont été déclarés vides
	       vides := true
	       for _, tas = range(ch.Tas){
	           if tas.Actif {
	               vides = false
	           }
	       }
	       if vides == true {

	       }
	   }
	*/
	return nil
}

// ************************** CRUD *******************************

/** 
    Insère un chantier plaquette en base
    + Crée et insère en base le(s) tas (crée un Tas par lieu de stockage)
    + Insère en base les liens UGs, parcelles, lieux-dits, fermiers
**/
func InsertPlaq(db *sqlx.DB, ch *Plaq, idsStockages, idsUG, idsLieudit, idsFermier []int) (idChantier int, err error) {
	query := `insert into plaq(
        datedeb,
        datefin,
        surface,
        granulo,
        exploitation,
        essence,
        fraisrepas,
        fraisreparation,
        notes
        ) values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	idChantier = int(0)
	err = db.QueryRow(
		query,
		ch.DateDebut,
		ch.DateFin,
		ch.Surface,
		ch.Granulo,
		ch.Exploitation,
		ch.Essence,
		ch.FraisRepas,
		ch.FraisReparation,
		ch.Notes).Scan(&idChantier)
	if err != nil {
		return idChantier, werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// tas - crée un tas par lieu de stockage sélectionné
	//
	for _, idStockage := range idsStockages {
		tas := NewTas(idStockage, idChantier, 0, true)
		_, err = InsertTas(db, tas)
		if err != nil {
			return idChantier, werr.Wrapf(err, "Erreur appel InsertTas()")
		}
	}
	//
	// insert associations avec UGs, Parcelles, Lieudits, Fermiers
	//
	err = insertLiensChantierUG(db, "plaq", idChantier, idsUG)
    if err != nil {
        return idChantier, werr.Wrapf(err, "Erreur appel insertLiensChantierUG()")
    }
	//
	err = insertLiensChantierParcelle(db, "plaq", idChantier, ch.LiensParcelles)
	if err != nil {
		return idChantier, werr.Wrapf(err, "Erreur appel insertLiensChantierParcelle()")
	}
	//
	err = insertLiensChantierLieudit(db, "plaq", idChantier, idsLieudit)
    if err != nil {
        return idChantier, werr.Wrapf(err, "Erreur appel insertLiensChantierLieudit()")
    }
	//
	err = insertLiensChantierFermier(db, "plaq", idChantier, idsFermier)
    if err != nil {
        return idChantier, werr.Wrapf(err, "Erreur appel insertLiensChantierFermier()")
    }
	//
	return idChantier, nil
}

/** 
    MAJ un chantier plaquette en base
    + Gère aussi les tas
    + MAJ en base les liens UGs, parcelles, lieux-dits, fermiers
    
    @param idsStockages ids tas APRÈS update
**/
func UpdatePlaq(db *sqlx.DB, ch *Plaq, idsStockages, idsUG, idsLieudit, idsFermier []int) (err error) {
	query := `update plaq set(
        datedeb,
        datefin,
        surface,
        granulo,
        exploitation,
        essence,
        fraisrepas, 
        fraisreparation,
        notes
        ) = ($1,$2,$3,$4,$5,$6,$7,$8,$9) where id=$10`
	_, err = db.Exec(
		query,
		ch.DateDebut,
		ch.DateFin,
		ch.Surface,
		ch.Granulo,
		ch.Exploitation,
		ch.Essence,
		ch.FraisRepas,
		ch.FraisReparation,
		ch.Notes,
		ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// tas
	//
	// on note AV les stockages associés au chantier avant update
	// on note AP les stockages associés au chantier après update
	// si AV et pas AP => supprimer tas AV
	// si AP et pas AV => créer tas AP
	// si AP et AV => ne rien faire
	idsStockageAP := idsStockages
	// calculer idsStockageAV à partir de la base
	idsStockageAV := []int{}
	query = "select id_stockage from tas where id_chantier=$1"
	err = db.Select(&idsStockageAV, query, ch.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query DB : "+query)
	}
	// si AV et pas AP => supprimer tas AV
	for _, av := range idsStockageAV {
		if !tiglib.InArrayInt(av, idsStockageAP) {
			// Attention, ne pas faire un DeleteTas() directement avec une query
			// car DeleteTas() a pour effet de supprimer les activités qui lui sont reliées.
			var idTasToDelete int
			query = "select id from tas where id_chantier=$1 and id_stockage=$2"
			err = db.Get(&idTasToDelete, query, ch.Id, av)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel Get(), query = "+query)
			}
			err = DeleteTas(db, idTasToDelete)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel DeleteTas()")
			}
		}
	}
	// si AP et pas AV => créer tas AP
	for _, ap := range idsStockageAP {
		if !tiglib.InArrayInt(ap, idsStockageAV) {
			tas := NewTas(ap, ch.Id, 0, true)
			_, err = InsertTas(db, tas)
			if err != nil {
				return werr.Wrapf(err, "Erreur appel InsertTas()")
			}
		}
	}
	//
	// update associations avec UGs, Parcelles, Lieudits, Fermiers
	//
	err = updateLiensChantierUG(db, "plaq", ch.Id, idsUG)
    if err != nil {
        return werr.Wrapf(err, "Erreur appel updateLiensChantierUG()")
    }
	//
	err = updateLiensChantierParcelle(db, "plaq", ch.Id, ch.LiensParcelles)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel updateLiensChantierParcelle()")
	}
	//
	err = updateLiensChantierLieudit(db, "plaq", ch.Id, idsLieudit)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel updateLiensChantierLieudit()")
	}
	//
	err = updateLiensChantierFermier(db, "plaq", ch.Id, idsFermier)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel updateLiensChantierFermier()")
	}
	//
	return nil
}

func DeletePlaq(db *sqlx.DB, id int) (err error) {
	var query string
	var ids []int
	var deletedId int
	//
	// delete transports associés à ce chantier
	//
	query = "select id from plaqtrans where id_chantier=$1"
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
	//
	// delete rangements associés à ce chantier
	//
	query = "select id from plaqrange where id_chantier=$1"
	err = db.Select(&ids, query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, deletedId = range ids {
		err = DeletePlaqRange(db, deletedId)
		if err != nil {
			return werr.Wrapf(err, "Erreur DeletePlaqRange()")
		}
	}
	//
	// delete opérations simples associées à ce chantier
	//
	query = "select id from plaqop where id_chantier=$1"
	err = db.Select(&ids, query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, deletedId = range ids {
		err = DeletePlaqOp(db, deletedId)
		if err != nil {
			return werr.Wrapf(err, "Erreur DeletePlaqOp()")
		}
	}
	//
	// delete tas associés à ce chantier
	//
	query = "select id from tas where id_chantier=$1"
	err = db.Select(&ids, query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	for _, deletedId = range ids {
		err = DeleteTas(db, deletedId)
		if err != nil {
			return werr.Wrapf(err, "Erreur DeletePlaqOp()")
		}
	}
	//
	// delete associations avec UGs, Parcelles, Lieudits, Fermiers
	//
    err = deleteLiensChantierUG(db, "plaq", id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierUG()")
	}
	//
    err = deleteLiensChantierParcelle(db, "plaq", id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierParcelle()")
	}
	//
    err = deleteLiensChantierLieudit(db, "plaq", id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierLieudit()")
	}
	//
    err = deleteLiensChantierFermier(db, "plaq", id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierFermier()")
	}
	//
	// delete le chantier, fait à la fin pour respecter les clés étrangères
	//
	query = "delete from plaq where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}
