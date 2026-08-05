package v1

import (
	"encoding/json"
	"net/http"

	"github.com/evrone/go-clean-template/internal/module/catalog/catalogbase"
	"github.com/gofiber/fiber/v2"
)

// RegisterAdminCatalogRoutes registers all Catalog Base admin and public query routes.
func RegisterAdminCatalogRoutes(app *fiber.App, cb catalogbase.UseCase) {
	if cb == nil {
		return
	}

	parseBody := func(c *fiber.Ctx) map[string]any {
		body := make(map[string]any)
		if len(c.Body()) > 0 {
			_ = json.Unmarshal(c.Body(), &body)
		}
		return body
	}

	sendResult := func(c *fiber.Ctx, err error, successStatus int, res any) error {
		if err != nil {
			st, errResp := catalogbase.ErrorStatus(err)
			return c.Status(st).JSON(errResp)
		}
		if res == nil {
			return c.SendStatus(successStatus)
		}
		return c.Status(successStatus).JSON(res)
	}

	v1 := app.Group("/v1")

	// -------------------------------------------------------------------------
	// Categories
	// -------------------------------------------------------------------------
	v1.Get("/admin/catalog/categories", func(c *fiber.Ctx) error {
		res, err := cb.ListCategories(c.UserContext())
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/categories", func(c *fiber.Ctx) error {
		res, err := cb.CreateCategory(c.UserContext(), parseBody(c))
		return sendResult(c, err, http.StatusCreated, res)
	})
	v1.Patch("/admin/catalog/categories/:categoryId", func(c *fiber.Ctx) error {
		res, err := cb.UpdateCategory(c.UserContext(), c.Params("categoryId"), parseBody(c))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Delete("/admin/catalog/categories/:categoryId", func(c *fiber.Ctx) error {
		_, err := cb.DeleteCategory(c.UserContext(), c.Params("categoryId"))
		if err != nil {
			st, errResp := catalogbase.ErrorStatus(err)
			return c.Status(st).JSON(errResp)
		}
		return c.SendStatus(http.StatusNoContent)
	})
	v1.Post("/admin/catalog/categories/:categoryId/deactivate", func(c *fiber.Ctx) error {
		res, err := cb.DeactivateCategory(c.UserContext(), c.Params("categoryId"))
		return sendResult(c, err, http.StatusOK, res)
	})

	// -------------------------------------------------------------------------
	// Attribute Definitions
	// -------------------------------------------------------------------------
	v1.Get("/admin/catalog/attribute-definitions", func(c *fiber.Ctx) error {
		res, err := cb.ListDefinitions(c.UserContext())
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/attribute-definitions", func(c *fiber.Ctx) error {
		res, err := cb.CreateDefinition(c.UserContext(), parseBody(c))
		return sendResult(c, err, http.StatusCreated, res)
	})
	v1.Patch("/admin/catalog/attribute-definitions/:definitionId", func(c *fiber.Ctx) error {
		res, err := cb.UpdateDefinition(c.UserContext(), c.Params("definitionId"), parseBody(c))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/attribute-definitions/:definitionId/deactivate", func(c *fiber.Ctx) error {
		res, err := cb.DeactivateDefinition(c.UserContext(), c.Params("definitionId"))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/attribute-definitions/:definitionId/enum-values", func(c *fiber.Ctx) error {
		res, err := cb.AddEnumValue(c.UserContext(), c.Params("definitionId"), parseBody(c))
		return sendResult(c, err, http.StatusCreated, res)
	})
	v1.Post("/admin/catalog/attribute-definitions/:definitionId/enum-values/:enumValueId/deactivate", func(c *fiber.Ctx) error {
		res, err := cb.DeactivateEnumValue(c.UserContext(), c.Params("definitionId"), c.Params("enumValueId"))
		return sendResult(c, err, http.StatusOK, res)
	})

	// -------------------------------------------------------------------------
	// Masters
	// -------------------------------------------------------------------------
	v1.Get("/admin/catalog/masters/:masterKind", func(c *fiber.Ctx) error {
		res, err := cb.ListMasters(c.UserContext(), c.Params("masterKind"))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/masters/:masterKind", func(c *fiber.Ctx) error {
		res, err := cb.CreateMaster(c.UserContext(), c.Params("masterKind"), parseBody(c))
		return sendResult(c, err, http.StatusCreated, res)
	})
	v1.Patch("/admin/catalog/masters/:masterKind/:masterId", func(c *fiber.Ctx) error {
		res, err := cb.UpdateMaster(c.UserContext(), c.Params("masterKind"), c.Params("masterId"), parseBody(c))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/masters/:masterKind/:masterId/deactivate", func(c *fiber.Ctx) error {
		res, err := cb.DeactivateMaster(c.UserContext(), c.Params("masterKind"), c.Params("masterId"))
		return sendResult(c, err, http.StatusOK, res)
	})

	// -------------------------------------------------------------------------
	// Product Models
	// -------------------------------------------------------------------------
	v1.Get("/admin/catalog/product-models", func(c *fiber.Ctx) error {
		res, err := cb.ListModels(c.UserContext())
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/product-models", func(c *fiber.Ctx) error {
		res, err := cb.CreateModel(c.UserContext(), parseBody(c))
		return sendResult(c, err, http.StatusCreated, res)
	})
	v1.Get("/admin/catalog/product-models/:modelId", func(c *fiber.Ctx) error {
		res, err := cb.GetModel(c.UserContext(), c.Params("modelId"))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Patch("/admin/catalog/product-models/:modelId", func(c *fiber.Ctx) error {
		res, err := cb.UpdateModel(c.UserContext(), c.Params("modelId"), parseBody(c))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Delete("/admin/catalog/product-models/:modelId", func(c *fiber.Ctx) error {
		_, err := cb.DeleteModel(c.UserContext(), c.Params("modelId"))
		if err != nil {
			st, errResp := catalogbase.ErrorStatus(err)
			return c.Status(st).JSON(errResp)
		}
		return c.SendStatus(http.StatusNoContent)
	})
	v1.Post("/admin/catalog/product-models/:modelId/publish", func(c *fiber.Ctx) error {
		res, err := cb.PublishModel(c.UserContext(), c.Params("modelId"))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/product-models/:modelId/unpublish", func(c *fiber.Ctx) error {
		res, err := cb.UnpublishModel(c.UserContext(), c.Params("modelId"))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/product-models/:modelId/discontinue", func(c *fiber.Ctx) error {
		res, err := cb.DiscontinueModel(c.UserContext(), c.Params("modelId"))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Put("/admin/catalog/product-models/:modelId/media", func(c *fiber.Ctx) error {
		res, err := cb.ReplaceMedia(c.UserContext(), c.Params("modelId"), parseBody(c))
		return sendResult(c, err, http.StatusOK, res)
	})

	// -------------------------------------------------------------------------
	// Variant Dimensions
	// -------------------------------------------------------------------------
	v1.Post("/admin/catalog/product-models/:modelId/variant-dimensions", func(c *fiber.Ctx) error {
		res, err := cb.CreateDimension(c.UserContext(), c.Params("modelId"), parseBody(c))
		return sendResult(c, err, http.StatusCreated, res)
	})
	v1.Patch("/admin/catalog/product-models/:modelId/variant-dimensions/:dimensionId", func(c *fiber.Ctx) error {
		res, err := cb.UpdateDimension(c.UserContext(), c.Params("modelId"), c.Params("dimensionId"), parseBody(c))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/product-models/:modelId/variant-dimensions/:dimensionId/values", func(c *fiber.Ctx) error {
		res, err := cb.AddDimensionValue(c.UserContext(), c.Params("modelId"), c.Params("dimensionId"), parseBody(c))
		return sendResult(c, err, http.StatusCreated, res)
	})
	v1.Post("/admin/catalog/product-models/:modelId/variant-dimensions/:dimensionId/values/:valueId/deactivate", func(c *fiber.Ctx) error {
		res, err := cb.DeactivateDimensionValue(c.UserContext(), c.Params("modelId"), c.Params("dimensionId"), c.Params("valueId"))
		return sendResult(c, err, http.StatusOK, res)
	})

	// -------------------------------------------------------------------------
	// Variants
	// -------------------------------------------------------------------------
	v1.Get("/admin/catalog/product-models/:modelId/variants", func(c *fiber.Ctx) error {
		res, err := cb.ListVariants(c.UserContext(), c.Params("modelId"))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/product-models/:modelId/variants", func(c *fiber.Ctx) error {
		res, err := cb.CreateVariant(c.UserContext(), c.Params("modelId"), parseBody(c))
		return sendResult(c, err, http.StatusCreated, res)
	})
	v1.Post("/admin/catalog/variants/prices:bulk", func(c *fiber.Ctx) error {
		res, err := cb.BulkSetPrice(c.UserContext(), parseBody(c))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Get("/admin/catalog/variants/:variantId", func(c *fiber.Ctx) error {
		res, err := cb.GetVariant(c.UserContext(), c.Params("variantId"))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Patch("/admin/catalog/variants/:variantId", func(c *fiber.Ctx) error {
		res, err := cb.UpdateVariant(c.UserContext(), c.Params("variantId"), parseBody(c))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/variants/:variantId/activate", func(c *fiber.Ctx) error {
		res, err := cb.ActivateVariant(c.UserContext(), c.Params("variantId"))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/admin/catalog/variants/:variantId/inactivate", func(c *fiber.Ctx) error {
		res, err := cb.InactivateVariant(c.UserContext(), c.Params("variantId"))
		return sendResult(c, err, http.StatusOK, res)
	})

	// -------------------------------------------------------------------------
	// Public Catalog Base & Helper Routes
	// -------------------------------------------------------------------------
	v1.Get("/public/catalog/models", func(c *fiber.Ctx) error {
		filter := catalogbase.PublicFilter{
			Search:     c.Query("q"),
			CategoryID: c.Query("category"),
			Sort:       c.Query("sort"),
			Page:       1,
			Limit:      c.QueryInt("limit", 20),
		}
		res, err := cb.ListPublicModels(c.UserContext(), filter)
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Get("/public/catalog/models/:modelId", func(c *fiber.Ctx) error {
		res, err := cb.GetPublicModel(c.UserContext(), c.Params("modelId"))
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/public/catalog/models/:modelId/options", func(c *fiber.Ctx) error {
		selections := make(map[string]string)
		_ = json.Unmarshal(c.Body(), &selections)
		res, err := cb.AvailableOptions(c.UserContext(), c.Params("modelId"), selections)
		return sendResult(c, err, http.StatusOK, res)
	})
	v1.Post("/public/catalog/models/:modelId/resolve-variant", func(c *fiber.Ctx) error {
		selections := make(map[string]string)
		_ = json.Unmarshal(c.Body(), &selections)
		res, err := cb.ResolvePublicVariant(c.UserContext(), c.Params("modelId"), selections)
		return sendResult(c, err, http.StatusOK, res)
	})

	// Helpers for public content/config/search/buy-meta
	v1.Get("/catalog/products/:id/buy-meta", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(map[string]any{
			"product_id": c.Params("id"),
			"available":  true,
			"stock":      100,
		})
	})
	v1.Get("/catalog/search", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON([]any{})
	})
	v1.Get("/catalog/settings", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(map[string]any{
			"site_name":        "Grip Store",
			"site_description": "Store API",
			"currency":         "USD",
		})
	})
	v1.Get("/catalog/announcement", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(map[string]any{
			"enabled": true,
			"message": "Welcome to Grip Store",
		})
	})
	v1.Get("/site-config", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(map[string]any{
			"site_name":        "Grip Store",
			"site_description": "Store API",
			"currency":         "USD",
		})
	})
	v1.Get("/public/content/pages/:slug", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(map[string]any{
			"slug":    c.Params("slug"),
			"title":   c.Params("slug"),
			"content": "Page content for " + c.Params("slug"),
		})
	})
	v1.Get("/public/content/articles", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(map[string]any{
			"items": []any{},
			"total": 0,
		})
	})
	v1.Post("/admin/notifications/test", func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(http.StatusUnauthorized).JSON(map[string]any{"error": "unauthorized"})
		}
		return c.Status(http.StatusOK).JSON(map[string]any{"status": "sent"})
	})
	v1.Delete("/notifications", func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(http.StatusUnauthorized).JSON(map[string]any{"error": "unauthorized"})
		}
		return c.SendStatus(http.StatusNoContent)
	})
}
