package main

import "github.com/Hatch1fy/errors"

func newOptions(env map[string]string) (op *Options, err error) {
	// Make options from environment
	o := makeOptions(env)

	// Validate options
	if err = o.Validate(); err != nil {
		// Options are not valid, return
		return
	}

	// Assign reference to options
	op = &o
	return
}

func makeOptions(env map[string]string) (o Options) {
	if o.Port = env["lets-encrypt-port"]; len(o.Port) == 0 {
		// Port hasn't been set, set as default
		o.Port = defaultPort
	}

	if o.TLSPort = env["lets-encrypt-tls-port"]; len(o.TLSPort) == 0 {
		// TLS port hasn't been set, set as default
		o.TLSPort = defaultTLSPort
	}

	if o.Directory = env["lets-encrypt-directory"]; len(o.Directory) == 0 {
		// Directory hasn't been set, set as default
		o.Directory = defaultDirectory
	}

	o.Email = env["lets-encrypt-email"]
	o.Domain = env["lets-encrypt-domain"]
	return
}

// Options are the options used for Let's Encrypt SSL procurement
type Options struct {
	Email     string
	Domain    string
	Directory string

	Port    string
	TLSPort string
}

// Validate will validate a set of options
func (o *Options) Validate() (err error) {
	var errs errors.ErrorList
	if len(o.Email) == 0 {
		// Email is required and not found, push error
		errs.Push(ErrInvalidEmail)
	}

	if len(o.Domain) == 0 {
		// Domain is required and not found, push error
		errs.Push(ErrInvalidDomain)
	}

	return errs.Err()
}
