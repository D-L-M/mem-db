package types


// Documents could be any format once they're parsed, so make a generic object
// type to define them
type Document map[string]interface {}


// JSON documents are generic key/value objects
type JsonDocument map[string]interface {}