package middleware

import (
	"net/http"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-gonic/gin"
)

// ClerkAuthMiddleware wraps Clerk’s WithHeaderAuthorization for Gin
func ClerkAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use Clerk's built-in middleware on the request
		handler := clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Pass request context back into Gin
			c.Request = c.Request.WithContext(r.Context())
			c.Next()
		}))

		// Run the handler
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
