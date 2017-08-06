package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/markwallsgrove/saml_federation_proxy/models"
	log "github.com/sirupsen/logrus"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoHandler func(http.ResponseWriter, *http.Request, *mgo.Session) interface{}
type handler func(http.ResponseWriter, *http.Request)

func notFoundError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func handleError(err error, w http.ResponseWriter) {
	if err == mgo.ErrNotFound {
		notFoundError(w)
	} else if err != nil {
		log.WithField("err", err).Error("api error")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// func apiEnableEntityDescriptorHandler(w http.ResponseWriter, r *http.Request, mongoSession *mgo.Session) interface{} {
// 	c := mongoSession.DB("fedproxy").C("entityDescriptors")
// 	en := getString(r, "enabled", "false") == "true"
// 	id := getString(r, "id", "")

// 	qry := bson.M{"entityid": id}
// 	change := mgo.Change{
// 		Update:    bson.M{"$set": bson.M{"enabled": en}},
// 		ReturnNew: true,
// 	}

// 	var ed models.EntityDescriptor
// 	_, err := c.Find(qry).Limit(1).Apply(change, &ed)

// 	handleError(err, w)
// 	return ed
// }

func apiEntityDescriptorHandler(w http.ResponseWriter, r *http.Request, mongoSession *mgo.Session) interface{} {
	c := mongoSession.DB("fedproxy").C("entityDescriptors")
	vars := mux.Vars(r)
	id := vars["id"]

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

func apiEntityDescriptorsHandler(w http.ResponseWriter, r *http.Request, session *mgo.Session) interface{} {
	os := getInt(r, "offset", 0, 1000000, 0)
	lmt := getInt(r, "limit", 1, 100, 10)

	entityDescriptors, err := models.PaginateEntityDescriptors(os, lmt, session)
	if err != nil {
		handleError(err, w)
		return nil
	}

	return entityDescriptors
}

func apiExportsHandler(w http.ResponseWriter, r *http.Request, session *mgo.Session) interface{} {
	exports, err := models.GetExports(session)
	handleError(err, w)
	return exports
}

func apiExportCreationHandler(w http.ResponseWriter, r *http.Request, session *mgo.Session) interface{} {
	var e models.Export
	err := models.Unmarshall(r.Body, &e)

	if err != nil {
		handleError(err, w)
		return nil
	}

	err = models.InsertExport(e, session)
	if err != nil {
		handleError(err, w)
		return nil
	}

	return nil
}

func apiPatchExportHandler(w http.ResponseWriter, r *http.Request, session *mgo.Session) interface{} {
	vars := mux.Vars(r)
	exportName := vars["id"]

	var patch models.ExportPatch
	err := models.Unmarshall(r.Body, &patch)
	if err != nil {
		handleError(err, w)
		return nil
	}

	export, err := models.PatchExport(exportName, patch, session)
	if err != nil {
		handleError(err, w)
		return nil
	}

	return export
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

	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	r.HandleFunc("/", htmlWrap(index, session)).Methods("GET")
	r.HandleFunc("/1/api/entitydescriptors", jsonWrap(apiEntityDescriptorsHandler, session)).Methods("GET")
	r.HandleFunc("/1/api/entitydescriptor/{id}", jsonWrap(apiEntityDescriptorHandler, session)).Methods("GET")
	r.HandleFunc("/1/api/exports", jsonWrap(apiExportsHandler, session)).Methods("GET")
	r.HandleFunc("/1/api/exports", jsonWrap(apiExportCreationHandler, session)).Methods("POST")
	r.HandleFunc("/1/api/exports/{id}", jsonWrap(apiPatchExportHandler, session)).Methods("PATCH")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
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
