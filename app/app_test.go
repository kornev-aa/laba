package app

import (
    "database/sql"
    "testing"
    "fyne-new/mocks"
    "go.uber.org/mock/gomock"
    _ "github.com/mattn/go-sqlite3"
)

func TestSaveMessage_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockRepository(ctrl)
    mockRepo.EXPECT().SaveMessage("тест").Return(nil)

    err := mockRepo.SaveMessage("тест")

    if err != nil {
        t.Errorf("SaveMessage вернул ошибку: %v", err)
    }
}

func TestGetMessages_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockRepository(ctrl)
    expectedMessages := []string{"привет", "мир"}

    mockRepo.EXPECT().GetMessages().Return(expectedMessages, nil)

    messages, err := mockRepo.GetMessages()

    if err != nil {
        t.Errorf("GetMessages вернул ошибку: %v", err)
    }
    if len(messages) != len(expectedMessages) {
        t.Errorf("ожидалось %d сообщений, получили %d", len(expectedMessages), len(messages))
    }
}

func TestGetMessages_Error(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockRepository(ctrl)
    expectedError := sql.ErrConnDone

    mockRepo.EXPECT().GetMessages().Return(nil, expectedError)

    messages, err := mockRepo.GetMessages()

    if err == nil {
        t.Error("ожидалась ошибка, но её нет")
    }
    if messages != nil {
        t.Error("ожидалось nil сообщений, но получили список")
    }
}

// проверяем реальное сохранение
func TestSQLiteRepo_SaveMessage_Real(t *testing.T) {
    // Создаем временную БД в памяти (быстро и не оставляет файлов)
    db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    // Создаем таблицу
    _, err = db.Exec(`CREATE TABLE messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT
    )`)
    if err != nil {
        t.Fatal(err)
    }

    repo := &SQLiteRepo{db: db}
    err = repo.SaveMessage("тестовое сообщение")
    
    if err != nil {
        t.Errorf("SaveMessage вернул ошибку: %v", err)
    }

    // Проверяем, что сообщение действительно сохранилось
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM messages WHERE text = ?", "тестовое сообщение").Scan(&count)
    if err != nil {
        t.Fatal(err)
    }
    if count != 1 {
        t.Errorf("Ожидалось 1 сообщение, получено %d", count)
    }
}

func TestSQLiteRepo_GetMessages_Real(t *testing.T) {
    db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT
    )`)
    if err != nil {
        t.Fatal(err)
    }

    // Добавляем тестовые данные
    _, err = db.Exec("INSERT INTO messages (text) VALUES (?)", "первое")
    if err != nil {
        t.Fatal(err)
    }
    _, err = db.Exec("INSERT INTO messages (text) VALUES (?)", "второе")
    if err != nil {
        t.Fatal(err)
    }

    repo := &SQLiteRepo{db: db}
    messages, err := repo.GetMessages()
    
    if err != nil {
        t.Errorf("GetMessages вернул ошибку: %v", err)
    }
    if len(messages) != 2 {
        t.Errorf("Ожидалось 2 сообщения, получено %d", len(messages))
    }
}
