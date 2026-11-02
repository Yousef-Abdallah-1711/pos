import {
  hexToRgb,
  rgbToHsl,
  hexToHsl,
  getLuminance,
  getContrastRatio,
  meetsWCAGAA,
  meetsWCAGAAA,
} from "./theme-colors";

describe("Color Utilities", () => {
  describe("hexToRgb", () => {
    it("should convert hex to RGB array", () => {
      expect(hexToRgb("#ffffff")).toEqual([255, 255, 255]);
      expect(hexToRgb("#000000")).toEqual([0, 0, 0]);
      expect(hexToRgb("#ff0000")).toEqual([255, 0, 0]);
    });

    it("should handle shorthand hex colors", () => {
      expect(hexToRgb("#fff")).toEqual([255, 255, 255]);
      expect(hexToRgb("#000")).toEqual([0, 0, 0]);
    });
  });

  describe("rgbToHsl", () => {
    it("should convert RGB to HSL string", () => {
      const hsl = rgbToHsl(255, 255, 255);
      expect(hsl).toMatch(/^\d+ \d+% \d+%$/);
    });

    it("should handle pure colors correctly", () => {
      // White
      const white = rgbToHsl(255, 255, 255);
      expect(white).toContain("100%");

      // Black
      const black = rgbToHsl(0, 0, 0);
      expect(black).toContain("0%");
    });
  });

  describe("hexToHsl", () => {
    it("should convert hex to HSL string", () => {
      const hsl = hexToHsl("#ffffff");
      expect(hsl).toMatch(/^\d+ \d+% \d+%$/);
    });

    it("should handle common colors", () => {
      expect(hexToHsl("#ffffff")).toBeDefined();
      expect(hexToHsl("#000000")).toBeDefined();
      expect(hexToHsl("#f97316")).toBeDefined();
    });
  });

  describe("getLuminance", () => {
    it("should calculate luminance for colors", () => {
      const whiteLum = getLuminance("#ffffff");
      const blackLum = getLuminance("#000000");
      expect(whiteLum).toBeGreaterThan(blackLum);
    });

    it("should return values between 0 and 1", () => {
      const lum = getLuminance("#ff0000");
      expect(lum).toBeGreaterThanOrEqual(0);
      expect(lum).toBeLessThanOrEqual(1);
    });

    it("white should have high luminance", () => {
      expect(getLuminance("#ffffff")).toBeGreaterThan(0.9);
    });

    it("black should have low luminance", () => {
      expect(getLuminance("#000000")).toBeLessThan(0.1);
    });
  });

  describe("getContrastRatio", () => {
    it("should calculate contrast ratio", () => {
      const ratio = getContrastRatio("#ffffff", "#000000");
      expect(ratio).toBeCloseTo(21, 0);
    });

    it("should return ratio of at least 1", () => {
      const ratio = getContrastRatio("#ffffff", "#ffffff");
      expect(ratio).toBeGreaterThanOrEqual(1);
    });

    it("white on black should be maximum contrast", () => {
      const ratio = getContrastRatio("#ffffff", "#000000");
      expect(ratio).toBeGreaterThan(20);
    });

    it("same color should have ratio of 1", () => {
      const ratio = getContrastRatio("#ffffff", "#ffffff");
      expect(ratio).toBe(1);
    });
  });

  describe("meetsWCAGAA", () => {
    it("should validate WCAG AA compliance (4.5:1)", () => {
      expect(meetsWCAGAA("#ffffff", "#000000")).toBe(true);
      expect(meetsWCAGAA("#000000", "#ffffff")).toBe(true);
    });

    it("should reject low contrast", () => {
      expect(meetsWCAGAA("#f0f0f0", "#ffffff")).toBe(false);
    });

    it("should handle text on background colors", () => {
      // Black text on white background should pass
      expect(meetsWCAGAA("#1f2937", "#ffffff")).toBe(true);
    });
  });

  describe("meetsWCAGAAA", () => {
    it("should validate WCAG AAA compliance (7:1)", () => {
      expect(meetsWCAGAAA("#ffffff", "#000000")).toBe(true);
      expect(meetsWCAGAAA("#000000", "#ffffff")).toBe(true);
    });

    it("should be more strict than AA", () => {
      const ratio = getContrastRatio("#ffffff", "#808080");
      const passesAA = ratio >= 4.5;
      const passesAAA = ratio >= 7;
      if (passesAA && !passesAAA) {
        expect(meetsWCAGAA("#ffffff", "#808080")).toBe(true);
        expect(meetsWCAGAAA("#ffffff", "#808080")).toBe(false);
      }
    });
  });

  describe("Realistic theme colors", () => {
    const themeColors = {
      primary: "#c2410c",
      secondary: "#0369a1",
      accent: "#fbbf24",
      background: "#ffffff",
      text: "#1f2937",
      border: "#e5e7eb",
    };

    it("text on background should meet WCAG AA", () => {
      expect(
        meetsWCAGAA(themeColors.text, themeColors.background)
      ).toBe(true);
    });

    it("primary on background should meet WCAG AA", () => {
      expect(
        meetsWCAGAA(themeColors.primary, themeColors.background)
      ).toBe(true);
    });

    it("secondary on background should meet WCAG AA", () => {
      expect(
        meetsWCAGAA(themeColors.secondary, themeColors.background)
      ).toBe(true);
    });
  });
});
