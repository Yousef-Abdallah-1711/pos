import { renderHook, waitFor } from "@testing-library/react";
import {
  useTheme,
  useThemeColors,
  useThemeHeader,
  useThemeFooter,
  useThemeIdentity,
  useThemeTypography,
  useThemeLoading,
  useThemeError,
} from "./use-theme";
import { useThemeStore } from "@/lib/store/theme-store";
import type { ThemeData } from "@/lib/store/theme-store";
import { getDefaultTheme } from "@/lib/utils/default-theme";

// Mock the theme store
jest.mock("@/lib/store/theme-store");

// Helper to create a complete mock store
type MockThemeStore = {
  currentTheme: ThemeData | null;
  isLoading: boolean;
  error: string | null;
  loadTheme: jest.Mock;
  setTheme: jest.Mock;
  setError: jest.Mock;
  clearError: jest.Mock;
};

const createMockStore = (overrides: Partial<MockThemeStore> = {}): MockThemeStore => ({
  currentTheme: null,
  isLoading: false,
  error: null,
  loadTheme: jest.fn().mockResolvedValue(undefined),
  setTheme: jest.fn(),
  setError: jest.fn(),
  clearError: jest.fn(),
  ...overrides,
});

let mockStoreState: MockThemeStore = createMockStore();
const mockUseThemeStore = useThemeStore as unknown as jest.Mock;

describe("Theme Hooks", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockStoreState = createMockStore();
    mockUseThemeStore.mockImplementation((selector?: (state: MockThemeStore) => unknown) =>
      selector ? selector(mockStoreState) : mockStoreState
    );
  });

  describe("useTheme", () => {
    it("should return theme data from store", () => {
      const theme = getDefaultTheme();
      mockStoreState = createMockStore({ currentTheme: theme });

      const { result } = renderHook(() => useTheme());

      expect(result.current.currentTheme).toEqual(theme);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("should indicate loading state", () => {
      mockStoreState = createMockStore({ isLoading: true });

      const { result } = renderHook(() => useTheme());

      expect(result.current.isLoading).toBe(true);
    });

    it("should return error state", () => {
      mockStoreState = createMockStore({ error: "Failed to load theme" });

      const { result } = renderHook(() => useTheme());

      expect(result.current.error).toBe("Failed to load theme");
    });
  });

  describe("useThemeColors", () => {
    it("should return color palette from theme", () => {
      const theme = getDefaultTheme();
      mockStoreState = createMockStore({ currentTheme: theme });

      const { result } = renderHook(() => useThemeColors());

      expect(result.current).toHaveProperty("primary");
      expect(result.current).toHaveProperty("secondary");
      expect(result.current).toHaveProperty("background");
      expect(result.current).toHaveProperty("text");
    });

  it("should return fallback colors and a loading flag if theme not loaded", () => {
      mockStoreState = createMockStore();

      const { result } = renderHook(() => useThemeColors());

    expect(result.current).toMatchObject({
      primary: "#f97316",
      secondary: "#0ea5e9",
      isLoading: true,
    });
    });
  });

  describe("useThemeHeader", () => {
    it("should return header configuration from theme", () => {
      const theme = getDefaultTheme();
      mockStoreState = createMockStore({ currentTheme: theme });

      const { result } = renderHook(() => useThemeHeader());

      expect(result.current).toBeDefined();
      expect(result.current).toHaveProperty("backgroundColor");
      expect(result.current).toHaveProperty("textColor");
    });
  });

  describe("useThemeFooter", () => {
    it("should return footer configuration from theme", () => {
      const theme = getDefaultTheme();
      mockStoreState = createMockStore({ currentTheme: theme });

      const { result } = renderHook(() => useThemeFooter());

      expect(result.current).toBeDefined();
      expect(result.current).toHaveProperty("backgroundColor");
    });
  });

  describe("useThemeIdentity", () => {
    it("should return brand identity from theme", () => {
      const theme = getDefaultTheme();
      mockStoreState = createMockStore({ currentTheme: theme });

      const { result } = renderHook(() => useThemeIdentity());

      expect(result.current).toBeDefined();
      expect(result.current).toHaveProperty("siteTitle");
    });
  });

  describe("useThemeTypography", () => {
    it("should return typography settings from theme", () => {
      const theme = getDefaultTheme();
      mockStoreState = createMockStore({ currentTheme: theme });

      const { result } = renderHook(() => useThemeTypography());

      expect(result.current).toBeDefined();
      expect(result.current).toHaveProperty("fontFamily");
      expect(result.current).toHaveProperty("baseFontSize");
    });
  });

  describe("useThemeLoading", () => {
    it("should return loading state", () => {
      mockStoreState = createMockStore({ isLoading: true });

      const { result } = renderHook(() => useThemeLoading());

      expect(result.current).toBe(true);
    });

    it("should return false when not loading", () => {
      mockStoreState = createMockStore({ isLoading: false });

      const { result } = renderHook(() => useThemeLoading());

      expect(result.current).toBe(false);
    });
  });

  describe("useThemeError", () => {
    it("should return error message", () => {
      mockStoreState = createMockStore({ error: "Theme not found" });

      const { result } = renderHook(() => useThemeError());

      expect(result.current).toBe("Theme not found");
    });

    it("should return null when no error", () => {
      mockStoreState = createMockStore({ error: null });

      const { result } = renderHook(() => useThemeError());

      expect(result.current).toBeNull();
    });
  });

  describe("Hook reactive updates", () => {
    it("should update when theme changes", () => {
      const theme1 = getDefaultTheme();
      const mockStore = {
        currentTheme: theme1,
        isLoading: false,
        error: null,
        loadTheme: jest.fn(),
        setTheme: jest.fn(),
        setError: jest.fn(),
        clearError: jest.fn(),
      };

      mockStoreState = mockStore;

      const { result, rerender } = renderHook(() => useTheme());

      expect(result.current.currentTheme).toEqual(theme1);

      // Simulate theme change
      mockStore.currentTheme = { ...theme1, name: "Updated Theme" };
      mockStoreState = mockStore;

      rerender();

      expect(result.current.currentTheme?.name).toBe("Updated Theme");
    });

    it("does not reload when unrelated store state changes", async () => {
      const store = createMockStore();
      mockStoreState = store;

      const { rerender } = renderHook(() => useTheme());
      await waitFor(() => expect(store.loadTheme).toHaveBeenCalledTimes(1));

      mockStoreState = { ...store, isLoading: true };
      rerender();
      mockStoreState = { ...store, isLoading: false, error: "Theme unavailable" };
      rerender();

      expect(store.loadTheme).toHaveBeenCalledTimes(1);
    });
  });
});
