/******************************************************************************
    Gestion des chantiers visités récemment

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-05-14 17:35:22+02:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
    "time"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type Recent struct {
	URL        string
	Label      string
	DateVisite time.Time
}

func GetRecents(db *sqlx.DB) ([]*Recent, error) {
	r := []*Recent{}
	query := "select * from recent order by datevisite desc"
	err := db.Select(&r, query)
	if err != nil {
		return r, werr.Wrapf(err, "Erreur query : "+query)
	}
	return r, nil
}

// Modifie l'historique
// si l'url est déjà en base, update la ligne en changeant la date de visite
// sinon ajoute une ligne.
// Ensuite efface les entrées les plus anciennes
func AddRecent(db *sqlx.DB, conf *Config, r *Recent) error {
    var err error
	var count int
	now := time.Now()
	_ = db.QueryRow("select count(*) from recent where url=$1", r.URL).Scan(&count)
    if count == 0 {
        query := `insert into recent(
            url,
            label,
            datevisite
            ) values($1,$2,$3)`
        _, err = db.Exec(
            query,
            r.URL,
            r.Label,
            now)
        if err != nil {
            return werr.Wrapf(err, "Erreur query : "+query)
        }
    } else {
        query := `update recent set datevisite=$1 where url=$2`
        _, err = db.Exec(query, now, r.URL)
        if err != nil {
            return werr.Wrapf(err, "Erreur query : "+query)
        }
    }
    // delete les plus anciens
	all := []*Recent{}
	query := "select * from recent order by datevisite"
	err = db.Select(&all, query)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
    query = "delete from recent where url=$1"
    for i, r := range all{
        if i < conf.NbRecent {
            continue
        }
        _, err := db.Exec(query, r.URL)
        if err != nil {
            return werr.Wrapf(err, "Erreur query : " + query)
        }
    }
	return nil
}
