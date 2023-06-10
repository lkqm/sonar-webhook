package sonar

type MeasuresComponentResponse struct {
	Component MeasuresComponent `json:"component"`
}

type MeasuresComponent struct {
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Qualifier   string    `json:"qualifier"`
	Measures    []Measure `json:"measures"`
}

// Measure 指标
type Measure struct {
	Metric    string   `json:"metric"`
	Value     string   `json:"value"`
	BestValue bool     `json:"bestValue,omitempty"`
	Period    Period  `json:"period"`
	Periods   []Period `json:"periods"`
}

type Period struct {
	Index     int    `json:"index"`
	Value     string `json:"value"`
	BestValue bool   `json:"bestValue,omitempty"`
}
