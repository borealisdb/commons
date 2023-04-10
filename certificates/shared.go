package certificates

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"net"
	"time"
)

type issuer struct {
	key  crypto.Signer
	cert *x509.Certificate
}

func makeRootCert() (*issuer, certsPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, certsPair{}, err
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, certsPair{}, err
	}
	skid, err := calculateSKID(privateKey.Public())
	if err != nil {
		return nil, certsPair{}, err
	}
	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "borealis root ca " + hex.EncodeToString(serial.Bytes()[:3]),
		},
		SerialNumber: serial,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(100, 0, 0),

		SubjectKeyId:          skid,
		AuthorityKeyId:        skid,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, privateKey.Public(), privateKey)
	if err != nil {
		return nil, certsPair{}, err
	}

	certPEM := new(bytes.Buffer)
	if err := pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return nil, certsPair{}, err
	}

	certPrivateKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(certPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return nil, certsPair{}, err
	}

	certificate, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, certsPair{}, err
	}

	return &issuer{
			key:  privateKey,
			cert: certificate,
		}, certsPair{
			Key:  certPrivateKeyPEM,
			Cert: certPEM,
		}, nil
}

func parseIPs(ipAddresses []string) ([]net.IP, error) {
	var parsed []net.IP
	for _, s := range ipAddresses {
		p := net.ParseIP(s)
		if p == nil {
			return nil, fmt.Errorf("invalid IP address %s", s)
		}
		parsed = append(parsed, p)
	}
	return parsed, nil
}

func calculateSKID(pubKey crypto.PublicKey) ([]byte, error) {
	spkiASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}
	_, err = asn1.Unmarshal(spkiASN1, &spki)
	if err != nil {
		return nil, err
	}
	skid := sha1.Sum(spki.SubjectPublicKey.Bytes)
	return skid[:], nil
}

type certsPair struct {
	Key  *bytes.Buffer
	Cert *bytes.Buffer
}

func sign(iss *issuer, domains []string, ipAddresses []string) (certsPair, error) {
	var cn string
	if len(domains) > 0 {
		cn = domains[0]
	} else if len(ipAddresses) > 0 {
		cn = ipAddresses[0]
	} else {
		return certsPair{}, fmt.Errorf("must specify at least one domain name or IP address")
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return certsPair{}, err
	}

	parsedIPs, err := parseIPs(ipAddresses)
	if err != nil {
		return certsPair{}, err
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return certsPair{}, err
	}
	template := &x509.Certificate{
		DNSNames:    domains,
		IPAddresses: parsedIPs,
		Subject: pkix.Name{
			CommonName: cn,
		},
		SerialNumber: serial,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(2, 0, 30),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, template, iss.cert, privateKey.Public(), iss.key)
	if err != nil {
		return certsPair{}, err
	}

	certPEM := new(bytes.Buffer)
	if err := pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return certsPair{}, err
	}

	certPrivateKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(certPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return certsPair{}, err
	}

	return certsPair{
		Key:  certPrivateKeyPEM,
		Cert: certPEM,
	}, nil
}

func main2() error {
	issuer, _, err := makeRootCert()
	if err != nil {
		return err
	}
	_, err = sign(issuer, []string{}, []string{})
	return err
}
