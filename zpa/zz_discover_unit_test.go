//go:build discover

package zpa

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dummyConfig(s map[string]*schema.Schema, depth int) map[string]any {
	out := map[string]any{}
	for k, v := range s {
		if !v.Required || depth > 4 {
			continue
		}
		out[k] = dummyValue(v, depth)
	}
	return out
}

func dummyValue(v *schema.Schema, depth int) any {
	switch v.Type {
	case schema.TypeString:
		return "1"
	case schema.TypeInt:
		return 1
	case schema.TypeBool:
		return true
	case schema.TypeList, schema.TypeSet:
		switch e := v.Elem.(type) {
		case *schema.Resource:
			return []any{dummyConfig(e.Schema, depth+1)}
		case *schema.Schema:
			return []any{dummyValue(e, depth+1)}
		}
	}
	return "1"
}

func record(fx *unitFixture, from int) string {
	var out []string
	for _, r := range fx.api.Server.AllRequests()[from:] {
		p := strings.Replace(r.Path, "/zpa/mgmtconfig/v1/admin/customers/123456789", "V1", 1)
		p = strings.Replace(p, "/zpa/mgmtconfig/v2/admin/customers/123456789", "V2", 1)
		if r.Query != "" {
			p += "?" + r.Query
		}
		out = append(out, r.Method+" "+p)
	}
	return strings.Join(out, "\n     ")
}

func TestDiscover(t *testing.T) {
	p := ZPAProvider()
	for _, full := range strings.Fields(os.Getenv("DISCOVER")) {
		kind, n, _ := strings.Cut(full, ":")
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("== %s PANIC %v\n", full, r)
				}
			}()
			fx := newUnitFixture(t)
			fx.catchAll("1", map[string]any{"id": "1", "name": "1"})
			ctx := context.Background()
			if kind == "D" {
				r := p.DataSourcesMap[n]
				fmt.Printf("== %s required=%v\n", full, dummyConfig(r.Schema, 0))
				d := fx.resourceData(r, "", map[string]any{"name": "1"})
				fmt.Printf("  READ(name) diags=%v id=%q\n     %s\n", r.ReadContext(ctx, d, fx.client), d.Id(), record(fx, 0))
				n0 := len(fx.api.Server.AllRequests())
				d = fx.resourceData(r, "1", map[string]any{})
				fmt.Printf("  READ(id) diags=%v id=%q\n     %s\n", r.ReadContext(ctx, d, fx.client), d.Id(), record(fx, n0))
				return
			}
			r := p.ResourcesMap[n]
			cfg := dummyConfig(r.Schema, 0)
			fmt.Printf("== %s required=%v\n", full, cfg)
			n0 := 0
			for _, step := range []struct {
				name string
				fn   func(context.Context, *schema.ResourceData, interface{}) []any
			}{} {
				_ = step
			}
			d := fx.resourceData(r, "", cfg)
			fmt.Printf("  CREATE diags=%v id=%q\n     %s\n", r.CreateContext(ctx, d, fx.client), d.Id(), record(fx, n0))
			n0 = len(fx.api.Server.AllRequests())
			d = fx.resourceData(r, "1", cfg)
			fmt.Printf("  READ diags=%v\n     %s\n", r.ReadContext(ctx, d, fx.client), record(fx, n0))
			n0 = len(fx.api.Server.AllRequests())
			if r.UpdateContext != nil {
				d = fx.resourceData(r, "1", cfg)
				fmt.Printf("  UPDATE diags=%v\n     %s\n", r.UpdateContext(ctx, d, fx.client), record(fx, n0))
			}
			n0 = len(fx.api.Server.AllRequests())
			d = fx.resourceData(r, "1", cfg)
			fmt.Printf("  DELETE diags=%v\n     %s\n", r.DeleteContext(ctx, d, fx.client), record(fx, n0))
		}()
	}
}
