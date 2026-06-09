package api

import (
	"log/slog"
	"net/http"

	"goadmin-backend/internal/cmd/api/routers"
	"goadmin-backend/internal/platform/httproute"
	"goadmin-backend/internal/platform/ratelimit"
)

// NewRouter returns a new router with all routes defined
func NewRouter(validator *OpenAPIValidator, handlers *Handlers, logger *slog.Logger) http.Handler {
	// use chi router
	router := routers.NewChiRouter(logger)

	// OpenAPI validator disabled — libopenapi-validator has bugs with
	// nested path parameters. Handlers validate their own inputs.
	// router.Use(validator.Middleware)

	// Rate limit auth endpoints: 20 req/min per IP, burst of 5
	authLimiter := ratelimit.New(20, 5)

	// public routes
	router.Get("/health", handlers.HealthHandler.healthCheck)
	router.Group(func(r httproute.Router) {
		r.Use(authLimiter.Middleware)
		r.Post("/auth/login", handlers.AuthHandler.Login)
		r.Post("/auth/signup", handlers.AuthHandler.Register)
		r.Post("/auth/forgot-password", handlers.AuthHandler.ForgotPassword)
	})
	router.Post("/auth/signin-with-google", handlers.AuthHandler.SignInWithGoogle)
	router.Post("/auth/refresh", handlers.AuthHandler.RefreshToken)
	router.Post("/auth/2fa/login", handlers.AuthHandler.Login2FA)
	router.Post("/auth/reset-password", handlers.AuthHandler.ResetPassword)

	// private routes (require authentication)
	router.Group(func(grt httproute.Router) {
		grt.Use(handlers.AuthHandler.Authenticator())

		grt.Route("/auth/logout", func(r httproute.Router) {
			r.Post("/", handlers.AuthHandler.Logout)
		})

		grt.Route("/auth/profile", func(r httproute.Router) {
			r.Get("/", handlers.AuthHandler.Profile)
		})
			grt.Route("/auth/2fa", func(r httproute.Router) {
				r.Get("/setup", handlers.AuthHandler.SetupTOTP)
				r.Post("/verify", handlers.AuthHandler.VerifyTOTP)
				r.Post("/disable", handlers.AuthHandler.DisableTOTP)
			})

		grt.Route("/v1/users", func(r httproute.Router) {
			r.Get("/", handlers.UserHandler.List)

			r.Route("/{userId}", func(r httproute.Router) {
				r.Get("/", handlers.UserHandler.GetByID)
				r.Patch("/", handlers.UserHandler.Update)
				r.Get("/roles", handlers.RBACHandler.ListRolesByUser)
				r.Post("/roles/{roleId}", handlers.RBACHandler.AssignRoleToUser)
				r.Delete("/roles/{roleId}", handlers.RBACHandler.UnassignRoleFromUser)
				r.Get("/permissions", handlers.RBACHandler.ListPermissionsByUser)
				r.Post("/permissions/{permissionId}", handlers.RBACHandler.AssignPermissionToUser)
				r.Delete("/permissions/{permissionId}", handlers.RBACHandler.UnassignPermissionFromUser)
			})
		})

		grt.Route("/v1/roles", func(r httproute.Router) {
			r.Get("/", handlers.RBACHandler.ListRoles)
			r.Post("/", handlers.RBACHandler.CreateRole)

			r.Route("/{id}", func(r httproute.Router) {
				r.Get("/", handlers.RBACHandler.GetRole)
				r.Patch("/", handlers.RBACHandler.UpdateRole)
				r.Delete("/", handlers.RBACHandler.DeleteRole)
				r.Get("/permissions", handlers.RBACHandler.ListPermissionsByRole)
				r.Post("/permissions/{permissionId}", handlers.RBACHandler.AddPermissionToRole)
				r.Delete("/permissions/{permissionId}", handlers.RBACHandler.RemovePermissionFromRole)
			})
		})

		grt.Route("/v1/permissions", func(r httproute.Router) {
			r.Get("/", handlers.RBACHandler.ListPermissions)
			r.Post("/", handlers.RBACHandler.CreatePermission)

			r.Route("/{id}", func(r httproute.Router) {
				r.Get("/", handlers.RBACHandler.GetPermission)
				r.Patch("/", handlers.RBACHandler.UpdatePermission)
				r.Delete("/", handlers.RBACHandler.DeletePermission)
			})
		})

		grt.Route("/v1/audit-logs", func(r httproute.Router) {
			r.Get("/", handlers.AuditHandler.List)
		})
	})

	return router
}
