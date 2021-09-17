/******************************************************************************
    Fichier servant à embarquer les fichiers de view/ dans l'exécutable compilé

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2021-05-11 15:35:22+01:00, Thierry Graff : Creation
********************************************************************************/
package view

import (
    "embed"
)

//go:embed *
var ViewFiles embed.FS
