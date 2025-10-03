package main

import (
	"log"
	"os"
	"project-z-backend/src/users"

	"github.com/gin-gonic/gin"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {

	users.InitDB()
	defer users.DB.Close()

	r := gin.Default()

	r.POST("/users", users.Register)
	r.GET("/users", users.UserInfo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	log.Println("Server running on port", port)
	r.Run(":" + port)
}
