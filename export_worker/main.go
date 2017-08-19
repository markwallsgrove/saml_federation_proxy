package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/beevik/etree"
	"github.com/markwallsgrove/saml_federation_proxy/models"
	dsig "github.com/russellhaering/goxmldsig"
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

func task(msgs <-chan amqp.Delivery, session *mgo.Session, ctx *dsig.SigningContext, forever chan bool) {
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
			ValidUntil:        &validUntil,
			EntityDescriptors: entityDescriptors,
		}

		// TODO: not so great
		xmlEncoded, err := xml.Marshal(entitiesDescriptor)
		if err != nil {
			log.WithError(err).Error("cannot encode entities descriptor")
			d.Reject(true)
			continue
		}
		readXMLDoc := etree.NewDocument()
		err = readXMLDoc.ReadFromBytes(xmlEncoded)
		failOnError(err, "cannot parse xml")

		elementToSign := readXMLDoc.Root()
		elementToSign.CreateAttr("ID", "id1234")

		signedElement, err := ctx.SignEnveloped(elementToSign)
		failOnError(err, "failed to sign envelop")

		var signedAssertionBuf []byte
		newDoc := etree.NewDocument()
		newDoc.SetRoot(signedElement)
		signedAssertionBuf, err = newDoc.WriteToBytes()
		failOnError(err, "failed to convert doc to bytes")

		log.WithField("payload", string(signedAssertionBuf)).Info("signed xml")

		exportResult := models.ExportResult{
			Name:    exportTask.Name,
			Payload: string(signedAssertionBuf),
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

	keyBytes, err := ioutil.ReadFile("/saml.pem")
	failOnError(err, "Failed to load signing key")

	certBytes, err := ioutil.ReadFile("/saml.crt")
	failOnError(err, "Failed to read certificate")

	keyPair, err := tls.X509KeyPair(certBytes, keyBytes)
	failOnError(err, "invalided to load keypair")

	keyStore := dsig.TLSCertKeyStore(keyPair)

	signingContext := dsig.NewDefaultSigningContext(keyStore)
	signingContext.Canonicalizer = dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("")
	err = signingContext.SetSignatureMethod(dsig.RSASHA256SignatureMethod)

	forever := make(chan bool)

	go task(msgs, session, signingContext, forever)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
