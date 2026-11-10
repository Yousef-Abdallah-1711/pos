'use client';

import React, { useState, useCallback } from 'react';
import { Loader, AlertCircle } from 'lucide-react';
import { ComponentSidebar, ThemeSection } from './component-sidebar';
import { SectionEditor } from './section-editor';
import { EnhancedPreview } from './enhanced-preview';
import { useTheme, useSectionContent, useSectionVisibility, useReorderSections } from '@/lib/hooks/use-theme';

interface ComponentManagerProps {
  restaurantId: number;
  primaryColor?: string;
  secondaryColor?: string;
  accentColor?: string;
  fontFamily?: string;
  logoUrl?: string;
}

export function ComponentManager({
  restaurantId,
  primaryColor = '#3b82f6',
  secondaryColor = '#10b981',
  accentColor = '#f59e0b',
  fontFamily = 'Inter',
  logoUrl = '',
}: ComponentManagerProps) {
  const [selectedSection, setSelectedSection] = useState<ThemeSection | null>(null);
  const [draggedSection, setDraggedSection] = useState<ThemeSection | null>(null);
  const [dropZoneIndex, setDropZoneIndex] = useState<number | null>(null);

  // Hooks
  const { theme, isLoading, error } = useTheme(restaurantId);
  const sectionVisibility = useSectionVisibility(restaurantId);
  const sectionContent = useSectionContent(restaurantId);
  const reorderSections = useReorderSections(restaurantId);

  const sections = theme?.sections || [];

  // Handle visibility toggle
  const handleVisibilityToggle = useCallback(
    async (sectionId: number, isVisible: boolean) => {
      try {
        await sectionVisibility.mutateAsync({ sectionId, isVisible });
      } catch (err) {
        console.error('Failed to toggle visibility:', err);
      }
    },
    [sectionVisibility]
  );

  // Handle section deletion (hide it)
  const handleDeleteSection = useCallback(
    async (sectionId: number) => {
      try {
        await sectionVisibility.mutateAsync({ sectionId, isVisible: false });
        if (selectedSection?.id === sectionId) {
          setSelectedSection(null);
        }
      } catch (err) {
        console.error('Failed to delete section:', err);
      }
    },
    [sectionVisibility, selectedSection]
  );

  // Handle section save
  const handleSectionSave = useCallback(
    async (updatedSection: ThemeSection) => {
      try {
        await sectionContent.mutateAsync({
          sectionId: updatedSection.id,
          title: updatedSection.title,
          subtitle: updatedSection.subtitle,
          description: updatedSection.description,
          content: updatedSection.content,
        });
        setSelectedSection(null);
      } catch (err) {
        console.error('Failed to save section:', err);
        throw err;
      }
    },
    [sectionContent]
  );

  // Handle drag start
  const handleDragStart = (section: ThemeSection) => {
    setDraggedSection(section);
  };

  // Handle drag over
  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
  };

  // Handle drop
  const handleDrop = useCallback(
    async (index: number) => {
      if (!draggedSection || !sections) return;

      try {
        const newSections = [...sections];
        const draggedIndex = newSections.findIndex(s => s.id === draggedSection.id);

        if (draggedIndex !== index) {
          // Reorder sections
          const [removed] = newSections.splice(draggedIndex, 1);
          newSections.splice(index, 0, removed);

          // Create order map
          const orderMap: Record<number, number> = {};
          newSections.forEach((section, idx) => {
            orderMap[section.id] = idx + 1;
          });

          // Call API
          await reorderSections.mutateAsync(orderMap);
        }
      } catch (err) {
        console.error('Failed to reorder sections:', err);
      } finally {
        setDraggedSection(null);
        setDropZoneIndex(null);
      }
    },
    [draggedSection, sections, reorderSections]
  );

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader className="h-8 w-8 animate-spin text-primary mr-2" />
        <p className="text-gray-600 dark:text-gray-400">Loading component manager...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Component Manager</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Customize your home page layout and content
        </p>
      </div>

      {/* Error Alert */}
      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex items-start gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <p className="text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      {/* Main Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left: Component Sidebar */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 lg:col-span-1 h-fit sticky top-6">
          <ComponentSidebar
            sections={sections}
            isLoading={isLoading}
            onSectionClick={setSelectedSection}
            onVisibilityToggle={handleVisibilityToggle}
            onDeleteSection={handleDeleteSection}
            onDragStart={handleDragStart}
            selectedSection={selectedSection}
          />
        </div>

        {/* Right: Enhanced Preview */}
        <div className="lg:col-span-2 space-y-4">
          <div>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">Live Preview</h2>
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
              Preview your home page on different devices with zoom and comparison modes
            </p>
          </div>
          <EnhancedPreview
            sections={sections}
            primaryColor={primaryColor}
            secondaryColor={secondaryColor}
            accentColor={accentColor}
            fontFamily={fontFamily}
            logoUrl={logoUrl}
          />
        </div>
      </div>

      {/* Section Editor Modal */}
      {selectedSection && (
        <SectionEditor
          section={selectedSection}
          isSaving={sectionContent.isPending}
          onSave={handleSectionSave}
          onClose={() => setSelectedSection(null)}
        />
      )}
    </div>
  );
}
