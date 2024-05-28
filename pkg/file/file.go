package file

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

func Csv2Struct(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening csv file: %s,", err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	r := csv.NewReader(file)
	records := make([][]string, 0)
	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) || errors.Is(err, csv.ErrBareQuote) || errors.Is(err, csv.ErrFieldCount) || errors.Is(err, csv.ErrQuote) {
			break
		}
		if err != nil {
			log.Printf("Error parsing csv: %s\n", err)
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func Json2Struct(path string, dest any) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening jsonFile file: %s, err: %s\n", path, err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(jsonFile)

	jsonParser := json.NewDecoder(jsonFile)
	if err = jsonParser.Decode(dest); err != nil {
		log.Printf("Error decoding json file: %s, err: %s\n", path, err)
		return err
	}

	return nil
}
