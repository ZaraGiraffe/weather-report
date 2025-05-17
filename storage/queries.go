package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InsertDubscriptionQuery(db *sql.DB, subscription *Subscription) error {
	_, err := db.Exec(fmt.Sprintf(
		`INSERT INTO subscriptions (email, city, created_at, updated_at, frequency_type, token, status) VALUES ('%s', '%s', %d, %d, %d, '%s', %d)`,
		subscription.Email,
		subscription.City,
		subscription.Created_at,
		subscription.Updated_at,
		subscription.Frequency_type,
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
		`SELECT id FROM subscriptions WHERE token = $1`, token,
	)

	subscription := &Subscription{}
	err := row.Scan(&subscription.Id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		log.Fatalf("ERROR: scan subscription query failed: %v", err)
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
		log.Fatalf("ERROR: update subscription status query failed: %v", err)
	}
	return nil
}

func DeleteSubscription(db *sql.DB, token string) error {
	_, err := db.Exec(fmt.Sprintf(`DELETE FROM subscriptions WHERE token = '%s'`, token))
	if err != nil {
		log.Fatalf("ERROR: delete subscription query failed: %v", err)
	}
	return nil
}

func GetSubscriptionByEmail(db *sql.DB, email string) (*Subscription, error) {
	row := db.QueryRow(
		`SELECT id FROM subscriptions WHERE email = $1`, email,
	)
	subscription := &Subscription{}
	err := row.Scan(&subscription.Id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		log.Fatalf("ERROR: scan subscription query failed: %v", err)
	}
	return subscription, nil
}