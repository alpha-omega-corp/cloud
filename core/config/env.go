package config

import (
	"flag"
	"fmt"
)

func GetEnv() string {
	cmd := flag.String("env", "", "")
	flag.Parse()
	fmt.Printf("enivonrment : \"%v\"\n", *cmd)
	
	return *cmd
}
