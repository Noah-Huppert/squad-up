package utils

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/fatih/structs"
)

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

func TestStructs_toMap_RespectsRecursionMax(t *testing.T) {
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

    // Call with recursion max 5 and expect to return error
    res, err := toMap(strct, 3, 0)

    a.Nil(res)
    a.Contains(err.Error(), "Recursed past level specified in recursionMax")
}