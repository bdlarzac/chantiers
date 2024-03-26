/*
	@copyright  BDL, Bois du Larzac
	@history    2019-09-26 17:41:35+02:00, Thierry Graff : Creation
*/
package main

import (
	"bdl.local/bdl/ctxt"
	// "bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/model"
	//"fmt"
	// "os"
	// "os/exec"
	// "regexp"
	// "sort"
	// "strings"
)

func main() {

	model.MustLoadEnv()
	ctxt.MustLoadConfig()
	ctxt.MustInitDB()
	ctx := ctxt.NewContext()

	liste_ugs(ctx)
}

