package utils

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/fatih/structs"
)

// Structs for recursion max tests
type ATest struct {
    BTest
}

type BTest struct {
    CTest
}

type CTest struct {
    DTest
}

type DTest struct {
    Str string
}

// Struct for DoesNotMarshalUnexportedFields test
type ExportTest struct {
    AField string
    bField string
    CField string
    dField string
}

// Struct for DoesNotMarshalIfJsonIgnoreTagPresent test
type JsonIgnoreTagTest struct {
    AField string
    BField string `json:"-"`
    CField string
}

// Struct and types for DoesNotTouchEmbeddedNonStructTypes test
type NonStructTestType string
type EmbedsNonStruct struct {
    NonStructTestType
    AField string
}

// Structs for PutsEmbeddedDupFieldsInSubField test
type ASameStructWithBSameStruct struct {
    BSameStruct
    AField string
    BField string
    CField string
}

type BSameStruct struct {
    AField string
    CField string
    DField string
}

// Test toMap function doesn't call itself recursively more than the recursiveMax arg specified
func TestStructs_toMap_ErrorWhenOverRecursionMax(t *testing.T) {
    a := assert.New(t)

    // Make object with depth of 3
    obj := ATest{
        BTest{
            CTest{
               DTest{
                    Str: "Str",
                },
            },
        },
    }

    strct := structs.New(obj)

    // Call with recursion max 3 and expect to return error
    res, err := toMap(strct, 3, 0)

    a.Nil(res)
    a.Contains(err.Error(), "Recursed past level specified in recursionMax")
}

// Test toMap function works ok when it calls itself recursively less than the recursiveMax arg specified
func TestStructs_toMap_OkWhenUnderRecursionMax(t *testing.T) {
    a := assert.New(t)

    // Make object with depth of 0
    obj := ATest{
        BTest{
            CTest{
               DTest{
                    Str: "Str",
                },
            },
        },
    }

    strct := structs.New(obj)

    // Call with recursion max 20 and expect no error
    res, err := toMap(strct, 20, 0)

    a.NotNil(res)
    a.Nil(err)
}

// Test toMap function does not marshall fields that are not exported
func TestStructs_toMap_DoesNotMarshalUnexportedFields(t *testing.T) {
    a := assert.New(t)

    // Create struct to test
    obj := ExportTest{
        AField: "afield",
        bField: "bfield",
        CField: "cfield",
        dField: "dfield",
    }

    strct := structs.New(obj)

    // Test
    res, err := toMap(strct, DefaultRecursionMaxDepth, 0)

    a.Equal(map[string]interface{}{
        "AField": "afield",
        "CField": "cfield",
    }, res)
    a.Nil(err)
}

// Test that toMap does not marshal fields which have the "json:-" tag
// As this designates that the field should not be marshaled
func TestStructs_toMap_DoesNotMarshalIfJsonIgnoreTagPresent(t *testing.T) {
    a := assert.New(t)

    // Make struct
    obj := JsonIgnoreTagTest{
        AField: "afield",
        BField: "bfield",
        CField: "cfield",
    }

    strct := structs.New(obj)

    // Test
    res, err := toMap(strct, DefaultRecursionMaxDepth, 0)

    a.Equal(map[string]interface{}{
        "AField": "afield",
        "CField": "cfield",
    }, res)
    a.Nil(err)
}

// Test that toMap doesn't do anything with non struct type embedded fields
func TestStructs_toMap_DoesNotTouchEmbeddedNonStructTypes(t *testing.T) {
    a := assert.New(t)

    // Create struct
    var tt NonStructTestType = "tt"
    obj := EmbedsNonStruct{
        NonStructTestType: tt,
        AField: "afield",
    }

    strct := structs.New(obj)

    // Test
    res, err := toMap(strct, DefaultRecursionMaxDepth, 0)

    a.Equal(map[string]interface{}{
        "AField": "afield",
    }, res)
    a.Nil(err)
}

// Test that toMap puts fields from an embedded struct in the main struct
func TestStructs_toMap_PutsEmbeddedFieldsInMainStruct(t *testing.T) {
    a := assert.New(t)

    // Make struct
    obj := ATest{
        BTest{
            CTest{
                DTest{
                    Str: "Str",
                },
            },
        },
    }

    strct := structs.New(obj)

    // Test
    res, err := toMap(strct, DefaultRecursionMaxDepth, 0)

    a.Equal(map[string]interface{}{
        "Str": "Str",
    }, res)
    a.Nil(err)
}


// Test that toMap puts fields that both the main struct and embedded struct have in a field with the embedded struct's
// name
func TestStructs_toMap_PutsEmbeddedDupFieldsInSubField(t *testing.T) {
    a := assert.New(t)

    // Create struct
    obj := ASameStructWithBSameStruct{
        BSameStruct: BSameStruct{
            AField: "afieldb",
            CField: "cfieldb",
            DField: "dfieldb",
        },
        AField: "afielda",
        BField: "bfielda",
        CField: "cfielda",
    }

    strct := structs.New(obj)

    // Test
    res, err := toMap(strct, DefaultRecursionMaxDepth, 0)

    a.Equal(map[string]interface{}{
        "BSameStruct": map[string]interface{}{
            "AField": "afieldb",
            "CField": "cfieldb",
        },
        "DField": "dfieldb",
        "AField": "afielda",
        "BField": "bfielda",
        "CField": "cfielda",
    }, res)
    a.Nil(err)
}