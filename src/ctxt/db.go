/******************************************************************************

    @copyright  Thierry Graff
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.

    @history    2019-09-26 23:41:06+02:00, Thierry Graff : Creation
********************************************************************************/

package ctxt

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var db *sqlx.DB

func init() {
	password := ""
	if config.Database.Password != "" {
		password = "password=" + config.Database.Password
	}
	connStr := fmt.Sprintf(
		"dbname=%s user=%s %s host=%s port=%s sslmode=%s",
		config.Database.DbName,
		config.Database.User,
		password,
		config.Database.Host,
		config.Database.Port,
		config.Database.SSLMode,
	)
	var err error

	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		LogError(err)
	}

	db.Exec(fmt.Sprintf(`set search_path='%s'`, config.Database.Schema))

	err = db.Ping()
	if err != nil {
		LogError(err)
	}
}
