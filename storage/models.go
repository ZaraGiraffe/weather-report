// This package is responsible for storing the models for the database
package storage

type Subscription struct {
	Id            int
	Email         string
	City          string
	CreatedAt     int64
	UpdatedAt     int64
	FrequencyType int
	Token         string
	Status        int
}
