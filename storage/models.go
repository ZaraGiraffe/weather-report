// This package is responsible for storing the models for the database
package storage

import (
	
)

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