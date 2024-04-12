package sui_types

import "github.com/thorli9527/sui-wallet-sdk/lib"

type DynamicFieldType struct {
	DynamicField  *lib.EmptyEnum `json:"DynamicField"`
	DynamicObject *lib.EmptyEnum `json:"DynamicObject"`
}

func (d DynamicFieldType) Tag() string {
	return ""
}

func (d DynamicFieldType) Content() string {
	return ""
}

type DynamicFieldName struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}
