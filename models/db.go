package models

import (
	"encoding/json"
	"io"
	"io/ioutil"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func UpdateExportResult(result ExportResult, session *mgo.Session) error {
	c := session.DB("fedproxy").C("exportResult")

	_, err := c.Upsert(
		bson.M{"name": result.Name},
		bson.M{"$set": bson.M{"payload": result.Payload}},
	)

	return err
}

func FindExportResult(name string, session *mgo.Session) (ExportResult, error) {
	c := session.DB("fedproxy").C("exportResult")

	var exportResult ExportResult
	err := c.Find(bson.M{"name": name}).One(&exportResult)
	return exportResult, err
}

func FindEntityDescriporsByName(names []string, session *mgo.Session) ([]*EntityDescriptor, error) {
	c := session.DB("fedproxy").C("entityDescriptors")

	var entityDescriptors []*EntityDescriptor

	err := c.Find(bson.M{"entityid": bson.M{"$in": names}}).All(&entityDescriptors)
	return entityDescriptors, err
}

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

func PatchExport(exportName string, patch ExportPatch, session *mgo.Session) (Export, error) {
	c := session.DB("fedproxy").C("exports")

	var export Export
	update := bson.M{}

	if patch.EntityDescriptors.Append != nil {
		update["$addToSet"] = bson.M{"entitydescriptors": bson.M{"$each": patch.EntityDescriptors.Append}}
	}

	if patch.EntityDescriptors.Delete != nil {
		update["$pullAll"] = bson.M{"entitydescriptors": patch.EntityDescriptors.Delete}
	}

	change := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}

	_, err := c.Find(bson.M{"name": exportName}).Apply(change, &export)

	if err != nil {
		return Export{}, err
	}

	return export, nil
}

func Unmarshall(r io.Reader, obj interface{}) error {
	b, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &obj)
	return err
}

type ExportPatch struct {
	EntityDescriptors PatchChange `json:"EntityDescriptors,omitempty"`
}

type PatchChange struct {
	Delete []string `json:"Delete,omitempty"`
	Append []string `json:"Append,omitempty"`
}

type Export struct {
	Name              string   `json:"Name"`
	EntityDescriptors []string `json:"EntityDescriptors"`
}

type ExportResult struct {
	Name    string `json:"Name"`
	Payload string `json:"Payload"`
}
