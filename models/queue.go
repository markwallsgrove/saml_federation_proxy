package models

type IngestJob struct {
	EntityID string
	XML      []byte
	Checksum []byte
}
