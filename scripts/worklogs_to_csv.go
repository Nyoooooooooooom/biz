package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type row struct {
	Date        string
	Hours       string
	Minutes     string
	Description string
	Ticket      string
	Source      string
}

var lineRE = regexp.MustCompile(`^\s*([0-9]{2}[-/][0-9]{2}[-/][0-9]{4})\s*-\s*([^\-]+?)\s*-\s*(.+?)\s*$`)
var ticketRE = regexp.MustCompile(`^\s*([0-9]{4,})\s*-\s*(.+)$`)

func parseDate(v string) (string, error) {
	v = strings.TrimSpace(v)
	for _, layout := range []string{"01-02-2006", "01/02/2006"} {
		if t, err := time.Parse(layout, v); err == nil {
			return t.Format("2006-01-02"), nil
		}
	}
	return "", fmt.Errorf("unsupported date format: %q", v)
}

func parseDurationToHoursAndMinutes(v string) (float64, int, error) {
	s := strings.ToLower(strings.TrimSpace(v))
	if s == "" {
		return 0, 0, fmt.Errorf("empty duration")
	}
	parts := strings.Fields(s)
	if len(parts) == 1 {
		num, unit, err := splitNumUnit(parts[0])
		if err != nil {
			return 0, 0, err
		}
		return toHoursMinutes(num, unit)
	}
	if len(parts) == 2 {
		num, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid duration number %q", parts[0])
		}
		return toHoursMinutes(num, parts[1])
	}
	if len(parts) == 4 {
		n1, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, 0, err
		}
		h1, m1, err := toHoursMinutes(n1, parts[1])
		if err != nil {
			return 0, 0, err
		}
		n2, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return 0, 0, err
		}
		h2, m2, err := toHoursMinutes(n2, parts[3])
		if err != nil {
			return 0, 0, err
		}
		totalMinutes := m1 + m2
		totalHours := h1 + h2 + float64(totalMinutes)/60.0
		return totalHours, totalMinutes, nil
	}
	return 0, 0, fmt.Errorf("unsupported duration format: %q", v)
}

func splitNumUnit(s string) (float64, string, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	i := 0
	for ; i < len(s); i++ {
		if (s[i] < '0' || s[i] > '9') && s[i] != '.' {
			break
		}
	}
	if i == 0 || i == len(s) {
		return 0, "", fmt.Errorf("invalid duration token: %q", s)
	}
	num, err := strconv.ParseFloat(s[:i], 64)
	if err != nil {
		return 0, "", err
	}
	return num, s[i:], nil
}

func toHoursMinutes(num float64, unit string) (float64, int, error) {
	u := strings.ToLower(strings.TrimSpace(unit))
	switch u {
	case "h", "hr", "hrs", "hour", "hours":
		minutes := int(num * 60.0)
		return float64(minutes) / 60.0, minutes, nil
	case "m", "min", "mins", "minute", "minutes":
		minutes := int(num)
		return float64(minutes) / 60.0, minutes, nil
	default:
		return 0, 0, fmt.Errorf("unsupported duration unit: %q", unit)
	}
}

func main() {
	in := flag.String("in", "fixtures/worklogs/raw.txt", "Input text file")
	out := flag.String("out", "fixtures/worklogs/worklogs_import.csv", "Output CSV file")
	source := flag.String("source", "legacy-text-log", "Source label")
	flag.Parse()

	fin, err := os.Open(*in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open input: %v\n", err)
		os.Exit(1)
	}
	defer fin.Close()

	fout, err := os.Create(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create output: %v\n", err)
		os.Exit(1)
	}
	defer fout.Close()

	w := csv.NewWriter(fout)
	defer w.Flush()

	if err := w.Write([]string{"Date", "Hours", "Minutes", "Description", "Ticket", "Source"}); err != nil {
		fmt.Fprintf(os.Stderr, "write header: %v\n", err)
		os.Exit(1)
	}

	s := bufio.NewScanner(fin)
	lineNo := 0
	count := 0
	for s.Scan() {
		lineNo++
		raw := strings.TrimSpace(s.Text())
		if raw == "" {
			continue
		}
		m := lineRE.FindStringSubmatch(raw)
		if len(m) != 4 {
			fmt.Fprintf(os.Stderr, "skip line %d (unrecognized format): %s\n", lineNo, raw)
			continue
		}

		dateISO, err := parseDate(m[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip line %d (bad date): %v\n", lineNo, err)
			continue
		}
		hours, minutes, err := parseDurationToHoursAndMinutes(m[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip line %d (bad duration): %v\n", lineNo, err)
			continue
		}

		desc := strings.TrimSpace(m[3])
		ticket := ""
		if tm := ticketRE.FindStringSubmatch(desc); len(tm) == 3 {
			ticket = strings.TrimSpace(tm[1])
			desc = strings.TrimSpace(tm[2])
		}

		r := row{
			Date:        dateISO,
			Hours:       fmt.Sprintf("%.2f", hours),
			Minutes:     strconv.Itoa(minutes),
			Description: desc,
			Ticket:      ticket,
			Source:      *source,
		}
		if err := w.Write([]string{r.Date, r.Hours, r.Minutes, r.Description, r.Ticket, r.Source}); err != nil {
			fmt.Fprintf(os.Stderr, "write row line %d: %v\n", lineNo, err)
			os.Exit(1)
		}
		count++
	}
	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scan input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("wrote %d rows to %s\n", count, *out)
}
