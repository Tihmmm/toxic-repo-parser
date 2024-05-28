package main

import (
	"errors"
	"log"
	"toxic-repos/pkg/file"
)

func parseBom(dest *result) error {
	switch bomFormat {
	case "cdxjson":
		var cdx cdxJson
		if err := file.Json2Struct(bom, &cdx); err != nil {
			log.Printf("Error parsing JSON file: %s", err)
			return err
		}
		dest.components = cdx.Components
	default:
		log.Println("Not implemented")
		return errors.ErrUnsupported
	}
	return nil
}
