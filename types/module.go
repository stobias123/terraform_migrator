package types

type Module struct {
	Name string `hcl:"name,label"`
	Attributes []AttributeConfig `hcl:"attribute,block"`
}