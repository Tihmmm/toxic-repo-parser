package main

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func parseDatasource(dest *result) error {
	var datasource string
	if datasourceLocal != "" {
		datasource = datasourceLocal
	} else {
		if err := downloadDatasource(datasourceRemote); err != nil {
			return err
		}
		datasource = datasourceFetchDest
	}

	switch datasourceType {
	case "sqlite3":
		return checkDB(dest, datasource)
	default:
		log.Println("Not implemented")
		os.Exit(2)
	}

	return nil
}

func downloadDatasource(url string) error {
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error downloading datasource: %s\n", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error downloading datasource. Response status code: %d\n", resp.StatusCode)
		return err
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			return
		}
	}(resp.Body)

	f, err := os.Create(datasourceFetchDest)
	if err != nil {
		log.Printf("Error saving datasource file: %s\n", err)
		return err
	}
	if _, err = io.Copy(f, resp.Body); err != nil {
		log.Printf("Error saving datasource file: %s\n", err)
		return err
	}

	return nil
}

func checkDB(res *result, datasource string) error {
	db, err := sql.Open("sqlite3", datasource)
	if err != nil {
		log.Printf("Error connecting to sqlite db: %s", err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			return
		}
	}(db)

	for _, v := range res.components {
		pack := v.Name

		r := new(repo)
		row := db.QueryRow("select * from repos where lower(name)=$1", strings.ToLower(pack))
		err := row.Scan(&r.id, &r.datetime, &r.problemType, &r.name, &r.commitLink, &r.description)
		if errors.Is(sql.ErrNoRows, err) {
			continue
		}
		if err != nil {
			log.Printf("Error querying sqlite db: %s", err)
			return err
		}
		res.toxicRepos = append(res.toxicRepos, *r)
	}

	return nil
}
