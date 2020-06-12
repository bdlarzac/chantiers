/******************************************************************************
    Structure de donnée recevant le contenu de config.yml
    Pourrait aussi être dans le package ctxt.
    Mis dans model pour éviter un problème de dépendance circulaire
    lorsqu'elle est utilisée dans model.

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-09-26, Thierry Graff : Creation
********************************************************************************/
package model

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DbName   string `yaml:"dbname"`
	} `yaml:"database"`
	Paths struct {
		LogicielFoncier string `yaml:"logiciel-foncier"`
		CoucheTypo      string `yaml:"couche-typo"`
	} `yaml:"paths"`
	PourcentagePerte float64   `yaml:"pourcentage-perte"`
	TVAExt           []float64 `yaml:"tva-ext"`
	TVABDL           struct {
		Livraison         float64 `yaml:"livraison"`
		VentePlaquettes   float64 `yaml:"vente-plaquettes"`
		BoisSurPied       float64 `yaml:"bois-sur-pied"`
		AutreValorisation float64 `yaml:"autre-valorisation"`
	} `yaml:"tva-bdl"`
	Facture struct {
		// metadata - pas affiché
		Auteur   string `yaml:"auteur"`
		Createur string `yaml:"createur"`
		// Infos affichées sur les factures
		Adresse string `yaml:"adresse"`
		Tel     string `yaml:"tel"`
		Email   string `yaml:"email"`
		SiteWeb string `yaml:"site-web"`
		Siret   string `yaml:"siret"`
		TVA     string `yaml:"tva"`
	} `yaml:"facture"`
	NbRecent int `yaml:"nb-recent"`
}
