package records

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"biz/internal/command"
	perr "biz/internal/platform/errors"
	recordsdomain "biz/internal/records"
	"github.com/spf13/cobra"
)

func resolveCollectionDBID(s *command.State, collection string) string {
	c := normalizeCollection(collection)
	if c == "" {
		return ""
	}
	if dbID, ok := s.Runtime.Config.Notion.Collections[c]; ok && strings.TrimSpace(dbID) != "" {
		return strings.TrimSpace(dbID)
	}
	switch c {
	case "invoice", "invoices":
		return s.Runtime.Config.Notion.InvoiceDBID
	default:
		// Allow passing a raw Notion database id directly.
		return strings.TrimSpace(collection)
	}
}

func normalizeCollection(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func parseProperties(rawData, rawFile string) (map[string]any, error) {
	data := strings.TrimSpace(rawData)
	file := strings.TrimSpace(rawFile)
	if data == "" && file == "" {
		return nil, perr.New(perr.KindValidation, "one of --data or --data-file is required")
	}
	if data != "" && file != "" {
		return nil, perr.New(perr.KindValidation, "provide only one of --data or --data-file")
	}
	if file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, perr.Wrap(perr.KindValidation, "failed to read data file", err)
		}
		data = strings.TrimSpace(string(b))
	}
	if data == "" {
		return nil, perr.New(perr.KindValidation, "data is required")
	}
	var props map[string]any
	if err := json.Unmarshal([]byte(data), &props); err != nil {
		return nil, perr.Wrap(perr.KindValidation, "invalid data json", err)
	}
	if len(props) == 0 {
		return nil, perr.New(perr.KindValidation, "data must contain at least one property")
	}
	return props, nil
}

func validatePropertiesAgainstSchema(cmd *cobra.Command, s *command.State, collection, dbID string, props map[string]any) error {
	schema, err := s.Runtime.Records.Schema(cmd.Context(), recordsdomain.SchemaRequest{
		Collection: collection,
		DBID:       dbID,
	})
	if err != nil {
		return err
	}
	for _, k := range propertyKeys(props) {
		if _, ok := schema.Properties[k]; !ok {
			return perr.New(perr.KindValidation, "property not found in schema: "+k)
		}
	}
	return nil
}

func ensureLastEditedMatch(cmd *cobra.Command, s *command.State, id, expected string) error {
	current, err := s.Runtime.Records.Get(cmd.Context(), recordsdomain.GetRequest{ID: id})
	if err != nil {
		return err
	}
	want, err := time.Parse(time.RFC3339, strings.TrimSpace(expected))
	if err != nil {
		return perr.Wrap(perr.KindValidation, "invalid --if-last-edited timestamp (must be RFC3339)", err)
	}
	if !current.LastEditedTime.UTC().Equal(want.UTC()) {
		return perr.New(
			perr.KindConflict,
			fmt.Sprintf("record changed: expected last_edited_time=%s actual=%s", want.UTC().Format(time.RFC3339), current.LastEditedTime.UTC().Format(time.RFC3339)),
		)
	}
	return nil
}

func diffProperties(current, next map[string]any) map[string]map[string]any {
	out := map[string]map[string]any{}
	for _, k := range propertyKeys(next) {
		nv := next[k]
		cv, ok := current[k]
		if !ok || !jsonValueEqual(cv, nv) {
			out[k] = map[string]any{
				"before": cv,
				"after":  nv,
			}
		}
	}
	return out
}

func jsonValueEqual(a, b any) bool {
	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)
	return string(ab) == string(bb)
}

func propertyKeys(props map[string]any) []string {
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
