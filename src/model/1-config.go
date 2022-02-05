/******************************************************************************
    Structure de donnée recevant le contenu de config.yml
    Pourrait aussi être dans le package ctxt.
    Mis dans model pour éviter un problème de dépendance circulaire
    lorsqu'elle est utilisée dans model.

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2019-09-26, Thierry Graff : Creation
********************************************************************************/
package model

type Config struct {
	Run struct {
		URL  string `yaml:"url"`
		Port string `yaml:"port"`
	}
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DbName   string `yaml:"dbname"`
		Schema   string `yaml:"schema"`
		SSLMode  string `yaml:"ssl-mode"`
		Backup   struct {
			Directory string `yaml:"directory"`
			CmdPgdump string `yaml:"cmd-pgdump"`
		} `yaml:"backup"`
	} `yaml:"database"`
	Paths struct {
		LogicielFoncier string `yaml:"logiciel-foncier"`
	} `yaml:"paths"`
	PourcentagePerte float64   `yaml:"pourcentage-perte"`
	DebutSaison      string    `yaml:"debut-saison"`
	TVAExt           []float64 `yaml:"tva-ext"`
	TVABDL           struct {
		Livraison           float64   `yaml:"livraison"`
		VentePlaquettes     float64   `yaml:"vente-plaquettes"`
		AutresValorisations []float64 `yaml:"autres-valorisations"`
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
	Affacture struct {
		Adresse string `yaml:"adresse"`
	} `yaml:"affacture"`
	NbRecent int `yaml:"nb-recent"`
	Dev      struct {
		SCTLData    string `yaml:"sctl-data"`
		SCTLAnalyse string `yaml:"sctl-analyse"`
	} `yaml:"dev"`
}
