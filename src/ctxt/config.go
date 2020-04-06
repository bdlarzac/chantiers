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

/*
type Config struct {
    Database struct{
        Host      string `yaml:"host"`
        Port      string `yaml:"port"`
        User      string `yaml:"user"`
        Password  string `yaml:"password"`
        DbName    string `yaml:"dbname"`
    } `yaml:"database"`
    Paths struct{
        LogicielFoncier string `yaml:"logiciel-foncier"`
        CoucheTypo      string `yaml:"couche-typo"`
    } `yaml:"paths"`
    PourcentagePerte    float64 `yaml:"pourcentage-perte"`
    TVAExt              []float64 `yaml:"tva-ext"`
    TVABDL struct {
        Livraison           float64 `yaml:"livraison"`
        VentePlaquettes     float64 `yaml:"vente-plaquettes"`
        BoisSurPied         float64 `yaml:"bois-sur-pied"`
        AutreValorisation   float64 `yaml:"autre-valorisation"`
    } `yaml:"tva-bdl"`
    Facture struct{
        // Concerne le document PDF généré ; pas affiché
        Auteur      string `yaml:"auteur"`
        Createur    string `yaml:"créateur"`
        // Infos affichées sur les factures
        Adresse     string `yaml:"adresse"`
        Tel         string `yaml:"tel"`
        Email       string `yaml:"email"`
        SiteWeb     string `yaml:"site-web"`
        Siret       string `yaml:"siret"`
        TVA         string `yaml:"tva"`
    } `yaml:"facture"`
}
*/

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
