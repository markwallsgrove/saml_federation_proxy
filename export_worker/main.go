package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/beevik/etree"
	"github.com/crewjam/saml"
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

		randomKeyStore := dsig.RandomKeyStoreForTest()
		ctx := dsig.NewDefaultSigningContext(randomKeyStore)

		_, certBytes, _ := randomKeyStore.GetKeyPair()
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			log.WithError(err).Error("cannot parse new cert")
		} else {
			buf := bytes.NewBufferString("")
			pem.Encode(buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
			log.WithField("cert", buf.String()).Info("cert")
		}

		entityDescriptors, err := models.FindEntityDescriporsByName(
			exportTask.EntityDescriptors,
			session,
		)

		id := "jlsdfjklfdjkl544534" // TODO: what is a good id?
		name := "https://fedproxy.com"
		validUntil := time.Now().Add(time.Duration(24 * time.Hour))

		entitiesDescriptor := saml.EntitiesDescriptor{
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

		doc := etree.NewDocument()
		err = doc.ReadFromBytes(xmlEncoded) //TODO:

		doc.CreateAttr("ID", id)

		el := doc.Root()

		sig, err := ctx.ConstructSignature(el, true)
		if err != nil {
			log.WithError(err).Error("Cannot construct signature")
			d.Reject(false)
			continue
		}

		ret := el.Copy()
		ret.Child = append(ret.Child, sig)
		// ret.Child = append(ret.Child, nil)
		// copy(ret.Child[1:], ret.Child[0:])
		// ret.Child[0] = sig

		doc.SetRoot(ret)

		// // Sign the element
		// signedElement, err := ctx.SignEnveloped(elementToSign)
		// if err != nil {
		// 	log.WithError(err).Error("cannot sign envelope")
		// 	d.Reject(false)
		// 	continue
		// }

		// Serialize the signed element. It is important not to modify the element
		// after it has been signed - even pretty-printing the XML will invalidate
		// the signature.
		// doc.SetRoot(signedElement)

		signedXML, err := doc.WriteToString()
		if err != nil {
			log.WithError(err).Error("cannot convert xml to string")
			d.Reject(false)
			continue
		}

		log.WithField("payload", signedXML).Info("signed xml")

		exportResult := models.ExportResult{
			Name:    exportTask.Name,
			Payload: string(signedXML),
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

	key, err := parseRsaPrivateKeyFromPemStr(keyBytes)
	failOnError(err, "Failed to parse RSA private key")

	certBytes, err := ioutil.ReadFile("/saml.crt")
	failOnError(err, "Failed to read certificate")

	block, _ := pem.Decode(certBytes)
	if block == nil {
		failOnError(errors.New("¯\\_(ツ)_/¯"), "Failed to decode certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		failOnError(err, "failed to parse certificate")
	}

	// cert, err := tls.X509KeyPair(certBytes, keyBytes)
	// failOnError(err, "Failed to load X509 key pair")

	// TODO: is there something else other than tls package that can be used?
	ks := &models.X509KeyStore{key, cert.Raw}

	var ksp *dsig.X509KeyStore
	ksp = new(dsig.X509KeyStore)
	*ksp = ks

	// Generate a key and self-signed certificate for signing
	// randomKeyStore := dsig.RandomKeyStoreForTest()
	ctx := &dsig.SigningContext{
		Hash:          crypto.SHA256,
		KeyStore:      *ksp,
		IdAttribute:   dsig.DefaultIdAttr,
		Prefix:        dsig.DefaultPrefix,
		Canonicalizer: dsig.MakeC14N11Canonicalizer(),
	}

	forever := make(chan bool)

	go task(msgs, session, ctx, forever)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
