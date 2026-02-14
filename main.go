package main

import (
	"database/sql"
	"log"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Открыть или создать файл базы данных
	db, err := sql.Open("sqlite3","/home/user/data/mydata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создать таблицу, если её нет
	createTableSQL := `CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	
	// Создать графический интерфейс
	myApp := app.New()
	myWindow := myApp.NewWindow("Заметки с базой данных")
	myWindow.Resize(fyne.NewSize(400, 300))
	
	// Поле для ввода текста
	input := widget.NewEntry()
	input.SetPlaceHolder("Введите сообщение...")
	
	// Список для отображения сообщений
	list := widget.NewLabel("Сохраненные сообщения:")
	output := widget.NewLabel("") // Сами сообщения
	
	// Функ. для обновления списка сообщений в интерфейсе
	refreshMessages := func() {
		rows, err := db.Query("SELECT text FROM messages ORDER BY id DESC")
		if err != nil {
			output.SetText("Ощибка загрузки")
			return
		}
		defer rows.Close()
		
		var messages string
		for rows.Next() {
			var text string
			rows.Scan(&text)
			messages += "- " + text + "\n"
		}
		if messages == "" {
			messages = "Пока нет записей."
		}
		output.SetText(messages)
	}
	
	// Кнопка для сохранения
	saveButton := widget.NewButton("Сохранить", func() {
		if input.Text == "" {
			return
		}
		_, err := db.Exec("INSERT INTO messages (text) VALUES (?)", input.Text)
		if err != nil {
			log.Println("Ошибка сохранения:", err)
			return
		}
		input.SetText("") // Очистить поле ввода
		refreshMessages() // Обновить список
	})
	
	// Сборка в окно
	content := container.NewVBox(
		input,
		saveButton,
		list,
		output,
	)
	
	myWindow.SetContent(content)
	
	// Загрузка сообщений при старте
	refreshMessages()
	
	// Запуск приложения
	myWindow.ShowAndRun()
}
