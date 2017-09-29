package main

import (
	"fmt"
	_ "github.com/Anson5555555Lee/dataFeed/matchers"
	"github.com/Anson5555555Lee/dataFeed/search"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	search.Run("president")
}
