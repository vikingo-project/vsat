package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"math/big"
	mrand "math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

func MD5(content string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(content))
	cipherBytes := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherBytes)
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}
func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func GenerateCert(hostname string) (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{hostname},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	template.DNSNames = append(template.DNSNames, hostname)
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	outCert := &bytes.Buffer{}
	outKey := &bytes.Buffer{}
	pem.Encode(outCert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	pem.Encode(outKey, pemBlockForKey(priv))

	return tls.X509KeyPair(outCert.Bytes(), outKey.Bytes())
}

func GenerateCertAndKey(host string) (string, string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{host},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	template.DNSNames = append(template.DNSNames, host)
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	outCert := &bytes.Buffer{}
	outKey := &bytes.Buffer{}
	pem.Encode(outCert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	pem.Encode(outKey, pemBlockForKey(priv))
	return outCert.String(), outKey.String(), nil
}

type network struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

func GetNetworks() ([]network, error) {
	var networks []network
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		// fmt.Println(iface.Name)
		addrs, err := iface.Addrs()
		if err != nil {
			// Interface has no address
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP

			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			// fmt.Println(ip)
			networks = append(networks, network{iface.Name, ip.String()})
			break
		}
	}
	return networks, nil
}

func AmAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("admin no")
		return false
	}
	fmt.Println("admin yes")
	return true
}

func IsDevMode() bool {
	devMode := os.Getenv("DEV")
	return devMode == "1" || devMode == "true"
}

func ExtractSettings(obj interface{}, data interface{}) error {
	err := mapstructure.Decode(data, &obj)
	if err != nil {
		return err
	}
	return nil
}

// EasyHash returns a string with length 6 bytes
func EasyHash(long bool) string {
	mrand.Seed(time.Now().UTC().UnixNano())
	tablePolynomial := crc32.MakeTable(0xEDB88320)
	hash := crc32.New(tablePolynomial)
	if _, err := io.Copy(hash, strings.NewReader(time.Now().String())); err != nil {
		return ""
	}
	hashInBytes := hash.Sum(nil)[:]
	out := hex.EncodeToString(hashInBytes) + fmt.Sprintf("%02x%02x", mrand.Intn(255), mrand.Intn(255))
	if long {
		for i := 0; i < 8; i++ {
			out += fmt.Sprintf("%02x", mrand.Intn(255))
		}
	}
	return out
}

func ExtractIP(s string) string {
	IP := ""
	parts := strings.Split(s, ":")
	if len(parts) > 2 {
		IP = strings.Join(parts[:len(parts)-1], ":")
	} else {
		IP = parts[0]
	}
	return IP
}

func PrintDebug(f string, args ...interface{}) {
	if IsDevMode() {
		log.Printf(f, args...)
	}
}
