package app

import "os"

func StartApplication() {
	wiring := newWiring()
	setupRoutes(wiring)
	runApplication()
}

func runApplication() {
	port := os.Getenv("ENV")
	if port == "" {
		port = "8085"
	}
	route.Run(":" + port)
}
