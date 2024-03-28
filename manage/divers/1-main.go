/*
Point d'entrée pour utiliser les fonctions de ce package.
Pas de gestion CLI, commenter les appels dans le code.

Utilisation :
ENV_CONFIG_FILE='../../config.env' APPLI_CONFIG_FILE='../../config.yml' go run *.go

@copyright  BDL, Bois du Larzac
@history    2024-03-26, Thierry Graff : Creation
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
)

func main() {

	model.MustLoadEnv()
	ctxt.MustLoadConfig()
	ctxt.MustInitDB()
	ctx := ctxt.NewContext()

	//liste_ugs(ctx)
	liste_parcelles(ctx)
}
