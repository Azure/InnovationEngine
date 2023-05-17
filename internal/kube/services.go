package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetAgentService(id string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "runner - " + id,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "runner",
				"id":  id,
			},
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     8080,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
}

func CreateAgentService(clientset *kubernetes.Clientset, service *corev1.Service) (*corev1.Service, error) {
	return clientset.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
}
