package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	saml "github.com/crewjam/saml"
	"github.com/gorilla/mux"
	"github.com/markwallsgrove/saml_federation_proxy/models"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type xmlContextHandler func(http.ResponseWriter, *http.Request, *Context) []byte
type contextHandler func(http.ResponseWriter, *http.Request, *Context) interface{}
type handler func(http.ResponseWriter, *http.Request)

type Context struct {
	Session *mgo.Session
	Channel *amqp.Channel
}

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

func apiEntityDescriptorHandler(w http.ResponseWriter, r *http.Request, context *Context) interface{} {
	c := context.Session.DB("fedproxy").C("entityDescriptors")
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return nil
	}

	var entityDescriptor saml.EntityDescriptor
	filter := bson.M{"entityid": id}
	err := c.Find(filter).One(&entityDescriptor)

	handleError(err, w)
	return entityDescriptor
}

func apiEntityDescriptorsHandler(w http.ResponseWriter, r *http.Request, context *Context) interface{} {
	os := getInt(r, "offset", 0, 1000000, 0)
	lmt := getInt(r, "limit", 1, 100, 10)

	entityDescriptors, err := models.PaginateEntityDescriptors(os, lmt, context.Session)
	if err != nil {
		handleError(err, w)
		return nil
	}

	return entityDescriptors
}

func apiExportsHandler(w http.ResponseWriter, r *http.Request, context *Context) interface{} {
	exports, err := models.GetExports(context.Session)
	handleError(err, w)
	return exports
}

func apiExportCreationHandler(w http.ResponseWriter, r *http.Request, context *Context) interface{} {
	var e models.Export
	err := models.Unmarshall(r.Body, &e)

	if err != nil {
		handleError(err, w)
		return nil
	}

	err = models.InsertExport(e, context.Session)
	if err != nil {
		handleError(err, w)
		return nil
	}

	return nil
}

func apiPatchExportHandler(w http.ResponseWriter, r *http.Request, context *Context) interface{} {
	vars := mux.Vars(r)
	exportName := vars["id"]

	var patch models.ExportPatch
	if err := models.Unmarshall(r.Body, &patch); err != nil {
		handleError(err, w)
		return nil
	}

	export, err := models.PatchExport(exportName, patch, context.Session)
	if err != nil {
		handleError(err, w)
		return nil
	}

	exportTask := models.ExportTask{
		Name:              exportName,
		EntityDescriptors: export.EntityDescriptors,
	}

	if err = models.QueueExportTask(exportTask, context.Channel); err != nil {
		handleError(err, w)
		return nil
	}

	return export
}

func apiExportPayloadHandler(w http.ResponseWriter, r *http.Request, context *Context) []byte {
	vars := mux.Vars(r)
	exportName := vars["id"]

	exportResult, err := models.FindExportResult(exportName, context.Session)

	if err != nil {
		handleError(err, w)
		return nil
	}

	return []byte(exportResult.Payload)
}

func index(w http.ResponseWriter, r *http.Request, context *Context) interface{} {
	http.ServeFile(w, r, "/public/index.html")
	return nil
}

func main() {
	log.SetLevel(log.DebugLevel)

	mongoConn := os.Getenv("MONGO_CONN")
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		log.WithField("err", err).Error("cannot connect to mongo")
		return
	}

	defer session.Close()

	queueConn := os.Getenv("QUEUE_CONN")
	conn, err := amqp.Dial(queueConn)
	if err != nil {
		log.WithError(err).Error("cannot connect to amqp")
		return
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.WithError(err).Error("cannot connect to amqp channel")
		return
	}

	defer ch.Close()

	context := &Context{
		Session: session,
		Channel: ch,
	}

	fs := http.FileServer(http.Dir("/public"))

	r := mux.NewRouter()
	r.HandleFunc("/", htmlWrap(index, context)).Methods("GET")
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	r.HandleFunc("/1/api/entitydescriptors", jsonWrap(apiEntityDescriptorsHandler, context)).Methods("GET")
	r.HandleFunc("/1/api/entitydescriptor/{id}", jsonWrap(apiEntityDescriptorHandler, context)).Methods("GET")
	r.HandleFunc("/1/api/exports", jsonWrap(apiExportsHandler, context)).Methods("GET")
	r.HandleFunc("/1/api/exports", jsonWrap(apiExportCreationHandler, context)).Methods("POST")
	r.HandleFunc("/1/api/exports/{id}", jsonWrap(apiPatchExportHandler, context)).Methods("PATCH")
	r.HandleFunc("/1/api/exports/{id}/payload", xmlWrap(apiExportPayloadHandler, context)).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:80",
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

func htmlWrap(t contextHandler, c *Context) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t(w, r, c)
	}
}

func jsonWrap(t contextHandler, c *Context) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		rst := t(w, r, c)
		if rst != nil {
			json.NewEncoder(w).Encode(rst)
		}
	}
}

func xmlWrap(t xmlContextHandler, c *Context) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		rst := t(w, r, c)
		w.Write(rst)
	}
}
