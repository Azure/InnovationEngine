package main

import (
	"context"
	"fmt"
	"net/http"
	"path"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Azure/InnovationEngine/internal/kube"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var PREFIX_BASE = "/api/ie"

func main() {
	fmt.Println("Hello, world!")
	server := echo.New()

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.GET(path.Join(PREFIX_BASE, "hello"), func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, world!"})
	})

	server.POST(PREFIX_BASE, func(c echo.Context) error {
		clientset, err := kube.GetKubernetesClient()

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
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
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, world!", "job": job.Name})
	})

	server.Logger.Fatal(server.Start("0.0.0.0:8080"))
}
