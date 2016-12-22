package utils

import "github.com/fatih/structs"

var RecursionMaxDepth = 20

// Converts a struct into a map with string keys and interface{} values.
// Takes a structs.Struct as input, this type was chosen to indicate that this method is for use on structs.
//
// Merges embedded types fields into main struct to mimic the implied "inheritance".
func toMap (s structs.Struct, recursionMax, recursionCounter int) *map[string]interface{} {
    var result map[string]interface{}

    // Loop over fields
    fields := s.Fields()
    for _, v := range fields {
        // Check to see if embedded type
        if v.IsEmbedded() == true {
            vals := toMap(v.Value(), recursionMax, recursionMax)
            // TODO: Make branch if v isn't embedded. Make sure map respects JSON tag for key name
        }
    }
}

// ToMap proxy caller. Allows user not specify recursion max depth and use default RecursionMaxDepth instead.
func ToMap (s structs.Struct) *map[string]interface{} {
    return toMap(s,RecursionMaxDepth, 0)
}
