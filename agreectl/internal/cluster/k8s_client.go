package cluster

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/applyconfigurations/core/v1"
	metav1config "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Secret struct {
	Bytes map[string][]byte
}

func (s *Secret) User() string {
	if s.Bytes == nil {
		return ""
	}
	return string(s.Bytes["user"])
}

func (s *Secret) Password() string {
	if s.Bytes == nil {
		return ""
	}
	return string(s.Bytes["password"])
}

func (s *Secret) DBName() string {
	if s.Bytes == nil {
		return ""
	}
	return string(s.Bytes["dbname"])
}

func (s *Secret) StringData() map[string]string {
	if s.Bytes == nil {
		return nil
	}
	result := make(map[string]string, len(s.Bytes))
	for k, v := range s.Bytes {
		result[k] = string(v)
	}
	return result
}

func (s *Secret) Data() map[string]string {
	return s.StringData()
}

type K8sClient interface {
	GetSecret(namespace, name string) (*Secret, error)
	UpsertSecret(namespace, name string, data map[string]string) error
	NodeIP() (string, error)
}

type realK8sClient struct {
	clientset *kubernetes.Clientset
	context   string
}

func NewK8sClient(context string) (K8sClient, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	if context != "" {
		configOverrides.CurrentContext = context
	}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}
	return &realK8sClient{
		clientset: clientset,
		context:   context,
	}, nil
}

func (c *realK8sClient) GetSecret(namespace, name string) (*Secret, error) {
	secret, err := c.clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &Secret{Bytes: secret.Data}, nil
}

func (c *realK8sClient) UpsertSecret(namespace, name string, data map[string]string) error {
	secretData := make(map[string][]byte, len(data))
	for k, v := range data {
		secretData[k] = []byte(v)
	}
	secret := "Secret"
	apiVersion := "v1"
	secretApplyConfig := v1.SecretApplyConfiguration{
		TypeMetaApplyConfiguration: metav1config.TypeMetaApplyConfiguration{
			Kind:       &secret,
			APIVersion: &apiVersion,
		},
		ObjectMetaApplyConfiguration: &metav1config.ObjectMetaApplyConfiguration{
			Name:      &name,
			Namespace: &namespace,
		},
		Data: secretData,
	}
	_, err := c.clientset.CoreV1().Secrets(namespace).Apply(context.Background(), &secretApplyConfig, metav1.ApplyOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *realK8sClient) NodeIP() (string, error) {
	nodes, err := c.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	for _, node := range nodes.Items {
		for _, cond := range node.Status.Conditions {
			if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
				for _, addr := range node.Status.Addresses {
					if addr.Type == corev1.NodeInternalIP {
						return addr.Address, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("no ready node with internal IP found")
}