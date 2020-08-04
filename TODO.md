# Todo

Misc notes for things upcoming, done, in progress.

## In Progress

- Make pagination size configurable

## Done

- Decided not to support positional arguments for the search command. Make circle back
  around eventually. The added logic gave no immediate benefit aside from ldapsearch
  parity. This means attributes and filter have to be passed as flags `--attribute`
  and `--filter`. Multiple attribute flags must be passed to generate a list of attrs.
- Added supported for Pagination control. Decided not to implement other controls atm.
