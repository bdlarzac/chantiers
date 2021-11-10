/******************************************************************************

    @copyright  BDL, Bois du Larzac
    @license    GPL
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
	connStr := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%s sslmode=%s",
		config.Database.DbName,
		config.Database.User,
		config.Database.Password,
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
