// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keto

namespaces: {
	User: {}
	Video: relations: {
		view: types: ["User", "Video#owner"]
		owner: types: ["User"]
	}
}
