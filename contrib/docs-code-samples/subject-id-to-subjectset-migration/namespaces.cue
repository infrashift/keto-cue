// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keto

namespaces: {
	User: {}
	ApiKey: {}
	File: {
		relations: viewer: types: ["User", "ApiKey"]
		permissions: view: {op: "direct", relation: "viewer"}
	}
}
