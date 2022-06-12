package L

// ApplyOpts dictates how the application will occur.
//
// ApplyOpts is used as an argument to *Config.Apply,
// for example in the call `trg.Apply(mod, &ApplyOpts{...})`.
// In the following, we call 'trg' the target configurations
// and 'mod' the modifying configuration.
type ApplyOpts struct {
	// whether to apply recursively to child loggers.
	Recursive bool `json:"recursive,omitempty"`

	// labels whose values are carried over to the result, independent
	// of the setting of RemoveAbsentLabels or whether a label is in
	// the key set of the labels of the modifying config (in Go, whether or
	// or not 'ok' is true after calling '_, ok := mod.Labels[label]').
	PreserveLabels map[string]bool `json:"preserveLabels,omitempty"`

	// If true, the labels not in the modifying config are removed
	// from the target config, unless they are specified in PreservedLabels
	// above.
	RemoveAbsentLabels bool `json:"removeAbsentLabels,omitempty"`
}

// Apply applies the configuration o to c.  Fields are copied over if they are
// not nil in o, otherwise left untouched.  Labels in o should not include the
// package name, but if they start with '.', they are expanded with the package
// name of 'c' when copied to c's Labels.
//
// if a label in 'c', with any package name stripped, is not in o, then it is
// removed from c.
func (c *Config) Apply(o *Config, opts *ApplyOpts) {
	if opts == nil {
		opts = &ApplyOpts{}
	}
	if o.E != nil {
		c.E = o.E
	}
	if o.W != nil {
		c.W = o.W
	}
	if o.F != nil {
		c.F = o.F
	}
	if o.Pre != nil {
		c.Pre = append([]Middleware{}, o.Pre...)
	}
	if o.Post != nil {
		c.Post = append([]Middleware{}, o.Post...)
	}
	if o.Labels == nil {
		return
	}
	if c.Labels == nil {
		c.Labels = make(map[string]int, len(o.Labels))
	}
	if opts.RemoveAbsentLabels {
		for k := range c.Labels {
			if _, ok := o.Labels[c.Localize(k)]; !ok {
				if _, ok := opts.PreserveLabels[c.Localize(k)]; !ok {
					delete(c.Labels, k)
				}
			}
		}
	}
	for k, v := range o.Labels {
		if opts.PreserveLabels[k] {
			continue
		}
		c.Labels[c.Unlocalize(k)] = v
	}
}
