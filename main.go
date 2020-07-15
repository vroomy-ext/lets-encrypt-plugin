package main

import (
	"fmt"

	"github.com/Hatch1fy/errors"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
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
	defaultPort      = "5001"
	defaultTLSPort   = "5002"
	defaultDirectory = "tls"
)

var (
	// Output writer
	out *scribe.Scribe = scribe.New("Let's Encrypt")
	// Default global registration options
	registrationOpts = registration.RegisterOptions{TermsOfServiceAgreed: true}
)

// Init will be called by vroomy on initialization before Configure
func Init(env map[string]string) (err error) {
	var opts *Options
	// Get options from environment variables
	if opts, err = newOptions(env); err != nil {
		err = fmt.Errorf("error parsing options: %v", err)
		return
	}

	var ok bool
	if ok, err = needsCertificate(opts.Directory); !ok || err != nil {
		err = fmt.Errorf("error checking certificate: %v", err)
		return
	}

	var u *User
	// Create a new user
	if u, err = newUser(opts.Email); err != nil {
		err = fmt.Errorf("error creating user: %v", err)
		return
	}

	// Create a new configuration
	config := lego.NewConfig(u)
	config.Certificate.KeyType = certcrypto.RSA2048

	var client *lego.Client
	// Initialize a new client
	if client, err = newClient(opts, config); err != nil {
		err = fmt.Errorf("error initializing client: %v", err)
		return
	}

	// Register user using Client
	if err = u.Register(client); err != nil {
		err = fmt.Errorf("error registering user \"%s\": %v", u.Email, err)
		return
	}

	// Make request
	request := makeRequest(opts.Domain)

	var certificates *certificate.Resource
	// Obtain certificates
	if certificates, err = client.Certificate.Obtain(request); err != nil {
		err = fmt.Errorf("error obtaining certificates: %v", err)
		return
	}

	// Save certificates to file
	if err = saveCertificates(opts.Directory, certificates); err != nil {
		err = fmt.Errorf("error saving certificates: %v", err)
		return
	}

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
