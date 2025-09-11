package lib

import (
	"fmt"
	"github.com/sil-org/nodeping-go-client"
	"os"
	"sort"
)

type UptimeResults struct {
	CheckLabels []string
	Uptimes map[string]float32
	StartTime int64
	EndTime int64
}

func GetContactGroupIDFromName(contactGroupName string, npClient *nodeping.NodePingClient) (string, error) {

	contactGroups, err := npClient.ListContactGroups()

	if err != nil {
		return "", err
	}

	if err != nil {
		fmt.Printf("Error retrieving contact groups: \n%s\n", err.Error())
		os.Exit(1)
	}

	cgID := ""
	for cgKey, cg := range contactGroups{
		if cg.Name == contactGroupName {
			cgID = cgKey
			break
		}
	}

	if cgID == "" {
		return "", fmt.Errorf("Could not find contact group with name \"%s\"\n", contactGroupName)
	}

	return cgID, nil
}

func GetCheckIDsAndLabels(
	contactGroupID string,
	npClient *nodeping.NodePingClient,
) ([]string, map[string]string, error) {

	checkIDs := map[string]string{}
	checkLabels := []string{}

	checks, err := npClient.ListChecks()

	if err != nil {
		return checkLabels, checkIDs, err
	}

	//fmt.Printf("First Check:\n%+v\n", checks[0])


	for _, check := range checks {
		// Notifications is a list of maps with the contactGroup ID as keys
		for _, notfctn := range check.Notifications {
			foundOne := false
			for nKey := range notfctn {
				if nKey == contactGroupID {
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
		//fmt.Printf("Got Uptime Response ...\n   %+v\n", nextUptime)
	}

	return uptimes
}


func GetUptimesForContactGroup(
	nodepingToken, contactGroupName, period string,
) (UptimeResults, error) {
	npClient, err := nodeping.New(nodeping.ClientConfig{Token: nodepingToken})
	emptyResults := UptimeResults{}

	if err != nil {
		err = fmt.Errorf("Error initializing cli: %s", err.Error())
		return emptyResults, err
	}

	cgID, err := GetContactGroupIDFromName(contactGroupName, npClient)

	if err != nil {
		return emptyResults, err
	}

	checkLabels, checkIDs, err := GetCheckIDsAndLabels(cgID, npClient)

	if err != nil {
		return emptyResults, err
	}

	start := int64(0)
	end := int64(0)

	if period != "" {
		periodObject := GetPeriodByName(period, 0)
		start = periodObject.From * 1000
		end = periodObject.To * 1000
	}

	uptimes := GetUptimesForChecks(checkIDs, start, end, npClient)
	uptimesByLabel := map[string]float32{}

	for _, label := range checkLabels {
		uptimesByLabel[label] = uptimes[checkIDs[label]]
	}

	results := UptimeResults{
		CheckLabels: checkLabels,
		Uptimes: uptimesByLabel,
		StartTime: start/1000,
		EndTime: end/1000,
	}

	return results, nil
}
