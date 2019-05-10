package analytics

// OptReporter is a short-lived Analytics to record a single metric (stripped of
// identifying information) and then disable itself. For use before user
// has opted in or out, and ONLY for use sending a metric of their choice.
type optReporter struct {
	a *remoteAnalytics
}

func newOptReporter(appName string, options ...Option) (*optReporter, error) {
	options = append(options,
		WithEnabled(true), // always enabled for first call
		WithUserID("anon"), WithMachineID("anon")) // anonymized
	a, err := NewRemoteAnalytics(appName, options...)
	if err != nil {
		return nil, err
	}

	return &optReporter{a}, nil
}

func (or *optReporter) incrOpt(c Opt) {
	or.a.Incr("analytics.opt", map[string]string{"choice": c.String()})
	or.a.enabled = false
}
