package agent

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/pygrum/siphon/internal/agent/controllers"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Fetcher struct {
	Agents []db.Agent
	Conn   *db.Conn
}

func NewFetcher() *Fetcher {
	f := &Fetcher{
		Conn: db.Initialize(),
	}
	f.Agents = f.Conn.Agents()
	return f
}

func (f *Fetcher) GetRecent() {
	for _, agent := range f.Agents {
		a := agent
		go func(agent *db.Agent) {
			resp, _ := f.BasicRequest(agent, "/samples", "", "")
			if resp == nil {
				return
			}
			body, _ := io.ReadAll(resp.Body)
			respObject := &controllers.SampleResponse{}
			if err := json.Unmarshal(body, respObject); err != nil {
				logger.Silentf("%v", err)
				return
			}
			if respObject.Status != controllers.StatusOk {
				logger.Silentf("request to (%s:%s) returned error %s", agent.Name, agent.AgentID, respObject.Status)
				return
			}
			if err := f.addSamples(respObject); err != nil {
				logger.Silentf("failed to add new samples: %v", err)
			}
		}(&a)
	}
}

func (f *Fetcher) addSamples(r *controllers.SampleResponse) error {
	for _, data := range r.Data {
		d := data
		d.ID = 0 // Unset ID (primary key) to let gorm create normal record
		go func(data db.Sample) {
			if f.Conn.SampleByHash(data.Hash) == nil {
				if err := f.Conn.Add(&data); err != nil {
					logger.Silentf("%v", err)
				}
			}
		}(d)
	}
	return nil
}

func (f *Fetcher) BasicRequest(a *db.Agent, endpoint, query, form string) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodGet, a.Endpoint+endpoint, strings.NewReader(form))
	if err != nil {
		logger.Silentf("unable to create new request for (%s:%s): %v", a.Name, a.AgentID, err)
		return nil, err
	}
	r.URL.RawQuery = query
	client, err := f.mTLSClient(a)
	if err != nil {
		logger.Silentf("failed to create mTLS client: %v", err)
		return nil, err
	}
	resp, err := client.Do(r)
	if err != nil {
		logger.Silentf("request failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func (f *Fetcher) mTLSClient(agent *db.Agent) (*http.Client, error) {
	serverCertFile := agent.CertPath
	certFile, keyFile := viper.GetString("cert_file"), viper.GetString("key_file")
	// Read server certificate file and add it to trusted certificate store (certPool). Right now I'm reading my cert instead of the servers
	caCert, err := os.ReadFile(serverCertFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            certPool,
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
				// function from github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc/electrum/electrum.go#L76-L111
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					certs := make([]*x509.Certificate, len(rawCerts))
					for i, asn1Data := range rawCerts {
						cert, err := x509.ParseCertificate(asn1Data)
						if err != nil {
							return fmt.Errorf("failed to parse certificate: %v", err)
						}
						certs[i] = cert
					}
					opts := x509.VerifyOptions{
						Roots:         certPool,
						CurrentTime:   time.Now(),
						DNSName:       "", // Skip hostname verification
						Intermediates: x509.NewCertPool(),
					}

					for i, cert := range certs {
						if i == 0 {
							continue
						}
						opts.Intermediates.AddCert(cert)
					}
					_, err := certs[0].Verify(opts)
					return err
				},
			},
		},
	}, nil
}

func (f *Fetcher) Download(agent *db.Agent, sha256Hash string) (io.ReadCloser, error) {
	q := url.Values{
		"sha256_hash": {sha256Hash},
	}
	resp, err := f.BasicRequest(agent, "/download", q.Encode(), "")
	if err != nil {
		return nil, err
	}
	if resp.Header.Get("Content-Type") == "application/json" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		respObject := &controllers.SampleResponse{}
		if err := json.Unmarshal(body, respObject); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("request to (%s:%s) returned error %s", agent.Name, agent.AgentID, respObject.Status)
	}
	return resp.Body, nil
}
