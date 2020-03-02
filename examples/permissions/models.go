package main

import (
	"time"

	dga "github.com/emicklei/dgraph-access"
)

type PermissionsInProject struct {
	dga.Node `json:",inline"`
	// Project        Project
	// Identity    CloudIdentity
	Permissions []string `json:"permissions"`
}

type CloudIdentity struct {
	dga.Node       `json:",inline"`
	Group          string `json:"group,omitempty"`
	User           string `json:"user,omitempty"`
	ServiceAccount string `json:"serviceAccount,omitempty"`
}

type Project struct {
	dga.Node `json:",inline"`
	Name     string `json:"project_name"`
}

type Version struct {
	dga.Node `json:",inline"`
	Snapshot time.Time `json:"snapshot"`
}
