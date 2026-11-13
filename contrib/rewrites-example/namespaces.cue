// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keto

namespaces: {
	User: relations: manager: ["User"]
	Group: relations: members: ["User", "Group"]
	Folder: {
		relations: {
			parents: ["File", "Folder"]
			viewers: ["Group#members"]
		}
		permissions: view: or: ["viewers", "parents->view"]
	}
	File: {
		relations: {
			parents: ["File", "Folder"]
			viewers: ["User", "Group#members"]
			owners: ["User", "Group#members"]
		}
		permissions: {
			view: or: ["parents->view", "viewers", "owners"]
			edit: "owners"
		}
	}
}
