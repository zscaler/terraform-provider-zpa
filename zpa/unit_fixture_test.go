package zpa

// Unit test fixture: runs resource and data source CRUD functions against an
// in-memory mock of the ZPA API (the zscaler-sdk-go OneAPI test harness), so
// provider logic (request building, state mapping, plan-time checks, error
// handling) is tested offline, without TF_ACC, credentials or a tenant.
//
// Unit tests are named Test<Resource>_<Behavior> and always run. Acceptance
// tests (TestAcc*, real API calls) still run only when TF_ACC is set.

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

type unitFixture struct {
	t      *testing.T
	api    *common.APITest
	client *Client
}

// newUnitFixture starts a mock ZPA API and returns a provider client wired to it.
// Requests to paths without a registered response get a 404.
func newUnitFixture(t *testing.T) *unitFixture {
	t.Helper()
	api := common.NewZPATest(t)
	client := &Client{Service: api.Service, policySetIDCache: make(map[string]string)}
	return &unitFixture{t: t, api: api, client: client}
}

// path builds a v1 mgmtconfig path for the mock customer, e.g. path("segmentGroup", "123").
func (f *unitFixture) path(resource string, parts ...string) string {
	return common.ZPAPath(f.api.CustomerID, resource, parts...)
}

// pathV2 builds a v2 mgmtconfig path for the mock customer.
func (f *unitFixture) pathV2(resource string, parts ...string) string {
	return common.ZPAv2Path(f.api.CustomerID, resource, parts...)
}

// on registers the response for every request matching method and path.
func (f *unitFixture) on(method, path string, response common.MockResponse) *unitFixture {
	f.api.On(method, path, response)
	return f
}

// calls returns how many requests the provider sent to method and path.
func (f *unitFixture) calls(method, path string) int {
	return f.api.Server.GetCallCount(method, path)
}

// requestBodies returns the decoded JSON bodies of the requests sent to method and path.
func (f *unitFixture) requestBodies(method, path string) []map[string]any {
	f.t.Helper()
	var bodies []map[string]any
	for _, r := range f.api.Server.AllRequests() {
		if r.Method != method || r.Path != path || len(r.Body) == 0 {
			continue
		}
		var body map[string]any
		if err := json.Unmarshal(r.Body, &body); err != nil {
			f.t.Fatalf("request body of %s %s is not a JSON object: %v", method, path, err)
		}
		bodies = append(bodies, body)
	}
	return bodies
}

// resourceData builds the ResourceData a CRUD function receives: raw is the
// configuration (or state) as a map, id the resource ID ("" for a create).
func (f *unitFixture) resourceData(r *schema.Resource, id string, raw map[string]any) *schema.ResourceData {
	f.t.Helper()
	d := schema.TestResourceDataRaw(f.t, r.Schema, raw)
	d.SetId(id)
	return d
}

// diff runs the plan for a resource with existing state and a new configuration,
// including CustomizeDiff, and returns the plan-time error, if any.
func (f *unitFixture) diff(r *schema.Resource, id string, state, config map[string]any) error {
	f.t.Helper()
	instance := f.resourceData(r, id, state).State()
	_, err := r.Diff(context.Background(), instance, terraform.NewResourceConfigRaw(config), f.client)
	return err
}

// updateData builds the ResourceData an update receives: the prior state plus
// the planned change to config, so HasChange and GetChange behave as in a plan.
func (f *unitFixture) updateData(r *schema.Resource, id string, state, config map[string]any) *schema.ResourceData {
	f.t.Helper()
	instance := f.resourceData(r, id, state).State()
	planned, err := r.Diff(context.Background(), instance, terraform.NewResourceConfigRaw(config), f.client)
	if err != nil {
		f.t.Fatalf("plan failed: %v", err)
	}
	d, err := schema.InternalMap(r.Schema).Data(instance, planned)
	if err != nil {
		f.t.Fatalf("building update data: %v", err)
	}
	return d
}

// validate runs Terraform's schema validation (ValidateFunc and friends) on a
// resource configuration and returns the first error summary ("" when valid).
func validate(r *schema.Resource, config map[string]any) string {
	for _, d := range r.Validate(terraform.NewResourceConfigRaw(config)) {
		if d.Severity == diag.Error {
			return d.Summary
		}
	}
	return ""
}

// terraformConfig wraps a raw configuration map for schema.Resource.Diff.
func terraformConfig(raw map[string]any) *terraform.ResourceConfig {
	return terraform.NewResourceConfigRaw(raw)
}

// body returns the JSON body of the single request sent to method and path.
func (f *unitFixture) body(method, path string) map[string]any {
	f.t.Helper()
	bodies := f.requestBodies(method, path)
	if len(bodies) != 1 {
		f.t.Fatalf("expected 1 %s %s request, got %d", method, path, len(bodies))
	}
	return bodies[0]
}

type crudFunc = func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics

// run calls a CRUD function and fails the test on an error diagnostic.
func (f *unitFixture) run(fn crudFunc, d *schema.ResourceData) {
	f.t.Helper()
	if diags := fn(context.Background(), d, f.client); diags.HasError() {
		f.t.Fatalf("unexpected error: %v", diags)
	}
}

// runErr calls a CRUD function and returns its first error summary ("" when none).
func (f *unitFixture) runErr(fn crudFunc, d *schema.ResourceData) string {
	f.t.Helper()
	for _, diagnostic := range fn(context.Background(), d, f.client) {
		if diagnostic.Severity == diag.Error {
			return diagnostic.Summary
		}
	}
	return ""
}

// v1Path and v2Path build mgmtconfig paths for the mock customer outside a fixture.
func v1Path(resource string, parts ...string) string {
	return common.ZPAPath(common.TestCustomerID, resource, parts...)
}

func v2Path(resource string, parts ...string) string {
	return common.ZPAv2Path(common.TestCustomerID, resource, parts...)
}

// cbiPath builds a cloud browser isolation config path for the mock customer.
func cbiPath(parts ...string) string {
	return path.Join(append([]string{"/zpa/cbiconfig/cbi/api/customers", common.TestCustomerID}, parts...)...)
}

// mockPolicySet registers the policy set of policyType (ID "ps-<policyType>")
// with no rules, for resources that look up or detach from policy rules.
func (f *unitFixture) mockPolicySet(policyType string) string {
	id := "ps-" + policyType
	f.on("GET", f.path("policySet", "policyType", policyType),
		common.SuccessResponse(policysetcontroller.PolicySet{ID: id, PolicyType: policyType}))
	f.on("GET", f.path("policySet", "rules", "policyType", policyType),
		common.SuccessResponse(common.ZPAList([]policysetcontroller.PolicyRule{})))
	return id
}

// policyRuleSpec returns the paths of rule "r1" in the policy set of
// policyType ("ps-<policyType>"): create and update use the v2 API when v2 is
// set; reads and deletes always use v1.
func policyRuleSpec(policyType string, v2 bool) crudSpec {
	ps := "ps-" + policyType
	spec := crudSpec{
		id:         "r1",
		createPath: v1Path("policySet", ps, "rule"),
		itemPath:   v1Path("policySet", ps, "rule", "r1"),
		listPath:   v1Path("policySet", "rules", "policyType", policyType),
		setup:      func(fx *unitFixture) { fx.mockPolicySet(policyType) },
	}
	if v2 {
		spec.createPath = v2Path("policySet", ps, "rule")
		spec.updatePath = v2Path("policySet", ps, "rule", "r1")
	}
	return spec
}

// apiErrorMarker is the reason returned by apiError; tests look for it in diagnostics.
const apiErrorMarker = "mock api failure"

// apiError is a 400 response: the SDK does not retry it.
func apiError() common.MockResponse {
	return common.SuccessResponseWithStatus(http.StatusBadRequest,
		map[string]any{"id": "invalid.request", "reason": apiErrorMarker})
}

// catchAll answers every request with obj: a path ending in id returns the
// object, CBI lists a JSON array, and any other path the paged list envelope,
// which also decodes as the object itself (for lookups such as policy sets).
func (f *unitFixture) catchAll(id string, obj any) {
	f.t.Helper()
	var fields map[string]any
	if err := json.Unmarshal([]byte(mustJSON(f.t, obj)), &fields); err != nil {
		f.t.Fatalf("catchAll: %v", err)
	}
	envelope := map[string]any{}
	for k, v := range fields {
		envelope[k] = v
	}
	envelope["list"] = []any{fields}
	envelope["totalPages"] = 1
	respond := func(r *http.Request, _ []byte) common.MockResponse {
		switch {
		case r.Method == "DELETE" || r.Method == "PUT" || r.Method == "PATCH":
			return common.NoContentResponse()
		case strings.HasSuffix(r.URL.Path, "/"+id):
			return common.SuccessResponse(fields)
		case strings.Contains(r.URL.Path, "/cbiconfig/"):
			return common.SuccessResponse([]any{fields})
		}
		return common.SuccessResponse(envelope)
	}
	for depth := 1; depth <= 14; depth++ {
		pattern := strings.Repeat("/*", depth)
		for _, method := range []string{"GET", "POST", "PUT", "PATCH", "DELETE"} {
			f.api.OnFunc(method, pattern, respond)
		}
	}
}

// crudSpec describes one resource for testCRUD: the configuration, the object
// the API returns, where it lives, and what must reach the request and state.
type crudSpec struct {
	resource *schema.Resource
	config   map[string]any
	id       string
	response any // API object returned by create and get

	createPath   string // POST
	itemPath     string // GET, PUT, DELETE
	updatePath   string // PUT, when it differs from itemPath
	deletePath   string // DELETE, when it differs from itemPath
	deleteMethod string // DELETE (default), or the method of a deactivate-style delete
	listPath     string // GET list, for resources that refresh from a list

	wantBody  map[string]any // subset expected in the create and update request bodies
	wantState map[string]any // subset expected in state after create

	setup        func(fx *unitFixture) // extra mocks (policy set lookups, dependencies)
	importPrefix string                // import ID prefix, e.g. "<namespace_id>/" for tag keys
	importLookup bool                  // the importer resolves IDs through the API too
	fullFields   map[string]any        // fields that must keep a specific shape in the fully populated object
	updateConfig func(config map[string]any)
	skipUpdate   bool // resources without an update API
}

// testCRUD runs the standard resource contract against the mock API:
// create sends the expected request and maps the response to state, read
// drops a resource that no longer exists, update sends a PUT and delete a DELETE.
func testCRUD(t *testing.T, spec crudSpec) {
	t.Helper()
	fixture := func(t *testing.T) *unitFixture {
		fx := newUnitFixture(t)
		if spec.setup != nil {
			spec.setup(fx)
		}
		return fx
	}
	// withObject registers the object at its item path and, for list-based reads, in the list.
	withObject := func(fx *unitFixture) {
		fx.on("GET", spec.itemPath, common.SuccessResponse(spec.response))
		if spec.listPath != "" {
			fx.on("GET", spec.listPath, common.SuccessResponse(common.ZPAList([]any{spec.response})))
		}
	}
	config := func() map[string]any {
		c := map[string]any{}
		for k, v := range spec.config {
			c[k] = v
		}
		return c
	}
	updatePath, deletePath, deleteMethod := spec.itemPath, spec.itemPath, "DELETE"
	if spec.deleteMethod != "" {
		deleteMethod = spec.deleteMethod
	}
	if spec.updatePath != "" {
		updatePath = spec.updatePath
	}
	if spec.deletePath != "" {
		deletePath = spec.deletePath
	}

	t.Run("Create", func(t *testing.T) {
		fx := fixture(t)
		fx.on("POST", spec.createPath, common.CreatedResponse(spec.response))
		withObject(fx)
		d := fx.resourceData(spec.resource, "", config())
		fx.run(spec.resource.CreateContext, d)
		if d.Id() != spec.id {
			t.Errorf("id = %q, want %q", d.Id(), spec.id)
		}
		assertSubset(t, "create request", fx.body("POST", spec.createPath), spec.wantBody)
		assertState(t, d, spec.wantState)
	})

	t.Run("ReadNotFoundRemovesFromState", func(t *testing.T) {
		fx := fixture(t)
		fx.on("GET", spec.itemPath, common.NotFoundResponse())
		if spec.listPath != "" {
			fx.on("GET", spec.listPath, common.SuccessResponse(common.ZPAList([]any{})))
		}
		d := fx.resourceData(spec.resource, spec.id, config())
		fx.run(spec.resource.ReadContext, d)
		if d.Id() != "" {
			t.Errorf("a resource deleted outside Terraform must be removed from state, id = %q", d.Id())
		}
	})

	if !spec.skipUpdate {
		t.Run("Update", func(t *testing.T) {
			fx := fixture(t)
			withObject(fx)
			fx.on("PUT", updatePath, common.NoContentResponse())
			c := config()
			if spec.updateConfig != nil {
				spec.updateConfig(c)
			}
			d := fx.updateData(spec.resource, spec.id, config(), c)
			fx.run(spec.resource.UpdateContext, d)
			assertSubset(t, "update request", fx.body("PUT", updatePath), spec.wantBody)
		})
	}

	t.Run("Delete", func(t *testing.T) {
		fx := fixture(t)
		withObject(fx)
		fx.on(deleteMethod, deletePath, common.NoContentResponse())
		d := fx.resourceData(spec.resource, spec.id, config())
		fx.run(spec.resource.DeleteContext, d)
		if n := fx.calls(deleteMethod, deletePath); n != 1 {
			t.Errorf("expected 1 %s %s, got %d", deleteMethod, deletePath, n)
		}
	})

	// API errors must reach the user: not swallowed, and read must not drop the resource.
	expectAPIError := func(t *testing.T, fx *unitFixture, fn crudFunc, d *schema.ResourceData) {
		t.Helper()
		if msg := fx.runErr(fn, d); !strings.Contains(msg, apiErrorMarker) {
			t.Errorf("expected the API error to be returned, got %q", msg)
		}
	}
	t.Run("CreateAPIError", func(t *testing.T) {
		fx := fixture(t)
		fx.on("POST", spec.createPath, apiError())
		d := fx.resourceData(spec.resource, "", config())
		expectAPIError(t, fx, spec.resource.CreateContext, d)
		if d.Id() != "" {
			t.Errorf("a failed create must not store an ID, got %q", d.Id())
		}
	})
	t.Run("ReadAPIError", func(t *testing.T) {
		fx := fixture(t)
		fx.on("GET", spec.itemPath, apiError())
		if spec.listPath != "" {
			fx.on("GET", spec.listPath, apiError())
		}
		d := fx.resourceData(spec.resource, spec.id, config())
		expectAPIError(t, fx, spec.resource.ReadContext, d)
		if d.Id() != spec.id {
			t.Errorf("a failed read must keep the resource in state, id = %q", d.Id())
		}
	})
	if !spec.skipUpdate {
		t.Run("UpdateAPIError", func(t *testing.T) {
			fx := fixture(t)
			withObject(fx)
			fx.on("PUT", updatePath, apiError())
			c := config()
			if spec.updateConfig != nil {
				spec.updateConfig(c)
			}
			d := fx.updateData(spec.resource, spec.id, config(), c)
			expectAPIError(t, fx, spec.resource.UpdateContext, d)
		})
	}
	t.Run("DeleteAPIError", func(t *testing.T) {
		fx := fixture(t)
		withObject(fx)
		fx.on(deleteMethod, deletePath, apiError())
		d := fx.resourceData(spec.resource, spec.id, config())
		expectAPIError(t, fx, spec.resource.DeleteContext, d)
	})

	// Reading an object with every field set maps it to state without error.
	t.Run("ReadFullObject", func(t *testing.T) {
		fx := fixture(t)
		full := populated(t, spec.response, spec.id)
		for k, v := range spec.fullFields {
			full[k] = v
		}
		fx.on("GET", spec.itemPath, common.SuccessResponse(full))
		if spec.listPath != "" {
			fx.on("GET", spec.listPath, common.SuccessResponse(common.ZPAList([]any{full})))
		}
		d := fx.resourceData(spec.resource, spec.id, config())
		fx.run(spec.resource.ReadContext, d)
		if d.Id() != spec.id {
			t.Errorf("id = %q, want %q", d.Id(), spec.id)
		}
	})

	// An update of a resource deleted outside Terraform must drop it from state
	// or fail; it must never report success with the resource still in state.
	if !spec.skipUpdate {
		t.Run("UpdateNotFound", func(t *testing.T) {
			fx := fixture(t)
			fx.on("GET", spec.itemPath, common.NotFoundResponse())
			fx.on("PUT", updatePath, common.NotFoundResponse())
			if spec.listPath != "" {
				fx.on("GET", spec.listPath, common.SuccessResponse(common.ZPAList([]any{})))
			}
			c := config()
			if spec.updateConfig != nil {
				spec.updateConfig(c)
			}
			d := fx.updateData(spec.resource, spec.id, config(), c)
			if msg := fx.runErr(spec.resource.UpdateContext, d); msg == "" && d.Id() != "" {
				t.Errorf("update of a deleted resource succeeded and kept it in state (id %q)", d.Id())
			}
		})
	}

	// With microtenant_id set, every request is scoped to the microtenant.
	if _, ok := spec.resource.Schema["microtenant_id"]; ok {
		t.Run("Microtenant", func(t *testing.T) {
			fx := newUnitFixture(t)
			fx.catchAll(spec.id, spec.response)
			c := config()
			c["microtenant_id"] = "mt1"
			for _, step := range []struct {
				name string
				fn   crudFunc
				d    *schema.ResourceData
			}{
				{"create", spec.resource.CreateContext, fx.resourceData(spec.resource, "", c)},
				{"read", spec.resource.ReadContext, fx.resourceData(spec.resource, spec.id, c)},
				{"update", spec.resource.UpdateContext, fx.updateData(spec.resource, spec.id, c, c)},
				{"delete", spec.resource.DeleteContext, fx.resourceData(spec.resource, spec.id, c)},
			} {
				if step.fn == nil || step.name == "update" && spec.skipUpdate {
					continue
				}
				before := len(fx.api.Server.AllRequests())
				if msg := fx.runErr(step.fn, step.d); msg != "" {
					t.Errorf("%s: %s", step.name, msg)
				}
				for _, r := range fx.api.Server.AllRequests()[before:] {
					if !strings.Contains(r.Query, "microtenantId=mt1") {
						t.Errorf("%s: %s %s is not scoped to the microtenant", step.name, r.Method, r.Path)
					}
				}
			}
		})
	}

	if spec.resource.Importer != nil {
		testImport(t, spec)
	}
}

// testImport imports by numeric ID and, when the object has a name, by name.
func testImport(t *testing.T, spec crudSpec) {
	t.Helper()
	r, id, obj := spec.resource, spec.id, spec.response
	importer := r.Importer.StateContext
	t.Run("ImportByID", func(t *testing.T) {
		importID := "72058304855015424"
		if spec.importLookup {
			importID = id
		}
		fx := newUnitFixture(t)
		fx.catchAll(id, obj)
		d := r.Data(nil)
		d.SetId(spec.importPrefix + importID)
		got, err := importer(context.Background(), d, fx.client)
		if err != nil || len(got) != 1 || got[0].Id() != importID {
			t.Errorf("import by ID: err=%v ids=%v, want %q", err, importedIDs(got), importID)
		}
	})
	var fields map[string]any
	_ = json.Unmarshal([]byte(mustJSON(t, obj)), &fields)
	name, _ := fields["name"].(string)
	if name == "" || reflect.ValueOf(importer).Pointer() == reflect.ValueOf(schema.ImportStatePassthroughContext).Pointer() {
		return
	}
	t.Run("ImportByUnknownName", func(t *testing.T) {
		fx := newUnitFixture(t)
		d := r.Data(nil)
		d.SetId(spec.importPrefix + "does-not-exist")
		if _, err := importer(context.Background(), d, fx.client); err == nil {
			t.Error("importing a name that does not exist must fail")
		}
	})
	t.Run("ImportByName", func(t *testing.T) {
		fx := newUnitFixture(t)
		fx.catchAll(id, obj)
		d := r.Data(nil)
		d.SetId(spec.importPrefix + name)
		got, err := importer(context.Background(), d, fx.client)
		if err != nil || len(got) != 1 || got[0].Id() != id {
			t.Errorf("import by name %q: err=%v ids=%v, want %q", name, err, importedIDs(got), id)
		}
	})
}

func importedIDs(ds []*schema.ResourceData) []string {
	var ids []string
	for _, d := range ds {
		ids = append(ids, d.Id())
	}
	return ids
}

// dataSourceSpec describes one data source for testDataSource.
type dataSourceSpec struct {
	resource   *schema.Resource
	listPath   string // list used by the name lookup
	listMethod string // GET (default) or POST for search endpoints
	plainList  bool   // the list endpoint returns a JSON array instead of the paged envelope
	itemPath   string // GET by ID; empty when the data source only looks up by name
	idFromList bool   // the ID lookup searches the list instead of reading the item
	items      any    // list returned by listPath ([]T); the match must not be the first item
	item       any    // object returned by itemPath
	name       string // name to look up; must match an item case-insensitively
	id         string // ID of the matching item
	extra      map[string]any
	setup      func(fx *unitFixture)
	fullFields map[string]any // fields that must keep a specific shape in the fully populated object

	wantState map[string]any
}

// testDataSource checks the name lookup picks the matching item, the ID
// lookup reads the item by ID, and an unknown name is an error.
func testDataSource(t *testing.T, spec dataSourceSpec) {
	t.Helper()
	fixture := func(t *testing.T) *unitFixture {
		fx := newUnitFixture(t)
		if spec.setup != nil {
			spec.setup(fx)
		}
		return fx
	}
	withExtra := func(c map[string]any) map[string]any {
		for k, v := range spec.extra {
			c[k] = v
		}
		return c
	}
	read := spec.resource.ReadContext
	listMethod := spec.listMethod
	if listMethod == "" {
		listMethod = "GET"
	}
	list := func(items []any) common.MockResponse {
		if spec.plainList {
			return common.SuccessResponse(items)
		}
		return common.SuccessResponse(common.ZPAList(items))
	}

	t.Run("ByName", func(t *testing.T) {
		fx := fixture(t)
		fx.on(listMethod, spec.listPath, list(toAnySlice(spec.items)))
		if spec.itemPath != "" {
			fx.on("GET", spec.itemPath, common.SuccessResponse(spec.item))
		}
		d := fx.resourceData(spec.resource, "", withExtra(map[string]any{"name": spec.name}))
		fx.run(read, d)
		if d.Id() != spec.id {
			t.Errorf("id = %q, want %q", d.Id(), spec.id)
		}
		assertState(t, d, spec.wantState)
	})

	if spec.itemPath != "" || spec.idFromList {
		t.Run("ByID", func(t *testing.T) {
			fx := fixture(t)
			if spec.idFromList {
				fx.on(listMethod, spec.listPath, list(toAnySlice(spec.items)))
			} else {
				fx.on("GET", spec.itemPath, common.SuccessResponse(spec.item))
			}
			d := fx.resourceData(spec.resource, spec.id, withExtra(map[string]any{}))
			fx.run(read, d)
			if d.Id() != spec.id {
				t.Errorf("id = %q, want %q", d.Id(), spec.id)
			}
			assertState(t, d, spec.wantState)
		})
	}

	t.Run("NameNotFound", func(t *testing.T) {
		fx := fixture(t)
		fx.on(listMethod, spec.listPath, list([]any{}))
		d := fx.resourceData(spec.resource, "", withExtra(map[string]any{"name": "does-not-exist"}))
		if fx.runErr(read, d) == "" {
			t.Error("expected an error for a name that does not exist")
		}
	})

	// Reading an object with every field set maps it to state without error.
	t.Run("FullObject", func(t *testing.T) {
		fx := fixture(t)
		sample := spec.item
		if sample == nil {
			sample = toAnySlice(spec.items)[0]
		}
		full := populated(t, sample, spec.id)
		full["name"] = spec.name
		for k, v := range spec.fullFields {
			full[k] = v
		}
		fx.on(listMethod, spec.listPath, list([]any{full}))
		if spec.itemPath != "" {
			fx.on("GET", spec.itemPath, common.SuccessResponse(full))
		}
		d := fx.resourceData(spec.resource, "", withExtra(map[string]any{"name": spec.name}))
		fx.run(read, d)
	})

	// API errors must be returned, not reported as "not found" or swallowed.
	t.Run("ByNameAPIError", func(t *testing.T) {
		fx := fixture(t)
		fx.on(listMethod, spec.listPath, apiError())
		d := fx.resourceData(spec.resource, "", withExtra(map[string]any{"name": spec.name}))
		if msg := fx.runErr(read, d); !strings.Contains(msg, apiErrorMarker) {
			t.Errorf("expected the API error to be returned, got %q", msg)
		}
	})
	if spec.itemPath != "" || spec.idFromList {
		t.Run("ByIDAPIError", func(t *testing.T) {
			fx := fixture(t)
			if spec.idFromList {
				fx.on(listMethod, spec.listPath, apiError())
			} else {
				fx.on("GET", spec.itemPath, apiError())
			}
			d := fx.resourceData(spec.resource, spec.id, withExtra(map[string]any{}))
			if msg := fx.runErr(read, d); !strings.Contains(msg, apiErrorMarker) {
				t.Errorf("expected the API error to be returned, got %q", msg)
			}
		})
	}
}

// assertSubset fails when want is not contained in got, compared as JSON:
// maps by key, lists by length with each wanted element matching one in got.
func assertSubset(t *testing.T, what string, got any, want map[string]any) {
	t.Helper()
	if len(want) == 0 {
		return
	}
	if !jsonSubset(normalizeJSON(t, got), normalizeJSON(t, want)) {
		t.Errorf("%s does not contain the expected fields\n got: %s\nwant: %s", what, mustJSON(t, expandSets(got)), mustJSON(t, want))
	}
}

// assertState checks the given attributes in state; sets are compared as lists.
func assertState(t *testing.T, d *schema.ResourceData, want map[string]any) {
	t.Helper()
	for key, w := range want {
		got := expandSets(d.Get(key))
		if !jsonSubset(normalizeJSON(t, got), normalizeJSON(t, w)) {
			t.Errorf("state %s = %s, want %s", key, mustJSON(t, got), mustJSON(t, w))
		}
	}
}

func jsonSubset(got, want any) bool {
	switch w := want.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			return false
		}
		for k, wv := range w {
			if !jsonSubset(g[k], wv) {
				return false
			}
		}
		return true
	case []any:
		g, ok := got.([]any)
		if !ok || len(g) != len(w) {
			return false
		}
		for _, wv := range w {
			found := false
			for _, gv := range g {
				if jsonSubset(gv, wv) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(got, want)
	}
}

func normalizeJSON(t *testing.T, v any) any {
	t.Helper()
	var out any
	if err := json.Unmarshal([]byte(mustJSON(t, expandSets(v))), &out); err != nil {
		t.Fatalf("normalize: %v", err)
	}
	return out
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(b)
}

// expandSets replaces *schema.Set values (also nested in state lists and maps) by their elements.
func expandSets(v any) any {
	switch x := v.(type) {
	case *schema.Set:
		return expandSets(x.List())
	case []any:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = expandSets(e)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, e := range x {
			out[k] = expandSets(e)
		}
		return out
	}
	return v
}

// populated returns a copy of obj (an SDK struct) with every field set: strings
// to a number-like value (valid as an ID, a count or an epoch timestamp), bools
// to true, one element in every slice and map, and every pointer allocated. The
// result is returned as JSON fields, with "id" set to id.
func populated(t *testing.T, obj any, id string) map[string]any {
	t.Helper()
	v := fillValue(reflect.TypeOf(obj), 0)
	var fields map[string]any
	if err := json.Unmarshal([]byte(mustJSON(t, v.Interface())), &fields); err != nil {
		t.Fatalf("populated: %v", err)
	}
	fields["id"] = id
	return fields
}

func fillValue(t reflect.Type, depth int) reflect.Value {
	v := reflect.New(t).Elem()
	if depth > 6 {
		return v
	}
	switch t.Kind() {
	case reflect.String:
		v.SetString("1700000000")
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(7)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(7)
	case reflect.Float32, reflect.Float64:
		v.SetFloat(7)
	case reflect.Ptr:
		v.Set(fillValue(t.Elem(), depth+1).Addr())
	case reflect.Slice:
		s := reflect.MakeSlice(t, 1, 1)
		s.Index(0).Set(fillValue(t.Elem(), depth+1))
		v.Set(s)
	case reflect.Map:
		m := reflect.MakeMap(t)
		m.SetMapIndex(fillValue(t.Key(), depth+1), fillValue(t.Elem(), depth+1))
		v.Set(m)
	case reflect.Interface:
		v.Set(reflect.ValueOf("1700000000"))
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).IsExported() {
				v.Field(i).Set(fillValue(t.Field(i).Type, depth+1))
			}
		}
	}
	return v
}

// toAnySlice converts a typed slice ([]T) to []any for the list envelope.
func toAnySlice(items any) []any {
	v := reflect.ValueOf(items)
	if v.Kind() != reflect.Slice {
		return nil
	}
	out := make([]any, v.Len())
	for i := range out {
		out[i] = v.Index(i).Interface()
	}
	return out
}
