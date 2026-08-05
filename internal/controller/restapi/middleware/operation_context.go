package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

// OperationIDKey is the context key for storing the canonical OpenAPI operationId.
const OperationIDKey = "operationId"

// openAPIParamRegex matches OpenAPI path parameter templates like {param}.
var openAPIParamRegex = regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)

// ConvertOpenAPIPathToFiber converts OpenAPI 3.0 path syntax ({param})
// to Fiber path syntax (:param).
func ConvertOpenAPIPathToFiber(openAPIPath string) string {
	return openAPIParamRegex.ReplaceAllString(openAPIPath, `:$1`)
}

// OperationBridge manages the canonical OperationID mapping derived from the OpenAPI specification.
type OperationBridge struct {
	// mapping stores key: "METHOD /path/pattern" -> operationId
	mapping map[string]string
	// operationIDs tracks used operationIds to prevent duplicates
	operationIDs map[string]string
}

// NewOperationBridge parses an OpenAPI spec and builds an immutable OperationBridge.
// Returns an error if duplicate operationIDs are detected during startup.
func NewOperationBridge(doc *openapi3.T) (*OperationBridge, error) {
	if doc == nil {
		return nil, fmt.Errorf("openapi spec document cannot be nil")
	}

	bridge := &OperationBridge{
		mapping:      make(map[string]string),
		operationIDs: make(map[string]string),
	}

	var basePath string
	if len(doc.Servers) > 0 && doc.Servers[0] != nil {
		basePath = strings.TrimSuffix(doc.Servers[0].URL, "/")
	}

	if doc.Paths != nil {
		for openAPIPath, pathItem := range doc.Paths.Map() {
			if pathItem == nil {
				continue
			}

			cleanPath := openAPIPath
			if !strings.HasPrefix(cleanPath, "/") {
				cleanPath = "/" + cleanPath
			}

			cleanFiberPath := ConvertOpenAPIPathToFiber(cleanPath)

			var fullFiberPath string
			if basePath != "" && !strings.HasPrefix(cleanPath, basePath) {
				fullFiberPath = ConvertOpenAPIPathToFiber(basePath + cleanPath)
			} else {
				fullFiberPath = cleanFiberPath
			}

			operations := map[string]*openapi3.Operation{
				http.MethodGet:    pathItem.Get,
				http.MethodPost:   pathItem.Post,
				http.MethodPut:    pathItem.Put,
				http.MethodPatch:  pathItem.Patch,
				http.MethodDelete: pathItem.Delete,
			}

			for method, op := range operations {
				if op == nil || op.OperationID == "" {
					continue
				}

				opID := op.OperationID

				// Detect duplicate operationId across spec
				if existingPath, exists := bridge.operationIDs[opID]; exists {
					return nil, fmt.Errorf("duplicate operationId detected at startup: '%s' defined at both '%s' and '%s %s'",
						opID, existingPath, method, openAPIPath)
				}

				bridge.operationIDs[opID] = fmt.Sprintf("%s %s", method, openAPIPath)

				// Index both full route pattern (with base path) and clean route pattern
				bridge.mapping[fmt.Sprintf("%s %s", method, fullFiberPath)] = opID
				bridge.mapping[fmt.Sprintf("%s %s", method, cleanFiberPath)] = opID
			}
		}
	}

	return bridge, nil
}

// Lookup derives the operationId for a given HTTP method and Fiber route pattern.
func (b *OperationBridge) Lookup(method, fiberRoutePattern string) (string, bool) {
	if b == nil || b.mapping == nil {
		return "", false
	}
	mapKey := fmt.Sprintf("%s %s", strings.ToUpper(method), fiberRoutePattern)
	opID, exists := b.mapping[mapKey]
	return opID, exists
}

// FiberMiddleware creates a Fiber handler that injects operationId into ctx.Locals.
// If a request hits a generated route group but operationId cannot be derived, it fails closed (500).
func (b *OperationBridge) FiberMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Ignore operational and documentation endpoints
		path := ctx.Path()
		if path == "/healthz" || path == "/metrics" || strings.HasPrefix(path, "/swagger") {
			return ctx.Next()
		}

		route := ctx.Route()
		if route == nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "missing operation id",
			})
		}

		opID, found := b.Lookup(ctx.Method(), route.Path)
		if !found {
			// Fail closed if operationId is missing for a managed API route
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "missing operation id",
			})
		}

		ctx.Locals(OperationIDKey, opID)
		return ctx.Next()
	}
}

// GetOperationID retrieves the operationId injected into Fiber context.
func GetOperationID(ctx *fiber.Ctx) string {
	if ctx == nil {
		return ""
	}
	val, ok := ctx.Locals(OperationIDKey).(string)
	if !ok {
		return ""
	}
	return val
}
