package main

import (
	"bytes"
	_ "database/sql"
	"demo3-gin/internal/controllers"
	"demo3-gin/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm"
	_ "io"
	_ "io/ioutil"
)

func main() {
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	//middleware log response body
	r.Use(ginBodyLogMiddleware())

	models.ConnectDataBase()

	r.GET("/users", controllers.GetAll)
	r.GET("/users/:id", controllers.GetById)
	r.POST("/users", controllers.Insert)
	r.PUT("/users/:id", controllers.Update)
	r.DELETE("/users/:id", controllers.Delete)
	r.Run("localhost:8080")
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func ginBodyLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString("\n"), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		fmt.Println("Response body: " + blw.body.String())
	}
}
