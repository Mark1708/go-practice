package main

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

// schema Схема таблиц
var schema = `
CREATE TABLE IF NOT EXISTS todo_item (
    id UUID DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
)
`

// TodoItem Сущность TO-DO листа
type TodoItem struct {
	ID        uuid.UUID `db:"id"`
	Content   string    `db:"content"`
	IsDone    bool      `db:"is_done"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (t *TodoItem) String() string {
	return fmt.Sprintf(
		"TodoItem { id: %s, content: %s, isDone: %v, createdAt: %s, updatedAt: %s }",
		t.ID, t.Content, t.IsDone, t.CreatedAt.Format(time.RFC3339), t.UpdatedAt.Format(time.RFC3339),
	)
}
