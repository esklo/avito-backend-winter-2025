package repository

import (
	"context"
	"fmt"
	"sort"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
)

func (r *repo) MakeTransfer(ctx context.Context, tx DB, senderID, receiverID int, amount int) error {
	db := r.getExecutor(tx)

	//todo: cleanup
	//start := time.Now()
	//_, err := db.Exec(ctx, `
	//    SELECT * FROM users
	//    WHERE id IN ($1, $2)
	//    ORDER BY id
	//    FOR UPDATE`, senderID, receiverID)
	//if err != nil {
	//	return err
	//}
	//
	//duration := time.Since(start)
	//if duration > 100*time.Millisecond {
	//	// log.Printf("SLOW SELECT FOR UPDATE: sender=%d receiver=%d duration=%v\n",
	//	//	senderID, receiverID, duration)
	//}

	_, err := db.Exec(ctx, `
		INSERT INTO transfers (sender_id, receiver_id, amount)
		VALUES ($1, $2, $3)
		ON CONFLICT (sender_id, receiver_id)
		DO UPDATE SET amount = transfers.amount + $3;
	`, senderID, receiverID, amount)
	if err != nil {
		return fmt.Errorf("make transfer: %w", err)
	}

	changes := []struct{ id, amount int }{
		{id: senderID, amount: -amount},
		{id: receiverID, amount: amount},
	}
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].id < changes[j].id
	})

	for _, change := range changes {
		_, err = db.Exec(ctx, `
		UPDATE users SET balance = balance + $2
		WHERE id = $1;
	`, change.id, change.amount)
		if err != nil {
			return fmt.Errorf("update balance %d: %w", change.amount, err)
		}
	}

	return nil
}

func (r *repo) ListTransactions(ctx context.Context, tx DB, userID int) (*model.CoinHistory, error) {
	db := r.getExecutor(tx)
	// todo: can it be simplified?
	rows, err := db.Query(ctx, `
		SELECT 
			'sent' as type,
			users.username,
			transfers.amount
		FROM transfers
		JOIN users ON users.id = transfers.receiver_id
		WHERE transfers.sender_id = $1
		
		UNION ALL
		
		SELECT 
			'received' as type,
			users.username,
			transfers.amount
		FROM transfers
		JOIN users ON users.id = transfers.sender_id
		WHERE transfers.receiver_id = $1;
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("select transactions: %w", err)
	}
	defer rows.Close()

	history := &model.CoinHistory{
		Received: make([]model.CoinsReceived, 0),
		Sent:     make([]model.CoinsSent, 0),
	}

	for rows.Next() {
		var (
			txType, username string
			amount           int
		)

		if err := rows.Scan(&txType, &username, &amount); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}

		switch txType {
		case "received":
			history.Received = append(history.Received, model.CoinsReceived{
				FromUser: username,
				Amount:   amount,
			})
		case "sent":
			history.Sent = append(history.Sent, model.CoinsSent{
				ToUser: username,
				Amount: amount,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transactions: %w", err)
	}

	return history, nil
}
