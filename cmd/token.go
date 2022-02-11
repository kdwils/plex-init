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

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "claim your plex server token and create a kubernetes secret",
	Long:  `claim your plex server token and create a kubernetes secret`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		host := viper.GetString("PLEX_HOST")
		scheme := viper.GetString("PLEX_SCHEME")

		plextoken := viper.GetString("PLEX_TOKEN")
		c := plex.NewClient(scheme, host, http.DefaultClient)
		claimToken, err := c.GetServerClaimToken(ctx, plextoken)
		if err != nil {
			log.Fatal(err)
		}

		kubePath := cmd.Flag("kube-config").Value.String()
		secretName := cmd.Flag("secret-name").Value.String()
		secretNamespace := cmd.Flag("namespace").Value.String()

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
	tokenCmd.Flags().String("secret-name", "plex", "name of the secret to create")
	tokenCmd.Flags().String("namespace", "default", "namespace to create the secret in")
	tokenCmd.Flags().String("kube-config", "", "path to kubeconfig")
	viper.SetDefault("PLEX_SCHEME", "https")
	viper.SetDefault("PLEX_HOST", "plex.tv")
}
