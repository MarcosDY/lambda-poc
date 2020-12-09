package util

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/spiffe/go-spiffe/v2/proto/spiffe/workload"
)

// SaveSvids persist X509 SVID into provided folder.
func SaveSvids(resp *workload.X509SVIDResponse, outputPath string) error {
	if len(resp.Svids) == 0 {
		return errors.New("workload response contains no svids")
	}

	federatedBundles := make(map[string][]*x509.Certificate)
	for federatedDomainID, federatedBundleDER := range resp.FederatedBundles {
		federatedBundle, err := x509.ParseCertificates(federatedBundleDER)
		if err != nil {
			return fmt.Errorf("failed to parse bundle for federated domain %q: %v", federatedDomainID, err)
		}
		if len(federatedBundle) == 0 {
			return fmt.Errorf("no certificates in bundle for federated domain %q", federatedDomainID)
		}
		federatedBundles[federatedDomainID] = federatedBundle
	}

	for i, respSVID := range resp.Svids {
		if err := saveSVID(i, respSVID, federatedBundles, outputPath); err != nil {
			return fmt.Errorf("failed to parse svid entry %d for spiffe id %q: %v", i, respSVID.SpiffeId, err)
		}
	}

	return nil
}

func saveSVID(i int, svid *workload.X509SVID, federatedBundles map[string][]*x509.Certificate, outputPath string) error {
	certificates, err := x509.ParseCertificates(svid.X509Svid)
	if err != nil {
		return err
	}
	if len(certificates) == 0 {
		return errors.New("no certificates found")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(svid.X509SvidKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}
	signer, ok := privateKey.(crypto.Signer)
	if !ok {
		return fmt.Errorf("private key is type %T, not crypto.Signer", privateKey)
	}

	bundle, err := x509.ParseCertificates(svid.Bundle)
	if err != nil {
		return fmt.Errorf("failed to parse trust bundle: %v", err)
	}
	if len(bundle) == 0 {
		return errors.New("no certificates in trust bundle")
	}

	svidPath := path.Join(outputPath, fmt.Sprintf("svid.%v.pem", i))
	keyPath := path.Join(outputPath, fmt.Sprintf("svid.%v.key", i))
	bundlePath := path.Join(outputPath, fmt.Sprintf("bundle.%v.pem", i))

	log.Printf("Writing SVID #%d to file %s.\n", i, svidPath)

	if err := writeCerts(svidPath, certificates); err != nil {
		return err
	}

	log.Printf("Writing key #%d to file %s.\n", i, keyPath)
	err = writeKey(keyPath, signer)
	if err != nil {
		return err
	}

	log.Printf("Writing bundle #%d to file %s.\n", i, bundlePath)
	err = writeCerts(bundlePath, bundle)
	if err != nil {
		return err
	}

	// sort and write the keys by trust domain so the output is consistent
	federatedDomains := make([]string, 0, len(federatedBundles))
	for trustDomain := range federatedBundles {
		federatedDomains = append(federatedDomains, trustDomain)
	}

	for j, trustDomain := range federatedDomains {
		bundlePath := path.Join(outputPath, fmt.Sprintf("federated_bundle.%d.%d.pem", i, j))
		log.Printf("Writing federated bundle #%d for trust domain %s to file %s.\n", j, trustDomain, bundlePath)
		err = writeCerts(bundlePath, federatedBundles[trustDomain])
		if err != nil {
			return err
		}
	}

	return nil
}
func writeCerts(filename string, certs []*x509.Certificate) error {
	pemData := []byte{}
	for _, cert := range certs {
		// TODO: demonstration purposes only, remove it
		log.Printf("SPIFFE ID: %q, NotAfter: %s", cert.URIs, cert.NotAfter.String())
		b := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		}
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	return writeFile(filename, pemData)
}

// writeKey takes a private key, formats as PEM, and writes it to filename
func writeKey(filename string, privateKey crypto.PrivateKey) error {
	data, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}
	b := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: data,
	}

	return ioutil.WriteFile(filename, pem.EncodeToMemory(b), 0600)
}

// writeFile creates or truncates filename, and writes data to it
func writeFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644) // nolint: gosec // expected permission for certificates
}
