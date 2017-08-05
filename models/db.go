package models

import (
	"encoding/json"
	"io"
	"io/ioutil"

	mgo "gopkg.in/mgo.v2"
)

func PaginateEntityDescriptors(os int, lmt int, session *mgo.Session) ([]EntityDescriptor, error) {
	// TODO: skip & limit is slow with large amount of data

	c := session.DB("fedproxy").C("entityDescriptors")

	var entityDescriptors []EntityDescriptor
	err := c.Find(nil).Sort("-_id").Skip(os).Limit(lmt).All(&entityDescriptors)
	return entityDescriptors, err
}

func GetExports(session *mgo.Session) ([]Export, error) {
	c := session.DB("fedproxy").C("exports")
	var exports []Export

	err := c.Find(nil).Sort("-_id").All(&exports)

	if exports == nil {
		exports = []Export{}
	}

	return exports, err
}

func InsertExport(e Export, session *mgo.Session) error {
	c := session.DB("fedproxy").C("exports")
	return c.Insert(e)
}

func Unmarshall(r io.Reader, obj interface{}) error {
	b, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &obj)
	return err
}

type Export struct {
	Name              string `json:"Name"`
	EntityDescriptors []EntityDescriptor
}
