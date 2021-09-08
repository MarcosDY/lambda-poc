package util

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/marcosdy/lambda-poc/secret"
)

// SaveSvid persist X509 SVID into provided folder.
func SaveSvid(svid *secret.Svid, outputPath string) error {
	log.Println("Getting spiffeID")
	block, _ := pem.Decode([]byte(svid.X509Svid))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Printf("Parse blocks: %v\n", err)
		return err
	}
	log.Printf("SPIFFE ID: %v", cert.URIs)

	log.Println("Writing /tmp/svid.pem.")
	log.Println(svid.X509Svid)

	if err := writeCerts("/tmp/svid.pem", svid.X509Svid); err != nil {
		return err
	}

	log.Println("Writing /tmp/svid.key")
	log.Println(svid.Bundle)
	if err := writeKey("/tmp/svid.key", svid.X509SvidKey); err != nil {
		return err
	}

	log.Println("Writing /tmp/bundle.pem")
	log.Println(svid.X509SvidKey)
	if err := writeCerts("/tmp/bundle.pem", svid.Bundle); err != nil {
		return err
	}

	// sort and write the keys by trust domain so the output is consistent
	federatedDomains := make([]string, 0, len(svid.FederatedBundles))
	for trustDomain := range svid.FederatedBundles {
		federatedDomains = append(federatedDomains, trustDomain)
	}

	for j, trustDomain := range federatedDomains {
		bundlePath := path.Join(outputPath, fmt.Sprintf("/tmp/federated_bundle.%d.pem", j))
		log.Printf("Writing federated bundle #%d for trust domain %s to file %s.\n", j, trustDomain, bundlePath)
		log.Printf("%+v \n", svid.FederatedBundles[trustDomain])
		if err := writeCerts(bundlePath, svid.FederatedBundles[trustDomain]); err != nil {
			return err
		}
	}

	return nil
}

func writeCerts(filename string, pemData string) error {
	return ioutil.WriteFile(filename, []byte(pemData), 0644) // nolint: gosec // expected permission for certificates
}

// writeKey takes a private key, formats as PEM, and writes it to filename
func writeKey(filename string, privateKey string) error {
	return ioutil.WriteFile(filename, []byte(privateKey), 0600)
}
