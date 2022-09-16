package dap

type SetBreakpointsListEvent struct {
	Event
	Arguments []SetBreakpointsArguments `json:"arguments"`
}

func (r *SetBreakpointsListEvent) GetEvent() *Event { return &r.Event }

type StoppedDataEvent struct {
	Event
	Body             StoppedEventBody              `json:"body"`
	ThreadsBody      ThreadsResponseBody           `json:"threads_body"`
	StackTraceBody   StackTraceResponseBody        `json:"stack_trace_body"`
	ScopesBody       ScopesResponseBody            `json:"scopes_body"`
	VariablesBodyMap map[int]VariablesResponseBody `json:"variables_body_map"`
}

func (r *StoppedDataEvent) GetEvent() *Event { return &r.Event }
