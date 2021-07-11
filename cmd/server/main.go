package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/AverageMarcus/cluster-api-provider-kind/api/v1alpha4"
	"github.com/AverageMarcus/cluster-api-provider-kind/pkg/kind"
)

var port = "3000"

// Start starts the Kind API server
func Start() error {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	logger := zap.New()
	kind := kind.New(logger)

	app.Post("/", func(c *fiber.Ctx) error {
		kindCluster := v1alpha4.KindCluster{}
		if err := c.BodyParser(&kindCluster); err != nil {
			logger.Error(err, "failed to parse incoming request")
			return err
		}

		if err := kind.CreateCluster(&kindCluster); err != nil {
			logger.Error(err, "failed to create Kind cluster")
			return err
		}

		return nil
	})

	app.Get("/:clusterName", func(c *fiber.Ctx) error {
		isReady, err := kind.IsReady(c.Params("clusterName"))
		if err != nil {
			logger.Error(err, "failed to check request status")
			return err
		}

		return c.JSON(isReady)
	})

	app.Get("/:clusterName/kubeconfig", func(c *fiber.Ctx) error {
		kubeconfig, err := kind.GetKubeConfig(c.Params("clusterName"))
		if err != nil {
			logger.Error(err, "failed to get kubeconfig")
			return err
		}
		return c.JSON(kubeconfig)
	})

	app.Delete("/:clusterName", func(c *fiber.Ctx) error {
		if err := kind.DeleteCluster(c.Params("clusterName")); err != nil {
			logger.Error(err, "failed to delete cluster")
			return err
		}

		return nil
	})

	return app.Listen(fmt.Sprintf(":%s", port))
}
