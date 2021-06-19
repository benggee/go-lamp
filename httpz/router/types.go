package router

import (
	"net/http"
	"time"
)

type (
	Route struct {
		Method  string
		Path    string
		Handler http.HandlerFunc
	}

	JwtConf struct {
		Enabled    bool
		Secret     string
		PrevSecret string
	}

	Signature struct {
		SignatureConf
		Enabled bool
	}

	SignatureConf struct {
		Strict      bool          `json:",default=false"`
		Expiry      time.Duration `json:",default=1h"`
		PrivateKeys []PrivateKeyConf
	}

	PrivateKeyConf struct {
		Fingerprint string
		KeyFile     string
	}

	Routers struct {
		Jwt       JwtConf
		Signature Signature
		Routes    []Route
	}
)
