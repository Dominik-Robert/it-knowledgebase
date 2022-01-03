package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dominik-robert/it-knowledgebase/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client             *mongo.Client
	database           *mongo.Database
	http_host          string
	http_port          string
	databaseSchemaName string
)

func init() {
	mongodb_host := GetEnvironment("MONGODB_HOST", "localhost")
	mongodb_port := GetEnvironment("MONGODB_PORT", "27017")
	mongodb_user := GetEnvironment("MONGODB_USER", "")
	mongodb_password := GetEnvironment("MONGODB_PASSWORD", "")
	mongodb_database := GetEnvironment("MONGODB_PASSWORD", "itknowledgebase")
	http_host = GetEnvironment("HTTP_HOST", "0.0.0.0")
	http_port = GetEnvironment("HTTP_PORT", "8080")
	databaseSchemaName = GetEnvironment("HTTP_PORT", "articles")

	var err error
	if mongodb_user != "" {
		credential := options.Credential{
			Username: mongodb_user,
			Password: mongodb_password,
		}
		client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://" + mongodb_host + ":" + mongodb_port).SetAuth(credential))
	} else {
		client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://" + mongodb_host + ":" + mongodb_port))
	}

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	database = client.Database(mongodb_database)
}

// test
func main() {
	defer client.Disconnect(context.TODO())
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.LoadHTMLGlob("templates/**/*")

	router.GET("/article/:id", func(c *gin.Context) {
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		article, err := GetArticles(bson.M{"_id": id}, options.Find())
		if err != nil {
			c.HTML(http.StatusOK, "err.html", gin.H{})
		} else {
			c.HTML(http.StatusOK, "detail.html", gin.H{
				"articles": article[0],
				"title":    "Main Webseite",
			})
		}
	})

	router.GET("/", func(c *gin.Context) {
		articles, err := GetArticles(bson.M{}, options.Find())

		if err != nil {
			c.HTML(http.StatusOK, "err.html", gin.H{})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"articles": articles,
				"title":    "Main Webseite",
			})
		}
	})

	router.StaticFS("/assets", http.Dir("assets/"))

	router.Run(http_host + ":" + http_port)
}

func GetEnvironment(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}

func GetArticles(filter bson.M, findOptions *options.FindOptions) ([]models.Article, error) {
	rows, err := database.Collection(databaseSchemaName).Find(context.TODO(), filter, findOptions)

	if err != nil {
		return nil, err
	}

	var article []models.Article
	err = rows.All(context.TODO(), &article)

	if err != nil {
		return nil, err
	}

	return article, nil
}
