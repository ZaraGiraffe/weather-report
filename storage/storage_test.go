package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"log"
)

func Test_InsertAndDeleteSubscription(t *testing.T) {
	db := NewStorageConnection("../test.config.json")

	subscription := Subscription{
		Email:     "test_insert@test.com",
		City:      "London",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		FrequencyType: 1,
		Token:         "insert_token",
		Status:        1,
	}

	err := InsertSubscriptionQuery(db, &subscription)
	assert.Nil(t, err)

	err = DeleteSubscriptionByToken(db, "insert_token")
	assert.Nil(t, err)
}

func Test_GetSubscription(t *testing.T) {
	db := NewStorageConnection("../test.config.json")

	newSubscription := Subscription{
		Email:     "test_get@test.com",
		City:      "London",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		FrequencyType: 1,
		Token:         "get_token",
		Status:        1,
	}

	err := InsertSubscriptionQuery(db, &newSubscription)
	assert.Nil(t, err)

	subscription, err := GetSubscriptionByToken(db, "get_token")
	log.Printf("subscription: %v", subscription)
	assert.Nil(t, err)
	assert.Equal(t, subscription.Email, "test_get@test.com")

	subscription, err = GetSubscriptionByEmail(db, "test_get@test.com")
	assert.Nil(t, err)
	assert.Equal(t, subscription.Token, "get_token")

	err = DeleteSubscriptionByToken(db, "get_token")
	assert.Nil(t, err)
}

func Test_UpdateSubscription(t *testing.T) {
	db := NewStorageConnection("../test.config.json")

	newSubscription := Subscription{
		Email:     "test_update@test.com",
		City:      "London",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		FrequencyType: 1,
		Token:         "update_token",
		Status:        0,
	}
	
	err := InsertSubscriptionQuery(db, &newSubscription)
	assert.Nil(t, err)

	err = UpdateSubscriptionStatus(db, "update_token", 1)
	assert.Nil(t, err)

	subscription, err := GetSubscriptionByToken(db, "update_token")
	assert.Nil(t, err)
	assert.Equal(t, subscription.Status, 1)

	err = DeleteSubscriptionByToken(db, "update_token")
	assert.Nil(t, err)
}
