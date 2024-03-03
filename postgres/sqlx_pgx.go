package main

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type CreateTodoItemDto struct {
	Content string
}

func main() {
	// Создаём подключение
	psqlInfoTemplate := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	db, err := sqlx.Connect(
		"pgx",
		fmt.Sprintf(psqlInfoTemplate, "localhost", 5432, "todo-user", "todo-pass", "todo"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Создаём таблицу
	db.MustExec(schema)

	// Начинаем транзакцию
	tx := db.MustBegin()

	tx.NamedExec("INSERT INTO todo_item (content) VALUES (:content)", &CreateTodoItemDto{Content: "Hello World"})
	_, err = db.NamedExec(`INSERT INTO todo_item (content) VALUES (:content)`,
		map[string]interface{}{
			"content": "Map parameter",
		})

	id := "dce060a7-e776-40f4-bbb4-59358e777cc1"
	tx.MustExec("INSERT INTO todo_item (id, content) VALUES ($1, $2)", id, "Get by ID")
	tx.MustExec("INSERT INTO todo_item (content) VALUES ($1)", "Привет мир")

	// Закрываем транзакцию
	tx.Commit()

	items := []TodoItem{}
	db.Select(&items, "SELECT * FROM todo_item ORDER BY created_at ASC")
	fmt.Printf("%s\n", items[0].String())

	item := TodoItem{}
	err = db.Get(&item, "SELECT * FROM todo_item WHERE id=$1", id)
	fmt.Printf("%s\n", item.String())

	item = TodoItem{}
	rows, err := db.Queryx("SELECT * FROM todo_item")
	for rows.Next() {
		err := rows.StructScan(&item)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s\n", item.String())
	}
}
