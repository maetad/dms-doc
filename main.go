package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gitlab.com/pakkaparn/dms-doc/doc"
	"gitlab.com/pakkaparn/dms-doc/mq"
	"gitlab.com/pakkaparn/dms-doc/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(".env.database", ".env"); err != nil {
		log.Fatal("Error loading environment file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", os.Getenv("DB_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("DB_PORT"), os.Getenv("DB_TIMEZONE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&doc.Doc{})
	db.AutoMigrate(&user.User{})

	r := gin.Default()

	docHandler := doc.NewDocHandler(db)

	r.GET("", docHandler.GetDoc)
	r.POST("", docHandler.CreateDoc)
	r.GET("/:id", docHandler.FindDoc)
	r.PUT("/:id", docHandler.UpdateDoc)
	r.DELETE("/:id", docHandler.DeleteDoc)

	mq.Received(db)

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
