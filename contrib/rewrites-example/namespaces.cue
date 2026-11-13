// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keto

namespaces: {
	User: relations: manager: types: ["User"]
	Group: relations: members: types: ["User", "Group"]
	Folder: {
		relations: {
			parents: types: ["File", "Folder"]
			viewers: types: ["Group#members"]
		}
		permissions: view: {
			op: "or"
			children: [
				{op: "direct", relation: "viewers"},
				{op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "view"},
			]
		}
	}
	File: {
		relations: {
			parents: types: ["File", "Folder"]
			viewers: types: ["User", "Group#members"]
			owners: types: ["User", "Group#members"]
		}
		permissions: {
			view: {
				op: "or"
				children: [
					{op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "view"},
					{op: "direct", relation: "viewers"},
					{op: "direct", relation: "owners"},
				]
			}
			edit: {op: "direct", relation: "owners"}
		}
	}
}
