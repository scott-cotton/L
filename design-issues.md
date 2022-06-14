# Design Issues

This document contains the major design issues and the decisions surrounding
them for L.  It is in the following format.  Subsections are issues.
Sub-sub-sections, except the last in a section,  are decisions made in support
of the associated issue.  The last sub-subsection is a status.

## Structured and Traditional logging interoperability.

Structured logging is pretty common, but makes no sense for interactive
commandlines or imposing on projects using one line log entries.

### Decision: handle this with formatters like zerolog

### Status 

- The log generation interface still imposes too much
structure for convenient usage in non-structured scenario.
- There is no mapping from unstructured to structured, only
the other way around.
- The table formatter suffices for rudimentary use.

## Support Runtime configurability.

### Decision: use an http endpoint

This is the only practical solution for ease of use.

### Decision: use RPC style

Document style is too complicated, what else is there?

### Decision: use JSONRPC

JSONRPC reduces the dependency burden of clients w.r.t.
rpc frameworks accross languages.  Roll your own in 
a few tens of lines of code is easier than managing 
an extra dependency and build-time configuration for users.

### Decision: use a simple set of labels

### Decision: provde a simple CLI

### Status
- seems to work so far.

## Context based logging

For structured logging, L should support contextual logging fields,
as this is a common pattern that relieves writing a lot of 
repeated code in servers.

### Decision

House an L object by context.

### Decision

Clone the object on return so it will not be modified if 
retrieved multiple times.

### Status

We can call 
```
obj := L.FromContext(ctx)
defer obj.{Log,Fatal}()
// or
defer func() {
	if err != nil {
		obj.Err(err)
	}
	obj.Log()
}()
```

For non-flow control levels, we use
```
traceObj := L.FromContextWith(ctx, Ltrace)
defer traceObj.Log() // now obj.Log will go to trace.
```

This seems to work.

## Levelled logging

### Decision: use configuration instead of fixed levels

There are too many different ways to do levels to attach them to
method/function names.

### Decision: use methods (.Fatal,.Log,.Err) for universal flow control

Flow control variation is always a concern for a package,
even when it corresponds to levels.  Levels which do not carry
with them implicit flow control semantics in any package (.Fatal,.Err)

For example, even a very high _trace_ level may need to handle errors,
because they are part of Go's flow control standard practices.  Why
not let it do Fatal as well?

### Decision: Make sure message generation is fast 

L Objects short circuit on nil, so this should work.

### Status:

- Short circuiting seems to work.
- Flow control is nice.
- Setting up different levels outside of flow control is reasonably concise for
  several use cases.





