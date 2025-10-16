# nodeping-cli
CLI for NodePing

# Usage

## CLI
```
Usage:
  nodeping-cli uptime [flags]

Flags:
  -g, --contact-group string   Name of the Contact Group to retrieve uptime data for.
  -h, --help                   help for uptime
  -p, --period period          Name of the time period to get uptime values for ...
                               [LastMonth LastYear ThisMonth ThisYear Today]
                               (default Period: Today)
```

## Library
The repo can also be imported and used in go packages using `go get github.com/sil-org/nodeping-cli`

Example usage:

```go
package example

import "github.com/sil-org/nodeping-cli"

func getUptime(token, group, periodName string) error {
	period, err := nodeping.GetPeriod(periodName)
	if err != nil {
		return err
	}

	uptimeResults, err := nodeping.GetUptimesForContactGroup(token, group, period)
	return err
}
```
