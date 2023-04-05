package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"snippetbox/pkg/models/postgres"
	"text/template"

	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "mydb"
)

// Добавляем поле templateCache в структуру зависимостей. Это позволит
// получить доступ к кэшу во всех обработчиках.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *postgres.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	// Определение нового флага из командной строки для настройки MySQL подключения.
	dns := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Чтобы функция main() была более компактной, мы поместили код для создания
	// пула соединений в отдельную функцию openDB(). Мы передаем в нее полученный
	// источник данных (DSN) из флага командной строки.

	db, err := openDB(dns)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Инициализируем новый кэш шаблона..
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// И добавляем его в зависимостях нашего
	// веб-приложения.

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &postgres.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск веб-сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dns)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
