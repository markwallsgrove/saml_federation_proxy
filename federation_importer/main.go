package main

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/crewjam/go-xmlsec"
	"github.com/markwallsgrove/saml_federation_proxy/models"
	"github.com/streadway/amqp"
)

func downloadFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func verifySignature(buff io.ReadCloser, key []byte) ([]byte, error) {
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

func calcChecksum(bytes []byte) []byte {
	hasher := md5.New()
	hasher.Write(bytes)
	return hasher.Sum(nil)
}

func createQueue(channel *amqp.Channel) (amqp.Queue, error) {
	return channel.QueueDeclare(
		"ingest_queue", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
}

func publishMessage(job *models.IngestJob, channel *amqp.Channel, queueName string) error {
	bytes, err := json.Marshal(job)
	if err != nil {
		log.Fatal("cannot marshal ingest job", err)
		return err
	}

	return channel.Publish(
		"",        // exchange
		queueName, // queue name
		false,     // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         bytes,
		},
	)
}

func getEntityDescriptors(entityURL string, pem []byte) (*models.EntitiesDescriptor, error) {
	buff, err := downloadFile(entityURL)

	if err != nil {
		log.Fatalln("error downloading file:", err)
		return nil, err
	}

	doc, err := verifySignature(buff, pem)

	if err != nil {
		log.Fatalln("signature error:", err)
		return nil, err
	}

	entitiesDescriptors, err := unmarshallXML(doc)
	if err != nil {
		log.Fatalln("unmarhsalling error:", err)
		return nil, err
	}

	return entitiesDescriptors, nil
}

func main() {
	// TODO: move to db
	ukfedCert := []byte(`-----BEGIN CERTIFICATE-----
MIIDxzCCAq+gAwIBAgIJAOwuoY8tsvYGMA0GCSqGSIb3DQEBCwUAMHoxCzAJBgNV
BAYTAkdCMUMwQQYDVQQKDDpVSyBBY2Nlc3MgTWFuYWdlbWVudCBGZWRlcmF0aW9u
IGZvciBFZHVjYXRpb24gYW5kIFJlc2VhcmNoMSYwJAYDVQQDDB1VSyBGZWRlcmF0
aW9uIE1ldGFkYXRhIFNpZ25lcjAeFw0xNDA4MjYxMjIwMjhaFw0zNzEyMzExMjIw
MjhaMHoxCzAJBgNVBAYTAkdCMUMwQQYDVQQKDDpVSyBBY2Nlc3MgTWFuYWdlbWVu
dCBGZWRlcmF0aW9uIGZvciBFZHVjYXRpb24gYW5kIFJlc2VhcmNoMSYwJAYDVQQD
DB1VSyBGZWRlcmF0aW9uIE1ldGFkYXRhIFNpZ25lcjCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAOqtfMvCmBuQudC4/jZFPYkHDNHFyp1FA3KJihIUXppF
vrecrO2wG5CpyqB1mZ+MlKf4jKcTMGBIXC2klD+FyrEdJMBhO6vRmJnNphg3uNZM
ks0NqIaZmtgc7e8435nMhqLHV95UK2oCLcT4gZrTaXa2vt9kukTOijB0KqDIfEG5
369EHXPItApAEeMlHebbWndl5n2I16nya/LeaoiU9qJ6sVz4xd1UtUesewrmYVKg
PA2JYEpovmnr13sTnGssai5Db/FkrE2NJ4Q4drbPYcwincUo/UXzrtuPclr+l3JE
gjtvDzPrBxxvK0S/gARrbKz5tk4LDLkYsj4PKlwVS+UCAwEAAaNQME4wHQYDVR0O
BBYEFE9HhBuMxrzBYOj1Kj/3gtzAgtUEMB8GA1UdIwQYMBaAFE9HhBuMxrzBYOj1
Kj/3gtzAgtUEMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBALJkjT3K
QL3w3xNfVe27nEOY44K2AZiu4IhqmRSslcyMnhnxrovEhLL3ieKFXQ+QFIkzVdR5
BcO3NrSIz5V6b+mHtr5IjqLFHzOzzjw/3i8LddGOsApJiav+JrU1CGJXCU4cwYDN
hAyfuAlrrEEL2lWMU1L1ZTzHsG1yWTfukfuvTftY5BwZ/dgANgIWwLDhvL6CAQZ3
g5XteFPyChU0Z7b3XAHdVNHDa2VzWSsSUDtSQZ9DyTuqSjZH1q2/qtdMcrbJpdMB
cndOf1pZRLzb6a+akIYi//1qO48HpB4wouH9gS3ZER+rNBhVWu301UYxoVI7o8mG
Yq7dENJce7lO9yE=
-----END CERTIFICATE-----`)

	// TODO: move to db
	entityURL := os.Getenv("ENTITY_URL")

	entitiesDescriptors, err := getEntityDescriptors(entityURL, ukfedCert)
	if err != nil {
		log.Fatal("cannot download entity descriptors", err)
		return
	}

	conn, err := amqp.Dial(os.Getenv("QUEUE_CONN"))
	if err != nil {
		log.Fatal("cannot connect to queue", err)
		return
	}

	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("cannot create channel in queue", err)
		return
	}

	queue, err := createQueue(channel)
	if err != nil {
		log.Fatal("cannot create queue")
		return
	}

	for _, entityDescriptor := range entitiesDescriptors.EntityDescriptor {
		entityID := entityDescriptor.EntityID
		entityDescriptorXML, _ := marshallEntityDescriptor(entityDescriptor)
		checksum := calcChecksum(entityDescriptorXML)

		ingestJob := &models.IngestJob{
			EntityID: entityID,
			XML:      entityDescriptorXML,
			Checksum: checksum,
		}

		err = publishMessage(ingestJob, channel, queue.Name)

		if err != nil {
			log.Fatal("cannot publish message", err)
		}
	}

	log.Println("fin.")
}
