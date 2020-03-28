package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/google/go-github/v30/github"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

func checkversions() (SupportedVersions, error) {
	client := github.NewClient(nil)

	taglist, _, err := client.Repositories.ListTags(context.Background(), "poseidon", "typhoon", nil)
	if err != nil {
		return SupportedVersions{}, err
	}
	var supported SupportedVersions
	for _, tag := range taglist {
		supported.Versions = append(supported.Versions, tag.GetName())
	}
	if err != nil {
		return SupportedVersions{}, err
	}
	return supported, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func cfgencrypt(key string, fileToEnc string) error {
	// Read in public key
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}

	recipient, err := readEntity(decoded)
	if err != nil {
		return err
	}

	f, err := os.Open(fileToEnc)
	if err != nil {
		return err
	}
	defer f.Close()

	dst, err := os.Create(fileToEnc + ".gpg")
	if err != nil {
		return err
	}
	defer dst.Close()
	err = encrypt([]*openpgp.Entity{recipient}, nil, f, dst)
	if err != nil {
		return err
	}
	err = os.Remove(fileToEnc)
	if err != nil {
		return err
	}
	return nil
}
func encrypt(recip []*openpgp.Entity, signer *openpgp.Entity, r io.Reader, w io.Writer) error {
	wc, err := openpgp.Encrypt(w, recip, signer, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return err
	}
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}
	return wc.Close()
}

func readEntity(key []byte) (*openpgp.Entity, error) {
	r := bytes.NewReader(key)

	block, err := armor.Decode(r)
	if err != nil {
		return nil, err
	}
	return openpgp.ReadEntity(packet.NewReader(block.Body))
}

// serviceAccount shows how to use a service account to authenticate.
func serviceAccount() error {
	// Download service account key per https://cloud.google.com/docs/authentication/production.
	// Set environment variable GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
	// This environment variable will be automatically picked up by the client.
	client, err := pubsub.NewClient(context.Background(), "your-project-id")
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	// Use the authenticated client.
	_ = client

	return nil
}
