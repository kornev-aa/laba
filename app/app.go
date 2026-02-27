package app

import "database/sql"

type Repository interface {
	SaveMessage(text string) error
	GetMessages() ([]string, error)
}

type SQLiteRepo struct {
	db *sql.DB
}

func (r *SQLiteRepo) SaveMessage(text string) error {
	_, err := r.db.Exec("INSERT INTO messages (text) VALUES (?)", text)
	return err
}

func (r *SQLiteRepo) GetMessages() ([]string, error) {
	rows, err := r.db.Query("SELECT text FROM messages ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var text string
		err := rows.Scan(&text)
		if err != nil {
			return nil, err
		}
		messages = append(messages, text)
	}
	return messages, nil
}
