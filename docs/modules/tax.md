# Tax Module

Path: `internal/tax`

## Purpose
Compute tax amount from subtotal and region.

## Current Behavior
- Uses `tax.rates` map from config.
- Uses `tax.default_region` when invoice region is empty.
- If `tax.required=true`, unknown region is validation error.

## Output
Returns:
- `rate`
- `amount`
- resolved `region`
