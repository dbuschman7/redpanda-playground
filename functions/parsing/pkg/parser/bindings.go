package parser

// The result of parsing is a slice of Bindings
type Bindings []Binding

// A Binding corresponds to “name = value”
type Binding struct {
	Name  string
	Value BindingValue
}

// BindingValue is a marker interface for the values in a Binding.
type BindingValue interface {
	IsBindingValue()
}

// BindingInt is a wrapper on int to implement the BindingValue interface.
type BindingInt int

// The marker method to be a BindingValue
func (BindingInt) IsBindingValue() {}

// BindingBool is a wrapper on bool to implement the BindingValue interface.
type BindingBool bool

// The marker method to be a BindingValue
func (BindingBool) IsBindingValue() {}

// BindingString is a wrapper on bool to implement the BindingValue interface.
type BindingString bool

// The marker method to be a BindingValue
func (BindingString) IsBindingValue() {}
