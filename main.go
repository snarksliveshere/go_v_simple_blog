package main

import (
	"blog/db/documents"
	"blog/models"
	"fmt"
	"html/template"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"github.com/tobyzxj/mgo"
)

var postsCollection *mgo.Collection
var counter int

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	//mongoConfig := "mongodb://root:example@localhost:27017/test?authMechanism=PLAIN"
	mongoConfig := "mongodb://localhost:27017/test"
	session, err := mgo.Dial(mongoConfig)
	if err != nil {
		panic(err)
	}
	postsCollection = session.DB("blog").C("posts")

	m := martini.Classic()

	unescFunc := template.FuncMap{"unescape": unescape}

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		Funcs:      []template.FuncMap{unescFunc},
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
	m.Get("/404", notFoundHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/delete/:id", deleteHandler)
	m.Post("/SavePost", savePostHandler)
	m.Post("/gethtml", getHtmlHandler)

	m.Run()
}

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := convertMarkDown(md)
	rnd.JSON(200, map[string]interface{}{"html": html})
}

func notFoundHandler(rnd render.Render) {
	rnd.HTML(200, "404", nil)
}

func indexHandler(rnd render.Render) {
	fmt.Println(counter)
	var postsDocuments []documents.PostDocument
	postsCollection.Find(nil).All(&postsDocuments)
	var posts []models.Post
	for _, doc := range postsDocuments {
		post := models.Post{
			Id:              doc.Id,
			Title:           doc.Title,
			ContentHTML:     doc.ContentHTML,
			ContentMarkdown: doc.ContentMarkdown,
		}
		posts = append(posts, post)
	}
	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	post := models.Post{}
	rnd.HTML(200, "write", post)
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

func deleteHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	if id == "" {
		rnd.Redirect("/404")
		//http.NotFound(w, r)
	}
	postsCollection.Remove(id)
	rnd.Redirect("/")
	//http.Redirect(w, r, "/", 302)
}

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]

	var postDocument documents.PostDocument
	err := postsCollection.FindId(id).One(&postDocument)

	if err != nil {
		// redirect to 404
		rnd.Redirect("/404")
		return
		//http.NotFound(w, r)
	}
	post := models.Post{
		Id:              postDocument.Id,
		Title:           postDocument.Title,
		ContentHTML:     postDocument.ContentHTML,
		ContentMarkdown: postDocument.ContentMarkdown,
	}
	rnd.HTML(200, "write", post)
	//rnd.HTML(200, "write", postDocument)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")

	contentHTML := convertMarkDown(contentMarkdown)
	postDocument := documents.PostDocument{id, title, contentHTML, contentMarkdown}

	if id != "" {
		postsCollection.UpdateId(id, postDocument)
	} else {
		id = GenerateId()
		postDocument.Id = id
		postsCollection.Insert(postDocument)
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
