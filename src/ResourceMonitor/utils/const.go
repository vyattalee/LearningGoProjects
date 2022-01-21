package utils

import "time"

const SecretKey = "secret"
const TokenDuration = 15 * time.Minute

const (
	ServerCertFile   = "cert/server-cert.pem"
	ServerKeyFile    = "cert/server-key.pem"
	ClientCACertFile = "cert/ca-cert.pem"
)

const (
	ClientCertFile   = "cert/client-cert.pem"
	ClientKeyFile    = "cert/client-key.pem"
	ServerCACertFile = "cert/ca-cert.pem" //client & server on different nodes, it is different ca-cert.pem
)
