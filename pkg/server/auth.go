package server

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

type KeyAuth func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error)

func SimpleKeyConfig(authorizedKeysMap map[string]bool) KeyAuth {
	return func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
		if authorizedKeysMap[string(pubKey.Marshal())] {
			log.Printf("User \"%s\" authenticated with PubKey.", c.User())
			return nil, nil
		}
		return nil, fmt.Errorf("unknown public key for %q", c.User())

	}
}
