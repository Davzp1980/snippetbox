package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	// Используем функцию mux.Handle() для регистрации обработчика для
	// всех запросов, которые начинаются с "/static/". Мы убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
