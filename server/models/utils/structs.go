package utils

import (
    "github.com/fatih/structs"
    "errors"
    "strings"
)

var DefaultRecursionMaxDepth = 20

// Converts a struct into a map with string keys and interface{} values.
// Takes a structs.Struct as input, this type was chosen to indicate that this method is for use on structs.
//
// Merges embedded type's fields into main struct to mimic the implied "inheritance".
// If an embedded type is not a struct than it will be treated like a normal KV pair.
//
// If an embedded type and the main struct have colliding keys the problematic keys from the embedded type will be placed
// in an object who's key is the name of the embedded type.
//
// Example given:
//
// Represents a security clearance to some company resources. Lets pretend it has a bunch of useful methods for
// accessing secure documents.
// type Clearance struct {
//     Id int // Id of clearance to access documents. Perhaps useful for revoking access.
//     Level int // Level of access
// }
//
// type Manager struct {
//     Clearance // Manager has embedded Clearance so that the useful methods of Clearance are easier to use.
//     Id int // User id
//     Name string // Name
//     Project string // Project
// }
//
// var s = Manager{Clearance{4534, 4}, 89243, "John Smith", "Secret project 3"}
// toMap(s) => {
//     Id: 89243,
//     Name: "John Smith",
//     Project: "Secret project 3",
//     Level: 4,
//     Clearance: {
//         Id: 4534,
//     },
// }
//
//  This method of merging is useful for presenting models which are meant to inherit properties.
//
// Note: Although I'm not sure this situation could actually occur given the semantics of Golang I will document it:
//     If an embedded field has multiple of a key that collides with the main struct the key's value will be the last
//     value in the struct. There is no special code that does this, this is just how the program runs (A loop over the
//     struct to put it into a map).
//
// In the same realm of impossibility as above this issue is documented:
//     If the main struct has a key with the same name as the name of an embedded type and a collided key needs to be
//     merged an error will be thrown.
//
// Returns error if recursed past level specified by recursionMax
//
// Returns converted map and error
func toMap (s *structs.Struct, recursionMax, recursionCounter int) (map[string]interface{}, error) {
    // Check recursion max
    recursionCounter++

    if (recursionCounter > recursionMax) {
        return nil, errors.New("Recursed past level specified in recursionMax")
    }

    // Result to fill later
    result := make(map[string]interface{}, 0)

    // Loop over fields
    fields := s.Fields()
    for _, field := range fields {
        // Check that field is exported
        if field.IsExported() == false {// If not then skip over
            continue
        }

        // Get field name
		fieldName := FieldName(*field)

		// Ignore field if FieldName returns empty string
		// This means a JSON tag was added to the field designating it should be ignored when marshalling
		if fieldName == "" {
			continue
		}

        // Check to see if field is embedded
        if field.IsEmbedded() == true {// If field is embedded
            // Check if embedded field is struct
		    if structs.IsStruct(field.Value()) {// If it is handle recursively
                // Call recursively to deal with this object
                embeddedVals, err := toMap(structs.New(field.Value()), recursionMax, recursionCounter)
                if err != nil {
                    return nil, errors.New("Error while processing embedded field \"" + field.Name() + "\": " + err.Error())
                }

                var embeddedMap map[string]interface{}

                // Merge map of embedded values with map of main struct
                for key, value := range embeddedVals {
                    // Check if key exists
                    if _, ok := result[key]; ok == true {
                        // If it does exist then put in object with key of embedded type's name
                        embeddedMap[key] = value
                    } else {
                        // If the key does not exist in the main object then set it
                        result[key] = value
                    }
                }

                // Attach embeddedMap onto main struct map
                // Check if main struct has field with same
                // If main struct does have field with name of embedded type throw error
                if _, ok := result[fieldName]; ok == true {
                    return nil, errors.New("[WTF] Main struct has separate key with name of embedded type: \"" + field.Name() + "\"")
                }

                // Set embedded values if they exist
                if len(embeddedMap) > 0 {
                    result[fieldName] = embeddedMap
                }
            }
        } else {// If field isn't embedded, aka normal
            result[FieldName(*field)] = field.Value()
        }
    }

    return result, nil
}

// FieldName returns the name of a field or an empty string if the field should be omitted.
//
// Field should be omitted if specified by "json" tag or not exported.
//
// Method follows rules of the "json" tag as described in encoding/json.Marshal except for the "string" flag.
func FieldName (field structs.Field) string {
    // Check that field is exported
    if field.IsExported() == false {// If not then return empty string
        return ""
    }

    // Get JSON tag
    jsonTag := field.Tag("json")

    // If "json" tag exists
    if jsonTag != "" {
        // Parse JSON tag
        jsonTagVals := strings.Split(jsonTag, ",")

        // If "json" tag has the value of "-" return empty string
        // Designates that field should not be marshaled
        if jsonTag == "-" {
            return ""
        }

        // Check if only name is given
        if strings.Contains(jsonTag, ",") == false {
            return jsonTag
        } else {// Parse into parts
            var jsonName string
            var jsonFlags []string

            if len(jsonTagVals) == 1 {// 1 value but comma provided = only flag provided, ex: ",omitempty"
                jsonFlags[0] = jsonTagVals[0]
            } else if len(jsonTagVals) >= 2 {// 2 or more values and comma provided = flag and name provided, ex: "name,omitempty"
                jsonName = jsonTagVals[0]
                jsonFlags = jsonTagVals[1:]
            }

            // Handle json flags
            // Even though we only handle one, build in ability to handle more in future
            // This design also allows us to have extra values of the "json" tag that we don't handle
            for _, flagVal := range jsonFlags {
                switch flagVal {
                case "omitempty":
                    // If "omitempty" json flag and empty field then return empty string
                    if field.IsZero() {
                        return ""
                    }
                }
            }

            // Now that all checks are complete return json name if exists
            if jsonName != "" {
                return jsonName
            }
        }
    }

    // If "json" tag exists but name is not provided or if "json" tag doesn't exist at all then return normal field name
    return field.Name()
}

// ToMap proxy caller. Allows user not specify recursion max depth and use default DefaultRecursionMaxDepth instead.
func ToMap (s *structs.Struct) (map[string]interface{}, error) {
    return toMap(s, DefaultRecursionMaxDepth, 0)
}
