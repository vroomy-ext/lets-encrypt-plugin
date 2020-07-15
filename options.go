package main

import "github.com/Hatch1fy/errors"

func newOptions(env map[string]string) (op *Options, err error) {
	o := makeOptions(env)
	if err = o.Validate(); err != nil {
		return
	}

	op = &o
	return
}

func makeOptions(env map[string]string) (o Options) {
	if o.Port = env["lets-encrypt-port"]; len(o.Port) == 0 {
		o.Port = defaultPort
	}

	if o.TLSPort = env["lets-encrypt-tls-port"]; len(o.TLSPort) == 0 {
		o.TLSPort = defaultTLSPort
	}

	if o.Directory = env["lets-encrypt-directory"]; len(o.Directory) == 0 {
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
		errs.Push(ErrInvalidEmail)
	}

	if len(o.Domain) == 0 {
		errs.Push(ErrInvalidDomain)
	}

	return errs.Err()
}
