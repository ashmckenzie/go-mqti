package littlefly

import "crypto/tls"

func NewTLSConfig(certFile, keyFile string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{cert},
	}
}
