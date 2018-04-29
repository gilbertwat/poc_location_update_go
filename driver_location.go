package main

import "github.com/jinzhu/gorm"

type DriverLocation struct {
	gorm.Model
	Driver string  `from:"driver_id" json:"driver_id" binding:"required`
	Lat    float64 `from:"lat" json:"lat" binding:"required`
	Lng    float64 `from:"lng" json:"lng" binding:"required`
}
