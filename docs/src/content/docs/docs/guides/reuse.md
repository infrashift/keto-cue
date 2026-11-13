---
title: Reuse & Comprehensions
description: Share rules between namespaces with CUE definitions, list comprehensions and field comprehensions.
---

A schema is CUE, so it can use the CUE language to remove repetition. The top
level of a schema is open. Declare helper values next to `namespaces`.

## Definitions

A definition starts with `#`. It is not part of the schema output, so it can
hold a shared value:

```cue
#editors: ["owners", "parents->edit"]

namespaces: {
	User: {}
	Folder: {
		relations: {owners: ["User"], parents: ["Folder"]}
		permissions: edit: or: #editors
	}
}
```

## List comprehensions

Use a comprehension to add the values of a list to another list:

```cue
#editors: ["owners", "parents->edit"]

namespaces: {
	User: {}
	Folder: {
		relations: {owners: ["User"], viewers: ["User"], parents: ["Folder"]}
		permissions: {
			edit: or: #editors
			view: or: ["viewers", "parents->view", for e in #editors {e}]
		}
	}
}
```

## Field comprehensions

A field comprehension makes one field for each value of a list. This schema
makes the permissions `read`, `write` and `delete`, and each one follows the
parent:

```cue
#actions: ["read", "write", "delete"]

namespaces: {
	User: {}
	Doc: {
		relations: {owners: ["User"], parents: ["Doc"]}
		permissions: {
			for a in #actions {(a): or: ["owners", "parents->\(a)"]}
		}
	}
}
```

`(a)` makes a field name from the value of `a`. `"\(a)"` puts the value in a
string.

## Shared namespaces

Embed a definition in a namespace to give several namespaces the same
relations:

```cue
#Owned: relations: {
	owners: ["User"]
	parents: ["Folder"]
}

namespaces: {
	User: {}
	Folder: relations: edit: ["User"]
	Doc: {
		#Owned
		permissions: edit: or: ["owners", "parents->edit"]
	}
	Sheet: {
		#Owned
		permissions: edit: "owners"
	}
}
```

A definition is closed. `#Owned & {permissions: ...}` is an error, because
`#Owned` does not declare `permissions`. An embedded definition accepts the
other fields of the struct.

## Hidden fields

A field that starts with `_` is hidden. It can hold a parameter of a
definition in the same package:

```cue
#via: {
	_rel: string
	or: ["owners", "parents->\(_rel)"]
}

namespaces: {
	User: {}
	Doc: {
		relations: {owners: ["User"], parents: ["Doc"]}
		permissions: view: #via & {_rel: "view"}
	}
}
```

## Standard library

A schema can import the CUE standard library, for example `strings` or
`list`. A schema cannot import other packages or modules. keto-cue builds the
schema in memory, so a schema cannot read local files or fetch modules.

```cue
import "list"

#readers: ["viewers"]
#writers: ["owners", "parents->edit"]

namespaces: {
	User: {}
	Doc: {
		relations: {owners: ["User"], viewers: ["User"], parents: ["Doc"]}
		permissions: {
			edit: or: #writers
			view: or: list.Concat([#readers, #writers])
		}
	}
}
```

## Errors in generated values

An error in a generated value points at the source of the value. When a
comprehension makes the same error for each entry, the compiler reports it
once.
