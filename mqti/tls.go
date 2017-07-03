package mqti

import "crypto/tls"

// NewTLSConfig ...
func NewTLSConfig(certFile, keyFile string) *tls.Config {
	var err error
	var cert tls.Certificate

	cert, err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{cert},
	}
}
