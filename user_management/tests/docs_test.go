package tests

import (
	"testing"

	"github.com/IntouchOpec/user_management/docs"
	"github.com/stretchr/testify/assert"
	"github.com/swaggo/swag"
)

func TestSwaggerInfo_Structure(t *testing.T) {
	// Test SwaggerInfo is properly initialized
	assert.NotNil(t, docs.SwaggerInfo)

	// Test SwaggerInfo fields
	assert.Equal(t, "1.0", docs.SwaggerInfo.Version)
	assert.Equal(t, "localhost:8080", docs.SwaggerInfo.Host)
	assert.Equal(t, "/api/v1", docs.SwaggerInfo.BasePath)
	assert.Contains(t, docs.SwaggerInfo.Schemes, "http")
	assert.Contains(t, docs.SwaggerInfo.Schemes, "https")
	assert.Equal(t, "User Management API", docs.SwaggerInfo.Title)
	assert.Contains(t, docs.SwaggerInfo.Description, "RESTful API for user management")
	assert.Equal(t, "swagger", docs.SwaggerInfo.InfoInstanceName)
	assert.Equal(t, "{{", docs.SwaggerInfo.LeftDelim)
	assert.Equal(t, "}}", docs.SwaggerInfo.RightDelim)
}

func TestSwaggerInfo_Template(t *testing.T) {
	// Test that SwaggerTemplate is not empty
	assert.NotEmpty(t, docs.SwaggerInfo.SwaggerTemplate)

	// Test that template contains expected sections
	template := docs.SwaggerInfo.SwaggerTemplate
	assert.Contains(t, template, "swagger")
	assert.Contains(t, template, "info")
	assert.Contains(t, template, "paths")
	assert.Contains(t, template, "definitions")
}

func TestSwaggerInfo_Registration(t *testing.T) {
	// Test that swagger info is registered
	// This tests the init() function indirectly

	// Get the registered spec
	spec := swag.GetSwagger(docs.SwaggerInfo.InstanceName())
	assert.NotNil(t, spec)

	// Verify it's the same instance
	assert.Equal(t, docs.SwaggerInfo, spec)
}

func TestSwaggerInfo_APIEndpoints(t *testing.T) {
	// Test that the swagger template contains expected API endpoints
	template := docs.SwaggerInfo.SwaggerTemplate

	// Check for health endpoint
	assert.Contains(t, template, "/health")

	// Check for user endpoints
	assert.Contains(t, template, "/users")
	assert.Contains(t, template, "/users/{id}")

	// Check for HTTP methods
	assert.Contains(t, template, "get")
	assert.Contains(t, template, "post")
	assert.Contains(t, template, "put")
	assert.Contains(t, template, "delete")
}

func TestSwaggerInfo_Definitions(t *testing.T) {
	// Test that swagger template contains model definitions
	template := docs.SwaggerInfo.SwaggerTemplate

	// Check for UserRequest model
	assert.Contains(t, template, "models.UserRequest")

	// Check for required fields
	assert.Contains(t, template, "name")
	assert.Contains(t, template, "email")
	assert.Contains(t, template, "age")

	// Check for optional fields
	assert.Contains(t, template, "phone")
	assert.Contains(t, template, "address")
	assert.Contains(t, template, "is_active")
}

func TestSwaggerInfo_ResponseSchemas(t *testing.T) {
	// Test that swagger template contains response schemas
	template := docs.SwaggerInfo.SwaggerTemplate

	// Check for success responses
	assert.Contains(t, template, "200")
	assert.Contains(t, template, "201")

	// Check for error responses
	assert.Contains(t, template, "400")
	assert.Contains(t, template, "404")
	assert.Contains(t, template, "500")

	// Check for response descriptions
	assert.Contains(t, template, "User created successfully")
	assert.Contains(t, template, "User updated successfully")
	assert.Contains(t, template, "User deleted successfully")
	assert.Contains(t, template, "Bad request")
	assert.Contains(t, template, "User not found")
}

func TestSwaggerInfo_Metadata(t *testing.T) {
	// Test swagger metadata
	template := docs.SwaggerInfo.SwaggerTemplate

	// Check for API metadata in the SwaggerInfo struct directly
	assert.Equal(t, "User Management API", docs.SwaggerInfo.Title)
	assert.Contains(t, docs.SwaggerInfo.Description, "RESTful API for user management")

	// Check for metadata in template
	assert.Contains(t, template, "API Support")
	assert.Contains(t, template, "MIT")

	// Check for contact information
	assert.Contains(t, template, "support@swagger.io")
	assert.Contains(t, template, "http://www.swagger.io/support")

	// Check for license
	assert.Contains(t, template, "https://opensource.org/licenses/MIT")
}
