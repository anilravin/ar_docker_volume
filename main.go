package main

import (
	"log"
	//	"flag"
	//	"fmt"
	//	"os"
	//	"path/filepath"
	//	"strings"

	"github.com/docker/go-plugins-helpers/volume"
)

func main() {
	driver := NewARVol()
	handler := volume.NewHandler(driver)
	//if err := handler.ServeUnix("root", "driver-example"); err != nil {
	if err := handler.ServeUnix("root", 0); err != nil {
		log.Fatalf("Error %v", err)
	}

	for {

	}
}
