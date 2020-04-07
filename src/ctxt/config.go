/******************************************************************************
    Chargement de config.yml

    contains package init()

    @copyright  BDL, Bois du Larzac
    @license    GPL
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
