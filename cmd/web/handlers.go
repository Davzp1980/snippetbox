package main

import (
	"errors"
	"net/http"
	"snippetbox/pkg/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Используем помощника render() для отображения шаблона.
	app.render(w, r, "home.page.html", &templateData{
		Snippets: s,
	})

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Страница не найдена.
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Используем помощника render() для отображения шаблона.
	app.render(w, r, "show.page.html", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}

		title := r.FormValue("title")     //"Не имей 100 рублей а имей 100 друзей"
		content := r.FormValue("content") //"Не имей 100 рублей а имей 100 друзей."
		expires := r.FormValue("expires") //"5"

		_, err = app.snippets.Insert(title, content, expires)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.ServeFile(w, r, "./ui/html/create.page.html")
		//http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", s), http.StatusSeeOther)
	}
}
