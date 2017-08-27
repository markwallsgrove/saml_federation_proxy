package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"

	xmlsec "github.com/crewjam/go-xmlsec"
	"github.com/markwallsgrove/saml_federation_proxy/models"
	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func parseRsaPrivateKeyFromPemStr(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM) // TODO: err
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// func getSignatureTemplate() []byte {
// 	return []byte(`
// 		<ds:Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
// 			<ds:SignedInfo>
// 				<ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#" />
// 				<ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256" />
// 				<ds:Reference URI="#_">
// 					<ds:Transforms>
// 						<ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature" />
// 						<ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#" />
// 					</ds:Transforms>
// 					<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256" />
// 					<ds:DigestValue></ds:DigestValue>
// 				</ds:Reference>
// 			</ds:SignedInfo>
// 			<ds:SignatureValue/>
// 			<ds:KeyInfo>
// 				<ds:KeyName />
// 			</ds:KeyInfo>
// 		</ds:Signature>
// 	`)
// }

func getSignature() models.Signature {
	return models.Signature{
		SignedInfo: models.SignedInfo{
			CanonicalizationMethod: models.CanonicalizationMethod{
				Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#",
			},
			SignatureMethod: models.SignatureMethod{
				Algorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
			},
			Reference: models.Reference{
				URI: "#_",
				DigestMethod: models.DigestMethod{
					Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
				},
				Transforms: models.Transforms{
					Transform: []models.Transform{
						models.Transform{Algorithm: "http://www.w3.org/2000/09/xmldsig#enveloped-signature"},
						models.Transform{Algorithm: "http://www.w3.org/2001/10/xml-exc-c14n#"},
					},
				},
			},
		},
	}
}

func signXml(key []byte, xml []byte) (error, []byte) {
	opts := xmlsec.SignatureOptions{
		XMLID: []xmlsec.XMLIDOption{{
			ElementName:      "EntitiesDescriptor",
			ElementNamespace: "urn:oasis:names:tc:SAML:2.0:metadata",
			AttributeName:    "ID",
		},
		}}

	if signedXml, err := xmlsec.Sign(key, xml, opts); err != nil {
		return err, nil
	} else {
		return nil, signedXml
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

		id := "jlsdfjklfdjkl544534" // TODO: what is a good id?
		name := "https://fedproxy.com"
		validUntil := time.Now().Add(time.Duration(24 * time.Hour))

		entitiesDescriptor := models.EntitiesDescriptor{
			ID:                &id,
			Name:              &name,
			Signature:         getSignature(),
			ValidUntil:        &validUntil,
			EntityDescriptors: entityDescriptors,
		}

		// TODO: not so great
		xml, err := xml.Marshal(entitiesDescriptor)
		if err != nil {
			log.WithError(err).Error("cannot encode entities descriptor")
			d.Reject(true)
			continue
		}

		log.WithField("xml", string(xml)).Info("signing xml")
		err, signedXml := signXml(key, xml)
		failOnError(err, "cannot sign XML")

		log.WithField("payload", string(signedXml)).Info("signed xml")

		exportResult := models.ExportResult{
			Name:    exportTask.Name,
			Payload: string(signedXml),
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
