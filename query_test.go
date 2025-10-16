package nodeping

import (
	"testing"

	"github.com/sil-org/nodeping-go-client"
)

func TestGetContactGroupIDFromName(t *testing.T) {
	npClient, _ := nodeping.New(nodeping.ClientConfig{Token: "TestToken"})

	npClient.MockResults = `
{
 "2018090614528ABCD-A-AAAA5":
    {"type":"group","customer_id":"2018090614528ABCD","name":"CGList1","members":["AAAA5","BBBB5","CCCC5","DDDD5","EEEE5"]},
 "2018090614528Abcd-B-BBBB5":
    {"type":"group","customer_id":"2018090614528ABCD","name":"CGList2","members":["FFFF5"]}
}
`

	resultsID, err := GetContactGroupIDFromName("CGList2", npClient)
	if err != nil {
		t.Error(err.Error())
		return
	}

	expectedID := "2018090614528Abcd-B-BBBB5"

	if resultsID != expectedID {
		t.Errorf("Wrong contact group ID. Expected %s, but got %s", expectedID, resultsID)
	}
}

func TestGetCheckIDsAndLabels(t *testing.T) {
	npClient, _ := nodeping.New(nodeping.ClientConfig{Token: "TestToken"})

	npClient.MockResults = `
{
  "2018090614528ABCD-MNOPQRST":
  {"_id":"2018090614528ABCD-MNOPQRST","customer_id":"2018090614528ABCD","label":"Example1","interval":5,
    "notifications":[
      {"AAAA5":
        {"schedule":"All","delay":5}}         
    ], 
    "runlocations":false,"type":"HTTP","status":"assigned","modified":1543260813141,"enable":"active","public":false,"dep":false,
    "parameters":
      {"target":"https://example1.org/","ipv6":false,"follow":false,"threshold":30,"sens":2},
    "created":1539715596937,"queue":"aaaaaaaa10","uuid":"aaaaaaa8-aaa4-aaa4-aaa4-aaaaaaaaaa11","firstdown":0,"state":1
  },
  "2018090614528ABCD-NOPQRSTU":
  {"_id":"2018090614528ABCD-NOPQRSTU","customer_id":"2018090614528ABCD","label":"Example2","interval":3,
    "notifications":[
      {"2018090614528ABCD-B-BBBB5":
        {"schedule":"All","delay":10}}
    ],
    "runlocations":false,"type":"HTTP","status":"assigned","modified":1543937160541,"enable":"active","public":false,"dep":false,
    "parameters":
      {"target":"https://example2.org/","ipv6":false,"follow":false,"threshold":30,"sens":2},
     "created":1539715552868,"queue":"bbbbbbbb10","uuid":"bbbbbbb8-bbb4-bbb4-bbb4-bbbbbbbbbb11","firstdown":0,"state":1
  },
  "2018090614528ABCD-OPQRSTUV":
  {"_id":"2018090614528ABCD-OPQRSTUV","customer_id":"2018090614528ABCD","label":"Example3","interval":1,
    "notifications":[
      {"2018090614528ABCD-C-CCCCC6":
        {"schedule":"All","delay":5}}
    ],
    "runlocations":false,"type":"HTTP","status":"assigned","modified":1539715504508,"enable":"active","public":false,"dep":false,
    "parameters":
      {"target":"https://example3.org/home","ipv6":false,"follow":false,"threshold":30,"sens":2},
    "created":1539715504508,"queue":"cccccccc10","uuid":"ccccccc8-ccc4-ccc4-ccc4-cccccccccc11","firstdown":0,"state":1
  },
  "2018090614528ABCD-PQRSTUVW":
  {"_id":"2018090614528ABCD-PQRSTUVW","customer_id":"2018090614528ABCD","label":"Example4","interval":1,
    "notifications":[
      {"2018090614528ABCD-B-BBBB5":
        {"schedule":"All","delay":5}},
      {"EEEE5":
        {"schedule":"All","delay":5}}
    ],
    "runlocations":false,"type":"HTTP","status":"assigned","modified":1543260719724,"enable":"active","public":false,"dep":false,
    "parameters":
       {"target":"https://example4.org/check","ipv6":false,"follow":false,"threshold":30,"sens":2},
    "created":1539715451787,"queue":"dddddddd10","uuid":"ddddddd8-ddd4-ddd4-ddd4-dddddddddd11","firstdown":0,"state":1
  }
}
`
	cgID := "2018090614528ABCD-B-BBBB5"

	checkLabels, checkIDs, _ := GetCheckIDsAndLabels(cgID, npClient)

	expectedLabels := []string{"Example2", "Example4"}
	if len(expectedLabels) != len(checkLabels) ||
		expectedLabels[0] != checkLabels[0] ||
		expectedLabels[1] != checkLabels[1] {

		t.Errorf("Got wrong list of check labels. \nExpected %v\n  But got %v", expectedLabels, checkLabels)
		return
	}

	resultsID := checkIDs[expectedLabels[1]]
	expectedID := "2018090614528ABCD-PQRSTUVW"

	if resultsID != expectedID {
		t.Errorf("Got wrong ID for %s. Expected %s, but got %s", expectedLabels[1], expectedID, resultsID)
	}
}

func TestGetUptimesForChecks(t *testing.T) {
	npClient, _ := nodeping.New(nodeping.ClientConfig{Token: "TestToken"})

	npClient.MockResults = `
{
  "2018-11":{"enabled":2592000000,"down":89390,"uptime":99.010},
  "2018-12":{"enabled":837810368,"down":80892,"uptime":99.012},
  "total":{"enabled":4744902919,"down":253073,"uptime":99.011}
}
`

	checkIDs := map[string]string{
		"check1": "c1ID",
		"check2": "c2ID",
	}
	uptimes := GetUptimesForChecks(checkIDs, 0, 0, npClient)
	expected := map[string]float32{
		"c1ID": 99.011,
		"c2ID": 99.011,
	}

	if len(uptimes) != 2 || uptimes["c1ID"] != expected["c1ID"] || uptimes["c2ID"] != expected["c2ID"] {
		t.Errorf("Got wrong uptime results. \nExpected %+v\n  but got %+v", expected, uptimes)
	}
}
