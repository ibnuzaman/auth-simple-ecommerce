package main

import (
	"github.com/ibnuzaman/auth-simple-ecommerce.git/cmd"
	"github.com/ibnuzaman/auth-simple-ecommerce.git/helpers"
)

func main() {
	helpers.SetupUpConfig()

	cmd.ServeHTTP()
}
