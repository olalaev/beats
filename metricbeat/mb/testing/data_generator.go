package testing

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
)

var (
	// Use `go test -data` to update files.
	dataFlag = flag.Bool("data", false, "Write updated data.json files")
)

// WriteEvent fetches a single event writes the output to a ./_meta/data.json
// file.
func WriteEvent(f mb.EventFetcher, t testing.TB) error {
	if !*dataFlag {
		t.Skip("skip data generation tests")
	}

	event, err := f.Fetch()
	if err != nil {
		return err
	}

	fullEvent := CreateFullEvent(f, event)
	WriteEventToDataJSON(t, fullEvent)
	return nil
}

// WriteEvents fetches events and writes the first event to a ./_meta/data.json
// file.
func WriteEvents(f mb.EventsFetcher, t testing.TB) error {
	if !*dataFlag {
		t.Skip("skip data generation tests")
	}

	events, err := f.Fetch()
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return fmt.Errorf("no events were generated")
	}

	fullEvent := CreateFullEvent(f, events[0])
	WriteEventToDataJSON(t, fullEvent)
	return nil
}

// CreateFullEvent builds a full event given the data generated by a MetricSet.
// This simulates the output of Metricbeat as if it were
// 2016-05-23T08:05:34.853Z and the hostname is host.example.com.
func CreateFullEvent(ms mb.MetricSet, metricSetData common.MapStr) beat.Event {
	startTime, err := time.Parse(time.RFC3339Nano, "2017-10-12T08:05:34.853Z")
	if err != nil {
		panic(err)
	}

	mbEvent := mb.TransformMapStrToEvent(ms.Module().Name(), metricSetData, nil)
	mbEvent.Namespace = ms.Registration().Namespace
	mbEvent.Timestamp = startTime
	mbEvent.Took = 115 * time.Microsecond
	mbEvent.Host = ms.Host()

	fullEvent := mbEvent.BeatEvent(ms.Module().Name(), ms.Name(), mb.AddMetricSetInfo)

	fullEvent.Fields["beat"] = common.MapStr{
		"name":     "host.example.com",
		"hostname": "host.example.com",
	}

	return fullEvent
}

// WriteEventToDataJSON writes the given event as "pretty" JSON to
// a ./_meta/data.json file. If the -data CLI flag is unset or false then the
// method is a no-op.
func WriteEventToDataJSON(t testing.TB, fullEvent beat.Event) {
	if !*dataFlag {
		return
	}

	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	fields := fullEvent.Fields
	fields["@timestamp"] = fullEvent.Timestamp

	output, err := json.MarshalIndent(&fullEvent.Fields, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	if err = ioutil.WriteFile(path+"/_meta/data.json", output, 0644); err != nil {
		t.Fatal(err)
	}
}
