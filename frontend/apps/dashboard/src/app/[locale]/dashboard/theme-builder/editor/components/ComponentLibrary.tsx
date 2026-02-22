'use client'

import { useState, useMemo } from 'react'
import { Plus, Trash2, GripVertical, Search, X, Eye, EyeOff } from 'lucide-react'
import { useThemeStore } from '@/hooks/useThemeStore'
import { ComponentConfig, ComponentType } from '@/types/theme'
import { componentRegistry } from '@pos-saas/ui-components/registry'
import { debugWarn } from '@/lib/debug'

// Helper function to extract label from translation objects
function extractLabel(label: any): string {
  if (typeof label === 'string') return label
  if (label && typeof label === 'object') {
    // Handle translation objects like {en: "...", ar: "..."}
    return label.en || label.ar || Object.values(label)[0] || ''
  }
  return ''
}

// Map registry component IDs to dashboard ComponentType
const REGISTRY_TO_COMPONENT_TYPE: Record<string, ComponentType> = {
  'hero-section': 'hero',
  'products-section': 'products',
  'why-us-section': 'why_us',
  'contact-section': 'contact',
  'testimonials-section': 'testimonials',
  'cta-section': 'cta',
  'custom-section': 'custom',
}

export function ComponentLibrary() {
  const theme = useThemeStore((state) => state.currentTheme)
  const addComponent = useThemeStore((state) => state.addComponent)
  const deleteComponent = useThemeStore((state) => state.deleteComponent)
  const updateComponent = useThemeStore((state) => state.updateComponent)
  const reorderComponents = useThemeStore((state) => state.reorderComponents)
  const selectComponent = useThemeStore((state) => state.selectComponent)
  const selectedComponentId = useThemeStore((state) => state.selectedComponentId)

  const [draggedItem, setDraggedItem] = useState<string | null>(null)
  const [showLibrary, setShowLibrary] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null)

  // Get all draggable components from registry
  const registryComponents = useMemo(() => {
    return componentRegistry.getDraggableComponents()
  }, [])

  // Get all unique categories from registry
  const allCategories = useMemo(() => {
    const categories = new Set<string>()
    registryComponents.forEach((comp) => {
      categories.add(comp.category)
    })
    return Array.from(categories).sort()
  }, [registryComponents])

  // Filter components based on search and category
  const filteredComponents = useMemo(() => {
    return registryComponents.filter((comp) => {
      // Category filter
      if (selectedCategory && comp.category !== selectedCategory) {
        return false
      }

      // Search filter (by name, tags, description)
      if (searchQuery.trim()) {
        const query = searchQuery.toLowerCase()
        const matchesName = comp.name.toLowerCase().includes(query)
        const matchesTags = comp.tags?.some((tag) => tag.toLowerCase().includes(query))
        const matchesDescription = comp.description?.toLowerCase().includes(query)
        return matchesName || matchesTags || matchesDescription
      }

      return true
    })
  }, [registryComponents, searchQuery, selectedCategory])

  if (!theme) return null
  const themeComponents = theme.components || []

  const handleAddComponent = (registryId: string) => {
    const registryComp = registryComponents.find((c) => c.id === registryId)
    if (!registryComp) return

    // Map registry ID to dashboard ComponentType
    const componentType = REGISTRY_TO_COMPONENT_TYPE[registryId] as ComponentType
    if (!componentType) {
      debugWarn(`Unknown component type mapping for ${registryId}`)
      return
    }

    const newComponent: ComponentConfig = {
      id: `component-${Date.now()}`,
      type: componentType,
      title: registryComp.name,
      displayOrder: themeComponents.length,
      enabled: true,
      config: registryComp.defaultConfig || {},
    }
    addComponent(newComponent)
  }

  const handleToggleComponent = (componentId: string) => {
    const target = themeComponents.find((comp) => String(comp.id) === componentId)
    if (!target) return
    updateComponent(componentId, { enabled: !target.enabled })
  }

  const handleDragStart = (id: string) => {
    setDraggedItem(id)
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
  }

  const handleDrop = (targetId: string) => {
    if (!draggedItem || draggedItem === targetId) return

    const draggedIdx = themeComponents.findIndex((c) => String(c.id) === draggedItem)
    const targetIdx = themeComponents.findIndex((c) => String(c.id) === targetId)

    if (draggedIdx < 0 || targetIdx < 0) return

    const newComponents = [...themeComponents]
    const [draggedComp] = newComponents.splice(draggedIdx, 1)
    newComponents.splice(targetIdx, 0, draggedComp)

    // ✅ Update display order for all components after reordering
    const updatedComponents = newComponents.map((comp, index) => ({
      ...comp,
      displayOrder: index,
    }))

    reorderComponents(updatedComponents)
    setDraggedItem(null)
  }

  return (
    <div className="space-y-6">
      {/* Current Components */}
      <div>
        <div className="flex items-center justify-between mb-3">
          <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300">
            Current Components ({themeComponents.length})
          </h4>
        </div>

        <div className="space-y-2 max-h-64 overflow-y-auto">
          {themeComponents.length === 0 ? (
            <p className="text-xs text-gray-500 dark:text-gray-400 py-4 text-center">
              No components added yet
            </p>
          ) : (
            themeComponents.map((component) => (
              <div
                key={component.id}
                draggable
                onDragStart={() => handleDragStart(String(component.id))}
                onDragOver={handleDragOver}
                onDrop={() => handleDrop(String(component.id))}
                onClick={() => selectComponent(String(component.id))}
                className={`p-3 rounded-lg border-2 transition cursor-move ${
                  selectedComponentId === String(component.id)
                    ? 'border-primary-600 bg-primary-50 dark:bg-primary-900/20'
                    : 'border-gray-200 dark:border-gray-600 hover:border-primary-300'
                } ${draggedItem === String(component.id) ? 'opacity-50' : ''}`}
              >
                <div className="flex items-center gap-2">
                  <GripVertical className="h-4 w-4 text-gray-400 flex-shrink-0" />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 dark:text-white truncate">
                      {extractLabel(component.title)}
                    </p>
                    <p className="text-xs text-gray-600 dark:text-gray-400 capitalize">
                      {component.type}
                    </p>
                  </div>
                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        handleToggleComponent(String(component.id))
                      }}
                      className="p-1 text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 transition"
                      title={component.enabled ? 'Hide component' : 'Show component'}
                  >
                    {component.enabled ? (
                      <Eye className="h-4 w-4" />
                    ) : (
                      <EyeOff className="h-4 w-4 opacity-50" />
                    )}
                  </button>
                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        deleteComponent(String(component.id))
                      }}
                      className="p-1 text-gray-400 hover:text-red-600 dark:hover:text-red-400 transition"
                    >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      {/* Component Library */}
      <div>
        <button
          onClick={() => setShowLibrary(!showLibrary)}
          className="w-full text-left px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition font-medium text-sm text-gray-900 dark:text-white"
        >
          {showLibrary ? '− ' : '+ '}Add Components
        </button>

        {showLibrary && (
          <div className="mt-3 space-y-3">
            {/* Search Input */}
            <div className="relative">
              <Search className="absolute left-3 top-2.5 h-4 w-4 text-gray-400" />
              <input
                type="text"
                placeholder="Search components..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-3 py-2 text-sm rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:border-primary-500"
              />
              {searchQuery && (
                <button
                  onClick={() => setSearchQuery('')}
                  className="absolute right-2 top-2.5 p-0.5 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                >
                  <X className="h-4 w-4" />
                </button>
              )}
            </div>

            {/* Category Filters */}
            {allCategories.length > 0 && (
              <div className="flex flex-wrap gap-2">
                <button
                  onClick={() => setSelectedCategory(null)}
                  className={`px-3 py-1 rounded-full text-xs font-medium transition ${
                    selectedCategory === null
                      ? 'bg-primary-600 text-white'
                      : 'bg-gray-200 dark:bg-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-500'
                  }`}
                >
                  All
                </button>
                {allCategories.map((category) => (
                  <button
                    key={category}
                    onClick={() => setSelectedCategory(category)}
                    className={`px-3 py-1 rounded-full text-xs font-medium transition capitalize ${
                      selectedCategory === category
                        ? 'bg-primary-600 text-white'
                        : 'bg-gray-200 dark:bg-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-500'
                    }`}
                  >
                    {category}
                  </button>
                ))}
              </div>
            )}

            {/* Component List */}
            <div className="space-y-2 max-h-96 overflow-y-auto">
              {filteredComponents.length === 0 ? (
                <p className="text-xs text-gray-500 dark:text-gray-400 py-4 text-center">
                  {registryComponents.length === 0
                    ? 'No components available'
                    : 'No components match your search'}
                </p>
              ) : (
                filteredComponents.map((comp) => (
                  <button
                    key={comp.id}
                    onClick={() => handleAddComponent(comp.id)}
                    className="w-full text-left px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 hover:bg-primary-50 dark:hover:bg-primary-900/20 hover:border-primary-300 dark:hover:border-primary-500 transition text-sm"
                  >
                    <div className="flex items-center gap-3">
                      <span className="text-lg">{comp.icon}</span>
                      <div className="flex-1 min-w-0">
                        <p className="font-medium text-gray-900 dark:text-white truncate">
                          {comp.name}
                        </p>
                        <p className="text-xs text-gray-600 dark:text-gray-400 truncate">
                          {comp.description}
                        </p>
                      </div>
                      <Plus className="h-4 w-4 text-gray-400 flex-shrink-0" />
                    </div>
                  </button>
                ))
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
