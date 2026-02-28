// Auth Store
export { useAuthStore } from './authStore'
export type { AuthState } from './authStore'

// POS Store
export { usePOSStore } from './posStore'
export type {
  PaymentMethod as POSPaymentMethod,
  POSAddItemProduct,
  POSCartItem,
  HeldOrder as POSHeldOrder,
  POSState,
} from './posStore'

// Preferences Store
export * from './preferencesStore'

// Theme Builder Store
export { useThemeBuilderStore } from './themeBuilderStore'
export type { ThemeBuilderState } from './themeBuilderStore'

// Theme Builder Hooks
export {
  useCurrentTheme,
  useThemeColors,
  useUpdateColors,
  useColorsWithUpdate,
  useColorValue,
  useThemeTypography,
  useUpdateTypography,
  useTypographyWithUpdate,
  useThemeHeader,
  useUpdateHeader,
  useHeaderWithUpdate,
  useThemeFooter,
  useUpdateFooter,
  useFooterWithUpdate,
  useThemeIdentity,
  useUpdateIdentity,
  useIdentityWithUpdate,
  useThemeComponents,
  useUpdateComponents,
  useGetComponentById,
  useCustomCss,
  useUpdateCustomCss,
  useIsDirty,
  useDirtyActions,
  useThemeLoading,
  useThemeSaving,
  useThemeError,
  useThemeSuccess,
  useThemeStatus,
  useThemeStatusActions,
  useThemeManagement,
  useThemeValidation,
  useThemeFormData,
  useHasThemeErrors,
  useThemeBuilder,
  useBulkThemeUpdate,
} from './themeBuilderHooks'

// Theme Builder Utilities
export {
  isValidHexColor,
  validateColorWithMessage,
  hexToRgb,
  rgbToHex,
  hexToHsl,
  hslToHex,
  getRelativeLuminance,
  getContrastRatio,
  meetsWCAGAA,
  meetsWCAGAAA,
  validateColorPalette,
  isValidFontSize,
  isValidLineHeight,
  isValidBorderRadius,
  validateTypography,
  validateHeader,
  validateFooter,
  generateSlug,
  sanitizeCss,
  mergeColors,
  cloneThemeData,
  getColorChanges,
  generatePaletteFromPrimary,
  validateCompleteTheme,
} from './themeBuilderUtils'
