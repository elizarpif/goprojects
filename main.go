// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
)

//не используя глобальный счетчик, посчитать количество пост-запросов

type Message struct {
	Msg string `json:"text"`
}

type Counter struct {
	n int
}

//POST
func (ctr *Counter) NewArticle(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} else {

		text := string(body)
		ctr.n++
		text = text + fmt.Sprint(ctr.n)

		resp_struct := &Message{Msg: text}

		fmt.Println(ctr.n)

		resp, err := json.Marshal(resp_struct)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {

			fmt.Println("x-req-id", r.Header.Get("x-req-id"))
			w.Write(resp)
		}

	}
	defer r.Body.Close()

}

func main() {

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	ctr := Counter{}

	r.Route("/articles", func(r chi.Router) {
		r.Post("/add", ctr.NewArticle)

	})

	http.ListenAndServe(":3000", r)

}
