'use client';

import React, { useState } from 'react';
import { Smartphone, Tablet, Monitor, ZoomIn, ZoomOut, Copy, Download } from 'lucide-react';
import { ManagerPreview } from './manager-preview';
import type { ThemeSection } from './component-sidebar';

interface EnhancedPreviewProps {
  sections: ThemeSection[];
  primaryColor: string;
  secondaryColor: string;
  accentColor: string;
  fontFamily: string;
  logoUrl: string;
}

type DeviceMode = 'mobile' | 'tablet' | 'desktop';

const DEVICE_SPECS = {
  mobile: {
    width: 375,
    height: 812,
    name: 'iPhone 15',
    icon: Smartphone,
  },
  tablet: {
    width: 768,
    height: 1024,
    name: 'iPad Pro',
    icon: Tablet,
  },
  desktop: {
    width: 1280,
    height: 800,
    name: 'Desktop (1280px)',
    icon: Monitor,
  },
};

export function EnhancedPreview({
  sections,
  primaryColor,
  secondaryColor,
  accentColor,
  fontFamily,
  logoUrl,
}: EnhancedPreviewProps) {
  const [deviceMode, setDeviceMode] = useState<DeviceMode>('desktop');
  const [zoomLevel, setZoomLevel] = useState(100);
  const [showComparisonMode, setShowComparisonMode] = useState(false);

  const currentDevice = DEVICE_SPECS[deviceMode];
  const previewWidth = (currentDevice.width * zoomLevel) / 100;
  const previewHeight = (currentDevice.height * zoomLevel) / 100;
  const zoomScale = zoomLevel / 100;
  const inverseScale = 100 / zoomScale;

  const handleZoomIn = () => {
    setZoomLevel(prev => Math.min(prev + 10, 200));
  };

  const handleZoomOut = () => {
    setZoomLevel(prev => Math.max(prev - 10, 50));
  };

  const handleResetZoom = () => {
    setZoomLevel(100);
  };

  const handleDownloadPreview = () => {
    // Create canvas-based screenshot
    const previewElement = document.getElementById('preview-content');
    if (previewElement) {
      // In a real app, use html2canvas library
      console.log('Downloading preview...');
      // html2canvas(previewElement).then(canvas => {
      //   const link = document.createElement('a');
      //   link.href = canvas.toDataURL('image/png');
      //   link.download = `preview-${deviceMode}-${Date.now()}.png`;
      //   link.click();
      // });
    }
  };

  const handleCopyPreviewLink = () => {
    const url = window.location.href;
    navigator.clipboard.writeText(url);
  };

  return (
    <div className="enhanced-preview space-y-4">
      {/* Toolbar */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
        <div className="flex flex-wrap items-center justify-between gap-4">
          {/* Device Selector */}
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Device:</span>
            <div className="flex gap-2">
              {(Object.entries(DEVICE_SPECS) as Array<[DeviceMode, typeof DEVICE_SPECS.mobile]>).map(
                ([mode, spec]) => {
                  const Icon = spec.icon;
                  return (
                    <button
                      key={mode}
                      onClick={() => {
                        setDeviceMode(mode);
                        setZoomLevel(100); // Reset zoom on device change
                      }}
                      className={`flex items-center gap-1.5 px-3 py-2 rounded-lg transition-colors ${
                        deviceMode === mode
                          ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 border border-blue-300 dark:border-blue-700'
                          : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                      }`}
                      title={`Preview on ${spec.name}`}
                    >
                      <Icon className="h-4 w-4" />
                      <span className="text-sm hidden sm:inline">{spec.name.split(' ')[0]}</span>
                    </button>
                  );
                }
              )}
            </div>
          </div>

          {/* Zoom Controls */}
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Zoom:</span>
            <button
              onClick={handleZoomOut}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
              title="Zoom out"
              disabled={zoomLevel <= 50}
            >
              <ZoomOut className="h-4 w-4 text-gray-600 dark:text-gray-400" />
            </button>
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300 min-w-[3rem] text-center">
              {zoomLevel}%
            </span>
            <button
              onClick={handleZoomIn}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
              title="Zoom in"
              disabled={zoomLevel >= 200}
            >
              <ZoomIn className="h-4 w-4 text-gray-600 dark:text-gray-400" />
            </button>
            <button
              onClick={handleResetZoom}
              className="px-3 py-1.5 text-sm bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded transition-colors text-gray-700 dark:text-gray-300"
            >
              Reset
            </button>
          </div>

          {/* Action Buttons */}
          <div className="flex items-center gap-2">
            <button
              onClick={handleCopyPreviewLink}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
              title="Copy preview link"
            >
              <Copy className="h-4 w-4 text-gray-600 dark:text-gray-400" />
            </button>
            <button
              onClick={handleDownloadPreview}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
              title="Download preview"
            >
              <Download className="h-4 w-4 text-gray-600 dark:text-gray-400" />
            </button>
            <button
              onClick={() => setShowComparisonMode(!showComparisonMode)}
              className={`px-3 py-1.5 text-sm rounded transition-colors ${
                showComparisonMode
                  ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              }`}
              title="Side-by-side comparison"
            >
              Compare
            </button>
          </div>
        </div>

        {/* Device Info */}
        <div className="mt-2 text-xs text-gray-500 dark:text-gray-400">
          {currentDevice.name} • {previewWidth.toFixed(0)}px × {previewHeight.toFixed(0)}px
        </div>
      </div>

      {/* Preview Area */}
      <div className="bg-gray-50 dark:bg-gray-900/50 rounded-lg border border-gray-200 dark:border-gray-700 p-6">
        {showComparisonMode ? (
          // Comparison Mode: Side-by-side views
          <div className="grid grid-cols-2 gap-4">
            {/* Original */}
            <div>
              <h3 className="text-xs font-semibold text-gray-600 dark:text-gray-400 mb-2">ORIGINAL</h3>
              <div
                id="preview-content"
                className="comparison-preview-frame bg-white rounded border border-gray-200 dark:border-gray-700 overflow-auto"
              >
                <ManagerPreview
                  sections={sections}
                  primaryColor={primaryColor}
                  secondaryColor={secondaryColor}
                  accentColor={accentColor}
                  fontFamily={fontFamily}
                  logoUrl={logoUrl}
                />
              </div>
            </div>

            {/* With Changes Indicator */}
            <div>
              <h3 className="text-xs font-semibold text-gray-600 dark:text-gray-400 mb-2">CURRENT</h3>
              <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded p-4 text-center text-sm text-blue-700 dark:text-blue-400">
                <p>No changes to compare yet</p>
                <p className="text-xs mt-2 text-blue-600 dark:text-blue-500">
                  Make changes to see side-by-side comparison
                </p>
              </div>
            </div>
          </div>
        ) : (
          // Single Preview Mode
          <div className="flex items-start justify-center overflow-x-auto">
            <div
              id="preview-content"
              className={`single-preview-frame bg-white dark:bg-gray-800 rounded border border-gray-300 dark:border-gray-600 shadow-lg overflow-hidden flex-shrink-0 ${
                deviceMode !== 'desktop' ? 'device-framed border-8 border-gray-900 dark:border-gray-950' : ''
              }`}
            >
              {/* Screen content with zoom applied */}
              <div className="single-preview-content">
                <ManagerPreview
                  sections={sections}
                  primaryColor={primaryColor}
                  secondaryColor={secondaryColor}
                  accentColor={accentColor}
                  fontFamily={fontFamily}
                  logoUrl={logoUrl}
                />
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Preview Info */}
      <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
        <div className="flex gap-2 text-sm">
          <div>
            <p className="font-medium text-blue-900 dark:text-blue-200">Preview Tips:</p>
            <ul className="text-blue-800 dark:text-blue-300 text-xs mt-1 space-y-1">
              <li>• Use device buttons to test on different screen sizes</li>
              <li>• Use zoom to see exact pixel dimensions</li>
              <li>• Click &quot;Compare&quot; to see changes side-by-side</li>
              <li>• Download preview to share with team</li>
            </ul>
          </div>
        </div>
      </div>

      <style jsx>{`
        .comparison-preview-frame {
          width: ${previewWidth}px;
          height: ${previewHeight}px;
          transform: scale(${zoomScale});
          transform-origin: top left;
        }

        .single-preview-frame {
          width: ${previewWidth}px;
          height: ${previewHeight}px;
        }

        .single-preview-content {
          overflow: auto;
          transform: scale(${zoomScale});
          transform-origin: top left;
          width: ${inverseScale}%;
          height: ${inverseScale}%;
        }

        .device-framed {
          border-radius: 12px;
          box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
        }
      `}</style>
    </div>
  );
}