'use client';

import React, { useState } from 'react';
import { Eye, EyeOff, Edit2, Trash2, GripVertical, ChevronDown } from 'lucide-react';
import { Button } from '@pos-saas/ui';
import type { ThemeSection } from '@/lib/hooks/use-theme';

export type { ThemeSection } from '@/lib/hooks/use-theme';

export interface ComponentSidebarProps {
  sections: ThemeSection[];
  isLoading: boolean;
  onSectionClick: (section: ThemeSection) => void;
  onVisibilityToggle: (sectionId: number, isVisible: boolean) => Promise<void>;
  onDeleteSection: (sectionId: number) => Promise<void>;
  onDragStart: (section: ThemeSection) => void;
  selectedSection: ThemeSection | null;
}

const SECTION_TYPE_NAMES: Record<string, string> = {
  hero: 'Hero Section',
  featured_items: 'Featured Items',
  why_choose_us: 'Why Choose Us',
  info_cards: 'Info Cards',
  cta: 'Call to Action',
  testimonials: 'Testimonials',
};

const SECTION_TYPE_COLORS: Record<string, string> = {
  hero: 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400',
  featured_items: 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400',
  why_choose_us: 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-400',
  info_cards: 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-400',
  cta: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400',
  testimonials: 'bg-pink-100 dark:bg-pink-900/30 text-pink-700 dark:text-pink-400',
};

interface SectionItemProps {
  section: ThemeSection;
  isSelected: boolean;
  onSectionClick: (section: ThemeSection) => void;
  onVisibilityToggle: (sectionId: number, isVisible: boolean) => Promise<void>;
  onDeleteSection: (sectionId: number) => Promise<void>;
  onDragStart: (section: ThemeSection) => void;
}

function SectionItem({
  section,
  isSelected,
  onSectionClick,
  onVisibilityToggle,
  onDeleteSection,
  onDragStart,
}: SectionItemProps) {
  const [isDeleting, setIsDeleting] = useState(false);
  const [isTogglingVisibility, setIsTogglingVisibility] = useState(false);

  const handleVisibilityToggle = async () => {
    try {
      setIsTogglingVisibility(true);
      await onVisibilityToggle(section.id, !section.is_visible);
    } finally {
      setIsTogglingVisibility(false);
    }
  };

  const handleDelete = async () => {
    if (confirm(`Are you sure you want to delete "${section.title || SECTION_TYPE_NAMES[section.section_type]}"?`)) {
      try {
        setIsDeleting(true);
        await onDeleteSection(section.id);
      } finally {
        setIsDeleting(false);
      }
    }
  };

  return (
    <div
      draggable
      onDragStart={() => onDragStart(section)}
      onClick={() => onSectionClick(section)}
      className={`p-3 rounded-lg border-2 cursor-move transition-all group ${
        isSelected
          ? 'border-primary bg-primary/5'
          : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
      } ${!section.is_visible ? 'opacity-50' : ''}`}
    >
      {/* Header Row */}
      <div className="flex items-center gap-2 mb-2">
        {/* Drag Handle */}
        <GripVertical className="h-4 w-4 text-gray-400 dark:text-gray-600 flex-shrink-0" />

        {/* Section Type Badge */}
        <span className={`text-xs font-medium px-2 py-1 rounded ${SECTION_TYPE_COLORS[section.section_type]}`}>
          {SECTION_TYPE_NAMES[section.section_type] || section.section_type}
        </span>

        {/* Visibility Toggle */}
        <button
          onClick={(e) => {
            e.stopPropagation();
            handleVisibilityToggle();
          }}
          disabled={isTogglingVisibility}
          className="ml-auto p-1 hover:bg-gray-200 dark:hover:bg-gray-700 rounded transition-colors"
          title={section.is_visible ? 'Hide section' : 'Show section'}
        >
          {section.is_visible ? (
            <Eye className="h-4 w-4 text-gray-600 dark:text-gray-400" />
          ) : (
            <EyeOff className="h-4 w-4 text-gray-400 dark:text-gray-600" />
          )}
        </button>
      </div>

      {/* Section Title */}
      <p className="text-sm font-medium text-gray-900 dark:text-white truncate mb-1">
        {section.title || `Section ${section.order}`}
      </p>

      {/* Section Subtitle */}
      {section.subtitle && (
        <p className="text-xs text-gray-600 dark:text-gray-400 truncate mb-2">
          {section.subtitle}
        </p>
      )}

      {/* Order Badge */}
      <div className="flex items-center justify-between mb-3">
        <span className="text-xs bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 px-2 py-1 rounded">
          Order: {section.order}
        </span>
      </div>

      {/* Action Buttons */}
      <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <button
          onClick={(e) => {
            e.stopPropagation();
            onSectionClick(section);
          }}
          className="flex-1 p-1 text-xs bg-blue-500 hover:bg-blue-600 text-white rounded flex items-center justify-center gap-1 transition-colors"
          title="Edit section"
        >
          <Edit2 className="h-3 w-3" />
          Edit
        </button>
        <button
          onClick={(e) => {
            e.stopPropagation();
            handleDelete();
          }}
          disabled={isDeleting}
          className="flex-1 p-1 text-xs bg-red-500 hover:bg-red-600 text-white rounded flex items-center justify-center gap-1 transition-colors disabled:opacity-50"
          title="Delete section"
        >
          <Trash2 className="h-3 w-3" />
          {isDeleting ? 'Deleting...' : 'Delete'}
        </button>
      </div>
    </div>
  );
}

export function ComponentSidebar({
  sections,
  isLoading,
  onSectionClick,
  onVisibilityToggle,
  onDeleteSection,
  onDragStart,
  selectedSection,
}: ComponentSidebarProps) {
  const [expandedCategories, setExpandedCategories] = useState<Record<string, boolean>>({
    hero: true,
    featured_items: true,
    why_choose_us: true,
    info_cards: true,
    cta: true,
    testimonials: true,
  });

  const toggleCategory = (type: string) => {
    setExpandedCategories(prev => ({
      ...prev,
      [type]: !prev[type],
    }));
  };

  // Group sections by type
  const sectionsByType: Record<string, ThemeSection[]> = {};
  sections.forEach(section => {
    if (!sectionsByType[section.section_type]) {
      sectionsByType[section.section_type] = [];
    }
    sectionsByType[section.section_type].push(section);
  });

  // Sort sections within each type by order
  Object.keys(sectionsByType).forEach(type => {
    sectionsByType[type].sort((a, b) => a.order - b.order);
  });

  if (isLoading) {
    return (
      <div className="p-6 text-center">
        <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full mx-auto mb-2" />
        <p className="text-sm text-gray-600 dark:text-gray-400">Loading sections...</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          Home Page Sections ({sections.length})
        </h3>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Drag to reorder • Click to edit • Toggle visibility
        </p>
      </div>

      {/* Sections by Type */}
      <div className="space-y-3">
        {Object.entries(sectionsByType).map(([type, typeSections]) => (
          <div key={type}>
            {/* Category Header */}
            <button
              onClick={() => toggleCategory(type)}
              className="w-full flex items-center gap-2 p-3 rounded-lg bg-gray-50 dark:bg-gray-700/50 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors mb-2"
            >
              <ChevronDown
                className={`h-4 w-4 text-gray-600 dark:text-gray-400 transition-transform ${
                  expandedCategories[type] ? '' : '-rotate-90'
                }`}
              />
              <span className="font-medium text-gray-700 dark:text-gray-300">
                {SECTION_TYPE_NAMES[type] || type}
              </span>
              <span className="ml-auto text-xs bg-gray-200 dark:bg-gray-600 text-gray-700 dark:text-gray-300 px-2 py-1 rounded">
                {typeSections.length}
              </span>
            </button>

            {/* Section Items */}
            {expandedCategories[type] && (
              <div className="space-y-2 pl-2">
                {typeSections.map(section => (
                  <SectionItem
                    key={section.id}
                    section={section}
                    isSelected={selectedSection?.id === section.id}
                    onSectionClick={onSectionClick}
                    onVisibilityToggle={onVisibilityToggle}
                    onDeleteSection={onDeleteSection}
                    onDragStart={onDragStart}
                  />
                ))}
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Empty State */}
      {sections.length === 0 && (
        <div className="text-center py-8">
          <p className="text-gray-600 dark:text-gray-400 text-sm">
            No sections yet. Add one to customize your home page!
          </p>
        </div>
      )}

      {/* Info Box */}
      <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3">
        <p className="text-xs text-blue-800 dark:text-blue-200">
          💡 <strong>Tip:</strong> Drag sections to reorder them on your home page. Hide sections you don't want to display.
        </p>
      </div>
    </div>
  );
}
