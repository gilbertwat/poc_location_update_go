package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
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

type updateChannels struct {
	M    sync.Mutex
	List *[]chan bool
}

func setupRouter(db *gorm.DB, lotteryUpdateChannels chan chan bool, number *int64) *gin.Engine {

	r := gin.Default()
	r.PATCH("/order", func(c *gin.Context) {
		currentLotteryID := atomic.LoadInt64(number)
		atomic.AddInt64(number, 1)
		fmt.Printf("curent lottery id: %d\n", currentLotteryID)
		myChannel := make(chan bool)
		defer close(myChannel)
		lotteryUpdateChannels <- myChannel
		isSuccess := <-myChannel
		c.JSON(200, gin.H{"currentId": currentLotteryID, "success": isSuccess})
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
	var number int64
	lotteryUpdateChannels := make(chan chan bool)
	defer close(lotteryUpdateChannels)
	list := make([]chan bool, 0)
	lotteryChannelList := updateChannels{sync.Mutex{}, &list}
	go func() {
		for {
			newChannel := <-lotteryUpdateChannels
			lotteryChannelList.M.Lock()
			list = append(*lotteryChannelList.List, newChannel)
			lotteryChannelList.List = &list
			lotteryChannelList.M.Unlock()
		}
	}()

	go func() {
		for {
			fmt.Println("First order")
			time.Sleep(time.Second * 5)
			fmt.Println("Times up!")
			totalNumberOfDriver := atomic.LoadInt64(&number)
			fmt.Printf("Total contestant: %d\n", totalNumberOfDriver)
			fmt.Printf("Total contestant: %d\n", len(*lotteryChannelList.List))
			if totalNumberOfDriver > 0 {
				successLotteryID := rand.Int63n(totalNumberOfDriver)
				fmt.Printf("Lucky guy: %d\n", successLotteryID)
				for i := int64(0); i < totalNumberOfDriver; i++ {
					(*lotteryChannelList.List)[i] <- successLotteryID == i
				}
			}
			number = 0
			newList := make([]chan bool, 0)
			lotteryChannelList.List = &newList
		}
	}()
	if err == nil {
		r := setupRouter(db, lotteryUpdateChannels, &number)
		r.Run() // listen and serve on 0.0.0.0:8080
	} else {
		fmt.Println(err)
	}
}
