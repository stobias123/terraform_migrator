package types
type AttributeAction string
const (
	// Add will add an attribute to the block.
	Add AttributeAction = "add"
	// Delete removes the attribute value
	Delete AttributeAction = "delete"
	// Update always replaces the attribute value.
	//TODO: Update or create?
	Update AttributeAction = "update"
	// Replace will only update the value when your source val is matched
	Replace AttributeAction = "replace"
)

type AttributeConfig struct {
	Name string `hcl:"name"`
	Action AttributeAction `hcl:"action"`
	OriginalValue *string `hcl:"original_value"`
	Value *string `hcl:"value"`
}