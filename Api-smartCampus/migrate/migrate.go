package main

import (
	models "github.com/jalil/Api-smartCampus/Models"
	"github.com/jalil/Api-smartCampus/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Tabsensor{})
}
