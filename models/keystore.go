package models

import "crypto/rsa"

type X509KeyStore struct {
	PrivateKey *rsa.PrivateKey
	Cert       []byte
}

func (ks *X509KeyStore) GetKeyPair() (*rsa.PrivateKey, []byte, error) {
	return ks.PrivateKey, ks.Cert, nil
}
