package main

import (
	"time"

	dga "github.com/emicklei/dgraph-access"
)

type PermissionsInProject struct {
	*dga.Node `json:",inline"`
	// Project        Project
	// ServiceAccount ServiceAccount
	// GroupOrUser    CloudIdentity
	Permissions []string `json:"permissions"`
}

type CloudIdentity struct {
	*dga.Node `json:",inline"`
	Group     string `json:"group"`
	User      string `json:"user"`
}

type Project struct {
	*dga.Node `json:",inline"`
	Name      string `json:"project_name"`
}

type ServiceAccount struct {
	*dga.Node `json:",inline"`
	Name      string `json:"serviceaccount_name"`
}

type Version struct {
	PermissionsInProjects []PermissionsInProject `json:"-"`
	Snapshot              time.Time              `json:"snapshot"`
}
