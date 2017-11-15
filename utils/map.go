package flatten


import (
    "../types"
)


// Flatten a map to a key/value map using dot notation to represent nested
// layers of keys
func FlattenDocumentToDotNotation(document types.Document) types.Document {

    flattenedMap := make(map[string]interface{})

    for key, value := range document {

        switch child := value.(type) {

            case types.Document:

                subMap := FlattenDocumentToDotNotation(child)
                
                for subKey, subValue := range subMap {
                    flattenedMap[key + "." + subKey] = subValue
                }
            
            default:
                flattenedMap[key] = value

        }

    }

    return flattenedMap

}