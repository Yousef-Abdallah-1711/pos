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
import { getDefaultTheme } from "@/lib/utils/default-theme";

// Mock the theme store
jest.mock("@/lib/store/theme-store");

// Helper to create a complete mock store
const createMockStore = (overrides = {}) => ({
  currentTheme: null,
  isLoading: false,
  error: null,
  loadTheme: jest.fn().mockResolvedValue(undefined),
  setTheme: jest.fn(),
  setError: jest.fn(),
  clearError: jest.fn(),
  ...overrides,
});

describe("Theme Hooks", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    // Set up default mock store for all tests
    (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore());
  });

  describe("useTheme", () => {
    it("should return theme data from store", () => {
      const theme = getDefaultTheme();
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore({ currentTheme: theme }));

      const { result } = renderHook(() => useTheme());

      expect(result.current.currentTheme).toEqual(theme);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("should indicate loading state", () => {
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore({ isLoading: true }));

      const { result } = renderHook(() => useTheme());

      expect(result.current.isLoading).toBe(true);
    });

    it("should return error state", () => {
      (useThemeStore as unknown as jest.Mock).mockReturnValue(
        createMockStore({ error: "Failed to load theme" })
      );

      const { result } = renderHook(() => useTheme());

      expect(result.current.error).toBe("Failed to load theme");
    });
  });

  describe("useThemeColors", () => {
    it("should return color palette from theme", () => {
      const theme = getDefaultTheme();
      (useThemeStore as unknown as jest.Mock).mockReturnValue(
        createMockStore({ currentTheme: theme })
      );

      const { result } = renderHook(() => useThemeColors());

      expect(result.current).toHaveProperty("primary");
      expect(result.current).toHaveProperty("secondary");
      expect(result.current).toHaveProperty("background");
      expect(result.current).toHaveProperty("text");
    });

  it("should return fallback colors and a loading flag if theme not loaded", () => {
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore());

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
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore({ currentTheme: theme }));

      const { result } = renderHook(() => useThemeHeader());

      expect(result.current).toBeDefined();
      expect(result.current).toHaveProperty("backgroundColor");
      expect(result.current).toHaveProperty("textColor");
    });
  });

  describe("useThemeFooter", () => {
    it("should return footer configuration from theme", () => {
      const theme = getDefaultTheme();
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore({ currentTheme: theme }));

      const { result } = renderHook(() => useThemeFooter());

      expect(result.current).toBeDefined();
      expect(result.current).toHaveProperty("backgroundColor");
    });
  });

  describe("useThemeIdentity", () => {
    it("should return brand identity from theme", () => {
      const theme = getDefaultTheme();
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore({ currentTheme: theme }));

      const { result } = renderHook(() => useThemeIdentity());

      expect(result.current).toBeDefined();
      expect(result.current).toHaveProperty("siteTitle");
    });
  });

  describe("useThemeTypography", () => {
    it("should return typography settings from theme", () => {
      const theme = getDefaultTheme();
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore({ currentTheme: theme }));

      const { result } = renderHook(() => useThemeTypography());

      expect(result.current).toBeDefined();
      expect(result.current).toHaveProperty("fontFamily");
      expect(result.current).toHaveProperty("baseFontSize");
    });
  });

  describe("useThemeLoading", () => {
    it("should return loading state", () => {
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore({ isLoading: true }));

      const { result } = renderHook(() => useThemeLoading());

      expect(result.current).toBe(true);
    });

    it("should return false when not loading", () => {
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore({ isLoading: false }));

      const { result } = renderHook(() => useThemeLoading());

      expect(result.current).toBe(false);
    });
  });

  describe("useThemeError", () => {
    it("should return error message", () => {
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore({ error: "Theme not found" }));

      const { result } = renderHook(() => useThemeError());

      expect(result.current).toBe("Theme not found");
    });

    it("should return null when no error", () => {
      (useThemeStore as unknown as jest.Mock).mockReturnValue(createMockStore({ error: null }));

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

      (useThemeStore as unknown as jest.Mock).mockReturnValue(mockStore);

      const { result, rerender } = renderHook(() => useTheme());

      expect(result.current.currentTheme).toEqual(theme1);

      // Simulate theme change
      mockStore.currentTheme = { ...theme1, name: "Updated Theme" };
      (useThemeStore as unknown as jest.Mock).mockReturnValue(mockStore);

      rerender();

      expect(result.current.currentTheme?.name).toBe("Updated Theme");
    });
  });
});
