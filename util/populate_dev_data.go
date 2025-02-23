// Copyright 2017 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"

	"github.com/web-platform-tests/wpt.fyi/shared"
	"github.com/web-platform-tests/wpt.fyi/shared/metrics"
)

var (
	project        = flag.String("project", "", "project ID used to connect to Datastore")
	datastoreHost  = flag.String("datastore_host", "", "Cloud Datastore emulator host")
	localHost      = flag.String("local_host", "localhost:8080", "local dev_appserver.py webapp host")
	remoteHost     = flag.String("remote_host", "wpt.fyi", "wpt.fyi host to fetch prod runs from")
	numRemoteRuns  = flag.Int("num_remote_runs", 10, "number of remote runs to copy from host to local environment")
	staticRuns     = flag.Bool("static_runs", false, "Include runs in the /static dir")
	remoteRuns     = flag.Bool("remote_runs", true, "Include copies of remote runs")
	seenTestRunIDs = mapset.NewSet()
	labels         = flag.String("labels", "", "Labels for which to fetch runs")
)

// populate_dev_data.go populates a local running webapp instance with some
// of the latest production entities, so that there's data to view.
//
// Usage (from util/):
// go run populate_dev_data.go
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	if *project != "" {
		os.Setenv("DATASTORE_PROJECT_ID", *project)
	}
	if *datastoreHost != "" {
		os.Setenv("DATASTORE_EMULATOR_HOST", *datastoreHost)
	}

	ctx := context.Background()
	shared.Clients.Init(ctx)
	defer shared.Clients.Close()

	log.Printf("Adding dev data to local emulator...")

	emptySecretToken := &shared.Token{}
	enabledFlag := &shared.Flag{Enabled: true}
	staticDataTime := time.Now()

	// Follow pattern established in run/*.py data collection code.
	const staticRunSHA = "24278ab61781de72ed363b866ae6b50b86822b27"
	summaryURLFmtString := "http://%s/static/%s/%s"
	chrome := shared.TestRun{
		ProductAtRevision: shared.ProductAtRevision{
			Product: shared.Product{
				BrowserName:    "chrome",
				BrowserVersion: "74.0",
				OSName:         "linux",
				OSVersion:      "3.16",
			},
			FullRevisionHash: staticRunSHA,
			Revision:         staticRunSHA[:10],
		},
		ResultsURL: fmt.Sprintf(summaryURLFmtString, *localHost, staticRunSHA[:10], "chrome[stable]-summary_v2.json.gz"),
		CreatedAt:  staticDataTime,
		TimeStart:  staticDataTime,
		Labels:     []string{"chrome", shared.StableLabel},
	}
	chromeExp := chrome
	chromeExp.BrowserVersion = "76.0"
	chromeExp.Labels = []string{"chrome", shared.ExperimentalLabel}
	chromeExp.ResultsURL = strings.Replace(chrome.ResultsURL, "[stable]", "[experimental]", -1)

	chrome_android := shared.TestRun{
		ProductAtRevision: shared.ProductAtRevision{
			Product: shared.Product{
				BrowserName:    "chrome_android",
				BrowserVersion: "74.0",
				OSName:         "android",
				OSVersion:      "9",
			},
			FullRevisionHash: staticRunSHA,
			Revision:         staticRunSHA[:10],
		},
		ResultsURL: fmt.Sprintf(summaryURLFmtString, *localHost, staticRunSHA[:10], "chrome_android[stable]-summary_v2.json.gz"),
		CreatedAt:  staticDataTime,
		TimeStart:  staticDataTime,
		Labels:     []string{"chrome_android", shared.StableLabel},
	}
	chrome_androidExp := chrome
	chrome_androidExp.BrowserVersion = "76.0"
	chrome_androidExp.Labels = []string{"chrome", shared.ExperimentalLabel}
	chrome_androidExp.ResultsURL = strings.Replace(chrome_android.ResultsURL, "[stable]", "[experimental]", -1)

	edge := chrome
	edge.BrowserName = "edge"
	edge.BrowserVersion = "18"
	edge.OSName = "windows"
	edge.OSVersion = "10"
	edge.ResultsURL = fmt.Sprintf(summaryURLFmtString, *localHost, staticRunSHA[:10], "edge[stable]-summary_v2.json.gz")
	edge.Labels = []string{"edge", shared.StableLabel}

	edgeExp := edge
	edgeExp.BrowserVersion = "20"
	edgeExp.ResultsURL = strings.Replace(edge.ResultsURL, "[stable]", "[experimental]", -1)
	edgeExp.Labels = []string{"edge", shared.ExperimentalLabel}

	firefox := chrome
	firefox.BrowserName = "firefox"
	firefox.BrowserVersion = "66"
	firefox.ResultsURL = fmt.Sprintf(summaryURLFmtString, *localHost, staticRunSHA[:10], "firefox[stable]-summary_v2.json.gz")
	firefox.Labels = []string{"firefox", shared.StableLabel}
	firefoxExp := firefox
	firefoxExp.BrowserVersion = "68.0"
	firefoxExp.Labels = []string{"firefox", shared.ExperimentalLabel}
	firefoxExp.ResultsURL = strings.Replace(firefox.ResultsURL, "[stable]", "[experimental]", -1)

	safari := chrome
	safari.BrowserName = "safari"
	safari.BrowserVersion = "12.1"
	safari.OSName = "mac"
	safari.OSName = "10.13"
	safari.ResultsURL = fmt.Sprintf(summaryURLFmtString, *localHost, staticRunSHA[:10], "safari[stable]-summary_v2.json.gz")
	safari.Labels = []string{"safari", shared.StableLabel}
	safariExp := safari
	safariExp.BrowserVersion = "81 preview"
	safariExp.Labels = []string{"safari", shared.ExperimentalLabel}
	safariExp.ResultsURL = strings.Replace(safari.ResultsURL, "[stable]", "[experimental]", -1)

	huawei_browser := chrome
	huawei_browser.BrowserName = "huawei_browser"
	huawei_browser.BrowserVersion = "Beta"
	huawei_browser.OSName = "openharmony"
	huawei_browser.OSName = "3.2.3"
	huawei_browser.ResultsURL = fmt.Sprintf(summaryURLFmtString, *localHost, staticRunSHA[:10], "huawei_browser[stable]-summary_v2.json.gz")
	huawei_browser.Labels = []string{"huawei browser", shared.StableLabel}
	huawei_browserExp := huawei_browser
	huawei_browserExp.BrowserVersion = "Beta"
	huawei_browserExp.Labels = []string{"huawei_browser", shared.ExperimentalLabel}
	huawei_browserExp.ResultsURL = strings.Replace(huawei_browser.ResultsURL, "[stable]", "[experimental]", -1)


	staticTestRuns := shared.TestRuns{
		chrome,
		chromeExp,
		chrome_android,
		chrome_androidExp,
		firefox,
		firefoxExp,
		safari,
		safariExp,
		huawei_browser,
		huawei_browserExp,
		edge,
		edgeExp,
	}
	labelRuns(staticTestRuns, "test", "static", shared.MasterLabel)

	timeZero := time.Unix(0, 0)
	// Follow pattern established in metrics/run/*.go data collection code.
	// Use unzipped JSON for local dev.
	const metricsURLFmtString = "/static/wptd-metrics/0-0/%s.json"
	staticTestRunMetadata := make([]interface{}, len(staticTestRuns))
	for i := range staticTestRuns {
		staticTestRunMetadata[i] = &staticTestRuns[i]
	}
	passRateMetadata := metrics.PassRateMetadata{
		TestRunsMetadata: metrics.TestRunsMetadata{
			StartTime: timeZero,
			EndTime:   timeZero,
			DataURL:   fmt.Sprintf(metricsURLFmtString, "pass-rates"),
		},
	}

	testRunKindName := "TestRun"
	passRateMetadataKindName := metrics.GetDatastoreKindName(
		metrics.PassRateMetadata{})

	log.Print("Adding local (empty) secrets...")
	store := shared.NewAppEngineDatastore(ctx, false)
	addSecretToken(store, "upload-token", emptySecretToken)
	addSecretToken(store, "github-wpt-fyi-bot-token", emptySecretToken)
	addSecretToken(store, "github-oauth-client-id", emptySecretToken)
	addSecretToken(store, "github-oauth-client-secret", emptySecretToken)
	addSecretToken(store, "secure-cookie-hashkey", &shared.Token{
		Secret: "a-very-secret-sixty-four-bytes!!a-very-secret-sixty-four-bytes!!",
	})
	addSecretToken(store, "secure-cookie-blockkey", &shared.Token{
		Secret: "a-very-secret-thirty-two-bytes!!",
	})

	log.Print("Adding flag defaults...")
	addFlag(store, "queryBuilder", enabledFlag)
	addFlag(store, "diffFilter", enabledFlag)
	addFlag(store, "diffFromAPI", enabledFlag)
	addFlag(store, "structuredQueries", enabledFlag)
	addFlag(store, "diffRenames", enabledFlag)
	addFlag(store, "paginationTokens", enabledFlag)

	log.Print("Adding uploader \"test\"...")
	addData(store, "Uploader", []interface{}{
		&shared.Uploader{Username: "test", Password: "123"},
	})

	if *staticRuns {
		log.Print("Adding local mock data (static/)...")
		for i, key := range addData(store, testRunKindName, staticTestRunMetadata) {
			staticTestRuns[i].ID = key.IntID()
		}
		stableRuns := shared.TestRuns{}
		defaultRuns := shared.TestRuns{}
		for _, run := range staticTestRuns {
			labels := run.LabelsSet()
			if labels.Contains(shared.StableLabel) {
				stableRuns = append(stableRuns, run)
			} else if labels.Contains("huawei_browser") || labels.Contains(shared.ExperimentalLabel) {
				defaultRuns = append(defaultRuns, run)
			}
		}
		stableInterop := passRateMetadata
		stableInterop.TestRunIDs = stableRuns.GetTestRunIDs()
		defaultInterop := passRateMetadata
		defaultInterop.TestRunIDs = defaultRuns.GetTestRunIDs()
		addData(store, passRateMetadataKindName, []interface{}{
			&stableInterop,
			&defaultInterop,
		})
	}

	if *remoteRuns {
		log.Print("Adding latest production TestRun data...")
		extraLabels := mapset.NewSet()
		if labels != nil {
			for _, s := range strings.Split(*labels, ",") {
				if s != "" {
					extraLabels.Add(s)
				}
			}
		}
		filters := shared.TestRunFilter{
			Labels:   extraLabels.Union(mapset.NewSetWith(shared.StableLabel)),
			MaxCount: numRemoteRuns,
		}
		copyProdRuns(store, filters)

		log.Print("Adding latest master TestRun data...")
		filters.Labels = extraLabels.Union(mapset.NewSetWith(shared.MasterLabel))
		copyProdRuns(store, filters)

		log.Print("Adding latest experimental TestRun data...")
		filters.Labels = extraLabels.Union(mapset.NewSetWith(shared.ExperimentalLabel))
		copyProdRuns(store, filters)

		log.Print("Adding latest beta TestRun data...")
		filters.Labels = extraLabels.Union(mapset.NewSetWith(shared.BetaLabel))
		copyProdRuns(store, filters)

		log.Print("Adding latest aligned Chrome/Chrome Android/Firefox/Safari/Huawei Browser experimental data...")
		filters.Labels = extraLabels.Union(mapset.NewSet(shared.MasterLabel))
		filters.Products, _ = shared.ParseProductSpecs("chrome[experimental]", "chrome_android[experimental]", "firefox[experimental]", "safari[experimental]","huawei_browser[experimental]")
		copyProdRuns(store, filters)

		log.Printf("Successfully copied a total of %v distinct TestRuns", seenTestRunIDs.Cardinality())

		log.Print("Adding latest production PendingTestRun...")
		copyProdPendingRuns(store, *numRemoteRuns)
	}

	log.Print("Adding test history data...")
	addFakeHistoryData(store)
}

func copyProdRuns(store shared.Datastore, filters shared.TestRunFilter) {
	for _, aligned := range []bool{false, true} {
		if aligned {
			filters.Aligned = &aligned
		}
		prodTestRuns, err := shared.FetchRuns(*remoteHost, filters)
		if err != nil {
			log.Print(err)
			continue
		}
		labelRuns(prodTestRuns, "prod")

		latestProductionTestRunMetadata := make([]interface{}, 0, len(prodTestRuns))
		for i := range prodTestRuns {
			if !seenTestRunIDs.Contains(prodTestRuns[i].ID) {
				seenTestRunIDs.Add(prodTestRuns[i].ID)
				latestProductionTestRunMetadata = append(latestProductionTestRunMetadata, &prodTestRuns[i])
			}
		}
		addData(store, "TestRun", latestProductionTestRunMetadata)
	}
}

func copyProdPendingRuns(store shared.Datastore, numRuns int) {
	pendingRuns, err := FetchPendingRuns(*remoteHost)
	if err != nil {
		log.Fatalf("Failed to fetch pending runs: %s", err.Error())
	}
	var castRuns []interface{}
	for i := range pendingRuns {
		castRuns = append(castRuns, &pendingRuns[i])
	}
	addData(store, "PendingTestRun", castRuns)
}

func labelRuns(runs []shared.TestRun, labels ...string) {
	for i := range runs {
		for _, label := range labels {
			runs[i].Labels = append(runs[i].Labels, label)
		}
	}
}

func addSecretToken(store shared.Datastore, id string, data interface{}) {
	key := store.NewNameKey("Token", id)
	if _, err := store.Put(key, data); err != nil {
		log.Fatalf("Failed to add %s secret: %s", id, err.Error())
	}
	log.Printf("Added %s secret", id)
}

func addFlag(store shared.Datastore, id string, data interface{}) {
	key := store.NewNameKey("Flag", id)
	if _, err := store.Put(key, data); err != nil {
		log.Fatalf("Failed to add %s flag: %s", id, err.Error())
	}
	log.Printf("Added %s flag", id)
}

func addData(store shared.Datastore, kindName string, data []interface{}) (keys []shared.Key) {
	keys = make([]shared.Key, len(data))
	for i := range data {
		keys[i] = store.NewIncompleteKey(kindName)
	}
	var err error
	if keys, err = store.PutMulti(keys, data); err != nil {
		log.Fatalf("Failed to add %s entities: %s", kindName, err.Error())
	}
	log.Printf("Added %v %s entities", len(data), kindName)
	return keys
}

// FetchPendingRuns fetches recent PendingTestRuns.
func FetchPendingRuns(wptdHost string) ([]shared.PendingTestRun, error) {
	url := "https://" + wptdHost + "/api/status"
	var pendingRuns []shared.PendingTestRun
	err := shared.FetchJSON(url, &pendingRuns)
	return pendingRuns, err
}

// TODO: Import real data here when staging is populated with real data
func addFakeHistoryData(store shared.Datastore) {
	// browser_name,browser_version,date,test_name,subtest_name,status
	devData := []map[string]string{
		{
			"run_id":       "5074677897101312",
			"date":         "2022-06-02T06:02:55.000Z",
			"test_name":    "example test name",
			"subtest_name": "",
			"status":       "OK",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2023-02-21T03:08:15.000Z",
			"test_name":    "example test name",
			"subtest_name": "",
			"status":       "TIMEOUT",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2023-03-30T20:27:32.000Z",
			"test_name":    "example test name",
			"subtest_name": "",
			"status":       "OK",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2022-06-02T06:02:55.000Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_1",
			"status":       "PASS",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2022-08-14T06:02:55.000Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_1",
			"status":       "FAIL",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2023-02-21T03:08:15.000Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_1",
			"status":       "NOTRUN",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2023-03-30T20:27:32.611Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_1",
			"status":       "FAIL",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2023-06-19T20:54:12.611Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_1",
			"status":       "PASS",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2022-06-02T06:02:55.000Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_2",
			"status":       "TIMEOUT",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2022-09-25T23:49:35.000Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_2",
			"status":       "PASS",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2023-02-21T03:08:15.000Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_2",
			"status":       "NOTRUN",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2023-03-30T20:27:32.611Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_2",
			"status":       "PASS",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2022-06-02T06:02:55.000Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_3",
			"status":       "PASS",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2023-02-21T03:08:15.000Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_3",
			"status":       "NOTRUN",
		},
		{
			"run_id":       "5074677897101312",
			"date":         "2023-03-30T20:27:32.611Z",
			"test_name":    "example test name",
			"subtest_name": "subtest_name_3",
			"status":       "PASS",
		},
	}

	browserMetadata := []map[string]string{
		{
			"browser": "chrome",
		},
		{
			"browser": "edge",
		},

		{
			"browser": "firefox",
		},

		{
			"browser": "safari",
		},
	}

	browserEntries := make([]interface{}, 0, len(devData))
	for _, metadata := range browserMetadata {
		for _, entry := range devData {
			testHistoryEntry := shared.TestHistoryEntry{
				BrowserName: metadata["browser"],
				RunID:       entry["run_id"],
				Date:        entry["date"],
				TestName:    entry["test_name"],
				SubtestName: entry["subtest_name"],
				Status:      entry["status"],
			}
			browserEntries = append(browserEntries, &testHistoryEntry)
		}
	}
	addData(store, "TestHistoryEntry", browserEntries)
}
