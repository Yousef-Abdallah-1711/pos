package main

import (
	"os"
	"regexp"
	"strconv"
	"testing"
)

type phase6RoutePermission struct {
	moduleID   int
	permission string
}

func TestPhase6Routes_AreWrappedWithExpectedRBACModules(t *testing.T) {
	content, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("failed to read main.go: %v", err)
	}

	routePermissions := make(map[string]phase6RoutePermission)
	re := regexp.MustCompile(`mux\.Handle\("([^"]+)",\s*wrapWithPermission\([^,]+,\s*([0-9]+),\s*"([A-Z]+)"\)\)`)
	matches := re.FindAllStringSubmatch(string(content), -1)
	for _, m := range matches {
		moduleID, convErr := strconv.Atoi(m[2])
		if convErr != nil {
			t.Fatalf("failed parsing module id %q for route %q: %v", m[2], m[1], convErr)
		}
		routePermissions[m[1]] = phase6RoutePermission{
			moduleID:   moduleID,
			permission: m[3],
		}
	}

	required := map[string]phase6RoutePermission{
		// Accounting (module 15)
		"GET /api/v1/accounting/accounts":                   {moduleID: 15, permission: "READ"},
		"POST /api/v1/accounting/accounts":                  {moduleID: 15, permission: "WRITE"},
		"DELETE /api/v1/accounting/accounts/{id}":           {moduleID: 15, permission: "DELETE"},
		"GET /api/v1/accounting/journal-entries":            {moduleID: 15, permission: "READ"},
		"POST /api/v1/accounting/journal-entries/{id}/post": {moduleID: 15, permission: "WRITE"},
		"POST /api/v1/accounting/journal-entries/{id}/void": {moduleID: 15, permission: "WRITE"},
		"GET /api/v1/accounting/ledger":                     {moduleID: 15, permission: "READ"},
		"GET /api/v1/accounting/trial-balance":              {moduleID: 15, permission: "READ"},
		"GET /api/v1/accounting/bank-accounts":              {moduleID: 15, permission: "READ"},
		"POST /api/v1/accounting/bank-accounts":             {moduleID: 15, permission: "WRITE"},
		"GET /api/v1/accounting/customers":                  {moduleID: 15, permission: "READ"},
		"POST /api/v1/accounting/customers":                 {moduleID: 15, permission: "WRITE"},
		"GET /api/v1/accounting/customer-invoices":          {moduleID: 15, permission: "READ"},
		"POST /api/v1/accounting/customer-invoices":         {moduleID: 15, permission: "WRITE"},
		"GET /api/v1/accounting/supplier-invoices":          {moduleID: 15, permission: "READ"},
		"POST /api/v1/accounting/supplier-invoices":         {moduleID: 15, permission: "WRITE"},
		"GET /api/v1/accounting/payments/received":          {moduleID: 15, permission: "READ"},
		"POST /api/v1/accounting/payments/received":         {moduleID: 15, permission: "WRITE"},
		"GET /api/v1/accounting/payments/made":              {moduleID: 15, permission: "READ"},
		"POST /api/v1/accounting/payments/made":             {moduleID: 15, permission: "WRITE"},

		// Supplier (module 16)
		"GET /api/v1/suppliers":         {moduleID: 16, permission: "READ"},
		"POST /api/v1/suppliers":        {moduleID: 16, permission: "WRITE"},
		"DELETE /api/v1/suppliers/{id}": {moduleID: 16, permission: "DELETE"},

		// Financial reports + tax (module 17)
		"GET /api/v1/reports/financial/profit-loss":      {moduleID: 17, permission: "READ"},
		"GET /api/v1/reports/financial/balance-sheet":    {moduleID: 17, permission: "READ"},
		"GET /api/v1/reports/financial/cash-flow":        {moduleID: 17, permission: "READ"},
		"GET /api/v1/reports/financial/trial-balance":    {moduleID: 17, permission: "READ"},
		"GET /api/v1/accounting/tax/periods":             {moduleID: 17, permission: "READ"},
		"POST /api/v1/accounting/tax/periods":            {moduleID: 17, permission: "WRITE"},
		"GET /api/v1/accounting/tax/vat-report":          {moduleID: 17, permission: "READ"},
		"POST /api/v1/accounting/tax/periods/{id}/close": {moduleID: 17, permission: "WRITE"},

		// Advanced HR modules 18/19/20
		"GET /api/v1/hr/performance":                   {moduleID: 18, permission: "READ"},
		"POST /api/v1/hr/performance":                  {moduleID: 18, permission: "WRITE"},
		"POST /api/v1/hr/performance/{id}/acknowledge": {moduleID: 18, permission: "WRITE"},
		"GET /api/v1/hr/training":                      {moduleID: 19, permission: "READ"},
		"POST /api/v1/hr/training":                     {moduleID: 19, permission: "WRITE"},
		"GET /api/v1/hr/training/expiring":             {moduleID: 19, permission: "READ"},
		"GET /api/v1/hr/insurance":                     {moduleID: 20, permission: "READ"},
		"POST /api/v1/hr/insurance":                    {moduleID: 20, permission: "WRITE"},
		"GET /api/v1/hr/contracts":                     {moduleID: 20, permission: "READ"},
		"POST /api/v1/hr/contracts":                    {moduleID: 20, permission: "WRITE"},
		"POST /api/v1/hr/contracts/{id}/renew":         {moduleID: 20, permission: "WRITE"},
		"GET /api/v1/hr/resignations":                  {moduleID: 20, permission: "READ"},
		"POST /api/v1/hr/resignations":                 {moduleID: 20, permission: "WRITE"},
		"GET /api/v1/hr/resignations/{id}/settlement":  {moduleID: 20, permission: "READ"},

		// HR reports (module 8 reporting)
		"GET /api/v1/reports/hr/headcount":           {moduleID: 8, permission: "READ"},
		"GET /api/v1/reports/hr/turnover":            {moduleID: 8, permission: "READ"},
		"GET /api/v1/reports/hr/absenteeism":         {moduleID: 8, permission: "READ"},
		"GET /api/v1/reports/hr/overtime":            {moduleID: 8, permission: "READ"},
		"GET /api/v1/reports/hr/leave-balances":      {moduleID: 8, permission: "READ"},
		"GET /api/v1/reports/hr/training-compliance": {moduleID: 8, permission: "READ"},
	}

	for route, expected := range required {
		actual, ok := routePermissions[route]
		if !ok {
			t.Fatalf("route %q is missing wrapWithPermission registration", route)
		}
		if actual.moduleID != expected.moduleID || actual.permission != expected.permission {
			t.Fatalf(
				"route %q expected module=%d permission=%s, got module=%d permission=%s",
				route,
				expected.moduleID,
				expected.permission,
				actual.moduleID,
				actual.permission,
			)
		}
	}
}
