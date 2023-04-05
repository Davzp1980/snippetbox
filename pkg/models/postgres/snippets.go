package postgres

import (
	"database/sql"
	"errors"
	"snippetbox/pkg/models"
	"strconv"
	"time"
)

// SnippetModel - Определяем тип который обертывает пул подключения sql.DB

type SnippetModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	var id int64
	//stmt := `INSERT INTO snippets (title, content, created, expires)
	// VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	ex, _ := strconv.Atoi(expires)
	err := m.DB.QueryRow(`INSERT INTO snippets (title, content,created,expires) 
		VALUES ($1,$2,$3,$4) returning id`, title, content, time.Now(), time.Now().AddDate(0, 0, ex)).Scan(&id)
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `select id,title, content,created,expires from snippets
	  	where expires > created and id=$1`

	s := &models.Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	// Пишем SQL запрос, который мы хотим выполнить.
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > now() ORDER BY created DESC LIMIT 10`
	// Используем метод Query() для выполнения нашего SQL запроса.
	// В ответ мы получим sql.Rows, который содержит результат нашего запроса.

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Откладываем вызов rows.Close(), чтобы быть уверенным, что набор результатов из sql.Rows
	// правильно закроется перед вызовом метода Latest(). Этот оператор откладывания
	// должен выполнится *после* проверки на наличие ошибки в методе Query().
	// В противном случае, если Query() вернет ошибку, это приведет к панике
	// так как он попытается закрыть набор результатов у которого значение: nil.
	defer rows.Close()

	var snippets []*models.Snippet
	// Используем rows.Next() для перебора результата. Этот метод предоставляем
	// первый а затем каждую следующею запись из базы данных для обработки
	// методом rows.Scan().
	for rows.Next() {
		s := &models.Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Добавляем структуру в срез.
		snippets = append(snippets, s)
	}
	// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
	// если в ходе работы у нас не возникла какая либо ошибка.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// Если все в порядке, возвращаем срез с данными.
	return snippets, nil
}
