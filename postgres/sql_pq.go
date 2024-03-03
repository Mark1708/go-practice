package main

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/lib/pq"
)

func main() {
	// Создаём подключение
	//psqlInfoTemplate := "postgresql://%s:%s@%s:%d/%s?sslmode=disable"
	psqlInfoTemplate := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	db, err := sql.Open(
		"postgres",
		//fmt.Sprintf(psqlInfoTemplate, "todo-user", "todo-pass", "localhost", 5432, "todo"),
		fmt.Sprintf(psqlInfoTemplate, "localhost", 5432, "todo-user", "todo-pass", "todo"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Создаём таблицу
	_, _ = db.Exec(schema)

	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	fmt.Println("\nPrepare statement and insert rows")
	stmt, err := tx.Prepare("INSERT INTO todo_item (content) VALUES ($1)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	contents := []string{"hello", "world"}
	for _, content := range contents {
		_, err = stmt.Exec(content)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("\nInsert row with param")
	id := "dce060a7-e776-40f4-bbb4-59358e777cc1"
	result, err := db.Exec("INSERT INTO todo_item (id, content) VALUES ($1, $2)",
		id, "Get by ID")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.LastInsertId()) // не поддерживается
	fmt.Println(result.RowsAffected()) // количество добавленных строк

	fmt.Println("\nInsert row and get entity id")
	var newId string
	db.QueryRow("INSERT INTO todo_item (content) VALUES ($1) returning id", "Привет мир").Scan(&newId)
	fmt.Println(newId) // А так можно

	// Закрываем транзакцию
	_ = tx.Commit()

	fmt.Println("\nSelect all")
	rows, err := db.Query("SELECT * FROM todo_item")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	items := []TodoItem{}

	for rows.Next() {
		p := TodoItem{}
		err := rows.Scan(&p.ID, &p.Content, &p.IsDone, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	for _, p := range items {
		fmt.Println(p.String())
	}

	fmt.Println("\nSelect field by id")
	var content string
	err = db.QueryRow("SELECT content FROM todo_item WHERE id = $1", newId).Scan(&content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(content)
}
