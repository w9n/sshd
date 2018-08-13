package cmd

import (
	"io/ioutil"
	"log"

	"github.com/nseps/sshd/pkg/server"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts a ssh server",
	Run: func(cmd *cobra.Command, args []string) {
		config := &ssh.ServerConfig{}

		certs := certificates{}
		authorizedKeysPath := cmd.Flag("auth-file").Value.String()
		if authorizedKeysPath != "" {
			certs.parseAuthFile(authorizedKeysPath)
		}
		config.PublicKeyCallback = server.SimpleKeyConfig(certs)

		address := *cmd.PersistentFlags().String("address", "0.0.0.0:9999", "address to listen on")
		err := server.SSHdListenAndServe(address, config)
		log.Fatalf("error %s", err)

	},
}

type certificates map[string]bool

func (c certificates) parseAuthFile(path string) error {
	authorizedKeysBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to load authorized_keys, err: %v", err)
	}

	authorizedKeysMap := map[string]bool{}
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			return err
		}
		authorizedKeysMap[string(pubKey.Marshal())] = true
		authorizedKeysBytes = rest

	}
	return nil
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().String("auth-file", "", "path to authorized keys")

}
