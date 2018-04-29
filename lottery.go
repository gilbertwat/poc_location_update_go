package main

import (
	"fmt"
	"math/rand"
	"time"
)

func lotteryTicketCollector(lotteryUpdateChannels *chan chan bool, lotteryChannelList *updateChannels) {
	for {
		newChannel := <-*lotteryUpdateChannels
		lotteryChannelList.M.Lock()
		list := append(*lotteryChannelList.List, newChannel)
		lotteryChannelList.List = &list
		lotteryChannelList.M.Unlock()
	}
}

func lotteryDrawer(lotteryChannelList *updateChannels) {
	for {
		fmt.Println("Lottery starts")
		time.Sleep(time.Second * 5)
		fmt.Println("Times up!")
		totalNumberOfDriver := len(*lotteryChannelList.List)
		fmt.Printf("Total contestant: %d\n", totalNumberOfDriver)
		if totalNumberOfDriver > 0 {
			successLotteryID := rand.Intn(totalNumberOfDriver)
			fmt.Printf("Lucky guy: %d\n", successLotteryID)
			for i := 0; i < totalNumberOfDriver; i++ {
				(*lotteryChannelList.List)[i] <- successLotteryID == i
			}
		}
		newList := make([]chan bool, 0)
		lotteryChannelList.List = &newList
	}
}
