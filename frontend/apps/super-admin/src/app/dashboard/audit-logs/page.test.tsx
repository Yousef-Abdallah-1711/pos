import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import AuditLogsPage from "./page";
import type { AuditLog } from "@/types";

const { getAuditLogsMock } = vi.hoisted(() => ({
  getAuditLogsMock: vi.fn(),
}));

vi.mock("@/lib/api", () => ({
  superAdminAPI: {
    getAuditLogs: getAuditLogsMock,
  },
}));

function createAuditLog(details?: unknown): AuditLog {
  return {
    id: 1,
    super_admin_id: 1,
    action: "LOGIN_SUCCESS",
    target_user_id: null,
    target_tenant_id: null,
    ip_address: "198.51.100.8",
    user_agent: "synthetic-browser",
    details,
    created_at: "2026-11-03T12:00:00.000Z",
  } as AuditLog;
}

function renderWithDetails(details?: unknown) {
  getAuditLogsMock.mockResolvedValue({
    data: [createAuditLog(details)],
    total: 1,
    page: 1,
    limit: 50,
    pages: 1,
  });

  return render(<AuditLogsPage />);
}

describe("AuditLogsPage", () => {
  beforeEach(() => {
    getAuditLogsMock.mockReset();
  });

  afterEach(() => {
    cleanup();
  });

  it.each([
    ["null", null],
    ["missing", undefined],
    ["non-object", "unexpected"],
    ["array", ["unexpected"]],
  ])("renders a log with %s details without crashing", async (_label, details) => {
    renderWithDetails(details);

    expect(await screen.findByRole("heading", { name: "LOGIN_SUCCESS" })).not.toBeNull();
    expect(screen.getByText("198.51.100.8")).not.toBeNull();
    expect(screen.queryByText("Details")).toBeNull();
  });

  it("keeps an empty details object empty and hides the details section", async () => {
    renderWithDetails({});

    expect(await screen.findByRole("heading", { name: "LOGIN_SUCCESS" })).not.toBeNull();
    expect(screen.queryByText("Details")).toBeNull();
  });

  it("renders ordinary details and redacts sensitive values", async () => {
    renderWithDetails({
      reason: "manual review",
      password: "do-not-render",
      nested: { access_token: "also-secret", tenant_id: 7 },
    });

    expect(await screen.findByText("Details")).not.toBeNull();
    expect(screen.getByText(/"reason":"manual review"/)).not.toBeNull();
    expect(screen.getByText(/"password":"\[REDACTED\]"/)).not.toBeNull();
    expect(screen.getByText(/"access_token":"\[REDACTED\]"/)).not.toBeNull();
    expect(screen.queryByText(/do-not-render|also-secret/)).toBeNull();
  });
});
