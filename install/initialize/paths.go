/******************************************************************************
    Chemins vers les répertoires utilisés par l'installation

    @copyright  BDL, Bois du
    @license    GPL
    @history    2019-11-05 05:50:37+01:00, Thierry Graff : Creation from a split
********************************************************************************/
package initialize
                                                                                                                                 
import (
	"path"
	"runtime"
	"bdl.local/bdl/ctxt"
)

// getCreateTableDir renvoie le chemin absolu vers le répertoire contenant
// les scripts de création des tables
func getCreateTableDir() string {
	_, filename, _, _ := runtime.Caller(0) // path to current go file
	return path.Join(path.Dir(path.Dir(filename)), "create-tables")
}

// getDataDir renvoie le chemin absolu vers le répertoire contenant
// des fichiers csv pour remplir certaines tables
// (tables "fixes", qui ne vont plus évoluer après remplissage initial).
func getDataDir() string {
	_, filename, _, _ := runtime.Caller(0) // path to current go file
	return path.Join(path.Dir(path.Dir(filename)), "data")
}

// getPrivateDir renvoie le chemin absolu vers le répertoire contenant
// des fichiers contenant des données personnelles
func getPrivateDir() string {
	_, filename, _, _ := runtime.Caller(0) // path to current go file
	return path.Join(path.Dir(path.Dir(filename)), "data-private")
}

// getSCTLDataDir renvoie le chemin absolu vers le répertoire contenant
// des exports de la base SCTL
func getSCTLDataDir(ctx *ctxt.Context, versionSCTL string) string {
    basedir := ctx.Config.Dev.SCTLData
	return path.Join(basedir, "csv-" + versionSCTL)
}
