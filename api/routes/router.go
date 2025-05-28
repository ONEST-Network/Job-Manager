package routes

import (
	"github.com/ONEST-Network/scheme-manager-adapter/api/handlers"
	"github.com/ONEST-Network/scheme-manager-adapter/api/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the application
func SetupRouter(baseHandler *handlers.BaseHandler) *gin.Engine {
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// API version 1
	v1 := router.Group("/api/v1")

	// Initialize handlers
	orgHandler := handlers.NewOrganizationHandler(baseHandler)
	schemeHandler := handlers.NewSchemeHandler(baseHandler)
	appHandler := handlers.NewApplicationHandler(baseHandler)
	
	setupOrganizationRoutes(v1, orgHandler)
	setupSchemeRoutes(v1, schemeHandler)
	setupApplicationRoutes(v1, appHandler)

	return router
}

// setupOrganizationRoutes configures all organization-related routes
func setupOrganizationRoutes(v1 *gin.RouterGroup, orgHandler *handlers.OrganizationHandler) {
	orgs := v1.Group("/organizations")
	{
		orgs.POST("", orgHandler.Create)
		orgs.GET("", orgHandler.List)
		orgs.GET("/:id", orgHandler.GetByID)
		orgs.GET("/api-key/:api_key", orgHandler.GetByAPIKey)
		orgs.PUT("/:id", orgHandler.Update)
		orgs.DELETE("/:id", orgHandler.Delete)
	}

	// Organization-specific scheme routes
	orgSchemes := v1.Group("/org-schemes")
	{
		orgSchemes.POST("/:org_id", orgHandler.CreateScheme)
		orgSchemes.GET("/:org_id", orgHandler.ListSchemes)
		orgSchemes.GET("/:org_id/:scheme_id", orgHandler.GetScheme)
	}
}

// setupSchemeRoutes configures all scheme-related routes
func setupSchemeRoutes(v1 *gin.RouterGroup, schemeHandler *handlers.SchemeHandler) {
	schemes := v1.Group("/schemes")
	{
		schemes.GET("/:id", schemeHandler.GetByID)
		schemes.PUT("/:id/status", schemeHandler.UpdateStatus)
		schemes.DELETE("/:id", schemeHandler.Delete)
	}
}

// setupApplicationRoutes configures all application-related routes
func setupApplicationRoutes(v1 *gin.RouterGroup, appHandler *handlers.ApplicationHandler) {
	// Scheme-specific application routes
	schemeApps := v1.Group("/scheme-applications")
	{
		schemeApps.POST("/:scheme_id", appHandler.Create)
		schemeApps.GET("/:scheme_id", appHandler.GetByScheme)
	}

	// General application routes
	applications := v1.Group("/applications")
	{
		applications.GET("/:id", appHandler.GetByID)
		applications.PUT("/:id/status", appHandler.UpdateStatus)
		applications.DELETE("/:id", appHandler.Delete)
	}
}