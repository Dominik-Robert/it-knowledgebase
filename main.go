package main

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dominik-robert/it-knowledgebase/models"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"github.com/sourcegraph/syntaxhighlight"
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
	databaseSchemaName = GetEnvironment("DATABASE_SCHEMA_NAME", "articles")

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

func replaceCodeParts(mdFile []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(mdFile))
	if err != nil {
		return "", err
	}

	// find code-parts via css selector and replace them with highlighted versions
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		oldCode := s.Text()
		formatted, err := syntaxhighlight.AsHTML([]byte(oldCode))
		if err != nil {
			log.Fatal(err)
		}
		s.SetHtml(string(formatted))
	})

	new, err := doc.Html()
	if err != nil {
		return "", err
	}

	// replace unnecessarily added html tags
	new = strings.Replace(new, "<html><head></head><body>", "", 1)
	new = strings.Replace(new, "</body></html>", "", 1)
	return new, nil
}

// test
func main() {
	defer client.Disconnect(context.TODO())
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.LoadHTMLGlob("templates/**/*")

	router.POST("/forms/newArticle", func(c *gin.Context) {
		title := c.PostForm("title")
		subtitle := c.PostForm("subtitle")
		contentMD := c.PostForm("content")
		html := blackfriday.MarkdownCommon([]byte(contentMD))
		replaced, err := replaceCodeParts(html)
		if err != nil {
			log.Fatal(err)
		}

		timestamp := time.Now().Unix()

		database.Collection(databaseSchemaName).InsertOne(context.TODO(), models.Article{
			Title:        title,
			Subtitle:     subtitle,
			ContentMD:    contentMD,
			Content:      template.HTML(replaced),
			Author:       []string{"Dominik Robert"},
			CreatedDate:  timestamp,
			ModifiedDate: timestamp,
		}, options.InsertOne())

		c.Redirect(http.StatusPermanentRedirect, "/")
	})

	router.GET("/admin/:site", func(c *gin.Context) {
		site := c.Param("site")
		c.HTML(http.StatusOK, site+".html", gin.H{})

	})

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
			c.HTML(http.StatusOK, "err.html", gin.H{
				"error": err,
			})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"articles": articles,
				"title":    "Main Webseite",
			})
		}
	})

	router.POST("/", func(c *gin.Context) {
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
