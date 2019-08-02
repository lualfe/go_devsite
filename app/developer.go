package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

type Desenvolvedor struct {
	Id          int
	UUID        string
	OauthuserID string
	Nome        string
	Skills      []int
	Stack       int
	Exp         int
	Email       string
	Senha       []byte
	Descricao   string
	Provider    string
	Avatar      Avatar
	Username    string
	Data        map[interface{}]interface{}
}

var Dev *Desenvolvedor

func Cadastrar(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "dev_session")
	email := r.FormValue("email")
	var err error
	senha, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("senha")), 14)
	nome := r.FormValue("nome")
	id := uuid.New()
	uuid := id.String()
	profile := strings.ToLower(Dev.Nome)
	profile = strings.Replace(profile, " ", "-", -1)
	profileURL := profile + "-" + uuid[:3]
	session.Values["username"] = profileURL
	Dev = &Desenvolvedor{
		Nome:     nome,
		Email:    email,
		Senha:    senha,
		UUID:     uuid,
		Username: profileURL,
	}
	statement := "INSERT INTO dev (nome, email, senha, dev_uuid, username) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		http.Error(w, "Cadastro não efetivado com sucesso", http.StatusInternalServerError)
	}
	defer stmt.Close()
	err = stmt.QueryRow(Dev.Nome, Dev.Email, Dev.Senha, Dev.UUID, Dev.Username).Scan(&Dev.Id)
	if err != nil {
		http.Error(w, "Cadastro falhou", http.StatusInternalServerError)
	}
	session.Values["authenticated"] = true
	session.Values["uuid"] = uuid
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("http://localhost:8080/user/%s/edit", profileURL), 302)
}

func handleLoginOauth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	switch action {
	case "login":
		if _, err := gothic.CompleteUserAuth(w, r); err == nil {
			http.Redirect(w, r, "http://localhost:8080/", 302)
			return
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	case "callback":
		session, _ := store.Get(r, "dev_session")
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
		}
		rows, err := Db.Query("SELECT COUNT(*) FROM dev WHERE oauthuserid = $1 AND provider = $2", user.UserID, user.Provider)
		if err != nil {
			log.Fatal(err)
		}
		var count int
		for rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				log.Fatal(err)
			}
		}
		if count != 0 {
			session.Values["authenticated"] = true
			http.Redirect(w, r, "http://localhost:8080/", 302)
			return
		}
		id := uuid.New()
		uuid := id.String()
		Dev = &Desenvolvedor{
			Nome:        user.Name,
			Email:       user.Email,
			Provider:    user.Provider,
			UUID:        uuid,
			OauthuserID: user.UserID,
		}
		statement := "INSERT INTO dev (nome, email, provider, dev_uuid, oauthuserid) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		stmt, err := Db.Prepare(statement)
		if err != nil {
			return
		}
		defer stmt.Close()
		err = stmt.QueryRow(Dev.Nome, Dev.Email, Dev.Provider, Dev.UUID, Dev.OauthuserID).Scan(&Dev.Id)
		session.Values["authenticated"] = true
		session.Values["UUID"] = uuid
		profile := strings.ToLower(Dev.Nome)
		profile = strings.Replace(profile, " ", "-", -1)
		profileURL := profile + "-" + uuid[:3]
		session.Values["username"] = profileURL
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("http://localhost:8080/user/%s/edit", profileURL), 302)
	}
}
