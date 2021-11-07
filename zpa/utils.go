package zpa

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ValidateStringFloatBetween(min, max float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		str, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string with value of float64", k))
			return
		}

		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			errors = append(errors, fmt.Errorf("expected type of %s to be float64: %v", k, err))
			return
		}

		if v < min || v > max {
			errors = append(errors, fmt.Errorf("expected %s to be in the range (%f - %f), got %f", k, min, max, v))
			return
		}

		return
	}
}

func SetToStringSlice(d *schema.Set) []string {
	list := d.List()
	return ListToStringSlice(list)
}

func SetToStringList(d *schema.ResourceData, key string) []string {
	setObj, ok := d.GetOk(key)
	if !ok {
		return []string{}
	}
	set, ok := setObj.(*schema.Set)
	if !ok {
		return []string{}
	}
	return SetToStringSlice(set)
}

func ListToStringSlice(v []interface{}) []string {
	if len(v) == 0 {
		return []string{}
	}

	ans := make([]string, len(v))
	for i := range v {
		switch x := v[i].(type) {
		case nil:
			ans[i] = ""
		case string:
			ans[i] = x
		}
	}

	return ans
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func MergeSchema(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	final := make(map[string]*schema.Schema)
	for _, s := range schemas {
		for k, v := range s {
			final[k] = v
		}
	}
	return final
}

func convertToListString(obj interface{}) []string {
	listI, ok := obj.([]interface{})
	if ok && len(listI) > 0 {
		list := make([]string, len(listI))
		for i, e := range listI {
			s, ok := e.(string)
			if ok {
				list[i] = e.(string)
			} else {
				log.Printf("[WARN] invalid type: %v\n", s)
			}
		}
		return list
	}
	return []string{}
}
