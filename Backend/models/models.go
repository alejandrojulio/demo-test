package models

import (
	"time"
)

type Stock struct {
	//gorm.Model
	Id         int
	Ticker     string    
	TargetFrom string   
	TargetTo   string    
	Company    string    
	Action     string    
	Brokerage  string    
	RatingFrom string   
	RatingTo   string    
	Time       time.Time 
}

type StockRecomendation struct {
	Id         int      
	Ticker     string    
	TargetFrom float64  
	TargetTo   float64 
	Company    string   
	Action     string    
	Brokerage  string    
	RatingFrom string    
	RatingTo   string    
	Time       time.Time 
	PotentialGain float64 
}

type ApiResponse struct {
	Items    []Stock `json:"items"`
	NextPage string `json:"next_page"`
}