package cmd

import (
	"bacon/pkg/config"
	"bacon/pkg/porkbun"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	createSymbol = "+"
	updateSymbol = "~"
	deleteSymbol = "-"
	keepSymbol   = "="
)

func newDeployCmd(client *porkbun.Client) *cobra.Command {
	var shouldCreate bool
	var shouldDelete bool
	var shouldUpdate bool

	deploy := &cobra.Command{
		Use:   "deploy <config-file>",
		Short: "Deploy records from a config file",
		Long: `Deploys DNS records from a YAML config file by updating records in place when
possible, then deleting old records and creating new records.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(client, args[0], shouldCreate, shouldDelete, shouldUpdate)
		},
	}

	deploy.Flags().BoolVarP(&shouldCreate, "create", "c", false, "create new records")
	deploy.Flags().BoolVarP(&shouldDelete, "delete", "d", false, "delete old records")
	deploy.Flags().BoolVarP(&shouldUpdate, "update", "u", false, "update changed records")

	return deploy
}

func deploy(client *porkbun.Client, configFile string, shouldCreate bool, shouldDelete bool, shouldUpdate bool) error {
	config, err := config.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("reading %v: %v", configFile, err)
	}

	from, err := client.AllRecords(config.Domain)
	if err != nil {
		return fmt.Errorf("fetching existing records: %v", err)
	}

	to := make([]porkbun.Record, len(config.Records))
	for i, record := range config.Records {
		to[i] = record.ToPorkbun()
	}

	added, removed, updated, unchanged := porkbun.DiffRecords(from, to)

	printRecords(removed, updated, added, unchanged)
	fmt.Println()

	if shouldDelete {
		for _, record := range removed {
			err := client.DeleteRecord(config.Domain, record)
			if err != nil {
				return fmt.Errorf("couldn't delete record: %v", err)
			}
		}
	}

	if shouldUpdate {
		for _, record := range updated {
			err := client.EditRecord(config.Domain, record)
			if err != nil {
				return fmt.Errorf("couldn't update record: %v", err)
			}
		}
	}

	if shouldCreate {
		for _, record := range added {
			_, err := client.CreateRecord(config.Domain, record)
			if err != nil {
				return fmt.Errorf("couldn't create record: %v", err)
			}
		}
	}

	clauses := []string{}
	if len(removed) > 0 {
		clauses = append(clauses, summarizeAction(shouldDelete, "deleted", "would delete", len(removed)))
	}
	if len(updated) > 0 {
		clauses = append(clauses, summarizeAction(shouldUpdate, "updated", "would update", len(updated)))
	}
	if len(added) > 0 {
		clauses = append(clauses, summarizeAction(shouldCreate, "created", "would create", len(added)))
	}
	if len(unchanged) > 0 {
		clauses = append(clauses, "kept "+countRecords(len(unchanged)))
	}

	if len(clauses) == 0 {
		fmt.Println("Summary: no records")
		return nil
	}

	summary := strings.Join(clauses, ", ")
	fmt.Println("Summary:", strings.ToUpper(summary[:1])+summary[1:])
	return nil
}

func printRecords(removed, updated, added, unchanged []porkbun.Record) {
	output := table.New("", "NAME", "TYPE", "CONTENT", "TTL", "PRI", "NOTES").
		WithHeaderFormatter(color.New(color.Underline).SprintfFunc()).
		WithWidthFunc(visibleWidth)

	for _, record := range removed {
		output.AddRow(recordRow(color.RedString(deleteSymbol), record)...)
	}
	for _, record := range updated {
		output.AddRow(recordRow(color.YellowString(updateSymbol), record)...)
	}
	for _, record := range added {
		output.AddRow(recordRow(color.GreenString(createSymbol), record)...)
	}
	for _, record := range unchanged {
		output.AddRow(recordRow(keepSymbol, record)...)
	}

	output.Print()
}

func recordRow(symbol string, record porkbun.Record) []any {
	return []any{symbol, record.Name, record.Type, record.Content, record.TTL, record.Priority, record.Notes}
}

var ansiEscapes = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Colored cells must not count their ANSI escapes toward column width.
func visibleWidth(text string) int {
	return utf8.RuneCountInString(ansiEscapes.ReplaceAllString(text, ""))
}

func summarizeAction(applied bool, appliedVerb string, plannedVerb string, count int) string {
	verb := plannedVerb
	if applied {
		verb = appliedVerb
	}
	return verb + " " + countRecords(count)
}

func countRecords(count int) string {
	if count == 1 {
		return "1 record"
	}
	return fmt.Sprintf("%d records", count)
}
