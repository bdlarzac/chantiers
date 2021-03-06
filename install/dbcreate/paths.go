/******************************************************************************
    Chemins vers les répertoires utilisés par l'installation

    @copyright  BDL, Bois du
    @license    GPL
    @history    2019-11-05 05:50:37+01:00, Thierry Graff : Creation from a split
********************************************************************************/
package dbcreate

import (
	"path"
	"runtime"

	"bdl.local/bdl/ctxt"
)

// GetCreateTableDir renvoie le chemin absolu vers le répertoire contenant
// les scripts de création des tables
func GetCreateTableDir() string {
	_, filename, _, _ := runtime.Caller(0) // path to current go file
	return path.Join(path.Dir(path.Dir(filename)), "db-sql-create")
}

// GetDataDir renvoie le chemin absolu vers le répertoire contenant
// des fichiers csv pour remplir certaines tables
// (tables "fixes", qui ne vont plus évoluer après remplissage initial).
func GetDataDir() string {
	_, filename, _, _ := runtime.Caller(0) // path to current go file
	return path.Join(path.Dir(path.Dir(filename)), "data")
}

// GetPrivateDir renvoie le chemin absolu vers le répertoire contenant
// des fichiers contenant des données personnelles
func GetPrivateDir() string {
	_, filename, _, _ := runtime.Caller(0) // path to current go file
	return path.Join(path.Dir(path.Dir(filename)), "data-private")
}

// GetSCTLDataDir renvoie le chemin absolu vers le répertoire contenant
// des exports de la base SCTL
func GetSCTLDataDir(ctx *ctxt.Context, versionSCTL string) string {
	basedir := ctx.Config.Dev.SCTLData
	return path.Join(basedir, "csv-"+versionSCTL)
}
