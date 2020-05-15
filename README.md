# wmclient

Common libraries for Windmill client apps.

[![Build Status](https://circleci.com/gh/tilt-dev/wmclient/tree/master.svg?style=shield)](https://circleci.com/gh/tilt-dev/wmclient)
[![GoDoc](https://godoc.org/github.com/tilt-dev/wmclient?status.svg)](https://godoc.org/github.com/tilt-dev/wmclient)

These libraries have Windmill-specific configuration baked in. They are probably
not useful for a wide audience. We publish them as open-source so that Windmill users
can inspect and verify the code running on their machine.

## Analytics

Windmill client apps may report usage to https://events.windmill.build, to help us
understand what features people use. This package implements logic for sending
these reports when the user opts in, and for storing the opt-in status.

We do not report any personally identifiable information. We do not report any
identifiable data about your code.

We do not share this data with anyone who is not an employee of Windmill
Engineering.  Data may be sent to third-party service providers like Datadog,
but only to help us analyze the data.

## License
Copyright 2018 Windmill Engineering

Licensed under [the Apache License, Version 2.0](LICENSE)
