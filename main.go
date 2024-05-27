package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Cdx struct {
	BomFormat   string `json:"bomFormat"`
	SpecVersion string `json:"specVersion"`
	Components  []struct {
		Name string `json:"name"`
		//Ref       string   `json:"ref"`
		//DependsOn []string `json:"dependsOn"`
	} `json:"components"`
}

type result struct {
	toxicRepos []repo
}

type repo struct {
	id          int
	datetime    string
	problemType string
	name        string
	commitLink  string
	description string
}

var blockers = []string{"malware", "ddos", "broken_assembly"}

func main() {
	db, err := sql.Open("sqlite3", "/tmp/toxic-repos.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	var res result
	var cdx Cdx
	err = parseJsonFile("/tmp/", "bom.json", &cdx)
	if err != nil {
		log.Fatalf("Error parsing JSON file: %v", err)
	}
	for _, v := range cdx.Components {
		//pack := getPackageName(v.Ref)
		pack := v.Name

		r := new(repo)
		err = dbCheckPack(pack, r, db)
		if errors.Is(sql.ErrNoRows, err) {
			//for _, d := range v.DependsOn {
			//	transPack := getPackageName(d)
			//	fmt.Println(transPack)
			//	err = dbCheckPack(transPack, r, db)
			//	if errors.Is(sql.ErrNoRows, err) {
			//		continue
			//	}
			//	if err != nil {
			//		panic(err)
			//	}
			//	res.toxicRepos = append(res.toxicRepos, *r)
			//}
			continue
		}
		if err != nil {
			panic(err)
		}
		res.toxicRepos = append(res.toxicRepos, *r)
	}

	if len(res.toxicRepos) > 0 && !checkForBlockers(&res) {
		fmt.Printf("Найдены токсичные компоненты:\n %v\n", res)
	}
}

func dbCheckPack(pack string, r *repo, db *sql.DB) error {
	row := db.QueryRow("select * from repos where lower(name)=$1", strings.ToLower(pack))
	return row.Scan(&r.id, &r.datetime, &r.problemType, &r.name, &r.commitLink, &r.description)
}

func checkForBlockers(toxicRepos *result) (passed bool) {
	for _, v := range toxicRepos.toxicRepos {
		if slices.Contains(blockers, v.problemType) {
			return false
		}
	}
	return true
}

func parseJsonFile(fileDir string, fileName string, dest any) error {
	jsonFile, err := os.Open(filepath.Join(fileDir, fileName))
	if err != nil {
		log.Printf("Error opening jsonFile file: %s, err: %s\n", fileName, err)
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
		log.Printf("Error decoding json file: %s, err: %s\n", fileName, err)
		return err
	}

	return nil
}

//func getPackageName(ref string) string {
//	s2, _, _ := strings.Cut(ref[strings.LastIndex(ref, "/")+1:], "@")
//	return s2
//}
