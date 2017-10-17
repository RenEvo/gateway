package main

import "fmt"
import "github.com/renevo/gateway"
import "os"

func main() {
	f, err := os.Open("./sample.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cfg, err := gateway.LoadConfiguration(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.String())
}
