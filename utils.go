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

func newClient(opts *Options, config *lego.Config) (client *lego.Client, err error) {
	// A client facilitates communication with the CA server.
	if client, err = lego.NewClient(config); err != nil {
		return
	}

	// We specify an http port of 5002 and an tls port of 5001 on all interfaces
	// because we aren't running as root and can't bind a listener to port 80 and 443
	// (used later when we attempt to pass challenges). Keep in mind that you still
	// need to proxy challenge traffic to port 5002 and 5001.
	if err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", opts.Port)); err != nil {
		return
	}

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

func saveCertificates(dir string, certificates *certificate.Resource) (err error) {
	if err = os.MkdirAll(dir, 0755); err != nil {
		return
	}

	if err = saveAndReplaceFile(path.Join(dir, "server.crt"), certificates.Certificate); err != nil {
		return
	}

	if err = saveAndReplaceFile(path.Join(dir, "server.key"), certificates.PrivateKey); err != nil {
		return
	}

	if err = saveAndReplaceFile(path.Join(dir, "server.csr"), certificates.CSR); err != nil {
		return
	}

	if err = saveAndReplaceFile(path.Join(dir, "server.url"), []byte(certificates.CertURL)); err != nil {
		return
	}

	return
}

func needsCertificate(dir string) (ok bool, err error) {
	filename := path.Join(dir, "server.crt")
	var bs []byte
	if bs, err = ioutil.ReadFile(filename); err != nil {
		ok = true
		err = nil
		return
	}

	var block *pem.Block
	if block, _ = pem.Decode([]byte(bs)); block == nil {
		err = fmt.Errorf("error parsing certificate PEM for source of %s", filename)
		return
	}

	var cert *x509.Certificate
	if cert, err = x509.ParseCertificate(block.Bytes); err != nil {
		err = fmt.Errorf("failed to parse certificate: " + err.Error())
		return
	}

	now := time.Now()
	ok = now.Before(cert.NotBefore) || now.After(cert.NotAfter)
	return
}
