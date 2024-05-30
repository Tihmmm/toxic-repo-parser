# toxic-repo-parser
```
Usage:
togore [-f <path_to_datasource> || -u <datasource_url>] [-b <comma,separated,blockers>] [-o <path_for_output_file>] -s <path_to_bom> -k <bom_format> -t <datasource_type> [flags]

Flags:
-s, --bom string                Path to BOM file for components to be parsed from (required)
-k, --bom-format string         BOM format. Supported formats: cdxjson(=CycloneDX json) (required) (default "cdxjson")
-f, --datasource-path string    Path to a toxic repos datasource
-t, --datasource-type string    Datasource type. Supported types: sqlite3 (required) (default "sqlite3")
-u, --datasource-url string     URL address to a remote toxic repos datasource. Defaults to a sqlite db as provided by github.com/toxic-repos/toxic-repos (default "https://raw.githubusercontent.com/toxic-repos/toxic-repos/main/data/sqlite/toxic-repos.sqlite3")
-b, --blockers strings          Comma separated list of problem types. If ANY of the provided problems is found, program exits with code 130. Default value is malware,ddos,broken_assembly (default [malware,ddos,broken_assembly])
-h, --help                      help for togore
-o, --output string             Path to output file. If not specified, program writes to stdout
-v, --version                   version for togore
```
Current version: beta