package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Translation struct {
	Language string `json:"language"`
	Greeting string `json:"greeting"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"greeting": "Hello, world!",
		})
	})

	r.POST("/translations/create", createTranslations)
	r.GET("/translate/:language", getTranslation)

	r.Run()
}

func getTranslation(c *gin.Context) {
	language := c.Param("language")
	db, err := gorm.Open("postgres", postgresConnectionString())
	if err != nil {
		log.Printf("Error connecting to PostgreSQL with connection string '%s': %s", postgresConnectionString(), err)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Could not connect to PostgreSQL database",
			},
		)
		return
	}
	defer db.Close()
	db.AutoMigrate(&Translation{})

	var translation Translation

	if err := db.Where(&Translation{Language: language}).First(&translation).Error; err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Could not find translation for language " + language,
			},
		)
		log.Println(err)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"greeting": translation.Greeting,
		})
	}
}

func createTranslations(c *gin.Context) {
	db, err := gorm.Open("postgres", postgresConnectionString())
	if err != nil {
		log.Printf("Error connecting to PostgreSQL with connection string '%s': %s", postgresConnectionString(), err)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Could not connect to PostgreSQL database",
			},
		)
		return
	}
	defer db.Close()
	db.AutoMigrate(&Translation{})

	translation := Translation{
		Language: "es",
		Greeting: "Hola, mundo!",
	}

	if err := db.Create(&translation).Error; err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Could not create initial translation",
			},
		)
		log.Println(err)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"translation_language": translation.Language,
		})
	}
}

func postgresConnectionString() string {
	return "host=" + postgresHost() + " port=" + postgresPort() + " user=" + postgresUser() + " dbname=" + postgresDbName() + " password=" + postgresDbPassword()
}

func postgresHost() string {
	return postgresConnectionValue("HOST")
}

func postgresPort() string {
	return postgresConnectionValue("PORT")
}

func postgresUser() string {
	return postgresConnectionValue("USER")
}

func postgresDbName() string {
	return postgresConnectionValue("NAME")
}

func postgresDbPassword() string {
	return postgresConnectionValue("PASSWORD")
}

func postgresConnectionValue(key string) string {
	return os.Getenv("DB_" + key)
}
