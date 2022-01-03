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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	database *mongo.Database
)

func init() {
	mongodb_host := GetEnvironment("MONGODB_HOST", "localhost")
	mongodb_port := GetEnvironment("MONGODB_PORT", "27017")
	mongodb_user := GetEnvironment("MONGODB_USER", "")
	mongodb_password := GetEnvironment("MONGODB_PASSWORD", "")
	mongodb_database := GetEnvironment("MONGODB_PASSWORD", "itknowledgebase")

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
	router.LoadHTMLGlob("templates/**/*")

	router.GET("/", func(c *gin.Context) {
		articles, err := GetArticles(bson.M{}, options.Find())

		if err != nil {
			c.HTML(http.StatusOK, "err.html", gin.H{})
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"articles": articles,
			"title":    "Main Webseite",
		})
	})
	router.GET("/detail/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail.html", gin.H{})
	})

	router.StaticFS("/assets", http.Dir("assets/"))

	router.Run("localhost:8080")
}

func GetEnvironment(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}

func GetArticles(filter bson.M, findOptions *options.FindOptions) ([]models.Article, error) {
	rows, err := database.Collection("articles").Find(context.TODO(), filter, findOptions)

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
