package main

import (
	"database/sql"
	"flag"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/openidConnect"
)

type templateHandler struct {
	once     sync.Once
	filename []string
	tpl      *template.Template
}

var TemplateData = make(map[string]interface{})

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		var files []string
		for _, file := range t.filename {
			files = append(files, filepath.Join("/Users/lucas/go/src/programathor/templates", file))
		}
		t.tpl = template.Must(template.ParseFiles(files...))
	})
	t.tpl.ExecuteTemplate(w, "layout", "")
}

func newRouter() *mux.Router {
	mux := mux.NewRouter()
	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("/Users/lucas/go/src/programathor/public"))))
	return mux
}

var Db *sql.DB
var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	var err error
	Db, err = sql.Open("postgres", "user=postgres dbname=programathor password="+os.Getenv("POSTGRES_SECRET")+" sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func main() {
	goth.UseProviders(
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:8080/auth/facebook/callback"),
	)
	openidConnect, _ := openidConnect.New(os.Getenv("OPENID_CONNECT_KEY"), os.Getenv("OPENID_CONNECT_SECRET"), "http://localhost:3000/auth/callback/openid-connect", os.Getenv("OPENID_CONNECT_DISCOVERY_URL"))
	if openidConnect != nil {
		goth.UseProviders(openidConnect)
	}
	m := make(map[string]string)
	m["facebook"] = "Facebook"
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var addr = flag.String("addr", ":8080", "Addr of application")
	flag.Parse()
	mux := newRouter()
	server := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}
	mux.HandleFunc("/auth/{provider}/{action}", handleLoginOauth)
	mux.Handle("/", &templateHandler{filename: []string{"layout.html", "navbar-loggedout.html", "main.html"}})
	mux.Handle("/cadastrar", &templateHandler{filename: []string{"layout.html", "navbar-loggedout.html", "cadastro.html"}}).Methods("GET")
	mux.HandleFunc("/cadastrar", Cadastrar).Methods("POST")
	server.ListenAndServe()
}
