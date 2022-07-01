package main

import (
	"encoding/json"
	"os"

	"github.com/microsoft/notation-cose/internal/version"
	"github.com/notaryproject/notation-go/plugin"
	"github.com/urfave/cli/v2"
)

var metadatCommand = &cli.Command{
	Name:   string(plugin.CommandGetMetadata),
	Usage:  "Get plugin metadata",
	Action: runGetMetaData,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:      "file",
			Usage:     "request json file",
			TakesFile: true,
			Hidden:    true,
		},
	},
}

var metadata []byte

func init() {
	var err error
	metadata, err = json.Marshal(plugin.Metadata{
		Name:                      "cose",
		Description:               "A notation plugin for COSE signature",
		Version:                   version.GetVersion(),
		URL:                       "https://github.com/microsoft/notation-cose",
		SupportedContractVersions: []string{plugin.ContractVersion},
		Capabilities:              []plugin.Capability{plugin.CapabilityEnvelopeGenerator},
	})
	if err != nil {
		panic(err)
	}
}

func runGetMetaData(ctx *cli.Context) error {
	os.Stdout.Write(metadata)
	return nil
}
