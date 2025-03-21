package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Chat struct {
	ID        string    `json:"id"`
	Users     []string  `json:"users"`
	Messages  []string  `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateChat(db *sql.DB, chat Chat) (Chat, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := uuid.NewV1()
	if err != nil {
		return Chat{}, err
	}
	chatId := "chat_" + id.String()

	query := "INSERT INTO chats (id, users, messages) VALUES ($1, $2, $3) RETURNING created_at"
	queryErr := db.QueryRowContext(ctx, query, chatId, chat.Users, chat.Messages).Scan(&chat.CreatedAt)
	if queryErr != nil {
		return Chat{}, queryErr
	}
	chat.ID = chatId
	return chat, nil
}

func GetChats(db *sql.DB) ([]Chat, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT id, users, messages, created_at, updated_at FROM chats;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.Users, &chat.Messages, &chat.CreatedAt, &chat.UpdatedAt); err != nil {
			return []Chat{}, err
		}
		chats = append(chats, chat)
	}
	if err := rows.Err(); err != nil {
		return []Chat{}, err
	}
	return chats, nil
}

func UpdateChatByID(db *sql.DB, chat Chat) (Chat, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	UPDATE chats 
	SET users = $1, messages = $2, updated_at = NOW() 
	WHERE id = $4 
	RETURNING id, users, messages, created_at, updated_at`
	var updatedChat Chat
	err := db.QueryRowContext(ctx, query, chat.Users, chat.Messages, chat.UpdatedAt, chat.ID).
		Scan(&updatedChat.ID, &updatedChat.Users, &updatedChat.Messages, &updatedChat.CreatedAt, &updatedChat.UpdatedAt)
	if err != nil {
		return Chat{}, err
	}
	return updatedChat, nil
}

func DeleteChatById(db *sql.DB, id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM chats WHERE id = $1"
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
