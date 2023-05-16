package main

import (
	"context"
	"fmt"
	"net/http"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Azure/InnovationEngine/internal/kube"
)

func main() {
	fmt.Println("Hello, world!")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})

	http.HandleFunc("/api/scenario", func(w http.ResponseWriter, r *http.Request) {
		clientset, err := kube.GetKubernetesClient()
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		}

		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "innovation-engine-",
			},
			Spec: batchv1.JobSpec{
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						RestartPolicy: v1.RestartPolicyNever,
						Containers: []v1.Container{
							{
								Name:  "runner",
								Image: "innovation-engine-runner",
								Command: []string{
									"runner",
								},
							},
						},
					},
				},
			},
		}

		job, err = clientset.BatchV1().Jobs("default").Create(context.TODO(), job, metav1.CreateOptions{})

		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		}

		fmt.Println(job)

	})

	http.ListenAndServe(":8080", nil)
}
