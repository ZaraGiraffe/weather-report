package storage

import (
	
)

var HOURLY_FREQUENCY int = 60
var DAILY_FREQUENCY int = 1440

var PENDING_STATUS int = 0
var CONFIRMED_STATUS int = 1

type Subscription struct {
	Id int
    Email string
    City string
    Created_at int64
    Updated_at int64
    Frequency_type int
    Token string
	Status int
}