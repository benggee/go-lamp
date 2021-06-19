package httpz

import "time"

type (
	HttpConf struct {
		Addr    string
		Timeout int64 `json:",default=5000"`
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
)
