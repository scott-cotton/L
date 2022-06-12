package rpc

import (
	"fmt"
	"regexp"

	"github.com/scott-cotton/L"
)

type ApplyParams struct {
	PkgPattern string       `json:"pkgPattern"`
	Opts       *L.ApplyOpts `json:"opts,omitempty"`
	Config     *L.Config    `json:"config"`
}

type ApplyResult []L.PackageConfig

func Apply(parms *ApplyParams) (ApplyResult, error) {
	if parms.PkgPattern == "" {
		parms.PkgPattern = ".*"
	}
	pkgRe, err := regexp.Compile(parms.PkgPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}
	var res ApplyResult
	walk := func(cfg *L.Config) {
		pkg := cfg.Package()
		if !pkgRe.MatchString(pkg) {
			return
		}
		cfg.Apply(parms.Config, parms.Opts)
		ocfg := cfg.Clone()
		labels := map[string]int{}
		for k, v := range ocfg.Labels {
			labels[ocfg.Localize(k)] = v
		}
		ocfg.Labels = labels
		res = append(res, L.PackageConfig{
			Config:  *ocfg,
			Package: pkg,
		})
	}
	L.Walk(walk)
	return res, nil
}
