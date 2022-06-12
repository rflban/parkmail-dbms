package main

import (
	"context"
	"fmt"
)

func main() {
	conf, _ := getConfig(context.Background())
	fmt.Printf("%+v\n", conf)
}
