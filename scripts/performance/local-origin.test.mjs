import assert from "node:assert/strict";
import test from "node:test";
import { requireLocalOrigin } from "./local-origin.mjs";

test("accepts HTTP loopback origins", () => {
  assert.equal(requireLocalOrigin("http://localhost:3000", "P2C_BASE_URL").origin, "http://localhost:3000");
  assert.equal(requireLocalOrigin("http://127.0.0.1:8080", "P2C_API_BASE_URL").hostname, "127.0.0.1");
  assert.equal(requireLocalOrigin("http://[::1]:8080", "P2C_API_BASE_URL").hostname, "[::1]");
});

test("rejects non-loopback hosts", () => {
  assert.throws(() => requireLocalOrigin("https://example.com", "P2C_BASE_URL"), /must use localhost/);
  assert.throws(() => requireLocalOrigin("http://192.168.1.10:3000", "P2C_BASE_URL"), /must use localhost/);
});

test("rejects credentials and URL components beyond the origin", () => {
  assert.throws(() => requireLocalOrigin("http://user:pass@localhost:3000", "P2C_BASE_URL"), /without credentials/);
  assert.throws(() => requireLocalOrigin("http://localhost:3000/api", "P2C_BASE_URL"), /without credentials/);
  assert.throws(() => requireLocalOrigin("http://localhost:3000/?token=x", "P2C_BASE_URL"), /without credentials/);
  assert.throws(() => requireLocalOrigin("file:///tmp/app", "P2C_BASE_URL"), /must use localhost/);
});
