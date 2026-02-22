'use client';

import React from 'react';

interface FontSelectorProps {
  selectedFont: string;
  onFontChange: (font: string) => void;
}

const FONT_FAMILIES = [
  { id: 'Inter', name: 'Inter', preview: 'sans-serif' },
  { id: 'Poppins', name: 'Poppins', preview: 'sans-serif' },
  { id: 'Roboto', name: 'Roboto', preview: 'sans-serif' },
  { id: 'Open Sans', name: 'Open Sans', preview: 'sans-serif' },
  { id: 'Playfair Display', name: 'Playfair Display', preview: 'serif' },
  { id: 'Lora', name: 'Lora', preview: 'serif' },
  { id: 'Merriweather', name: 'Merriweather', preview: 'serif' },
  { id: 'Montserrat', name: 'Montserrat', preview: 'sans-serif' },
  { id: 'Source Code Pro', name: 'Source Code Pro', preview: 'monospace' },
  { id: 'IBM Plex Mono', name: 'IBM Plex Mono', preview: 'monospace' },
];

export function FontSelector({ selectedFont, onFontChange }: FontSelectorProps) {
  const fontFamilyRules = FONT_FAMILIES.map(
    (font, index) => `.font-item-${index} { font-family: '${font.id}', sans-serif; }`
  ).join('\n');

  const selectedFontEscaped = selectedFont.replace(/'/g, "\\'");

  return (
    <div className="font-selector space-y-3">
      <div className="grid grid-cols-1 gap-2">
        {FONT_FAMILIES.map((font, index) => (
          <button
            key={font.id}
            onClick={() => onFontChange(font.id)}
            className={`p-3 rounded border-2 transition-all text-left ${
              selectedFont === font.id
                ? 'border-primary bg-primary/5'
                : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
            }`}
          >
            <div className={`font-item-${index} text-sm font-medium text-gray-900 dark:text-white`}>
              {font.name}
            </div>
            <div className="text-xs text-gray-600 dark:text-gray-400 mt-1">
              {font.preview === 'sans-serif' && 'Clean & modern'}
              {font.preview === 'serif' && 'Elegant & traditional'}
              {font.preview === 'monospace' && 'Technical & minimal'}
            </div>
          </button>
        ))}
      </div>

      {/* Font Preview */}
      <div className="mt-4 space-y-2">
        <p className="text-sm font-medium text-gray-700 dark:text-gray-300">Preview</p>
        <div className="selected-font-preview p-4 bg-gray-50 dark:bg-gray-700 rounded border border-gray-200 dark:border-gray-600">
          <p className="text-lg font-medium text-gray-900 dark:text-white">
            The quick brown fox jumps
          </p>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            1234567890 !@#$%^&*()
          </p>
        </div>
      </div>

      {/* Font Categories */}
      <div className="mt-4 space-y-2">
        <p className="text-sm font-medium text-gray-700 dark:text-gray-300">Categories</p>
        <div className="flex flex-wrap gap-2">
          <span className="px-2 py-1 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs rounded">
            Sans-serif
          </span>
          <span className="px-2 py-1 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs rounded">
            Serif
          </span>
          <span className="px-2 py-1 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs rounded">
            Monospace
          </span>
        </div>
      </div>

      <style jsx>{`
        .selected-font-preview {
          font-family: '${selectedFontEscaped}', sans-serif;
        }

        ${fontFamilyRules}
      `}</style>
    </div>
  );
}