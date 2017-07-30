package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/markwallsgrove/saml_federation_proxy/models"

	mgo "gopkg.in/mgo.v2"
)

type mongoHandler func(http.ResponseWriter, *http.Request, *mgo.Session)
type handler func(http.ResponseWriter, *http.Request)

func getInt(r *http.Request, name string, min int, max int, def int) int {
	num, _ := strconv.Atoi(r.URL.Query().Get(name))

	if num < min || num > max {
		return def
	}

	return num
}

func apiEntityDescriptorHandler(w http.ResponseWriter, r *http.Request, mongoSession *mgo.Session) {
	c := mongoSession.DB("fedproxy").C("entityDescriptors")
	os := getInt(r, "offset", 0, 1000000, 0)
	lmt := getInt(r, "limit", 1, 100, 10)

	// TODO: skip & limit is slow with large amount of data
	query := c.Find(nil).Sort("-_id").Skip(os).Limit(lmt)

	var entityDescriptors []models.EntityDescriptor
	if err := query.All(&entityDescriptors); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	} else {
		json.NewEncoder(w).Encode(entityDescriptors)
	}
}

func index(w http.ResponseWriter, r *http.Request, mongoSession *mgo.Session) {
	http.ServeFile(w, r, "/public/index.html")
}

func htmlWrap(t mongoHandler, s *mgo.Session) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t(w, r, s)
	}
}

func jsonWrap(t mongoHandler, s *mgo.Session) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		t(w, r, s)
	}
}

func main() {
	session, err := mgo.Dial("mongodb")
	if err != nil {
		log.Fatal("cannot connect to mongo: ", err)
		return
	}

	defer session.Close()

	fs := http.FileServer(http.Dir("/public"))

	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", htmlWrap(index, session))
	http.HandleFunc("/1/api/entitydescriptor", jsonWrap(apiEntityDescriptorHandler, session))
	http.ListenAndServe(":8080", nil)
}
