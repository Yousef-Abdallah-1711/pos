import { describe, expect, it } from "vitest";
import { getSafeAuditDetails } from "./audit-log-details";

describe("getSafeAuditDetails", () => {
  it.each([null, undefined, "invalid", 42, ["unexpected"]])(
    "ignores non-object details: %s",
    (details) => {
      expect(getSafeAuditDetails(details)).toBeNull();
    }
  );

  it("keeps ordinary fields and redacts sensitive fields recursively", () => {
    expect(
      getSafeAuditDetails({
        reason: "manual review",
        password: "do-not-render",
        metadata: [{ access_token: "also-secret", tenant_id: 7 }],
      })
    ).toEqual({
      reason: "manual review",
      password: "[REDACTED]",
      metadata: [{ access_token: "[REDACTED]", tenant_id: 7 }],
    });
  });

  it("returns an empty object without inventing details", () => {
    expect(getSafeAuditDetails({})).toEqual({});
  });
});
