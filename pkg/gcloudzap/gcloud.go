package gcloudzap

import gcl "cloud.google.com/go/logging"

// GoogleCloudLogger encapsulates the important methods of gcl.Logger
type GoogleCloudLogger interface {
	Flush() error
	Log(e gcl.Entry)
}
