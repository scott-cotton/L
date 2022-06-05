# L

L is a minimalist project-aware structured logger for Go.

## Why?

There are indeed already too many logging libraries out there.  It makes a lot
of sense to continue using them if they serve you well, and maybe even if they
don't given the effort of changing something as low level as logging.

However, logging has a evolved a lot since the standard library logging
interface.  First, structured logging has become the norm for newer
applications, perhaps due to the nature of logging in cloud distributed systems
such as Kubernetes.  Second, as observed by
[capnslog](https://github.com/coreos/capnslog), logging should serve the project
it uses well: the main entry point should determine the logging configuration of
its imports, and moreover do so in a way that is runtime configurable.

Unfortunately, all the libraries out there seem to use the idea of the global
logger, where this state can be changed, and often is, inside library code.


L provides

- Efficient structured logging, with reasonable syntax for creating nested structures
  incrementally.
- Formatters and a simple formatting interface.
- The ability for the main entry point to determine and set up configuration
  for all Loggers in all packages it depends on.
- The ability for a package to define its own default configuration.
- Support for _hooks_ aka MiddleWare.
- An admin interface for manipulating labels in a Go process.

## Working with L

To use L in some package, simply import L and call L.New(cfg) with
your desired configuration.  This should normally done at the top level
of each package once.

## Basic structured logging

Structured logging is about providing structured output, which is 
typically fairly unstructured actually and refers to anything that 
can be viewed as a `map[string]any`, where `any` is any ground type
or another `map[string]any` or a slice of values, i.e. Go's representation
of free form json objects.

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

// slower, but easier on the eyes
myL.Fmt(`{"k": %d, %q: [ "%s", %d ] }`, 3, "key", "a", 4) 

```

## Custom Log project-native types

L provides a means to avoid copying your native types at 
all their points they are used and logged, as this clearly is tedious.

The idea is simple:  define a json Unmarshaler for a type and then call
```
v := &T{...}
obj.JSON(v)
```

Note that there are libraries to autogenerate marshalers, and also
one can easily enough define an alternate marshaler for logging specific purposes.

The distinction that L makes here is that if `v` does not implement json.Marshaler,
then the program will fail to compile.  The restriction helps to guarantee that logging
is fast.



### Overriding the configuration at the entry point

The entry point can set the configuration for all
imported packages, transitively by calling `L.Apply(...)`:
```
import "github.com/scott-cotton/L"

L.Apply(MyAppConfig())
```

In this context, the meaning of the labels is as follows:
a label key is a regular expression which 
- a label value is an indication of whether or not labels matched 
with the corresponding key should be kept.  It can contain package names
- deletion takes precedence over insertion, so setting the value to
  false in the label map will always delete the corresponding matched
  keys.

This allows for the application of labels accross all packages following
a project policy.

Applying a config allows setting middleware for pre- and post-processing.
nil values for the middleware, .E, .F, and .W attributes will not overwrite
package-specified ones.


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

[this test file](levels_test.go) contains a full working example.

This mechanism is dynamic.  At runtime, you can set the logging to a given
level.  For example, you can automatically increase the level if the frequency
of errors goes above some threshold.

L provides more general middleware for filtering logging.
