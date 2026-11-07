'use client'

import { useState } from 'react'
import { Plus, Trash2, GripVertical, Eye, EyeOff, Zap, Box, MessageSquare, FileText, Star, Layout } from 'lucide-react'
import { useThemeStore } from '@/hooks/useThemeStore'
import { ComponentType, ComponentConfig } from '@/types/theme'

// Component type to icon mapping
const COMPONENT_ICONS: Record<ComponentType, React.ReactNode> = {
  hero: <Zap className="h-4 w-4" />,
  products: <Box className="h-4 w-4" />,
  why_us: <Star className="h-4 w-4" />,
  contact: <MessageSquare className="h-4 w-4" />,
  testimonials: <FileText className="h-4 w-4" />,
  cta: <Layout className="h-4 w-4" />,
  custom: <Box className="h-4 w-4" />,
}

const COMPONENT_TYPES: { type: ComponentType; label: string; description: string }[] = [
  { type: 'hero', label: 'Hero Section', description: 'Large banner with headline' },
  { type: 'products', label: 'Featured Products', description: 'Product showcase' },
  { type: 'why_us', label: 'Why Choose Us', description: 'Key features section' },
  { type: 'contact', label: 'Contact Section', description: 'Contact form and info' },
  { type: 'testimonials', label: 'Testimonials', description: 'Customer reviews' },
  { type: 'cta', label: 'Call to Action', description: 'Action button section' },
  { type: 'custom', label: 'Custom HTML', description: 'Custom content block' },
]

export function ComponentBuilder() {
  const theme = useThemeStore((state) => state.currentTheme)
  const addComponent = useThemeStore((state) => state.addComponent)
  const deleteComponent = useThemeStore((state) => state.deleteComponent)
  const selectComponent = useThemeStore((state) => state.selectComponent)
  const selectedComponentId = useThemeStore((state) => state.selectedComponentId)
  const reorderComponents = useThemeStore((state) => state.reorderComponents)

  const [showAddMenu, setShowAddMenu] = useState(false)
  const [draggedItem, setDraggedItem] = useState<string | null>(null)

  // Toggle component enabled/disabled state
  const handleToggleComponent = (id: string, e: React.MouseEvent) => {
    e.stopPropagation()
    if (theme) {
      const component = theme.components.find(c => c.id === id)
      if (component) {
        const updatedComponent = { ...component, enabled: !component.enabled }
        const updatedComponents = theme.components.map(c => c.id === id ? updatedComponent : c)
        reorderComponents(updatedComponents)
      }
    }
  }

  if (!theme) {
    return (
      <div className="flex items-center justify-center h-full text-gray-500">
        Select a theme to add components
      </div>
    )
  }

  const handleAddComponent = (type: ComponentType) => {
    const newComponent: ComponentConfig = {
      id: `component-${Date.now()}`,
      type,
      title: `${COMPONENT_TYPES.find((c) => c.type === type)?.label || type}`,
      displayOrder: theme.components.length, // was `order`, a field ComponentConfig does not have
      enabled: true,
      settings: {},
    }
    addComponent(newComponent)
    setShowAddMenu(false)
  }

  const handleDragStart = (id: string) => {
    setDraggedItem(id)
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
  }

  const handleDrop = (targetId: string) => {
    if (!draggedItem || draggedItem === targetId) return

    const draggedIdx = theme.components.findIndex((c) => c.id === draggedItem)
    const targetIdx = theme.components.findIndex((c) => c.id === targetId)

    const newComponents = [...theme.components]
    const [draggedComp] = newComponents.splice(draggedIdx, 1)
    newComponents.splice(targetIdx, 0, draggedComp)

    reorderComponents(newComponents)
    setDraggedItem(null)
  }

  return (
    <div className="h-full flex flex-col bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Components</h2>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {theme.components.length} component{theme.components.length !== 1 ? 's' : ''}
            </p>
          </div>

          <div className="relative">
            <button
              onClick={() => setShowAddMenu(!showAddMenu)}
              className="flex items-center gap-2 px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg text-sm font-medium transition"
            >
              <Plus className="h-4 w-4" />
              Add
            </button>

            {/* Add Menu Dropdown */}
            {showAddMenu && (
              <div className="absolute right-0 mt-2 w-56 bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-lg shadow-lg z-10">
                <div className="p-2">
                  {COMPONENT_TYPES.map((comp) => (
                    <button
                      key={comp.type}
                      onClick={() => handleAddComponent(comp.type)}
                      className="w-full text-left px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition"
                    >
                      <p className="text-sm font-medium text-gray-900 dark:text-white">
                        {comp.label}
                      </p>
                      <p className="text-xs text-gray-600 dark:text-gray-400">{comp.description}</p>
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Components List */}
      <div className="flex-1 overflow-y-auto p-4 space-y-2">
        {theme.components.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-gray-500 text-center">
            <p className="text-sm">No components added yet</p>
            <p className="text-xs mt-2">Click "Add" to start building</p>
          </div>
        ) : (
          theme.components.map((component, index) => (
            <div
              key={component.id}
              draggable
              onDragStart={() => handleDragStart(component.id)}
              onDragOver={handleDragOver}
              onDrop={() => handleDrop(component.id)}
              onClick={() => selectComponent(component.id)}
              className={`p-3 rounded-lg border-2 transition cursor-move ${
                !component.enabled ? 'opacity-60' : ''
              } ${
                selectedComponentId === component.id
                  ? 'border-primary-600 bg-primary-50 dark:bg-primary-900/20'
                  : 'border-gray-200 dark:border-gray-600 hover:border-primary-300 dark:hover:border-primary-700'
              } ${draggedItem === component.id ? 'opacity-50' : ''}`}
            >
              <div className="flex items-start gap-2">
                {/* Drag Handle */}
                <div className="text-gray-400 mt-1 flex-shrink-0 cursor-grab active:cursor-grabbing">
                  <GripVertical className="h-4 w-4" />
                </div>

                {/* Component Icon */}
                <div className="text-primary-600 dark:text-primary-400 mt-1 flex-shrink-0">
                  {COMPONENT_ICONS[component.type as ComponentType] || <Box className="h-4 w-4" />}
                </div>

                {/* Component Info */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <p className="text-sm font-medium text-gray-900 dark:text-white truncate">
                      {component.title}
                    </p>
                    <span className="text-xs px-2 py-0.5 bg-gray-200 dark:bg-gray-600 text-gray-700 dark:text-gray-300 rounded">
                      #{index + 1}
                    </span>
                  </div>
                  <p className="text-xs text-gray-600 dark:text-gray-400 capitalize">
                    {component.type.replace('_', ' ')}
                  </p>
                </div>

                {/* Toggle Visibility */}
                <button
                  onClick={(e) => handleToggleComponent(component.id, e)}
                  aria-label={component.enabled ? 'Disable component' : 'Enable component'}
                  title={component.enabled ? 'Disable component' : 'Enable component'}
                  className={`flex-shrink-0 transition ${
                    component.enabled
                      ? 'text-green-600 hover:text-green-700 dark:text-green-400 dark:hover:text-green-300'
                      : 'text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'
                  }`}
                >
                  {component.enabled ? <Eye className="h-4 w-4" /> : <EyeOff className="h-4 w-4" />}
                </button>

                {/* Delete Button */}
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    deleteComponent(component.id)
                  }}
                  aria-label={`Delete ${component.title}`}
                  title="Delete component"
                  className="text-gray-400 hover:text-red-600 dark:hover:text-red-400 flex-shrink-0 transition"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  )
}
