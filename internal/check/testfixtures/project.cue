// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keto

namespaces: {
	User: {}
	Project: {
		relations: {
			owner: ["User"]
			developer: ["User"]
		}
		permissions: {
			isOwner: "owner"
			isOwnerOrDeveloper: or: ["owner", "developer"]
			writeCollaborator: "isOwner"
			readCollaborator:  "isOwnerOrDeveloper"
			deleteProject:     "isOwner"
			writeProject:      "isOwnerOrDeveloper"
			readProject:       "isOwnerOrDeveloper"
		}
	}
}
