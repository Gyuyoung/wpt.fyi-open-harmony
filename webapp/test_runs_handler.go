// Copyright 2017 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package webapp

import (
	"net/http"
	"time"

	"github.com/web-platform-tests/wpt.fyi/shared"
)

// testRunsHandler handles GET/POST requests to /test-runs
func testRunsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET is supported.", http.StatusMethodNotAllowed)
		return
	}
	filter, err := parseTestRunsUIFilter(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	RenderTemplate(w, r, "test-runs.html", filter)
}

// parseTestRunsUIFilter parses the standard TestRunFilter, as well as the extra
// pr param.
func parseTestRunsUIFilter(r *http.Request) (filter testRunUIFilter, err error) {
	q := r.URL.Query()
	testRunFilter, err := shared.ParseTestRunFilterParams(q)
	if err != nil {
		return filter, err
	}

	pr, err := shared.ParsePRParam(q)
	if err != nil {
		return filter, err
	}

	// isDefault := testRunFilter.IsDefaultQuery() && pr == nil
	isDefault := false
	if isDefault {
		// Get runs from a week ago, onward, by default.
		aWeekAgo := time.Now().Truncate(time.Hour*24).AddDate(0, 0, -7)
		testRunFilter.From = &aWeekAgo
		testRunFilter = testRunFilter.MasterOnly()
	} else if testRunFilter.MaxCount == nil {
		oneHundred := 100
		testRunFilter.MaxCount = &oneHundred
	}
	filter = convertTestRunUIFilter(testRunFilter)
	if pr != nil {
		filter.PR = pr
	}

	return filter, nil
}
