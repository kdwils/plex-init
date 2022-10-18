package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/kyledwilson/plex-init/kubernetes"
	"github.com/kyledwilson/plex-init/plex"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	kubePath        string
	secretName      string
	secretNamespace string
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "claim your plex server token and create a kubernetes secret",
	Long:  `claim your plex server token and create a kubernetes secret`,
	Run: func(cmd *cobra.Command, args []string) {
		c := plex.NewClient(http.DefaultClient)
		ctx := context.Background()
		plextoken := viper.GetString("PLEX_TOKEN")
		claimToken, err := c.GetServerClaimToken(ctx, plextoken)
		if err != nil {
			log.Fatal(err)
		}

		kc, err := kubernetes.NewClient(kubePath, secretNamespace)
		if err != nil {
			log.Fatal(err)
		}

		secret, err := kc.NewPlexSecret(ctx, secretName, secretNamespace, claimToken)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("secret %s created in %s namespace", secret.Name, secret.Namespace)
	},
}

func init() {
	claimCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().StringVar(&kubePath, "kube-config", "", "path to kubeconfig")
	tokenCmd.Flags().StringVar(&secretName, "secretName", "plex", "path to kubeconfig")
	tokenCmd.Flags().StringVar(&secretNamespace, "namespace", "default", "namespace to create the plex token secret in")
}
