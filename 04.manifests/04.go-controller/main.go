package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config := getConfig()

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	log.Println("Starting pod watcher...")

	// reconnect loop (important)
	for {
		err := watchPods(clientset)
		log.Printf("Watch ended: %v. Reconnecting in 3s...", err)
		time.Sleep(3 * time.Second)
	}
}

func getConfig() *rest.Config {
	// Try in-cluster first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config
	}

	// Fallback to local kubeconfig
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Cannot load kubeconfig: %v", err)
	}

	return config
}

func watchPods(clientset *kubernetes.Clientset) error {
	watcher, err := clientset.CoreV1().Pods(metav1.NamespaceAll).Watch(context.Background(), metav1.ListOptions{LabelSelector: "project=ng-voice"})
	if err != nil {
		return err
	}

	for event := range watcher.ResultChan() {

		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			continue // avoid panic on tombstones
		}

		switch event.Type {

		case watch.Added:
			log.Printf("[CREATED] %s/%s | Phase=%s",
				pod.Namespace, pod.Name, pod.Status.Phase)

		case watch.Modified:
			log.Printf("[UPDATED] %s/%s | Phase=%s Ready=%v",
				pod.Namespace, pod.Name,
				pod.Status.Phase, isPodReady(pod))

		case watch.Deleted:
			log.Printf("[DELETED] %s/%s",
				pod.Namespace, pod.Name)
		}
	}

	return nil // channel closed → reconnect
}

func isPodReady(pod *corev1.Pod) bool {
	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
