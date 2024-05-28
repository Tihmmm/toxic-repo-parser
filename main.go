package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"log"
	"os"
	"slices"
)

type component struct {
	Name string `json:"name"`
}

type cdxJson struct {
	BomFormat   string      `json:"bomFormat"`
	SpecVersion string      `json:"specVersion"`
	Components  []component `json:"components"`
}

type result struct {
	toxicRepos []repo
	components []component
}

type repo struct {
	id          int
	datetime    string
	problemType string
	name        string
	commitLink  string
	description string
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

var (
	bomFormats      = []string{"cdxjson"}
	datasourceTypes = []string{"sqlite3", "json", "csv"}
	defaultBlockers = []string{"malware", "ddos", "broken_assembly"}
)

const (
	defaultBomFormat        = "cdxjson"
	defaultDatasourceRemote = "https://raw.githubusercontent.com/toxic-repos/toxic-repos/main/data/sqlite/toxic-repos.sqlite3"
	defaultDatasourceType   = "sqlite3"
	datasourceFetchDest     = "togore-ds"
)

var (
	bom              string
	bomFormat        string
	datasourceLocal  string
	datasourceRemote string
	datasourceType   string
	blockers         *[]string
	output           string

	rootCmd = &cobra.Command{
		Use:     "togore [-f <path_to_datasource> || -u <datasource_url>] [-b <comma,separated,blockers>] [-o <path_for_output_file>] -s <path_to_bom> -k <bom_format> -t <datasource_type>",
		Version: "omega",
		Short:   "",
		Long:    ``,
		Run:     run,
	}
)

func init() {
	rootCmd.Flags().StringVarP(&bom, "bom", "s", "", "Path to BOM file for components to be parsed from (required)")
	rootCmd.MarkFlagRequired("bom")
	rootCmd.Flags().StringVarP(&bomFormat, "bom-format", "k", defaultBomFormat, "BOM format. Supported formats: cdxjson(=CycloneDX json) (required)")
	rootCmd.MarkFlagRequired("bom-format")
	rootCmd.Flags().StringVarP(&datasourceLocal, "datasource-path", "f", "", "Path to a toxic repos datasource")
	rootCmd.Flags().StringVarP(&datasourceRemote, "datasource-url", "u", defaultDatasourceRemote, "URL address to a remote toxic repos datasource. Defaults to a sqlite db as provided by github.com/toxic-repos/toxic-repos")
	rootCmd.MarkFlagsMutuallyExclusive("datasource-path", "datasource-url")
	rootCmd.Flags().StringVarP(&datasourceType, "datasource-type", "t", defaultDatasourceType, "Datasource type. Supported types: sqlite3 (required)")
	rootCmd.MarkFlagRequired("datasource-type")
	blockers = rootCmd.Flags().StringSliceP("defaultBlockers", "b", defaultBlockers, "Comma separated list of problem types. If ANY of the provided problems is found, program exits with code 130. Default value is malware,ddos,broken_assembly")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Path to output file. If not specified, program writes to stdout")
}

func run(cmd *cobra.Command, args []string) {
	if !argsValid() {
		log.Println("Unsupported arguments")
		os.Exit(1)
	}

	var res result
	if err := parseBom(&res); err != nil {
		os.Exit(1)
	}

	if err := parseDatasource(&res); err != nil {
		os.Exit(1)
	}

	if len(res.toxicRepos) > 0 {
		log.Printf("Found toxic repos:\n %v\n", res.toxicRepos)
	}
	if containsBlockers(res.toxicRepos) {
		os.Exit(130)
	}
}

func argsValid() bool {
	return slices.Contains(bomFormats, bomFormat) && slices.Contains(datasourceTypes, datasourceType)
}

func containsBlockers(toxicRepos []repo) bool {
	for _, v := range toxicRepos {
		if slices.Contains(*blockers, v.problemType) {
			return true
		}
	}
	return false
}
