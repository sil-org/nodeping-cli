package lib

import (
	"fmt"
	"sort"

	"github.com/sil-org/nodeping-go-client"
)

type UptimeResults struct {
	CheckLabels []string
	Uptimes     map[string]float32
	StartTime   int64
	EndTime     int64
}

func GetContactGroupIDFromName(contactGroupName string, npClient *nodeping.NodePingClient) (string, error) {
	contactGroups, err := npClient.ListContactGroups()
	if err != nil {
		return "", fmt.Errorf("error retrieving contact groups: %w", err)
	}

	cgID := ""
	for cgKey, cg := range contactGroups {
		if cg.Name == contactGroupName {
			cgID = cgKey
			break
		}
	}

	if cgID == "" {
		return "", fmt.Errorf(`contact group not found with name: "%s"`, contactGroupName)
	}

	return cgID, nil
}

func GetCheckIDsAndLabels(id string, client *nodeping.NodePingClient) ([]string, map[string]string, error) {
	checkIDs := map[string]string{}
	checkLabels := []string{}

	checks, err := client.ListChecks()
	if err != nil {
		return checkLabels, checkIDs, err
	}

	for _, check := range checks {
		// Notifications is a list of maps with the contactGroup ID as keys
		for _, notification := range check.Notifications {
			foundOne := false
			for nKey := range notification {
				if nKey == id {
					checkIDs[check.Label] = check.ID
					checkLabels = append(checkLabels, check.Label)
					foundOne = true
					break
				}
			}
			if foundOne {
				break
			}
		}
	}

	sort.Strings(checkLabels)
	return checkLabels, checkIDs, nil
}

func GetUptimesForChecks(
	checkIDs map[string]string,
	start, end int64,
	npClient *nodeping.NodePingClient,
) map[string]float32 {
	uptimes := map[string]float32{}

	for _, checkID := range checkIDs {
		nextUptime, err := npClient.GetUptime(checkID, start, end)
		if err != nil {
			fmt.Printf("Error getting uptime for check ID %s.\n%s\n", checkID, err.Error())
			continue
		}
		uptimes[checkID] = nextUptime["total"].Uptime
	}

	return uptimes
}

func GetUptimesForContactGroup(token, group string, period Period) (UptimeResults, error) {
	var emptyResults UptimeResults
	npClient, err := nodeping.New(nodeping.ClientConfig{Token: token})
	if err != nil {
		return emptyResults, fmt.Errorf("error initializing cli: %w", err)
	}

	cgID, err := GetContactGroupIDFromName(group, npClient)
	if err != nil {
		return emptyResults, err
	}

	checkLabels, checkIDs, err := GetCheckIDsAndLabels(cgID, npClient)
	if err != nil {
		return emptyResults, err
	}

	start := period.From.Unix() * 1000
	end := period.To.Unix() * 1000

	uptimes := GetUptimesForChecks(checkIDs, start, end, npClient)
	uptimesByLabel := map[string]float32{}

	for _, label := range checkLabels {
		uptimesByLabel[label] = uptimes[checkIDs[label]]
	}

	results := UptimeResults{
		CheckLabels: checkLabels,
		Uptimes:     uptimesByLabel,
		StartTime:   start / 1000,
		EndTime:     end / 1000,
	}

	return results, nil
}
