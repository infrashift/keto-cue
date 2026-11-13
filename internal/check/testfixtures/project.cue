// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keto

namespaces: {
	User: {}
	Project: {
		relations: {
			owner: types: ["User"]
			developer: types: ["User"]
		}
		permissions: {
			isOwner: {op: "direct", relation: "owner"}
			isOwnerOrDeveloper: {
				op: "or"
				children: [
					{op: "direct", relation: "owner"},
					{op: "direct", relation: "developer"},
				]
			}
			writeCollaborator: {op: "direct", relation: "isOwner"}
			readCollaborator: {op: "direct", relation: "isOwnerOrDeveloper"}
			deleteProject: {op: "direct", relation: "isOwner"}
			writeProject: {op: "direct", relation: "isOwnerOrDeveloper"}
			readProject: {op: "direct", relation: "isOwnerOrDeveloper"}
		}
	}
}
