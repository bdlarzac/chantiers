/******************************************************************************

    @copyright  Thierry Graff
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.

    @history    2019-09-26 23:41:06+02:00, Thierry Graff : Creation
********************************************************************************/

package ctxt

import (
	"log"
	"strings"

	"bdl.local/bdl/model"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var db *sqlx.DB

// ajoute le schema a l'url
// et supprime sslmode=prefer mis automatiquement par Scalingo
func ajusteDbURL(url string, schema string) string {
	dbURL := strings.ReplaceAll(url, "sslmode=prefer", "") // non supporté par lib/pq
	if strings.HasPrefix(dbURL, "postgres") && !strings.Contains(dbURL, "search_path") {
		if !strings.Contains(dbURL, "?") {
			dbURL += "?"
		}
		dbURL += "&search_path=" + schema
	}
	return dbURL
}

func MustInitDB() {
	var err error

	dbURL := ajusteDbURL(model.SERVER_ENV.DATABASE_URL, model.SERVER_ENV.SCHEMA)

	db, err = sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Connexion DB impossible : %v", err)
	}

	//TODO: ici faire upgrade versions
}
