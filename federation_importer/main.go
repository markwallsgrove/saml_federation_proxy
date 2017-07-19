package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/crewjam/go-xmlsec"
	"github.com/markwallsgrove/saml_federation_proxy/models"
)

// const ukfedURL = "http://metadata.ukfederation.org.uk/ukfederation-metadata.xml"
const ukfedURL = "http://localhost:8181/ukfederation-metadata.xml"

func downloadFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func verifySignature(buff io.ReadCloser, keyPath string) ([]byte, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalln("error loading key from path", err)
		return nil, err
	}

	doc, err := ioutil.ReadAll(buff)
	if err != nil {
		log.Fatalln("error reading xml buffer", err)
		return nil, err
	}

	defer buff.Close()

	err = xmlsec.Verify(key, doc, xmlsec.SignatureOptions{
		XMLID: []xmlsec.XMLIDOption{{
			ElementName:      "EntitiesDescriptor",
			ElementNamespace: "urn:oasis:names:tc:SAML:2.0:metadata",
			AttributeName:    "ID",
		},
		},
	})

	return doc, err
}

func unmarshallXML(buff []byte) (*models.EntitiesDescriptor, error) {
	entitiesDescriptor := new(models.EntitiesDescriptor)
	err := xml.Unmarshal(buff, entitiesDescriptor)
	return entitiesDescriptor, err
}

func marshallEntityDescriptor(entityDescriptor *models.EntityDescriptor) ([]byte, error) {
	bytes, err := xml.Marshal(entityDescriptor)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func calcChecksum(bytes []byte) string {
	hasher := md5.New()
	hasher.Write(bytes)
	return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
	// TODO: cli args (queue location, url)
	// TODO: queue system
	buff, err := downloadFile(ukfedURL)

	if err != nil {
		log.Fatalln("error downloading file:", err)
		return
	}

	doc, err := verifySignature(buff, "/home/smoky/workspace/golang/src/github.com/markwallsgrove/saml_federation_proxy/ukfederation.pem")

	if err != nil {
		log.Fatalln("signature error:", err)
		return
	}

	entitiesDescriptors, err := unmarshallXML(doc)
	if err != nil {
		log.Fatalln("unmarhsalling error:", err)
		return
	}

	for _, entityDescriptor := range entitiesDescriptors.EntityDescriptor {
		entityID := entityDescriptor.EntityID
		entityDescriptorXML, _ := marshallEntityDescriptor(entityDescriptor)
		checksum := calcChecksum(entityDescriptorXML)
		fmt.Println(entityID, checksum)
	}

	log.Println("fin.")
}
