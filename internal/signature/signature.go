package signature

import (
	"context"
	"crypto"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/microsoft/notation-cose/pkg/cose"
	"github.com/notaryproject/notation-go"
	"github.com/notaryproject/notation-go/crypto/cryptoutil"
	"github.com/notaryproject/notation-go/crypto/timestamp"
	"github.com/notaryproject/notation-go/plugin"
)

func Sign(ctx context.Context, req *plugin.GenerateEnvelopeRequest) (*plugin.GenerateEnvelopeResponse, error) {
	signer, opts, err := getSignerWithOptions(req.KeyID)
	if err != nil {
		return nil, plugin.RequestError{
			Code: plugin.ErrorCodeValidation,
			Err:  fmt.Errorf("invalid request input: %w", err),
		}
	}
	var sig []byte
	sig, err = signer.Sign(ctx, req.Payload, opts)
	if err != nil {
		return nil, plugin.RequestError{
			Code: plugin.ErrorCodeGeneric,
			Err:  err,
		}
	}
	return &plugin.GenerateEnvelopeResponse{
		SignatureEnvelope:     sig,
		SignatureEnvelopeType: req.SignatureEnvelopeType,
		Annotations:           nil,
	}, nil
}

func getSignerWithOptions(keyInfo string) (*cose.Signer, notation.SignOptions, error) {
	// parse options
	var opts notation.SignOptions
	items := strings.SplitN(keyInfo, ":", 3)
	if len(items) < 2 {
		return nil, opts, errors.New("missing signing key pair")
	}
	keyPath := items[0]
	certPath := items[1]
	var tsEndpoint string
	if len(items) > 2 {
		tsEndpoint = items[2]
	}

	// read key / cert pair
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, opts, err
	}
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, opts, err
	}
	keyPair, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, opts, err
	}

	// parse cert
	certs, err := cryptoutil.ParseCertificatePEM(certPEM)
	if err != nil {
		return nil, opts, err
	}

	// construct signer
	privateKey, ok := keyPair.PrivateKey.(crypto.Signer)
	if !ok {
		return nil, opts, errors.New("unsupported private key")
	}
	signer, err := cose.NewSigner(privateKey, certs)
	if err != nil {
		return nil, opts, err
	}

	// hack: refine options
	// notation#feat-kv-extensibility uses an older version of notation-go-lib,
	// which does not support TSA in options.
	if tsEndpoint != "" {
		tsa, err := timestamp.NewHTTPTimestamper(nil, tsEndpoint)
		if err != nil {
			return nil, opts, err
		}
		opts.TSA = tsa
	}
	return signer, opts, nil
}
