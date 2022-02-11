package kubernetes

import (
	"context"
	"os"

	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type SecretClient struct {
	Client coreV1Types.SecretInterface
}

func NewClient(kubePath, namespace string) (SecretClient, error) {
	var sc SecretClient

	config, err := getConfig(kubePath)
	if err != nil {
		return sc, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return sc, err
	}

	sc.Client = clientset.CoreV1().Secrets(namespace)
	return sc, nil
}

func (s SecretClient) NewPlexSecret(ctx context.Context, secretName, namespace, claimToken string) (*coreV1.Secret, error) {
	secret := &coreV1.Secret{
		Type: coreV1.SecretTypeOpaque,
		ObjectMeta: v1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"token": []byte(claimToken),
		},
	}

	_, err := s.Client.Get(ctx, secretName, v1.GetOptions{})
	if err != nil {
		return s.Client.Create(ctx, secret, v1.CreateOptions{})
	}

	return s.Client.Update(ctx, secret, v1.UpdateOptions{})

}

func getConfig(kubePath string) (*rest.Config, error) {
	if _, err := os.Stat("/.dockerenv"); err == nil && kubePath == "" {
		return rest.InClusterConfig()
	}

	return clientcmd.BuildConfigFromFlags("", kubePath)
}
