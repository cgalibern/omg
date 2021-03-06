package cluster

import "opensvc.com/opensvc/util/timestamp"

type (
	// SchedulerThreadEntry describes a task queued for execution by the
	// opensvc scheduler thread.
	SchedulerThreadEntry struct {
		Action string      `json:"action"`
		Csum   string      `json:"csum"`
		Expire timestamp.T `json:"expire"`
		Path   string      `json:"path"`
		Queued timestamp.T `json:"queued"`
		Rid    string      `json:"rid"`
	}

	// SchedulerThreadStatus describes the OpenSVC daemon scheduler thread
	// state, which is responsible for executing node and objects scheduled
	// jobs.
	SchedulerThreadStatus struct {
		ThreadStatus
		Delayed []SchedulerThreadEntry `json:"delayed"`
	}
)
