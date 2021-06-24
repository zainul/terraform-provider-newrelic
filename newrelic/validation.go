package newrelic

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func float64Gte(gte float64) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(float64)
		if !ok {
			return diag.Errorf("expected type of %s to be float64", k)
		}

		if v >= gte {
			return nil
		}

		return diag.Errorf("expected %s to be greater than or equal to %v, got %v", k, gte, v)
	}
}

func intInSlice(valid []int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(int)
		if !ok {
			return diag.Errorf("expected type of %s to be int", k)
		}

		for _, p := range valid {
			if v == p {
				return nil
			}
		}

		return diag.Errorf("expected %s to be one of %v, got %v", k, valid, v)
	}
}

// float64AtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type float64 and is at least min (inclusive)
func float64AtLeast(min float64) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(float64)
		if !ok {
			return diag.Errorf("expected type of %s to be float64", k)
		}

		if v < min {
			return diag.Errorf("expected %s to be at least %f, got %f", k, min, v)
		}

		return nil
	}
}
