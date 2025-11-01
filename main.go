package main

import (
	"github.com/ibnuzaman/auth-simple-ecommerce.git/cmd"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
)

// @title Auth Simple Ecommerce API
// @version 1.0
// @description Authentication microservice for e-commerce platform with JWT
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:9000
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	helpers.SetupUpConfig()

	helpers.SetupLogger()

	helpers.SetupPostgreSQL()

	cmd.ServeHTTP()
}
