package generator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/spf13/viper"
	"math/big"
	"net"
	"os"
)

const (
	AgentIDLength = 16
	AgentIDPrefix = "AA"
)

var charSet = []byte("0123456789ABCDEF")

func RandID() string {
	ret := make([]byte, AgentIDLength-len(AgentIDPrefix))
	for i := 0; i < AgentIDLength-len(AgentIDPrefix); i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		ret[i] = charSet[num.Int64()]
	}
	return AgentIDPrefix + string(ret)
}

func IsAgentID(s string) bool {
	return len(s) == AgentIDLength && s[:2] == "AA"
}

func Generate() error {
	name := viper.GetString("name")
	goos := viper.GetString("os")
	arch := viper.GetString("arch")
	iFace := viper.GetString("host")
	port := viper.GetString("port")
	outFile := viper.GetString("outfile")
	siphonCert := viper.GetString("cert_file")

	certData, err := os.ReadFile(siphonCert)
	if err != nil {
		logger.Fatalf("could not read siphon certificate from %s: %v", siphonCert, err)
	}
	agentID := RandID()
	if len(name) == 0 {
		name = agentID
	}
	mainPath := viper.GetString("src_path")

	builder := NewBuilder("go", goos, arch)
	builder.AddSrcPath(mainPath)
	builder.SetFlags(
		Flag{"main.AgentID", agentID},
		Flag{"main.Interface", iFace},
		Flag{"main.Port", port},
		// Add certificate data base64 encoded to agent, so it is added to its rootCAs
		Flag{"main.ClientCertData", base64.StdEncoding.EncodeToString(certData)},
	)
	builder.SetOutFile(outFile)
	// Exits if unsuccessful
	builder.Build()

	conn := db.Initialize()
	agent := &db.Agent{
		AgentID:  agentID,
		Name:     name,
		Endpoint: fmt.Sprintf("https://%s/api", net.JoinHostPort(iFace, port)),
	}
	return conn.Add(agent)
}
