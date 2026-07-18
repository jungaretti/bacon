package deployment

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonReporter struct {
	writer io.Writer
}

type jsonOutput struct {
	Summary jsonSummary  `json:"summary"`
	Records []jsonResult `json:"records"`
}

type jsonSummary struct {
	Deleted     int `json:"deleted"`
	WouldDelete int `json:"wouldDelete"`
	Updated     int `json:"updated"`
	WouldUpdate int `json:"wouldUpdate"`
	Created     int `json:"created"`
	WouldCreate int `json:"wouldCreate"`
	Kept        int `json:"kept"`
}

type jsonResult struct {
	Action Action     `json:"action"`
	DryRun bool       `json:"dryRun"`
	Record jsonRecord `json:"record"`
}

type jsonRecord struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      string `json:"ttl"`
	Priority string `json:"priority"`
	Notes    string `json:"notes"`
}

func (reporter *jsonReporter) Preview(results []Operation) error {
	return nil
}

func (reporter *jsonReporter) Report(results []Operation) error {
	output := jsonOutput{Records: make([]jsonResult, len(results))}
	for i, result := range results {
		switch result.Action {
		case Delete:
			if result.DryRun {
				output.Summary.WouldDelete++
			} else {
				output.Summary.Deleted++
			}
		case Update:
			if result.DryRun {
				output.Summary.WouldUpdate++
			} else {
				output.Summary.Updated++
			}
		case Create:
			if result.DryRun {
				output.Summary.WouldCreate++
			} else {
				output.Summary.Created++
			}
		case Keep:
			output.Summary.Kept++
		}

		record := result.Record
		output.Records[i] = jsonResult{
			Action: result.Action,
			DryRun: result.DryRun,
			Record: jsonRecord{
				Name:     record.Name,
				Type:     record.Type,
				Content:  record.Content,
				TTL:      record.TTL,
				Priority: record.Priority,
				Notes:    record.Notes,
			},
		}
	}

	encoded, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(reporter.writer, string(encoded))
	return err
}
