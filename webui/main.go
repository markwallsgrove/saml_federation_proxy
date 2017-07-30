package main

import (
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

type mongoHandler func(http.ResponseWriter, *http.Request, *mgo.Session)
type handler func(http.ResponseWriter, *http.Request)

func index(w http.ResponseWriter, r *http.Request, mongoSession *mgo.Session) {
	// c := mongoSession.DB("fedproxy").C("entityDescriptors")

	// var entityDescriptors []models.EntityDescriptor
	// query := c.Find(nil).Sort("-_id").Limit(10)

	// if err := query.All(&entityDescriptors); err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte("500 - Something bad happened!"))
	// 	return
	// }

	// fmt.Fprintf(w, "Hi there, I love!")
	// for _, entityDescriptor := range entityDescriptors {
	// 	fmt.Fprintf(w, "<br /> entity descriptor: %s", entityDescriptor.EntityID)
	// }

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
	http.ListenAndServe(":8080", nil)
}
