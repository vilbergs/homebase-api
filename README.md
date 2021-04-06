# Homebase API

Api service for homebase devices and applications

## Running locally

`deployctl run --no-check --libs=ns,fetchevent ./zones.ts`

**Notes:** The `--no-check` flag is required for now there some problems with Deno.readFileSync in the type checks. 