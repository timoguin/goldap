package cmd

import (
	"crypto/tls"
	"os"

	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Help text for the search command
var SearchHelp = map[string]string{
	"Use":   "search",
	"Short": "Execute an LDAP search",

	"Long": `
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

`,
}

// The search command as defined by Cobra
var searchCmd = &cobra.Command{
	Use:   SearchHelp["Use"],
	Short: SearchHelp["Short"],
	Long:  SearchHelp["Long"],
	Run:   SearchCmdRun,
}

// Variables passed to the LDAP search query
var (
	aliasDeref string
	attribute  []string
	bindDN     string
	controls   []ldap.Control
	filter     string
	ldapURI    string
	pagingSize int
	password   string
	scope      string
	searchBase string
	sizeLimit  int
	startTLS   bool
	timeLimit  int
	typesOnly  bool
)

// Maps string names from the --scope CLI flag to constants
var ScopeMapping = map[string]int{
	"_default": ldap.ScopeWholeSubtree,
	"base":     ldap.ScopeBaseObject,
	"one":      ldap.ScopeSingleLevel,
	"sub":      ldap.ScopeWholeSubtree,
}

// Maps string names from the --alias-deref CLI flag to constants
var AliasDerefMapping = map[string]int{
	"_default": ldap.NeverDerefAliases,
	"never":    ldap.NeverDerefAliases,
	"always":   ldap.DerefAlways,
	"search":   ldap.DerefFindingBaseObj,
	"find":     ldap.NeverDerefAliases,
}

// Look up an integer value from a string map, fall back to _default key
func LookupConst(s string, m map[string]int) int {
	v, found := m[s]
	if !found {
		v = m["_default"]
	}
	return v
}

// Represents all the values passed to NewSearchRequest and NewSearchRequestWithPaging
type LDAPNewSearchRequestInput struct {
	SearchBase string // Create alias: BaseDN / base-dn
	Scope      int
	AliasDeref int
	SizeLimit  int
	TimeLimit  int
	TypesOnly  bool
	Filter     string
	Attributes []string
	Controls   []ldap.Control
}

// Execute the search
func SearchCmdRun(cmd *cobra.Command, args []string) {
	Logger.Debug("test in search command")

	// Prepare input for search request
	searchInput := LDAPNewSearchRequestInput{
		SearchBase: viper.GetString("search-base"),
		SizeLimit:  viper.GetInt("size-limit"),
		TimeLimit:  viper.GetInt("time-limit"),
		Scope:      LookupConst(viper.GetString("scope"), ScopeMapping),
		AliasDeref: LookupConst(viper.GetString("alias-deref"), AliasDerefMapping),
		TypesOnly:  viper.GetBool("types-only"),
		Filter:     viper.GetString("filter"),
		Controls:   nil,

		// Append any --attribute flag values with any defined in config
		Attributes: append(
			viper.GetStringSlice("attributes"),
			viper.GetStringSlice("attribute")...,
		),
	}

	// Connect to the LDAP server and start a session
	ldapURI = viper.GetString("ldap-uri")

	session, err := ldap.DialURL(ldapURI)
	if err != nil {
		Logger.Fatalw("Failed to dial URI",
			"uri", ldapURI,
			"err", err,
		)
		os.Exit(1)
	}
	defer session.Close()

	// Authenticate the client via Simple BIND
	bindDN = viper.GetString("bind-dn")
	password = viper.GetString("password")

	if err := session.Bind(bindDN, password); err != nil {
		Logger.Fatalw("Failed to perform bind",
			"binddn", bindDN,
			"err", err,
		)
		os.Exit(1)
	}

	// Reconnect with TLS
	if viper.GetBool("start-tls") {
		err = session.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			Logger.Fatalw("Failed to connect with TLS", "err", err)
			os.Exit(1)
		}
	}

	// Prepare the LDAP search inputs
	searchRequest := ldap.NewSearchRequest(
		searchInput.SearchBase,
		searchInput.Scope,
		searchInput.AliasDeref,
		searchInput.SizeLimit,
		searchInput.TimeLimit,
		searchInput.TypesOnly,
		searchInput.Filter,
		searchInput.Attributes,
		searchInput.Controls,
	)

	// Perform the search
	pagingSize := uint32(viper.GetInt("paging-size"))
	searchResult, err := session.SearchWithPaging(searchRequest, pagingSize)
	if err != nil {
		Logger.Fatalw("Failed to perform search",
			"searchRequest", searchRequest,
			"err", err,
		)
		os.Exit(1)
	}

	// Print search results
	Logger.Infow("Performed search", "entries", searchResult.Entries)

	// // Marshal JSON and print
	// var jsonData []byte
	// jsonData, err = json.Marshal(searchResult.Entries)
	// if err != nil {
	// 	logger.Fatalw("Failed to marshal results to JSON", "err", err)
	// }
	// fmt.Println(string(jsonData))
}

// Initialize the search command: parse CLI flags and other configuration
func init() {
	rootCmd.AddCommand(searchCmd)

	// Search command flags
	searchCmd.PersistentFlags().StringVarP(&ldapURI, "ldap-uri", "H", "", "URI of the LDAP host")
	searchCmd.PersistentFlags().StringVarP(&bindDN, "bind-dn", "D", "", "LDAP binddn used for authentication")
	searchCmd.PersistentFlags().StringVarP(&password, "password", "w", "", "LDAP password")
	searchCmd.Flags().StringVarP(&searchBase, "search-base", "b", "", "The starting point for the search (LDAP searchbase)")
	searchCmd.Flags().StringVarP(&scope, "scope", "s", "sub", "The search scope (base | one | sub)")
	searchCmd.Flags().StringVarP(&aliasDeref, "alias-deref", "a", "never", "How to handle alias derefering (never | always | search | find)")
	searchCmd.Flags().IntVarP(&sizeLimit, "size-limit", "z", 0, "Limit the number of records returned")
	searchCmd.Flags().IntVarP(&timeLimit, "time-limit", "l", 60, "Number of seconds to wait for the query to return (default: 60)")
	searchCmd.Flags().BoolVarP(&typesOnly, "types-only", "A", false, "Only return attribute names (not values)")
	searchCmd.Flags().StringVarP(&filter, "filter", "f", "(objectClass=*)", "The filter string for the query")
	searchCmd.Flags().StringSliceVarP(&attribute, "attribute", "", []string{}, "Specify multiple attributes by passing this flag multiple times")
	searchCmd.Flags().IntVarP(&pagingSize, "paging-size", "p", 100, "Number of records to return with pagination")
	searchCmd.Flags().BoolVarP(&startTLS, "start-tls", "Z", false, "Connect using TLS")

	// Bind all command flags to Viper config
	if err := viper.BindPFlags(searchCmd.Flags()); err != nil {
		Logger.Fatalf("Error binding config flags with viper: %s", err)
		os.Exit(1)
	}
}
