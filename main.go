package main

import (
	"blog/models"
	"fmt"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

var posts map[string]*models.Post
var counter int

func main() {
	posts = make(map[string]*models.Post, 0)

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		//Funcs:           nil,
		//Delims:          render.Delims{},
		Charset:    "UTF-8",
		IndentJSON: true,
		//IndentXML:       false,
		//PrefixJSON:      nil,
		//PrefixXML:       nil,
		//HTMLContentType: "",
	}))

	counter = 0
	m.Use(func(r *http.Request) {
		if r.URL.Path == "/write" {
			counter++
		}
	})

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))
	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit", editHandler)
	m.Post("/SavePost", savePostHandler)
	m.Get("/delete", deleteHandler)
	m.Get("/404", notFoundHandler)

	m.Run()
}

func notFoundHandler(rnd render.Render) {
	rnd.HTML(200, "404", posts)
}
func indexHandler(rnd render.Render) {
	fmt.Println(counter)
	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	rnd.HTML(200, "write", nil)
}

//func writeHandler(w http.ResponseWriter, r *http.Request) {
//	t, err := template.ParseFiles("templates/write.html", "templates/layout.html", "templates/footer.html")
//	if err != nil {
//		fmt.Fprintf(w, err.Error())
//		return
//	}
//
//	t.ExecuteTemplate(w, "write", nil)
//}

func deleteHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		rnd.Redirect("/404")
		//http.NotFound(w, r)
	}
	delete(posts, id)
	rnd.Redirect("/")
	//http.Redirect(w, r, "/", 302)
}

func editHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	post, found := posts[id]
	if !found {
		// redirect to 404
		rnd.Redirect("/404")
		return
		//http.NotFound(w, r)
	}
	rnd.HTML(200, "write", post)
	//t.ExecuteTemplate(w, "write", post)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")
	var post *models.Post

	if id != "" {
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id = GenerateId()
		post := models.NewPost(id, title, content)
		posts[post.Id] = post
	}

	rnd.Redirect("/")
}

//func savePostHandler(w http.ResponseWriter, r *http.Request) {
//	id := r.FormValue("id")
//	title := r.FormValue("title")
//	content := r.FormValue("content")
//	var post *models.Post
//
//	if id != "" {
//		post = posts[id]
//		post.Title = title
//		post.Content = content
//	} else {
//		id = GenerateId()
//		post := models.NewPost(id, title, content)
//		posts[post.Id] = post
//	}
//
//	http.Redirect(w, r, "/", 302)
//}
