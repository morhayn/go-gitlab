//
// Copyright 2021, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestCreateLabel(t *testing.T) {
	mux, client := setup(t)

	mux.HandleFunc("/api/v4/projects/1/labels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"id":1, "name": "MyLabel", "color" : "#11FF22", "priority": 2}`)
	})

	// Create new label
	l := &CreateLabelOptions{
		Name:     Ptr("MyLabel"),
		Color:    Ptr("#11FF22"),
		Priority: Ptr(2),
	}
	label, _, err := client.Labels.CreateLabel("1", l)
	if err != nil {
		t.Fatal(err)
	}
	want := &Label{ID: 1, Name: "MyLabel", Color: "#11FF22", Priority: 2}
	if !reflect.DeepEqual(want, label) {
		t.Errorf("Labels.CreateLabel returned %+v, want %+v", label, want)
	}
}

func TestDeleteLabelbyID(t *testing.T) {
	mux, client := setup(t)

	mux.HandleFunc("/api/v4/projects/1/labels/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	// Delete label
	_, err := client.Labels.DeleteLabel("1", "1", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteLabelbyName(t *testing.T) {
	mux, client := setup(t)

	mux.HandleFunc("/api/v4/projects/1/labels/MyLabel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	// Delete label
	label := &DeleteLabelOptions{
		Name: Ptr("MyLabel"),
	}

	_, err := client.Labels.DeleteLabel("1", "MyLabel", label)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateLabel(t *testing.T) {
	mux, client := setup(t)

	mux.HandleFunc("/api/v4/projects/1/labels/MyLabel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		fmt.Fprint(w, `{"id":1, "name": "New Label", "color" : "#11FF23" , "description":"This is updated label", "priority": 42}`)
	})

	// Update label
	l := &UpdateLabelOptions{
		NewName:     Ptr("New Label"),
		Color:       Ptr("#11FF23"),
		Description: Ptr("This is updated label"),
		Priority:    Ptr(42),
	}

	label, resp, err := client.Labels.UpdateLabel("1", "MyLabel", l)

	if resp == nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	want := &Label{ID: 1, Name: "New Label", Color: "#11FF23", Description: "This is updated label", Priority: 42}

	if !reflect.DeepEqual(want, label) {
		t.Errorf("Labels.UpdateLabel returned %+v, want %+v", label, want)
	}
}

func TestSubscribeToLabel(t *testing.T) {
	mux, client := setup(t)

	mux.HandleFunc("/api/v4/projects/1/labels/5/subscribe", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{  "id" : 5, "name" : "kind/bug", "color" : "#d9534f", "description": "Bug reported by user", "open_issues_count": 1, "closed_issues_count": 0, "open_merge_requests_count": 1, "subscribed": true,"priority": null}`)
	})

	label, _, err := client.Labels.SubscribeToLabel("1", "5")
	if err != nil {
		t.Fatal(err)
	}
	want := &Label{ID: 5, Name: "kind/bug", Color: "#d9534f", Description: "Bug reported by user", OpenIssuesCount: 1, ClosedIssuesCount: 0, OpenMergeRequestsCount: 1, Subscribed: true}
	if !reflect.DeepEqual(want, label) {
		t.Errorf("Labels.SubscribeToLabel returned %+v, want %+v", label, want)
	}
}

func TestUnsubscribeFromLabel(t *testing.T) {
	mux, client := setup(t)

	mux.HandleFunc("/api/v4/projects/1/labels/5/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.Labels.UnsubscribeFromLabel("1", "5")
	if err != nil {
		t.Fatal(err)
	}
}

func TestListLabels(t *testing.T) {
	mux, client := setup(t)

	mux.HandleFunc("/api/v4/projects/1/labels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{  "id" : 5, "name" : "kind/bug", "color" : "#d9534f", "description": "Bug reported by user", "open_issues_count": 1, "closed_issues_count": 0, "open_merge_requests_count": 1, "subscribed": true,"priority": null}]`)
	})

	o := &ListLabelsOptions{
		ListOptions: ListOptions{
			Page:    1,
			PerPage: 10,
		},
	}
	label, _, err := client.Labels.ListLabels("1", o)
	if err != nil {
		t.Log(err.Error() == "invalid ID type 1.1, the ID must be an int or a string")
	}
	want := []*Label{{ID: 5, Name: "kind/bug", Color: "#d9534f", Description: "Bug reported by user", OpenIssuesCount: 1, ClosedIssuesCount: 0, OpenMergeRequestsCount: 1, Subscribed: true}}
	if !reflect.DeepEqual(want, label) {
		t.Errorf("Labels.ListLabels returned %+v, want %+v", label, want)
	}
}

func TestGetLabel(t *testing.T) {
	mux, client := setup(t)

	mux.HandleFunc("/api/v4/projects/1/labels/5", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{  "id" : 5, "name" : "kind/bug", "color" : "#d9534f", "description": "Bug reported by user", "open_issues_count": 1, "closed_issues_count": 0, "open_merge_requests_count": 1, "subscribed": true,"priority": null}`)
	})

	label, _, err := client.Labels.GetLabel("1", 5)
	if err != nil {
		t.Log(err)
	}
	want := &Label{ID: 5, Name: "kind/bug", Color: "#d9534f", Description: "Bug reported by user", OpenIssuesCount: 1, ClosedIssuesCount: 0, OpenMergeRequestsCount: 1, Subscribed: true}
	if !reflect.DeepEqual(want, label) {
		t.Errorf("Labels.GetLabel returned %+v, want %+v", label, want)
	}
}
