package domain

// Method /
type Method struct {
	Name          string
	ParameterList []MethodValue
	ResultList    []MethodValue
}

// MethodValue /
type MethodValue struct {
	Name string
	Type string
}
