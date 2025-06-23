package types

type ActivityReport struct {
	ActivityName string                 `json:"name"`
	Inputs       map[string]interface{} `json:"input,omitempty"`
	Outputs      *interface{}           `json:"output,omitempty"`
	Error        map[string]interface{} `json:"error,omitempty"`
}

type LinkReport struct {
	LinkName string `json:"linkName,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
}

type FlowReport struct {
	Name             string                 `json:"testName,omitempty"`
	ActivityReport   []ActivityReport       `json:"activities"`
	LinkReport       []LinkReport           `json:"links,omitempty"`
	FlowErrorHandler FlowErrorHandler       `json:"errorHandler,omitempty"`
	SubFlow          map[string]interface{} `json:"subFlow,omitempty"`
}

type Report struct {
	Trigger *Trigger    `json:"trigger"`
	Flows   *FlowReport `json:"flows"`
}

type Trigger struct {
	ID       string                 `json:"id,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	Handler  Handler                `json:"handler"`
}

type Handler struct {
	FlowName string `json:"flow_name"`
}

type FlowErrorHandler struct {
	ActivityReport []ActivityReport `json:"activities"`
	LinkReport     []LinkReport     `json:"links,omitempty"`
}

type DebugOptions struct {
	Op                  int
	ReturnID            bool
	FlowURI             string
	PreservedInstanceId string
	InitStepId          int
	ExecOptions         *DebugExecOptions
	Rerun               bool
	OriginalInstanceId  string
	DetachExecution     bool
}

type DebugExecOptions struct {
	Interceptor *Interceptor
}
