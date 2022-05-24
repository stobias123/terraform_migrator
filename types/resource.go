package types

// Resource will parse a terraform resource block
type Resource struct {
	Name string `hcl:"name,label"`
	Attributes []AttributeConfig `hcl:"attribute,block"`
}

