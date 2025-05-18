// This package is responsible for storing the queries for the database
package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InsertSubscriptionQuery(db *sql.DB, subscription *Subscription) error {
	_, err := db.Exec(fmt.Sprintf(
		`INSERT INTO subscriptions (email, city, created_at, updated_at, frequency_type, token, status) VALUES ('%s', '%s', %d, %d, %d, '%s', %d)`,
		subscription.Email,
		subscription.City,
		subscription.CreatedAt,
		subscription.UpdatedAt,
		subscription.FrequencyType,
		subscription.Token,
		subscription.Status,
	))
	if err != nil {
		log.Printf("ERROR: insert subscription query failed: %v", err)
		return err
	}
	return nil
}

func GetSubscriptionByToken(db *sql.DB, token string) (*Subscription, error) {
	row := db.QueryRow(
		`SELECT id, email, token, city, created_at, updated_at, frequency_type, status FROM subscriptions WHERE token = $1`, token,
	)

	subscription := &Subscription{}
	err := row.Scan(&subscription.Id, &subscription.Email, &subscription.Token, &subscription.City, &subscription.CreatedAt, &subscription.UpdatedAt, &subscription.FrequencyType, &subscription.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		log.Printf("ERROR: get subscription by token query failed: %v", err)
		return nil, err
	}
	return subscription, nil
}

func UpdateSubscriptionStatus(db *sql.DB, token string, status int) error {
	_, err := db.Exec(fmt.Sprintf(
		`UPDATE subscriptions SET status = %d WHERE token = '%s'`,
		status,
		token,
	))
	if err != nil {
		log.Printf("ERROR: update subscription status query failed: %v", err)
		return err
	}
	return nil
}

func DeleteSubscriptionByToken(db *sql.DB, token string) error {
	_, err := db.Exec(fmt.Sprintf(`DELETE FROM subscriptions WHERE token = '%s'`, token))
	if err != nil {
		log.Printf("ERROR: delete subscription query failed: %v", err)
		return err
	}
	return nil
}

func GetSubscriptionByEmail(db *sql.DB, email string) (*Subscription, error) {
	row := db.QueryRow(
		`SELECT id, email, token, city, created_at, updated_at, frequency_type, status FROM subscriptions WHERE email = $1`, email,
	)
	subscription := &Subscription{}
	err := row.Scan(&subscription.Id, &subscription.Email, &subscription.Token, &subscription.City, &subscription.CreatedAt, &subscription.UpdatedAt, &subscription.FrequencyType, &subscription.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		log.Printf("ERROR: get subscription by email query failed: %v", err)
		return nil, err
	}
	return subscription, nil
}

func GetAllSubscriptionsWithTimeConstraint(db *sql.DB, timeConstraint int64, frequencyType int) ([]*Subscription, error) {
	rows, err := db.Query(
		`SELECT * FROM subscriptions WHERE updated_at > $1 AND frequency_type = $2`,
		timeConstraint,
		frequencyType,
	)
	if err != nil {
		log.Printf("ERROR: get all subscriptions with time constraint query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	subscriptions := []*Subscription{}
	for rows.Next() {
		subscription := &Subscription{}
		err := rows.Scan(&subscription.Id, &subscription.Email, &subscription.City, &subscription.CreatedAt, &subscription.UpdatedAt, &subscription.FrequencyType, &subscription.Token, &subscription.Status)
		if err != nil {
			log.Printf("ERROR: get all subscriptions with time constraint query failed: %v", err)
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}

func UpdateSubscriptionLastSent(db *sql.DB, token string, lastSent int64) error {
	_, err := db.Exec(fmt.Sprintf(
		`UPDATE subscriptions SET updated_at = %d WHERE token = '%s'`,
		lastSent,
		token,
	))
	if err != nil {
		log.Printf("ERROR: update subscription was deleted: %v", err)
		return nil
	}
	return nil
}
