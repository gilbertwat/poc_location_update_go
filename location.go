package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DriverLocation struct {
	gorm.Model
	Driver string  `from:"driver_id" json:"driver_id" binding:"required`
	Lat    float64 `from:"lat" json:"lat" binding:"required`
	Lng    float64 `from:"lng" json:"lng" binding:"required`
}

func setupRouter(db *gorm.DB, ch chan string) *gin.Engine {

	r := gin.Default()
	r.PUT("/locations", func(c *gin.Context) {
		var driverLocation DriverLocation
		if err := c.ShouldBindJSON(&driverLocation); err == nil {
			db.Create(&driverLocation)
			time.Sleep(10 * time.Second)
			c.JSON(200, driverLocation)
			ch <- fmt.Sprint(driverLocation.ID)
		} else {
			//c.Status(422)
			c.JSON(422, gin.H{"error": err.Error()})
		}
	})

	return r
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=gorm password=postgres sslmode=disable")
	defer db.Close()
	db.AutoMigrate(&DriverLocation{})
	ch := make(chan string)
	go func() {
		for {
			msg := <-ch
			fmt.Println(msg)
		}
	}()
	if err == nil {
		r := setupRouter(db, ch)
		r.Run() // listen and serve on 0.0.0.0:8080
	} else {
		fmt.Println(err)
	}
}
