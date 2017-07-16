package main

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/crewjam/go-xmlsec"
)

const ukfedURL = "http://metadata.ukfederation.org.uk/ukfederation-metadata.xml"

type EntitiesDescriptors struct {
	EntityDescriptor []EntityDescriptor `xml:"EntityDescriptor"`
	Name             string             `xml:"Name,attr"`
}

type EntityDescriptor struct {
	ID               string           `xml:"ID,attr"`
	EntityID         string           `xml:"EntityID,attr"`
	IDPSSODescriptor IDPSSODescriptor `xml:"IDPSSODescriptor"`
	SPSSODescriptor  SPSSODescriptor  `xml:"SPSSODescriptor"`
}

type IDPSSODescriptor struct {
	ProtocolSupportEnumeration string        `xml:"protocolSupportEnumeration,attr"`
	KeyDescriptor              KeyDescriptor `xml:"KeyDescriptor"`
}

type SPSSODescriptor struct {
	ProtocolSupportEnumeration string `xml:"protocolSupportEnumeration,attr"`
}

type KeyDescriptor struct {
	KeyInfo KeyInfo `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
}

type KeyInfo struct {
	X509Data X509Data `xml:"http://www.w3.org/2000/09/xmldsig# X509Data"`
}

type X509Data struct {
	X509Certificate X509Certificate `xml:"http://www.w3.org/2000/09/xmldsig# X509Certificate"`
}

type X509Certificate struct {
}

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

func unmarshallXML(buff []byte) (*EntitiesDescriptors, error) {
	entitiesDescriptor := new(EntitiesDescriptors)
	err := xml.Unmarshal(buff, entitiesDescriptor)
	return entitiesDescriptor, err
}

func main() {
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

	for _, entitiesDescriptor := range entitiesDescriptors.EntityDescriptor {
		log.Println(entitiesDescriptor.IDPSSODescriptor.ProtocolSupportEnumeration)
		// TODO: broken
		log.Println(entitiesDescriptor.IDPSSODescriptor.KeyDescriptor.KeyInfo.X509Data.X509Certificate)
	}

	log.Println("fin.")
}
