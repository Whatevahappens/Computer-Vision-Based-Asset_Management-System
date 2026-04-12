package server

import (
	"backend/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.SetupRoutes(r)
	r.Run(":8080")
}
