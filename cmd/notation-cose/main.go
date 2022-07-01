package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/microsoft/notation-cose/internal/version"
	"github.com/notaryproject/notation-go/plugin"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "notation-cose",
		Usage:   "COSE Plugin for Notation",
		Version: version.GetVersion(),
		Commands: []*cli.Command{
			signCommand,
			verifyCommand,
			metadatCommand,
		},
	}
	if err := app.Run(os.Args); err != nil {
		var reer plugin.RequestError
		if !errors.As(err, &reer) {
			err = plugin.RequestError{
				Code: plugin.ErrorCodeGeneric,
				Err:  err,
			}
		}
		data, _ := json.Marshal(err)
		os.Stderr.Write(data)
		os.Exit(1)
	}
}
