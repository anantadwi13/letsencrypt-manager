package internal

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type certbot struct {
	config Config
}

func NewCertbot(config Config) CertificateManager {
	return &certbot{config: config}
}

func (c *certbot) GetAll(ctx context.Context) ([]*Certificate, error) {
	var (
		separatorStr = "- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -"
		notFoundStr  = "No certificates found."
		//foundStr       = "Found the following certs:"
		certNameStr    = "  Certificate Name: "
		serialNumStr   = "    Serial Number: "
		keyTypeStr     = "    Key Type: "
		domainsStr     = "    Domains: "
		expiryDateStr  = "    Expiry Date: "
		certPathStr    = "    Certificate Path: "
		privKeyPathStr = "    Private Key Path: "
	)

	response, err := runCommand(ctx, "certbot", "certificates")
	if err != nil {
		return nil, err
	}
	//resByte, err := os.ReadFile("./temp/example.response")
	//if err != nil {
	//	return nil, err
	//}
	//response = string(resByte)
	resSplit := strings.Split(response, "\n")
	isOpened := false
	var certs []*Certificate
	var currentCert *Certificate
	for _, str := range resSplit {
		if str == "" {
			continue
		}
		if strings.Contains(str, separatorStr) {
			isOpened = !isOpened
		}
		if isOpened {
			if str == notFoundStr {
				return nil, nil
			}
			if strings.Contains(str, certNameStr) {
				if currentCert != nil {
					certs = append(certs, currentCert)
				}
				currentCert = new(Certificate)
				certName := strings.ReplaceAll(str, certNameStr, "")
				currentCert.Name = certName
				continue
			}
			if currentCert == nil {
				continue
			}
			if strings.Contains(str, serialNumStr) {
				serialNumber := strings.ReplaceAll(str, serialNumStr, "")
				currentCert.SerialNumber = serialNumber
				continue
			}
			if strings.Contains(str, keyTypeStr) {
				keyType := strings.ReplaceAll(str, keyTypeStr, "")
				currentCert.KeyType = keyType
				continue
			}
			if strings.Contains(str, domainsStr) {
				domains := strings.ReplaceAll(str, domainsStr, "")
				currentCert.Domains = strings.Split(domains, " ")
				continue
			}
			if strings.Contains(str, expiryDateStr) {
				expiryDate := strings.ReplaceAll(str, expiryDateStr, "")
				date, err := time.Parse("2006-01-02 15:04:05-07:00", expiryDate[:25])
				if err != nil {
					return nil, err
				}
				currentCert.ExpiryDate = date
				continue
			}
			if strings.Contains(str, certPathStr) {
				certPath := strings.ReplaceAll(str, certPathStr, "")
				cert, err := readFile(certPath)
				if err != nil {
					return nil, err
				}
				currentCert.Public = cert
				continue
			}
			if strings.Contains(str, privKeyPathStr) {
				privKeyPath := strings.ReplaceAll(str, privKeyPathStr, "")
				privKey, err := readFile(privKeyPath)
				if err != nil {
					return nil, err
				}
				currentCert.Private = privKey
				continue
			}
		} else {
			if currentCert != nil {
				certs = append(certs, currentCert)
			}
		}
	}
	return certs, nil
}

func (c *certbot) Get(ctx context.Context, domain string) (*Certificate, error) {
	certs, err := c.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var cert *Certificate
	for _, c := range certs {
		if c.Name == domain {
			cert = c
			break
		}
	}
	return cert, nil
}

func (c *certbot) Add(ctx context.Context, email string, domain string, altDomains ...string) (*Certificate, error) {
	args := []string{
		"certonly", "--webroot", "-n", "-m", email, "--agree-tos", "-w", c.config.PublicStaticPath(), "-d", domain,
	}
	for _, altDomain := range altDomains {
		args = append(args, "-d", altDomain)
	}
	cmd, err := runCommand(ctx, "certbot", args...)
	if err != nil {
		raw := strings.Split(cmd, "\n")
		for _, s := range raw {
			log.Println(s)
		}
		return nil, err
	}
	cert, err := c.Get(ctx, domain)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func (c *certbot) Delete(ctx context.Context, domain string) error {
	cmd, err := runCommand(ctx, "certbot", "delete", "--cert-name", domain, "-n")
	if err != nil {
		raw := strings.Split(cmd, "\n")
		for _, s := range raw {
			log.Println(s)
		}
		return err
	}
	return nil
}

func (c *certbot) RenewAll(ctx context.Context) error {
	cmd, err := runCommand(ctx, "certbot", "renew")
	if err != nil {
		raw := strings.Split(cmd, "\n")
		for _, s := range raw {
			log.Println(s)
		}
		return err
	}
	return nil
}

func (c *certbot) Renew(ctx context.Context, domain string) error {
	cmd, err := runCommand(ctx, "certbot", "renew", "--cert-name", domain)
	if err != nil {
		raw := strings.Split(cmd, "\n")
		for _, s := range raw {
			log.Println(s)
		}
		return err
	}
	return nil
}

func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	raw, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	resBytes, err := io.ReadAll(raw)
	if err != nil {
		return "", err
	}
	errBytes, err := io.ReadAll(stdErr)
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return string(errBytes), err
	}
	return string(resBytes), nil
}

func readFile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
