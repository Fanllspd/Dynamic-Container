package k3s

import (
	"context"
	"fmt"
	"k3s-client/global"
	"k3s-client/utils"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func InitKubernetesClient() {
	// 解码 CA 证书
	caCert, err := utils.LoadBase64EncodedFile(global.App.Config.Kubernetes.CAFile)
	if err != nil {
		log.Fatalf("Error loading CA certificate: %v", err)
	}

	// 解码客户端证书
	clientCert, err := utils.LoadBase64EncodedFile(global.App.Config.Kubernetes.CertFile)
	if err != nil {
		log.Fatalf("Error loading client certificate: %v", err)
	}

	// 解码客户端密钥
	clientKey, err := utils.LoadBase64EncodedFile(global.App.Config.Kubernetes.KeyFile)
	if err != nil {
		log.Fatalf("Error loading client key: %v", err)
	}

	// 创建自定义的 rest.Config
	kubeConfig := &rest.Config{
		Host:    global.App.Config.Kubernetes.ApiServer,
		Timeout: time.Duration(global.App.Config.Kubernetes.Timeout) * time.Second,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   caCert,
			CertData: clientCert,
			KeyData:  clientKey,
		},
	}

	// 创建客户端集
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	// 列出所有的 Pod
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}

	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Printf("Pod Name: %s, Namespace: %s\n", pod.Name, pod.Namespace)
	}
}
