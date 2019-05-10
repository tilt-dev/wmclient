package analytics

import "fmt"

// OptReporter is a short-lived Analytics to record a single metric (stripped of
// identifying information) and then disable itself. For use before user
// has opted in or out, and ONLY for use sending a metric of their choice.
type OptReporter struct {
	a *remoteAnalytics
	used bool
}

func NewOptReporter(appName string, options ...Option) (*OptReporter, error) {
	options = append(options,
		WithEnabled(true), // always enabled for first call
		WithUserID("anon"), WithMachineID("anon")) // anonymized
	a, err := NewRemoteAnalytics(appName, options...)
	if err != nil {
		return nil, err
	}

	return &OptReporter{a: a}, nil
}

func (or *OptReporter) incrOpt(c Opt) error {
	if or.used {
		return fmt.Errorf("optReporter already used, can't incr opt: %s", c.String())
	}
	or.a.Incr("analytics.opt", map[string]string{"choice": c.String()})

	or.used = true
	or.a.enabled = false // disable analytics, just to be safe.
	return nil
}
