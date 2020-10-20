package util

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	MongoDB struct {
		ConnectionString  string
		DBName            string
		ArticleCollection string
		LinkCollection    string
	}
	NLP struct {
		Domain string
	}
}

var Config config

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("ERROR READING CONFIG FILE")
	}
	Config = config{}
	Config.MongoDB.ConnectionString = viper.GetString("MONGO_CONNECTION_STRING")
	Config.MongoDB.DBName = viper.GetString("MONGO_DB_NAME")
	Config.MongoDB.ArticleCollection = viper.GetString("MONGO_ARTICLE_COLLECTION")
	Config.MongoDB.LinkCollection = viper.GetString("MONGO_LINK_COLLECTION")

	Config.NLP.Domain = viper.GetString("NLP_DOMAIN")
}
