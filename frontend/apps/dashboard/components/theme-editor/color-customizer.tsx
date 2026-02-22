'use client';

import React, { useState } from 'react';
import { Copy, Check } from 'lucide-react';

interface ColorCustomizerProps {
  primaryColor: string;
  secondaryColor: string;
  accentColor: string;
  onPrimaryChange: (color: string) => void;
  onSecondaryChange: (color: string) => void;
  onAccentChange: (color: string) => void;
}

interface ColorPreset {
  name: string;
  primary: string;
  secondary: string;
  accent: string;
}

const COLOR_PRESETS: ColorPreset[] = [
  {
    name: 'Modern Blue',
    primary: '#3b82f6',
    secondary: '#10b981',
    accent: '#f59e0b',
  },
  {
    name: 'Warm Sunset',
    primary: '#dc2626',
    secondary: '#ea580c',
    accent: '#fbbf24',
  },
  {
    name: 'Ocean Fresh',
    primary: '#0369a1',
    secondary: '#0891b2',
    accent: '#06b6d4',
  },
  {
    name: 'Forest Green',
    primary: '#15803d',
    secondary: '#22c55e',
    accent: '#84cc16',
  },
  {
    name: 'Purple Dream',
    primary: '#7c3aed',
    secondary: '#a855f7',
    accent: '#d946ef',
  },
  {
    name: 'Dark Elegance',
    primary: '#1f2937',
    secondary: '#374151',
    accent: '#9ca3af',
  },
];

interface ColorInputProps {
  label: string;
  color: string;
  onChange: (color: string) => void;
}

function ColorInput({ label, color, onChange }: ColorInputProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(color);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const isValidHex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/.test(color);

  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
        {label}
      </label>
      <div className="flex items-center gap-2">
        {/* Color Input */}
        <div className="flex-1 flex items-center gap-2">
          <input
            type="color"
            value={color}
            onChange={(e) => onChange(e.target.value)}
            className="h-10 w-16 rounded border border-gray-300 dark:border-gray-600 cursor-pointer"
          />
          <input
            type="text"
            value={color}
            onChange={(e) => onChange(e.target.value)}
            placeholder="#000000"
            className={`flex-1 px-3 py-2 border rounded font-mono text-sm ${
              isValidHex
                ? 'border-gray-300 dark:border-gray-600'
                : 'border-red-500 dark:border-red-400'
            } bg-white dark:bg-gray-700 text-gray-900 dark:text-white`}
          />
        </div>

        {/* Copy Button */}
        <button
          onClick={handleCopy}
          className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
          title="Copy color code"
        >
          {copied ? (
            <Check className="h-4 w-4 text-green-600" />
          ) : (
            <Copy className="h-4 w-4 text-gray-600 dark:text-gray-400" />
          )}
        </button>
      </div>

      {!isValidHex && (
        <p className="text-xs text-red-600 dark:text-red-400">
          Invalid hex color. Use format: #RRGGBB or #RGB
        </p>
      )}
    </div>
  );
}

export function ColorCustomizer({
  primaryColor,
  secondaryColor,
  accentColor,
  onPrimaryChange,
  onSecondaryChange,
  onAccentChange,
}: ColorCustomizerProps) {
  const presetColorRules = COLOR_PRESETS.map(
    (preset, index) => `
      .preset-${index}-primary { background-color: ${preset.primary}; }
      .preset-${index}-secondary { background-color: ${preset.secondary}; }
      .preset-${index}-accent { background-color: ${preset.accent}; }
    `
  ).join('\n');

  return (
    <div className="color-customizer space-y-6">
      {/* Color Inputs */}
      <div className="space-y-4">
        <ColorInput
          label="Primary Color"
          color={primaryColor}
          onChange={onPrimaryChange}
        />
        <ColorInput
          label="Secondary Color"
          color={secondaryColor}
          onChange={onSecondaryChange}
        />
        <ColorInput
          label="Accent Color"
          color={accentColor}
          onChange={onAccentChange}
        />
      </div>

      {/* Color Preview */}
      <div className="space-y-2">
        <p className="text-sm font-medium text-gray-700 dark:text-gray-300">Preview</p>
        <div className="grid grid-cols-3 gap-2">
          <div
            className="preview-color preview-primary h-16 rounded border border-gray-200 dark:border-gray-700"
            title="Primary"
          />
          <div
            className="preview-color preview-secondary h-16 rounded border border-gray-200 dark:border-gray-700"
            title="Secondary"
          />
          <div
            className="preview-color preview-accent h-16 rounded border border-gray-200 dark:border-gray-700"
            title="Accent"
          />
        </div>
      </div>

      {/* Color Presets */}
      <div className="space-y-2">
        <p className="text-sm font-medium text-gray-700 dark:text-gray-300">Presets</p>
        <div className="space-y-2">
          {COLOR_PRESETS.map((preset, index) => (
            <button
              key={preset.name}
              onClick={() => {
                onPrimaryChange(preset.primary);
                onSecondaryChange(preset.secondary);
                onAccentChange(preset.accent);
              }}
              className="w-full flex items-center gap-3 p-3 rounded border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors text-left"
            >
              <div className="flex gap-1">
                <div className={`w-6 h-6 rounded border border-gray-200 preset-${index}-primary`} />
                <div className={`w-6 h-6 rounded border border-gray-200 preset-${index}-secondary`} />
                <div className={`w-6 h-6 rounded border border-gray-200 preset-${index}-accent`} />
              </div>
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                {preset.name}
              </span>
            </button>
          ))}
        </div>
      </div>
      <style jsx>{`
        .preview-primary {
          background-color: ${primaryColor};
        }

        .preview-secondary {
          background-color: ${secondaryColor};
        }

        .preview-accent {
          background-color: ${accentColor};
        }

        ${presetColorRules}
      `}</style>
    </div>
  );
}
