package docs

import "github.com/gofiber/fiber/v2"

// swaggerJSON is a minimal OpenAPI 3.0 document that describes the API.
// You can regenerate or extend it as your API evolves.
const swaggerJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "URL Shortener API",
    "version": "1.0.0",
    "description": "Simple URL Shortener service with rate limiting and redirects."
  },
  "servers": [
    { "url": "/" }
  ],
  "paths": {
    "/v1/healthz": {
      "get": {
        "summary": "Liveness probe",
        "responses": {
          "200": {
            "description": "OK",
            "content": { "text/plain": { "schema": { "type": "string", "example": "ok" } } }
          }
        }
      }
    },
    "/v1/readyz": {
      "get": {
        "summary": "Readiness probe",
        "responses": {
          "200": {
            "description": "Ready",
            "content": { "text/plain": { "schema": { "type": "string", "example": "ready" } } }
          }
        }
      }
    },
    "/v1/shorten": {
      "post": {
        "summary": "Shorten a URL",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/ShortenRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Shortened successfully",
            "content": { "application/json": { "schema": { "$ref": "#/components/schemas/ShortenResponse" } } }
          },
          "400": { "description": "invalid_url or bad_request" },
          "409": { "description": "conflict (custom alias already exists)" },
          "500": { "description": "internal" }
        }
      }
    },
    "/{code}": {
      "get": {
        "summary": "Resolve and redirect by code",
        "parameters": [
          {
            "name": "code",
            "in": "path",
            "required": true,
            "schema": { "type": "string" },
            "description": "Short code generated for the original URL"
          }
        ],
        "responses": {
          "302": { "description": "Found, redirects to original URL" },
          "404": { "description": "Not found" },
          "500": { "description": "Internal error (settings retrieval failure)" }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "ShortenRequest": {
        "type": "object",
        "required": ["url"],
        "properties": {
          "url": { "type": "string", "format": "uri", "example": "https://example.com/long" },
          "custom_alias": { "type": "string", "nullable": true, "example": "my-custom" }
        }
      },
      "ShortenResponse": {
        "type": "object",
        "properties": {
          "code": { "type": "string", "example": "abc123" },
          "short_url": { "type": "string", "format": "uri", "example": "http://localhost:3000/abc123" }
        }
      }
    }
  }
}`

// uiHTML renders Swagger UI via CDN, loading the spec from /v1/docs/swagger.json
const uiHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Swagger UI - URL Shortener</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '/v1/docs/swagger.json',
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis],
    });
  };
  </script>
</body>
</html>`

// SwaggerJSON serves the OpenAPI specification in JSON format.
func SwaggerJSON(c *fiber.Ctx) error {
	c.Type("json")
	return c.SendString(swaggerJSON)
}

// SwaggerUI serves a minimal Swagger UI that loads the local swagger.json.
func SwaggerUI(c *fiber.Ctx) error {
	c.Type("html")
	return c.SendString(uiHTML)
}
