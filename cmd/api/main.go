package main

import (
	"net/http"
	"path"

	"github.com/Azure/InnovationEngine/internal/kube"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	BASE_ROUTE        = "/api"
	HEALTH_ROUTE      = path.Join(BASE_ROUTE, "health")
	EXECUTION_ROUTE   = path.Join(BASE_ROUTE, "execute")
	DEPLOYMENTS_ROUTE = path.Join(BASE_ROUTE, "deployments")
)

func main() {
	server := echo.New()

	// Setup middleware.
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.GET(HEALTH_ROUTE, func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "OK"})
	})

	server.POST(EXECUTION_ROUTE, func(c echo.Context) error {
		clientset, err := kube.GetKubernetesClient()

		id := uuid.New().String()

		// Create deployment
		deployment := kube.GetAgentDeployment(id)
		_, err = kube.CreateAgentDeployment(clientset, deployment)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		// Create service
		service := kube.GetAgentService(id)
		_, err = kube.CreateAgentService(clientset, service)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"deployment": deployment.Name, "service": service.Name})
	})

	server.Logger.Fatal(server.Start(":8080"))
}
