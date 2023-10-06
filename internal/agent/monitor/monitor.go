package monitor

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pygrum/siphon/internal/agent/controllers"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Scout struct {
	Agent        *controllers.Agent
	Logger       *logger.Logger
	MonitorPaths []monitorPath
}

type monitorPath struct {
	Path      string
	Recursive bool
}

var watcher *fsnotify.Watcher

func NewScout(agent *controllers.Agent) (*Scout, error) {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	var folders []monitorPath
	if err := viper.UnmarshalKey("monitor_folders", &folders); err != nil {
		return nil, err
	}
	l, err := logger.NewLogger("./agent_monitor.log")
	if err != nil {
		return nil, err
	}
	return &Scout{
		Agent:        agent,
		MonitorPaths: folders,
		Logger:       l,
	}, nil
}

func (s *Scout) Start() error {
	for _, d := range s.MonitorPaths {
		if err := addDirectory(&d); err != nil {
			return err
		}
	}
	go s.start()
	return nil
}

func (s *Scout) RunTLS() error {
	s.Agent.CertFile = viper.GetString("cert_file")
	s.Agent.KeyFile = viper.GetString("key_file")

	cert := s.Agent.ClientCertData

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	tlsConfig := &tls.Config{
		ClientCAs:          certPool,
		InsecureSkipVerify: true,
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
	}
	return s.Agent.RunTLS(tlsConfig)
}

func (s *Scout) start() {
	for {
		select {
		case event := <-watcher.Events:
			if event.Has(fsnotify.Create) {
				s.Logger.Write(logger.Sinfof("file created: %s", event.Name))
				if err := s.Add(event.Name); err != nil {
					s.Logger.Write(logger.Serrorf("failed to add %s to database: %v", event.Name, err))
				}
			} else if event.Has(fsnotify.Write) {
				s.Logger.Write(logger.Sinfof("file written: %s", event.Name))
				if err := s.Add(event.Name); err != nil {
					s.Logger.Write(logger.Serrorf("failed to add %s to database: %v", event.Name, err))
				}
			} else if event.Has(fsnotify.Rename) {
				s.Logger.Write(logger.Sinfof("file renamed: %s", event.Name))
			} else if event.Has(fsnotify.Remove) {
				s.Logger.Write(logger.Sinfof("file removed: %s", event.Name))
			} else {
				s.Logger.Write(logger.Sinfof("file mode changed: %s", event.Name))
			}
			break
		case event := <-watcher.Errors:
			s.Logger.Write(logger.Serrorf("error: %s", event.Error()))
			break
		}
	}
}

func (s *Scout) Add(file string) error {
	conn := s.Agent.Conn
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	// Don't save empty files
	if fi.Size() == 0 {
		return nil
	}
	// If file caching is on, save copy of file to agent folder (protected)
	if viper.GetBool("cache") {
		b, _ := os.ReadFile(file)
		cwd, _ := os.Getwd()
		newFileName := filepath.Join(cwd, filepath.Base(file))
		_ = os.WriteFile(newFileName, b, 0600)
		file = newFileName
	}

	sample := &db.Sample{
		// Set creation time manually
		Model: gorm.Model{
			CreatedAt: fi.ModTime(),
		},
		Name:       filepath.Base(file),
		Path:       file,
		FileType:   filepath.Ext(file)[1:], // Skip the leading '.'
		FileSize:   uint(fi.Size()),
		Source:     s.Agent.ID,
		Hash:       fileHash(file),
		UploadTime: fi.ModTime(),
	}
	return conn.Add(sample)
}

func fileHash(file string) string {
	f, _ := os.Open(file)
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func addDirectory(mp *monitorPath) error {
	if !mp.Recursive {
		return watcher.Add(mp.Path)
	}
	// Add folders to watch recursively
	return filepath.Walk(mp.Path, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}
