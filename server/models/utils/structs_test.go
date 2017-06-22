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
    str string
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

func TestStructs_toMap_ErrorWhenOverRecursionMax(t *testing.T) {
    a := assert.New(t)

    // Make object with depth of 3
    obj := ATest{
        BTest{
            CTest{
               DTest{
                    str: "str",
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

func TestStructs_toMap_OkWhenUnderRecursionMax(t *testing.T) {
    a := assert.New(t)

    // Make object with depth of 0
    obj := DTest{
        str: "str",
    }

    strct := structs.New(obj)

    // Call with recursion max 20 and expect no error
    res, err := toMap(strct, 20, 0)

    a.NotNil(res)
    a.Nil(err)
}

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

// Tests that toMap does not marshal fields which have the "json:-" tag
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