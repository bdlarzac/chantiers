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

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Configuration "métier", venant de config.yml
type Config struct {
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
}

// Configuration spécifique au déploiement
// en utilisant les variables d'environnement ou config.env
type serverEnv struct {
	DATABASE_URL      string
	DATABASE_SCHEMA   string
	DATABASE_USER     string
	DATABASE_PASSWORD string
	DATABASE_HOST     string
	DATABASE_PORT     string
	DATABASE_DBNAME   string
	RUN_SERVER_ADDR   string // http.Server
	PORT              string // RUN_PORT
	RUN_MODE          string
	BACKUP_DIR        string
	CMD_PGDUMP        string
}

var SERVER_ENV serverEnv

// MustLoadEnv charge les variables d'environnement spécifiques au serveur
// Dans l'ordre :
// - si CONFIG_FILE existe dans l'environnement : lecture du fichier en question
// - si non lecture config.env s'il existe
// Pas d'erreur s'il n'y a pas de fichier de conf
// - les variables d'environnement : elles sont prioritaires !
func MustLoadEnv() {
	configFile := os.Getenv("ENV_CONFIG_FILE")
	if configFile == "" {
		configFile = "config.env"
	}
	err := godotenv.Load(configFile)
	if err != nil {
		log.Printf("Chargement env %s : %s\n", configFile, err)
	} else {
		log.Printf("Chargement env %s OK\n", configFile)
	}
	SERVER_ENV = serverEnv{
		DATABASE_URL:      os.Getenv("DATABASE_URL"),
		DATABASE_SCHEMA:   os.Getenv("DATABASE_SCHEMA"),
		DATABASE_USER:     os.Getenv("DATABASE_USER"),
		DATABASE_PASSWORD: os.Getenv("DATABASE_PASSWORD"),
		DATABASE_HOST:     os.Getenv("DATABASE_HOST"),
		DATABASE_PORT:     os.Getenv("DATABASE_PORT"),
		DATABASE_DBNAME:   os.Getenv("DATABASE_DBNAME"),
		RUN_SERVER_ADDR:   os.Getenv("RUN_SERVER_ADDR"),
		PORT:              os.Getenv("PORT"), // RUN_PORT
		RUN_MODE:          os.Getenv("RUN_MODE"),
		CMD_PGDUMP:        os.Getenv("CMD_PGDUMP"),
		BACKUP_DIR:        os.Getenv("BACKUP_DIR"),
	}
}
