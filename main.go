package main

import (
	"fmt"

	"github.com/Hatch1fy/errors"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/lego"
	legolog "github.com/go-acme/lego/log"
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
	defaultPort      = "80"
	defaultTLSPort   = "443"
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
		// Error parsing options, return
		err = fmt.Errorf("error parsing options: %v", err)
		return
	}

	// Create log wrapper
	var lw logWrapper
	// Set log wrapper as our lego logger
	legolog.Logger = &lw

	var ok bool
	// Check to see if a new certificate is needed
	if ok, err = needsCertificate(opts.Directory); err != nil {
		// Error encountered while checking certificate, return
		err = fmt.Errorf("error checking certificate: %v", err)
		return
	} else if !ok {
		// No certificate needed, return
		return
	}

	out.Notification("Certificate is expired (or expiring soon), executing renewal process")

	var u *User
	// Create a new user
	if u, err = newUser(opts.Email); err != nil {
		// Error creating new user, return
		err = fmt.Errorf("error creating user: %v", err)
		return
	}

	// Create a new configuration
	config := lego.NewConfig(u)
	config.Certificate.KeyType = certcrypto.RSA2048

	var client *lego.Client
	// Initialize a new client
	if client, err = newClient(opts, config); err != nil {
		// Error initializing new client, return
		err = fmt.Errorf("error initializing client: %v", err)
		return
	}

	out.Success("Client created")

	// Register user using Client
	if err = u.Register(client); err != nil {
		// Error registering user, return
		err = fmt.Errorf("error registering user \"%s\": %v", u.Email, err)
		return
	}

	out.Success("User registered")

	// Make request
	request := makeRequest(opts.Domain)

	var certificates *certificate.Resource
	// Obtain certificates
	if certificates, err = client.Certificate.Obtain(request); err != nil {
		// Error obtaining certificates, return
		err = fmt.Errorf("error obtaining certificates: %v", err)
		return
	}

	out.Success("Certificates obtained")

	// Save certificates to file
	if err = saveCertificates(opts.Directory, certificates); err != nil {
		// Error saving certificates, return
		err = fmt.Errorf("error saving certificates: %v", err)
		return
	}

	out.Success("Certificate renewal process complete")
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
