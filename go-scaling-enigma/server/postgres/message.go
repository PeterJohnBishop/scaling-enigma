package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Message struct {
	ID        string    `json:"id"`
	Chat      string    `json:"chat"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Media     []string  `json:"media"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateMessage(db *sql.DB, message Message) (Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := uuid.NewV1()
	if err != nil {
		return Message{}, err
	}
	msgId := "msg_" + id.String()

	query := "INSERT INTO messages (id, chat, sender, text, media) VALUES ($1, $2, $3, $4, $5) RETURNING created_at"
	err = db.QueryRowContext(ctx, query, msgId, message.Chat, message.Sender, message.Text, message.Media).Scan(&message.CreatedAt)
	if err != nil {
		return Message{}, err
	}
	return message, nil
}

func GetMessages(db *sql.DB) ([]Message, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT id, chat, sender, text, media, created_at FROM messages;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Chat, &msg.Sender, &msg.Text, &msg.Media, &msg.CreatedAt); err != nil {
			return []Message{}, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return []Message{}, err
	}
	return messages, nil
}

func DeleteMessageById(db *sql.DB, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM messages WHERE id = $1"
	res, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}
	return nil
}
