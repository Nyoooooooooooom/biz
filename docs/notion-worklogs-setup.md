# Notion Worklogs Setup (Invoice-Ready)

This guide replaces plaintext worklog notes with structured Notion tables that feed your invoice process.

## 1) Create Notion Databases

Create two databases in Notion:

1. `Clients`
2. `Worklogs`

Optional third database if you want invoice grouping in Notion:

3. `Invoices`

## 2) `Clients` Database Schema

Recommended properties:

- `Name` (Title)
- `Billing Rate` (Number)
- `Currency` (Select: USD, EUR, GBP)
- `Tax Region` (Select/Text)
- `Active` (Checkbox)

## 3) `Worklogs` Database Schema

Recommended properties:

- `Title` (Title) — short summary (can be first 60 chars of description)
- `Date` (Date)
- `Hours` (Number)
- `Minutes` (Number)
- `Description` (Rich text)
- `Ticket` (Rich text)
- `Client` (Relation -> `Clients`)
- `Status` (Select: Draft, Approved, Ready to Invoice, Invoiced)
- `Source` (Select/Text)
- `Billable` (Checkbox)
- `Rate Override` (Number, optional)
- `Effective Rate` (Formula)
- `Amount` (Formula)

`Effective Rate` formula:

```notion
if(empty(prop("Rate Override")), prop("Client").first().prop("Billing Rate"), prop("Rate Override"))
```

`Amount` formula:

```notion
round(prop("Hours") * prop("Effective Rate") * 100) / 100
```

## 4) Import Existing Plaintext Logs

Use the provided converter in this repo:

```bash
cd /Users/main-user/Documents/Projects/biz
go run ./scripts/worklogs_to_csv.go \
  -in ./fixtures/worklogs/raw.txt \
  -out ./fixtures/worklogs/worklogs_import.csv
```

Then in Notion:

1. Open `Worklogs` database
2. `...` menu -> `Merge with CSV`
3. Select `fixtures/worklogs/worklogs_import.csv`
4. Map CSV columns to properties

## 5) Link Worklogs to Invoice Flow

Use this Notion filter for invoice selection:

- `Status` is `Ready to Invoice`
- `Billable` is checked
- `Date` is within target billing period

Your `biz invoice` flow can then query this subset and aggregate into line items.

## 6) Connect `biz` CLI to Notion

Set config in `/Users/main-user/Documents/Projects/biz/config.yaml`:

```yaml
notion:
  token: "<YOUR_NOTION_INTEGRATION_TOKEN>"
  invoice_db_id: "<INVOICE_DB_ID>"
invoice:
  source: notion
```

Also share the target Notion databases with your integration.

## 7) Quick Connectivity Check

```bash
cd /Users/main-user/Documents/Projects/biz
./bin/biz --config ./config.yaml invoice list --status ready --json
```

If connected, you should get `"code": "OK"` and live Notion records.
