package main

import "link-shortener/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
