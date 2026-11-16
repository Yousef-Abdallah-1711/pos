import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import "@testing-library/jest-dom/vitest";
import { afterEach, describe, expect, it } from "vitest";
import { PreferenceControls } from "@/components/layout/PreferenceControls";
import { PreferencesProvider, useSuperAdminPreferences } from "./PreferencesProvider";

function CurrentPreferences() {
  const { locale, theme } = useSuperAdminPreferences();
  return <output>{`${locale}:${theme}`}</output>;
}

afterEach(() => {
  cleanup();
  localStorage.clear();
  document.documentElement.className = "";
  document.documentElement.lang = "en";
  document.documentElement.dir = "ltr";
});

describe("Super Admin preferences", () => {
  it("restores and persists Arabic direction and explicit theme choices", async () => {
    localStorage.setItem("servio-super-admin-preferences", JSON.stringify({ locale: "ar", theme: "dark" }));
    render(
      <PreferencesProvider>
        <CurrentPreferences />
        <PreferenceControls />
      </PreferencesProvider>,
    );

    await waitFor(() => {
      expect(document.documentElement.lang).toBe("ar");
      expect(document.documentElement.dir).toBe("rtl");
      expect(document.documentElement).toHaveClass("dark");
    });
    expect(screen.getByText("ar:dark")).toBeVisible();

    fireEvent.click(screen.getByRole("button", { name: "التغيير إلى الإنجليزية" }));
    await waitFor(() => expect(document.documentElement.dir).toBe("ltr"));
    fireEvent.click(screen.getByRole("button", { name: "Switch to light mode" }));
    await waitFor(() => expect(document.documentElement).not.toHaveClass("dark"));
    expect(JSON.parse(localStorage.getItem("servio-super-admin-preferences") || "{}")).toEqual({ locale: "en", theme: "light" });
  });
});
