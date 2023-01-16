/******************************************************************************
    Chargement de config.yml

    Contains package init()

    @copyright  Thierry Graff
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.

    @history    2019-09-26, Thierry Graff : Creation
********************************************************************************/
package ctxt

import (
	"io/ioutil"
	"log"
	"os"

	"bdl.local/bdl/model"
	"gopkg.in/yaml.v3"
)

var config *model.Config

func MustLoadConfig() {

	configFile := os.Getenv("APPLI_CONFIG_FILE")
	if configFile == "" {
		configFile = "../config.yml"
	}

	y, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("Chargement conf %s : %s\n", configFile, err)
	} else {
		log.Printf("Chargement conf %s OK\n", configFile)
	}

	err = yaml.Unmarshal(y, &config)
	if err != nil {
		panic(err)
	}
}
