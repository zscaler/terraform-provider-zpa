package common

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Converts camelCase to snake_case with exceptions
func camelToSnake(name string) string {
	if name == "ID" || name == "Id" {
		return "id"
	}
	if name == "MicroTenantID" || name == "MicroTenantId" {
		return "microtenant_id"
	}
	if name == "MicroTenantName" {
		return "microtenant_name"
	}

	if name == "DNSQueryType" {
		return "dns_query_type"
	}
	var result strings.Builder
	for i, char := range name {
		if unicode.IsUpper(char) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(char))
	}
	return result.String()
}

// Converts field names based on specific rules
func convertFieldName(fieldName string) string {
	return camelToSnake(fieldName)
}

func StructToSchema(v interface{}) map[string]*schema.Schema {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	schemaMap := make(map[string]*schema.Schema)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tf")
		if tag == "" {
			continue
		}

		fieldName := convertFieldName(field.Name)

		s := &schema.Schema{
			Type:     getSchemaType(field.Type),
			Optional: tagContains(tag, "optional"),
			Required: tagContains(tag, "required"),
			Computed: tagContains(tag, "computed"),
		}

		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
			s.Elem = &schema.Resource{
				Schema: StructToSchema(reflect.New(field.Type.Elem()).Interface()),
			}
		}

		schemaMap[fieldName] = s
	}
	return schemaMap
}

func getSchemaType(t reflect.Type) schema.ValueType {
	switch t.Kind() {
	case reflect.String:
		return schema.TypeString
	case reflect.Bool:
		return schema.TypeBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return schema.TypeInt
	case reflect.Slice:
		if t.Elem().Kind() == reflect.String {
			return schema.TypeList
		}
		return schema.TypeSet
	default:
		return schema.TypeString
	}
}

func tagContains(tag, key string) bool {
	return strings.Contains(tag, key)
}

func DataToStructPointer(d *schema.ResourceData, v interface{}) {
	rv := reflect.ValueOf(v).Elem()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldName := convertFieldName(rv.Type().Field(i).Name)

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if val, ok := d.GetOk(fieldName); ok {
				field.SetString(val.(string))
			}
		case reflect.Bool:
			if val, ok := d.GetOk(fieldName); ok {
				field.SetBool(val.(bool))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if val, ok := d.GetOk(fieldName); ok {
				field.SetInt(int64(val.(int)))
			}
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.Struct {
				if val, ok := d.GetOk(fieldName); ok {
					list := val.(*schema.Set).List()
					slice := reflect.MakeSlice(field.Type(), len(list), len(list))
					for i := 0; i < len(list); i++ {
						item := reflect.New(field.Type().Elem()).Elem()
						mapDataToStruct(list[i].(map[string]interface{}), item)
						slice.Index(i).Set(item)
					}
					field.Set(slice)
				}
			} else {
				if val, ok := d.GetOk(fieldName); ok {
					list := val.([]interface{})
					slice := reflect.MakeSlice(field.Type(), len(list), len(list))
					for i := 0; i < len(list); i++ {
						slice.Index(i).Set(reflect.ValueOf(list[i]))
					}
					field.Set(slice)
				}
			}
		}
	}
}

func StructToData(v interface{}, d *schema.ResourceData) error {
	rv := reflect.ValueOf(v).Elem()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldName := convertFieldName(rv.Type().Field(i).Name)

		switch field.Kind() {
		case reflect.String:
			d.Set(fieldName, field.String())
		case reflect.Bool:
			d.Set(fieldName, field.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			d.Set(fieldName, field.Int())
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.Struct {
				slice := make([]interface{}, field.Len())
				for i := 0; i < field.Len(); i++ {
					mapData := make(map[string]interface{})
					structToMap(field.Index(i), mapData)
					slice[i] = mapData
				}
				d.Set(fieldName, slice)
			} else {
				slice := make([]interface{}, field.Len())
				for i := 0; i < field.Len(); i++ {
					slice[i] = field.Index(i).Interface()
				}
				d.Set(fieldName, slice)
			}
		}
	}
	return nil
}

func mapDataToStruct(data map[string]interface{}, v reflect.Value) {
	for key, value := range data {
		field := v.FieldByName(strings.Title(key))
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		}
	}
}

func structToMap(v reflect.Value, data map[string]interface{}) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := convertFieldName(v.Type().Field(i).Name)
		data[fieldName] = field.Interface()
	}
}
