package main

import (
	"fmt"
	"urlshorter/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()

	fmt.Printf("%#v\n", cfg) // удалить во время прода

	//TODO: init config:cleanenv

	//TODO: init logger: slog

	//TODO: init storage: sqlite

	//TODO: init router: chi, "chi render"

	//TODO: run server
}
