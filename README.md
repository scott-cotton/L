# L

L is a minimalist project-aware structured logger for Go.

## Why?

As observed by [capnslog](https://github.com/coreos/capnslog), logging should
serve the project it uses well: the main entry point should determine the
logging configuration of its imports, and moreover do so in a way that is
runtime configurable.

If projects A and B organise their logging differently, and project C would
like to use A and B and yet manage the logging, L can be used to resolve the
differences between A and B and C's desired logging setup.  In other words, if
you would like your project to be used in other projects with distinct logging
styles, and you find L worthy enough in its current nascent state, the L can do
that.

L permits flexible per-project configuration in a way that can be overriden by
a main entrypoint, viewed, and (in progress) manipulated via an authenticated
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
imported packages, transitively by calling `L.Apply(...)`:
```
import "github.com/scott-cotton/L"

L.Apply(MyAppConfig())
```


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
of errors goes above some threshold.

L provides more general middleware for filtering logging.
