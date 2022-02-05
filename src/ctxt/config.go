/******************************************************************************
    Chargement de config.yml

    Contains package init()

    @copyright  Thierry Graff
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    
    @history    2019-09-26, Thierry Graff : Creation
********************************************************************************/
package ctxt

import (
	"bdl.local/bdl/model"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var config *model.Config

func init() {
	y, err := ioutil.ReadFile("../config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(y, &config)
	if err != nil {
		panic(err)
	}
}
