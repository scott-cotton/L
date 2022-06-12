# L

L is a minimalist project-aware structured logger for Go.

[![GoDoc](https://godoc.org/github.com/scott-cotton/L?status.svg)](https://godoc.org/github.com/scott-cotton/L)

## Why?

As observed by [capnslog](https://github.com/coreos/capnslog), logging should
serve the project it uses well: the main entry point should determine the
logging configuration of its imports, and moreover do so in a way that is
runtime configurable.

If projects A and B organise their logging differently, and project C would
like to use A and B and yet manage the logging, L can be used to resolve the
differences between A and B and C's desired logging setup.  In other words, if
you would like your project to be used in other projects with distinct logging
styles, and you find L worthy enough in its current nascent state, then L can do
that.

L permits flexible per-project configuration in a way that can be overriden by
a main entrypoint, viewed, and manipulated via an authenticated
RPC service accessible via an HTTP gateway at runtime.

L provides

- Efficient structured logging, with reasonable syntax for creating nested structures
  incrementally.
- Formatters and a simple formatting interface.
- The ability for the main entry point to determine and set up configuration
  for all Loggers in all packages it depends on.
- The ability for a package to define its own default configuration.
- Support for MiddleWare.
- An HMAC authenticated jsonrpc service exposed via HTTP for viewing (and, coming soon:
manipulating logs). 

## Working with L

To use L in some package, simply import L and call L.New(cfg) with
your desired configuration.  This should normally done at the top level
of each package once.

If you would like to further facilitate the use of your project within
projects with different logging styles/criteria, then simply describe
your logging setup and make this description part of your release process.

Other projects can then use L to set your logging to their conventions.

## Basic structured logging

Structured logging is about providing structured output, which is typically
fairly unstructured actually and refers to anything that can be viewed as a
`map[string]any`, where `any` is any ground type or another `map[string]any` or
a slice of values, i.e. Go's representation of free form json objects.

```go
// log a field {"k": 77}
var myL = L.New(L.NewConfig())

myL.Dict().Field("k", 77).Log()

// add a bunch of stuff in a chain
myL.Dict().Field("k", 3).Field("mini", "kube").Field("who", "me").Log()

// add a tree structure '{"k": 3, "k2": [ "a", 4 ] }'
myL.Dict().Field("k", 3).Set(
	"k2", 
	myL.Array().Str("a").Int(4)).
	Log()

// log a string
myl.Str("hello").Log()

// add any fields you can with defer
func F() {
	obj := myL.Dict()
	defer obj.Log()
	//
	obj.Field("func", "F")
	//
	//
}
```

### Overriding the configuration at the entry point

The entry point can manipulate the configuration for all
imported packages, transitively by calling [`L.ApplyConfig(...)`](https://pkg.go.dev/github.com/scott-cotton/L#ApplyConfig)`:
```
import "github.com/scott-cotton/L"

L.ApplyConfig(MyAppConfig(), &L.ApplyOpts{Recursive: true})
```

`ApplyConfig` can either overwrite the configuration of all
`L` loggers or set specific variables.

More dynamic and fine-grained control is available via 
[`Walk`](https://pkg.go.dev/github.com/scott-cotton/L#Walk).


### Levelled Logging

Traditionally, levelled logging uses an interface where the level is made explicit
by a function call at each log site.  For example, we often see

```
log.Info().Msg("hello")
```

or similar.  Logging can happen at very high frequency in terms of lines of code,
leading to a whole lot ".\<level\>()" repeats.

Another problem with existing levelled logging options is that one must choose
the levels for a whole project, making inter-project dependencies difficult when
different projects use different levels.

L proposes to treat this differently.  L uses labels which can be turned on or
off at the call site and at main entry points, together with middle ware to
have efficient levelled logging. 

The idea is that for a project which wants to define its own levels, the project
defines a project specific level key and associated integer values.

[this test file](https://github.com/scott-cotton/L/blob/main/level_test.go)
contains a full working example.

This mechanism is dynamic.  At runtime, you can set the logging to a given
level.  For example, you can automatically increase the level if the frequency
of errors goes above some threshold.  The L implementation eliminates the 
processing time of construction of loggable objects when filtering messages
in this way.

L provides more general middleware for filtering logging, with the same
performance considerations.

## Middleware

L provides middleware for filtering, handling special events, adding
custom automatic fields, etc.

Middleware is extremely powerful.  In addition to automatic fields,
it can be used for filtering, to adaptively change the log level,
to sample logs, to set up alerts, to set up metrics such as prometheus
or expvar, etc.

Middleware is invoked when the associated logger is locked, and has
full access to the associated log config labels, see below.

## Labels

Labels are the bread and butter of configuring and manipulating L loggers.

Labels are first of all just a `map[string]int`, associated with every `Config`
object.  There is a distinct config object for every logger, so we have a 1-1
relation between label maps and loggers.

Label keys concretely form a single global namespace for configuring loggers.
These keys, however, can be package qualified and used with or without
knowledge of the package names in play.

A label that starts with a '.' is considered a pattern for \<pkgname\>.  This
way, when setting labels or searching for loggers' configurations, one can
either find all loggers with a given package-qualified label, or search for
global labels, or search for labels within a set of packages.

Labels are integer valued, which is adapted (thus far) to most uses and remains
a simple atomic type.  This simplicity also aids in the specification of a 
logger.

Labels are available to middleware for reading and writing, so they can be used
to auto-monitor error rates or to dynamically trigger increased verbosity
localized to a specific functionality.

While many projects use a full fledged monitoring solution such as prometheus, 
many projects are not suited to depending on a 3rd party monitoring
service.  L can work out of the box, integrating your monitoring and remote
debugability in a single library.

## RPC Service

Please see [the design doc](https://github.com/scott-cotton/L/blob/main/rpc/design.md)



