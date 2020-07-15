package main

import (
	"fmt"

	"github.com/Hatch1fy/errors"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/challenge/http01"
	"github.com/go-acme/lego/challenge/tlsalpn01"
	"github.com/go-acme/lego/lego"
	"github.com/hatchify/scribe"
	"github.com/vroomy/common"
)

const (
	// ErrInvalidEmail is returned when the let's encrypt email is missing from the environment values
	ErrInvalidEmail = errors.Error("cannot create SSL certificate without 'lets-encrypt-email' environment value")
	// ErrInvalidDomain is returned when the let's encrypt domain is missing from the environment values
	ErrInvalidDomain = errors.Error("cannot create SSL certificate without 'lets-encrypt-domain' environment value")
)

const (
	defaultPort    = "5001"
	defaultTLSPort = "5002"
)

var out *scribe.Scribe = scribe.New("Let's Encrypt")

// Init will be called by vroomy on initialization before Configure
func Init(env map[string]string) (err error) {
	var opts *Options
	// Get options from environment variables
	if opts, err = newOptions(env); err != nil {
		return
	}

	var u *User
	// Create a new user
	if u, err = newUser(opts.Email); err != nil {
		return
	}

	// Create a new configuration
	config := lego.NewConfig(u)
	config.Certificate.KeyType = certcrypto.RSA2048

	var client *lego.Client
	// Initialize a new client
	if client, err = newClient(config, opts); err != nil {
		return
	}

	// Register user using Client
	if err = u.Register(client); err != nil {
		return
	}

	// Make request
	request := makeRequest(opts.Domain)

	var certificates *certificate.Resource
	// Obtain certificates
	if certificates, err = client.Certificate.Obtain(request); err != nil {
		return
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	fmt.Printf("%#v\n", certificates)
	return
}

// Backend returns the underlying backend to the plugin
func Backend() interface{} {
	return nil
}

// Load will be called by vroomy after plugin initialization
func Load(p common.Plugins) (err error) {

	return
}

func newClient(config *lego.Config, opts *Options) (client *lego.Client, err error) {
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
