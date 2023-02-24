/*
*****************************************************************************

	Fichier servant à embarquer les fichiers de static/ dans l'exécutable compilé

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.

	@history    2021-05-11 15:31:29+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package static

import (
	"embed"
)

//go:embed *
var StaticFiles embed.FS
