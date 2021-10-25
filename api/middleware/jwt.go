package middleware

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
)

var publicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAnPq4vHhlaGYzGtEnEXe1
vQoFBCBeCmOoIXUya1TpVbrscLRrBXLzDdqBnOyNZfAFvzEj1ghZTlpecATiL6C6
O8PAVdFF82Jf+8NmbMzw2a1AUjhtxLrvxOmqwg/yn7e2wLol9xnml4bbGr0iIszK
wNEPHUlDvBR4APYjH9DkPDpG+wYRUKuoIKNEEQf/uhUEyJdn1Bx+1ge5m1n91Ibo
K0Y2wn6mR85RHc5t+50iGQTXz7xJwPjSZQcoZjSB4T0WL2CjdsHqyVxjX0L3TuF8
6VnohzkSLDuEMfciuRi+VTDKawcMvDUoijtbJHPe9iZqpa7LeLFk2cxfBei8waOI
FCC1Sh4YZWpnXpgcnQ5KIT0yKq2WVIBSsCUBfUCwM0QaiaPjyIkYlBKIvPKjfHj5
s7hWTDC2yBU4npTWhi/y57kCgYbOJjj3SEy8Kb23VYF85TvVEnDSahmooMi6642n
c+CCk5UUv0bASlkMRAH8UupA7ZDSUbwz7ZfAOWvz/qGZRcuuPWa6c/doxaK5hgff
hjy0qHQj/rbT08qC6OF3vop9z25aljgMsszamwQJKM4hCLTdfwJeFJ8MK/5BwzB9
3jnqcRyu7toZPQxOuJp1Ckd2xqMB1DsmzNFiWDvH59FoK3UboxdF4zzZGbgqDXBv
tLLHT5MavkIhYxaAh33D8s8CAwEAAQ==
-----END PUBLIC KEY-----`

// JWT is Fully Configured JWT Middleware to Support CWBI-Auth
var JWT = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningMethod: "RS512",
	KeyFunc: func(t *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	},
})

// JWTMock is JWT Middleware
var JWTMock = middleware.JWTWithConfig(
	middleware.JWTConfig{
		SigningKey: []byte("mock"),
	},
)