package tool

import (
	"fmt"
	"io"
)

// HelpCategory is a bitmask for help categories.
// Translation of the CURLHELP_ defines from tool_help.h.
type HelpCategory uint

const (
	HelpAuth       HelpCategory = 1 << iota
	HelpConnection              // 1 << 1
	HelpCurl                    // 1 << 2
	HelpDeprecated              // 1 << 3
	HelpDNS                     // 1 << 4
	HelpFile                    // 1 << 5
	HelpFTP                     // 1 << 6
	HelpGlobal                  // 1 << 7
	HelpHTTP                    // 1 << 8
	HelpIMAP                    // 1 << 9
	HelpImportant               // 1 << 10
	HelpLDAP                    // 1 << 11
	HelpOutput                  // 1 << 12
	HelpPOP3                    // 1 << 13
	HelpPost                    // 1 << 14
	HelpProxy                   // 1 << 15
	HelpSCP                     // 1 << 16
	HelpSFTP                    // 1 << 17
	HelpSMTP                    // 1 << 18
	HelpSSH                     // 1 << 19
	HelpTelnet                  // 1 << 20
	HelpTFTP                    // 1 << 21
	HelpTimeout                 // 1 << 22
	HelpTLS                     // 1 << 23
	HelpUpload                  // 1 << 24
	HelpVerbose                 // 1 << 25
	HelpAll = 0xfffffff
)

// HelpText holds the text for a single command-line option.
// Translation of the C `helptxt` struct.
type HelpText struct {
	Option      string
	Description string
	Categories  HelpCategory
}

// CategoryDescriptor holds the name and description of a help category.
// Translation of the C `category_descriptors` struct.
type CategoryDescriptor struct {
	Option      string
	Description string
	Category    HelpCategory
}

// categories is the translated list of help categories from tool_help.c.
var categories = []CategoryDescriptor{
	{"auth", "Authentication methods", HelpAuth},
	{"connection", "Manage connections", HelpConnection},
	{"curl", "The command line tool itself", HelpCurl},
	{"http", "HTTP and HTTPS protocol", HelpHTTP},
	// Add other categories as needed
}

// helptext is a placeholder for the full list of options from tool_hugehelp.c.
// We use a small sample here to build and test the printing logic.
var helptext = []HelpText{
	{"-d, --data <data>", "HTTP POST data", HelpHTTP | HelpPost | HelpImportant},
	{"-H, --header <header>", "Pass custom header to server", HelpHTTP | HelpImportant},
	{"-I, --head", "Show document info only", HelpHTTP | HelpFTP | HelpFile},
	{"-L, --location", "Follow redirects", HelpHTTP | HelpImportant},
	{"-o, --output <file>", "Write to file instead of stdout", HelpImportant},
	{"-u, --user <user:password>", "Server user and password", HelpAuth | HelpImportant},
	{"-v, --verbose", "Make the operation more talkative", HelpVerbose | HelpImportant},
}

// PrintHelp prints help information for a given category.
// This is a translation of the C `tool_help` function.
func PrintHelp(writer io.Writer, category string) {
	cols := GetTerminalColumns() // Assumes this function exists from previous translations.

	if category == "" {
		fmt.Fprintln(writer, "Usage: curl [options...] <url>")
		printCategory(writer, HelpImportant, cols)
		fmt.Fprintln(writer, "\nThis is not the full help; this menu is split into categories.")
		fmt.Fprintln(writer, `Use "curl --help [category]" to get help for one category.`)
		return
	}

	if category == "all" {
		printCategory(writer, HelpAll, cols)
		return
	}

	if category == "category" {
		fmt.Fprintln(writer, "Available categories:")
		for _, cat := range categories {
			fmt.Fprintf(writer, "  %-12s %s\n", cat.Option, cat.Description)
		}
		return
	}

	for _, cat := range categories {
		if cat.Option == category {
			fmt.Fprintf(writer, "%s: %s\n", cat.Option, cat.Description)
			printCategory(writer, cat.Category, cols)
			return
		}
	}

	fmt.Fprintf(writer, "Unknown help category: %s\n", category)
}

// printCategory formats and prints all options belonging to a given category.
// This is a translation of the C `print_category` function.
func printCategory(writer io.Writer, category HelpCategory, cols int) {
	var maxOptLen int
	for _, ht := range helptext {
		if ht.Categories&category != 0 {
			if len(ht.Option) > maxOptLen {
				maxOptLen = len(ht.Option)
			}
		}
	}
	// Add some padding
	maxOptLen += 2

	for _, ht := range helptext {
		if ht.Categories&category != 0 {
			// Basic formatting, not as complex as the C version's column calculation
			line := fmt.Sprintf(" %-*s  %s", maxOptLen, ht.Option, ht.Description)
			if len(line) > cols {
				line = line[:cols]
			}
			fmt.Fprintln(writer, line)
		}
	}
}