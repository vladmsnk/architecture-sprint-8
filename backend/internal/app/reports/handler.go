package reports

import (
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"net/http"
	"reports-api/internal/app/auth"
	report_getter "reports-api/internal/usecase/reports"
)

type Handler struct {
	report_getter.Getter
}

func New(keySet jwk.Set, reportGetter report_getter.Getter) *gin.Engine {
	h := Handler{
		Getter: reportGetter,
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.Use(auth.AuthMiddleware(keySet))

	r.GET("/reports", h.GetReport)

	return r
}
