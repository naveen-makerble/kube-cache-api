package main

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Pod struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
}

type UserCache struct {
	sync.RWMutex
	Data      []Pod
	Timestamp int64 // Unix timestamp (seconds)
}

var cache = struct {
	sync.RWMutex
	Users map[string]*UserCache
}{
	Users: make(map[string]*UserCache),
}

const cacheTTL = 60 // seconds

func NewKubeClient() (*kubernetes.Clientset, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fallback to local kubeconfig
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return kubernetes.NewForConfig(config)
}

func GetOrCreateUserCache(userID string, clientset *kubernetes.Clientset) *UserCache {
	cache.Lock()
	defer cache.Unlock()

	uc, exists := cache.Users[userID]
	now := time.Now().Unix()
	if exists && now-uc.Timestamp < cacheTTL {
		return uc
	}

	// Fetch from Kubernetes (all namespaces for demo; you can scope by user/namespace)
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	podList := make([]Pod, 0, len(pods.Items))
	if err == nil {
		for _, p := range pods.Items {
			podList = append(podList, Pod{
				Name:      p.Name,
				Namespace: p.Namespace,
				Status:    string(p.Status.Phase),
			})
		}
	}

	if !exists {
		uc = &UserCache{}
		cache.Users[userID] = uc
	}
	uc.Data = podList
	uc.Timestamp = now
	return uc
}
