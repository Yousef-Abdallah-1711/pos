'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Loader, AlertCircle, Save } from 'lucide-react';
import { Button } from '@pos-saas/ui';
import { ColorCustomizer } from './color-customizer';
import { LogoUploader } from './logo-uploader';
import { FontSelector } from './font-selector';
import { ThemePreview } from './theme-preview';
import { useTheme } from '@/lib/hooks/use-theme';

interface ThemeEditorProps {
  restaurantId: number;
}

export function ThemeEditor({ restaurantId }: ThemeEditorProps) {
  const router = useRouter();
  const { theme, isLoading, updateTheme, isSaving } = useTheme(restaurantId);

  const [formData, setFormData] = useState({
    primaryColor: '#3b82f6',
    secondaryColor: '#10b981',
    accentColor: '#f59e0b',
    logoUrl: '',
    fontFamily: 'Inter',
  });

  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // Load theme data when it arrives
  useEffect(() => {
    if (theme) {
      setFormData({
        primaryColor: theme.primary_color || '#3b82f6',
        secondaryColor: theme.secondary_color || '#10b981',
        accentColor: theme.accent_color || '#f59e0b',
        logoUrl: theme.logo_url || '',
        fontFamily: theme.font_family || 'Inter',
      });
    }
  }, [theme]);

  const handleColorChange = (key: 'primaryColor' | 'secondaryColor' | 'accentColor', value: string) => {
    setFormData(prev => ({ ...prev, [key]: value }));
    setError(null);
  };

  const handleLogoChange = (url: string) => {
    setFormData(prev => ({ ...prev, logoUrl: url }));
  };

  const handleFontChange = (font: string) => {
    setFormData(prev => ({ ...prev, fontFamily: font }));
  };

  const handleSave = async () => {
    try {
      setError(null);
      await updateTheme({
        primary_color: formData.primaryColor,
        secondary_color: formData.secondaryColor,
        accent_color: formData.accentColor,
        logo_url: formData.logoUrl,
        font_family: formData.fontFamily,
      });
      setSuccess('Theme saved successfully!');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save theme');
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader className="h-8 w-8 animate-spin text-primary mr-2" />
        <p className="text-gray-600 dark:text-gray-400">Loading theme...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Theme Editor</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">Customize your restaurant&apos;s appearance</p>
        </div>
      </div>

      {/* Error Alert */}
      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex items-start gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <p className="text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      {/* Success Alert */}
      {success && (
        <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4 flex items-start gap-3">
          <div className="h-5 w-5 bg-green-600 dark:bg-green-400 rounded-full flex-shrink-0 mt-0.5" />
          <p className="text-green-800 dark:text-green-200">{success}</p>
        </div>
      )}

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Panel - Controls */}
        <div className="lg:col-span-1">
          <div className="space-y-6">
            {/* Color Customizer */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">Colors</h2>
              <ColorCustomizer
                primaryColor={formData.primaryColor}
                secondaryColor={formData.secondaryColor}
                accentColor={formData.accentColor}
                onPrimaryChange={(color) => handleColorChange('primaryColor', color)}
                onSecondaryChange={(color) => handleColorChange('secondaryColor', color)}
                onAccentChange={(color) => handleColorChange('accentColor', color)}
              />
            </div>

            {/* Logo Uploader */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">Logo</h2>
              <LogoUploader
                logoUrl={formData.logoUrl}
                onLogoChange={handleLogoChange}
              />
            </div>

            {/* Font Selector */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">Typography</h2>
              <FontSelector
                selectedFont={formData.fontFamily}
                onFontChange={handleFontChange}
              />
            </div>

            {/* Save Button */}
            <Button
              onClick={handleSave}
              disabled={isSaving}
              className="w-full"
              size="lg"
            >
              {isSaving ? (
                <>
                  <Loader className="h-4 w-4 mr-2 animate-spin" />
                  Saving...
                </>
              ) : (
                <>
                  <Save className="h-4 w-4 mr-2" />
                  Save Theme
                </>
              )}
            </Button>
          </div>
        </div>

        {/* Right Panel - Preview */}
        <div className="lg:col-span-2">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden sticky top-6">
            <div className="p-6 border-b border-gray-200 dark:border-gray-700">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Live Preview</h2>
              <p className="text-sm text-gray-600 dark:text-gray-400">See changes in real-time</p>
            </div>
            <ThemePreview
              primaryColor={formData.primaryColor}
              secondaryColor={formData.secondaryColor}
              accentColor={formData.accentColor}
              logoUrl={formData.logoUrl}
              fontFamily={formData.fontFamily}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
