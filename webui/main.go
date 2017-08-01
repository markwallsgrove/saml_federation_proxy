package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/markwallsgrove/saml_federation_proxy/models"
	log "github.com/sirupsen/logrus"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoHandler func(http.ResponseWriter, *http.Request, *mgo.Session) interface{}
type handler func(http.ResponseWriter, *http.Request)

func handleError(err error, w http.ResponseWriter) {
	if err == mgo.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	} else if err != nil {
		log.WithField("err", err).Error("api error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	}
}

func apiEnableEntityDescriptorHandler(w http.ResponseWriter, r *http.Request, mongoSession *mgo.Session) interface{} {
	c := mongoSession.DB("fedproxy").C("entityDescriptors")
	en := getString(r, "enabled", "false") == "true"
	id := getString(r, "id", "")

	qry := bson.M{"entityid": id}
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"enabled": en}},
		ReturnNew: true,
	}

	var ed models.EntityDescriptor
	_, err := c.Find(qry).Limit(1).Apply(change, &ed)

	handleError(err, w)
	return ed
}

func apiEntityDescriptorHandler(w http.ResponseWriter, r *http.Request, mongoSession *mgo.Session) interface{} {
	c := mongoSession.DB("fedproxy").C("entityDescriptors")
	id := getString(r, "id", "")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return nil
	}

	var entityDescriptor models.EntityDescriptor
	filter := bson.M{"entityid": id}
	err := c.Find(filter).One(&entityDescriptor)

	handleError(err, w)
	return entityDescriptor
}

func apiEntityDescriptorsHandler(w http.ResponseWriter, r *http.Request, mongoSession *mgo.Session) interface{} {
	c := mongoSession.DB("fedproxy").C("entityDescriptors")
	os := getInt(r, "offset", 0, 1000000, 0)
	lmt := getInt(r, "limit", 1, 100, 10)

	// TODO: skip & limit is slow with large amount of data
	// TODO: not returning .Enabled :(
	var entityDescriptors []models.EntityDescriptor
	err := c.Find(nil).Sort("-_id").Skip(os).Limit(lmt).All(&entityDescriptors)

	for _, entityDescriptor := range entityDescriptors {
		x, _ := json.Marshal(entityDescriptor)
		log.WithField("EntityDescriptor", string(x)).Debug("entity descriptor")
	}

	handleError(err, w)
	return entityDescriptors
}

func index(w http.ResponseWriter, r *http.Request, mongoSession *mgo.Session) interface{} {
	http.ServeFile(w, r, "/public/index.html")
	return nil
}

func main() {
	log.SetLevel(log.DebugLevel)

	session, err := mgo.Dial("mongodb")
	if err != nil {
		log.WithField("err", err).Error("cannot connect to mongo")
		return
	}

	defer session.Close()

	fs := http.FileServer(http.Dir("/public"))

	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", htmlWrap(index, session))
	http.HandleFunc("/1/api/entitydescriptors", jsonWrap(apiEntityDescriptorsHandler, session))
	http.HandleFunc("/1/api/entitydescriptor", jsonWrap(apiEntityDescriptorHandler, session))
	http.HandleFunc("/1/api/entityDescriptor/toggle", jsonWrap(apiEnableEntityDescriptorHandler, session))
	http.ListenAndServe(":8080", nil)
}

func getString(r *http.Request, name string, def string) string {
	if vl := r.URL.Query().Get(name); vl != "" {
		return vl
	}

	return def
}

func getInt(r *http.Request, name string, min int, max int, def int) int {
	num, _ := strconv.Atoi(getString(r, name, ""))

	if num < min || num > max {
		return def
	}

	return num
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
		rst := t(w, r, s)
		if rst != nil {
			json.NewEncoder(w).Encode(rst)
		}
	}
}
