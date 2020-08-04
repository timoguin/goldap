# goldap

This tool is intended to be a clone of the standard LDAP CLI tool suite written in
Golang. It is fully configurable via CLI flags, environment variables, or config files.

It currently only supports LDAP search functionality.

## Installing

If you have a properly-configured Golang development environment, you can quickly build
and install from source with `go get`:

```
$ go get github.com/timoguin/goldap
```

Binaries are currently not published, but they will be published with GitHub releases
in the future.

If you clone the repo, you can build and install with `make install`.

## Usage

Help documentation is built-in to the CLI itself. You can pass `-h` to any command to
get help output.

```
$ goldap search -h
===========================================================================================
|                                                                                         |
| Clone of the ldapsearch tool.                                                           |
|                                                                                         |
| The short format of the flags are identical between the two. Long names are formatted   |
| to be more readable, while the familiar names from ldapsearch are supported as aliases. |
| The filter must be passed via the --filter flag, and attributes must be passed via one  |
| or more --attribute flags.                                                              |
|                                                                                         |  
===========================================================================================

Usage:
  goldap search [flags]

Flags:
  -a, --alias-deref string   How to handle alias derefering (never | always | search | find) (default "never")
      --attribute strings    Specify multiple attributes by passing this flag multiple times
  -D, --bind-dn string       LDAP binddn used for authentication
  -f, --filter string        The filter string for the query (default "(objectClass=*)")
  -h, --help                 help for search
  -H, --ldap-uri string      URI of the LDAP host
  -p, --paging-size int      Number of records to return with pagination (default 100)
  -w, --password string      LDAP password
  -s, --scope string         The search scope (base | one | sub) (default "sub")
  -b, --search-base string   The starting point for the search (LDAP searchbase)
  -z, --size-limit int       Limit the number of records returned
  -Z, --start-tls            Connect using TLS
  -l, --time-limit int       Number of seconds to wait for the query to return (default: 60) (default 60)
  -A, --types-only           Only return attribute names (not values)

Global Flags:
  -c, --config-file string   Path to the config file
  -d, --debug                Enable debug logs
```

## Configuration

All configuration options can be specifed in one or more ways. Here is the order of
precendence for determining the final configuration of any one command:

- CLI flags
- Environment variables prefixed with `GOLDAP_`
- Config file

All 3 of these configuration sources are merged together in that order of precendence.
This allows setting default configuration options via a file, but using environment
variables and flags to customize the specific command.

For example, you can put the `ldap-uri` and `bind-dn` into a config file to set the
default LDAP endpoint and your LDAP Bind DN (username for simple auth), while setting
your password via the `GOLDAP_PASSWORD` environment variable.

### Config File

The config file location can be customized via the `--config-file` flag. By default,
the CLI searches for the config file in the following order of precedence:

- `/etc/goldap/`
- `$HOME/`
- `./`

The config file can be written in multiple formats. The CLI uses the file extension to
determine how to parse the config. Use whichever of the following you are most
comfortable with:

| Format          | File Name      | File Extension               |
|-----------------|----------------|------------------------------|
| JSON            | .goldap        | json                         |
| TOML            | .goldap        | toml                         |
| YAML            | .goldap        | yaml / yml                   |
| HCL             | .goldap        | hcl                          |
| INI             | .goldap        | ini                          |
| envfile         | dotfile or env | None                         |
| Java properties | .goldap        | .properties / .props / .prop |

If a config file is found that has no file extension, the CLI attempts to parse it as
YAML.
