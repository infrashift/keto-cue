// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keto

namespaces: {
	User: {}
	ApiKey: {}
	File: {
		relations: viewer: ["User", "ApiKey"]
		permissions: view: "viewer"
	}
}
