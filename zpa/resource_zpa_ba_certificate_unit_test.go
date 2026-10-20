package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
)

func TestResourceBaCertificate_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource: resourceBaCertificate(),
		config: map[string]any{
			"name":        "portal.acme.com",
			"description": "Portal certificate",
			"cert_blob":   "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----",
		},
		id: "c1",
		response: bacertificate.BaCertificate{
			ID: "c1", Name: "portal.acme.com", Description: "Portal certificate", Certificate: "-----BEGIN CERTIFICATE-----",
		},
		createPath: v1Path("certificate"),
		itemPath:   v1Path("certificate", "c1"),
		// the certificate and key are uploaded as certBlob
		wantBody: map[string]any{
			"name":     "portal.acme.com",
			"certBlob": "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----",
		},
		wantState: map[string]any{"name": "portal.acme.com", "description": "Portal certificate"},
		// certificates are immutable: any change replaces the resource
		skipUpdate: true,
	})
}
