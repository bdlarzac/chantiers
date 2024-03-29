/*
Fichier servant à embarquer les fichiers de view/ dans l'exécutable compilé

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.

@history    2021-05-11 15:35:22+01:00, Thierry Graff : Creation
*/
package view

import (
	"embed"
)

//go:embed common/*.js
var ViewFiles embed.FS

//go:embed *.html common/*.html
var TemplatesFiles embed.FS
