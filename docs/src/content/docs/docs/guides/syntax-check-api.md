---
title: Syntax-Check API
description: Check a CUE schema over HTTP or gRPC, and the limits that protect the server.
---

The syntax-check API compiles one CUE schema and returns its errors. It does
not change the schema of the server.

## HTTP

The API listens on port 4469 by default. Send the schema as the request body:

```shell
curl -X POST --data-binary @schema.cue http://127.0.0.1:4469/opl/syntax/check
```

A valid schema gives an empty result:

```json
{}
```

A schema with errors gives one entry for each error. Lines and columns start
at 1:

```json
{
  "errors": [
    {
      "message": "namespace \"A\" did not declare relation \"nope\"",
      "start": { "Line": 3, "column": 18 },
      "end": { "Line": 3, "column": 24 }
    }
  ]
}
```

The endpoint keeps its OPL name for API compatibility.

## gRPC

`SyntaxService.Check` takes the schema as the `content` field and returns the
same errors. The `--syntax-remote` flag of the CLI sets its address.

## Limits

The API accepts schemas from clients. These limits protect the server:

| Limit | Value | HTTP | gRPC |
|---|---|---|---|
| Request body | 1 MiB | 413 | `ResourceExhausted` |
| Concurrent checks | The number of CPUs. An extra request waits up to 10 seconds. | 429 | `ResourceExhausted` |
| Time of one check | 10 seconds | 408 | `DeadlineExceeded` |

## Security

Go cannot stop a running CUE evaluation. A timed-out check keeps its slot
until the evaluation ends. CUE comprehensions and functions such as
`list.Range` can use much memory within these limits.

Put the syntax-check listener behind access control when untrusted clients
can reach it.
