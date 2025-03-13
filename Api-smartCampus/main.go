package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jalil/Api-smartCampus/controllers"
	"github.com/jalil/Api-smartCampus/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}

func main() {
	r := gin.Default()
	r.GET("/tabsensor", controllers.TabsensorIndex)
	r.Run() // listen and serve on 0.0.0.0:8080
}
