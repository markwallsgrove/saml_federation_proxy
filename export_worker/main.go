package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"io/ioutil"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/crewjam/saml"
	"github.com/ma314smith/signedxml"
	"github.com/markwallsgrove/saml_federation_proxy/models"
	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func task(msgs <-chan amqp.Delivery, session *mgo.Session, key []byte, forever chan bool) {
	for d := range msgs {
		var exportTask models.ExportTask
		err := models.Unmarshall(bytes.NewReader(d.Body), &exportTask)

		if err != nil {
			log.WithField("error", err).Error("cannot process export task")
			d.Reject(true)
			continue
		}

		entityDescriptors, err := models.FindEntityDescriporsByName(
			exportTask.EntityDescriptors,
			session,
		)

		id := "jlsdfjklfdjkl544534"
		name := "https://fedproxy.com"
		validUntil := time.Now().Add(time.Duration(24 * time.Hour))

		entitiesDescriptor := saml.EntitiesDescriptor{
			// XmlnsDs:      "http://www.w3.org/2000/09/xmldsig#",
			// XmlnsIdpdisc: "urn:oasis:names:tc:SAML:profiles:SSO:idp-discovery-protocol",
			// XmlnsInit:    "urn:oasis:names:tc:SAML:profiles:SSO:request-init",
			// XmlnsMdattr:  "urn:oasis:names:tc:SAML:metadata:attribute",
			// XmlnsMdrpi:   "urn:oasis:names:tc:SAML:metadata:rpi",
			// XmlnsMdui:    "urn:oasis:names:tc:SAML:metadata:ui",
			// XmlnsRemd:    "http://refeds.org/metadata",
			// XmlnsSaml:    "urn:oasis:names:tc:SAML:2.0:assertion",
			// XmlnsShibmd:  "urn:mace:shibboleth:metadata:1.0",
			// XmlnsXsi:     "http://www.w3.org/2001/XMLSchema-instance",
			// Xmlns:            "urn:oasis:names:tc:SAML:2.0:metadata",
			ID:                &id,
			Name:              &name,
			ValidUntil:        &validUntil,
			EntityDescriptors: entityDescriptors,
		}

		// xmlEncoded, err := xml.MarshalIndent(entitiesDescriptor, "  ", "    ")
		xmlEncoded, err := xml.Marshal(entitiesDescriptor)
		if err != nil {
			log.WithError(err).Error("cannot encode entities descriptor")
			d.Reject(true)
			continue
		}

		// TODO: cannot find start node
		// _, err = xmlsec.Sign(key, xmlEncoded, xmlsec.SignatureOptions{
		// 	XMLID: []xmlsec.XMLIDOption{{
		// 		ElementName:      "EntitiesDescriptor",
		// 		ElementNamespace: "urn:oasis:names:tc:SAML:2.0:metadata",
		// 		AttributeName:    "ID",
		// 	}},
		// })
		signer, err := signedxml.NewSigner(string(xmlEncoded))

		if err != nil {
			log.WithError(err).Error("cannot sign entities descriptor")
			// d.Reject(false)
			// continue
		} else {
			p, _ := pem.Decode(key)
			if p == nil {
				log.Error("no pem block found")
			} else {
				pk, err := x509.ParsePKCS1PrivateKey(p.Bytes)
				if err != nil {
					log.WithError(err).Error("cannot parse private key with pem bytes")
				} else {
					ep, err := signer.Sign(pk)
					if err != nil {
						log.WithError(err).Error("cannot sign xml")
					} else {
						log.WithField("payload", ep).Info("signed xml payload")
					}
				}
			}
		}

		log.WithFields(log.Fields{
			"xml": string(xmlEncoded),
			"key": string(key),
		}).Info("signed xml")

		exportResult := models.ExportResult{
			Name:    exportTask.Name,
			Payload: string(xmlEncoded),
		}

		if err := models.UpdateExportResult(exportResult, session); err != nil {
			log.WithError(err).Error("cannot upsert export result")
			d.Reject(true)
			continue
		}

		d.Ack(false)
	}
}

func main() {
	session, err := mgo.Dial("mongodb")
	if err != nil {
		log.WithField("err", err).Fatal("cannot connect to mongo")
		return
	}

	defer session.Close()

	queueConn := os.Getenv("QUEUE_CONN")
	conn, err := amqp.Dial(queueConn)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"descriptor_exporter_queue", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	key, err := ioutil.ReadFile("/saml.pem")
	failOnError(err, "Failed to load signing key")

	forever := make(chan bool)

	go task(msgs, session, key, forever)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
