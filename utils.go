package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/challenge/http01"
	"github.com/go-acme/lego/challenge/tlsalpn01"
	"github.com/go-acme/lego/lego"
)

func newClient(opts *Options, config *lego.Config) (client *lego.Client, err error) {
	// A client facilitates communication with the CA server.
	if client, err = lego.NewClient(config); err != nil {
		return
	}

	// Set HTTP provider
	if err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", opts.Port)); err != nil {
		return
	}

	// Set TLS provider
	if err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", opts.TLSPort)); err != nil {
		return
	}

	return
}

func makeRequest(domains ...string) (r certificate.ObtainRequest) {
	r.Domains = domains
	r.Bundle = true
	return
}

func saveFile(name string, bs []byte) (err error) {
	var f *os.File
	if f, err = os.Create(name + ".tmp"); err != nil {
		return
	}
	defer f.Close()

	if _, err = f.Write(bs); err != nil {
		return
	}

	return
}

func saveAndReplaceFile(name string, bs []byte) (err error) {
	if err = saveFile(name, bs); err != nil {
		return
	}

	return replaceFile(name)
}

func replaceFile(name string) (err error) {
	return os.Rename(name+".tmp", name)
}

func saveCertificates(dir string, certificates *certificate.Resource) (err error) {
	if err = os.MkdirAll(dir, 0755); err != nil {
		// Error encountered while creating directory
		return
	}

	// Save SSL certificate
	if err = saveAndReplaceFile(path.Join(dir, "server.crt"), certificates.Certificate); err != nil {
		return
	}

	// Save SSL key
	if err = saveAndReplaceFile(path.Join(dir, "server.key"), certificates.PrivateKey); err != nil {
		return
	}

	// Save SSL CSR
	if err = saveAndReplaceFile(path.Join(dir, "server.csr"), certificates.CSR); err != nil {
		return
	}

	// Save SSL certificate meta URL file
	if err = saveAndReplaceFile(path.Join(dir, "server.url"), []byte(certificates.CertURL)); err != nil {
		return
	}

	return
}

func needsCertificate(dir string) (ok bool, err error) {
	var bs []byte
	// Set filename value as server.crt within the target directory
	filename := path.Join(dir, "server.crt")
	// Read target file
	if bs, err = ioutil.ReadFile(filename); err != nil {
		// We encountered an error reading the file, which most likely means that the file
		// does not exist. We can safely assume that a certificate is needed.
		return true, nil
	}

	var cert *x509.Certificate
	// Parse existing SSL certificate
	if cert, err = parseCertificate(filename, bs); err != nil {
		// We encountered an error parsing the certificate, return
		return
	}

	// Check to see if the certificate is expired, if expired - we need a certificate
	ok = isCertificateExpired(cert)
	return
}

func parseCertificate(filename string, bs []byte) (cert *x509.Certificate, err error) {
	var block *pem.Block
	// Decode PEM block
	if block, _ = pem.Decode([]byte(bs)); block == nil {
		err = fmt.Errorf("error parsing certificate PEM for source of %s", filename)
		return
	}

	// Parse certificate from PEM block
	if cert, err = x509.ParseCertificate(block.Bytes); err != nil {
		err = fmt.Errorf("failed to parse certificate: " + err.Error())
		return
	}

	return
}

func isCertificateExpired(cert *x509.Certificate) (expired bool) {
	now := time.Now()
	if now.Before(cert.NotBefore) {
		// Certificate is not yet valid, return true
		return true
	}

	if now.After(cert.NotAfter) {
		// Certificate has expired, return true
		return true
	}

	// Set value for a week from now
	nextWeek := now.Add(time.Hour * 24 * 7)
	// Check to see if certificate expires a week from now
	return nextWeek.After(cert.NotAfter)
}
