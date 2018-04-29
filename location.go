package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type updateChannels struct {
	M    sync.Mutex
	List *[]chan bool
}

func setupRouter(db *gorm.DB, resultBroadcaster chan chan bool) *gin.Engine {

	r := gin.Default()
	r.PATCH("/order", func(c *gin.Context) {
		myChannel := make(chan bool)
		defer close(myChannel)
		resultBroadcaster <- myChannel
		isSuccess := <-myChannel
		c.JSON(200, gin.H{"success": isSuccess})
	})

	r.PUT("/locations", func(c *gin.Context) {
		var driverLocation DriverLocation
		if err := c.ShouldBindJSON(&driverLocation); err == nil {
			db.Create(&driverLocation)
			c.JSON(200, driverLocation)
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
	resultBroadcaster := make(chan chan bool)
	defer close(resultBroadcaster)
	resultSubscriptors := make([]chan bool, 0)
	lotteryChannelList := updateChannels{sync.Mutex{}, &resultSubscriptors}

	go lotteryTicketCollector(&resultBroadcaster, &lotteryChannelList)
	go lotteryDrawer(&lotteryChannelList)

	if err == nil {
		r := setupRouter(db, resultBroadcaster)
		r.Run() // listen and serve on 0.0.0.0:8080
	} else {
		fmt.Println(err)
	}
}
