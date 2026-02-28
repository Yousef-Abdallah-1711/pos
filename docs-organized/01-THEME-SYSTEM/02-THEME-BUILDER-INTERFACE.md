# 🎨 2. THEME BUILDER INTERFACE & SETTINGS

**Last Updated**: 2025-01-09
**Status**: ✅ Complete & Fully Functional
**Version**: 1.0

---

## TABLE OF CONTENTS

1. [Overview](#overview)
2. [Theme Builder UI Structure](#theme-builder-ui-structure)
3. [Presets Gallery](#presets-gallery)
4. [Theme Editor Pages](#theme-editor-pages)
5. [Settings Panels](#settings-panels)
6. [How Each Theme Looks Unique](#how-each-theme-looks-unique)
7. [Component Management](#component-management)
8. [Live Preview System](#live-preview-system)
9. [File Locations](#file-locations)
10. [Features & Capabilities](#features--capabilities)

---

## OVERVIEW

The **Theme Builder** is a comprehensive dashboard interface that allows users to:

1. ✅ **Browse** 16 pre-designed theme presets
2. ✅ **Create** themes by duplicating presets
3. ✅ **Customize** colors, typography, headers, and footers
4. ✅ **Manage** components (enable/disable, reorder)
5. ✅ **Preview** changes in real-time
6. ✅ **Save** themes to database
7. ✅ **Publish** themes to public website
8. ✅ **Support** bilingual content (EN/AR)

### Architecture

```
┌──────────────────────────────────────────────────────────┐
│          THEME BUILDER DASHBOARD                         │
├──────────────────────────────────────────────────────────┤
│ URL: /dashboard/theme-builder                            │
│ Auth: Required (JWT token)                               │
│ Responsive: Mobile, Tablet, Desktop                      │
│ Languages: English (LTR), Arabic (RTL)                   │
└──────────────────────────────────────────────────────────┘
       │
       ├─ /presets (Presets Gallery)
       │  - Browse 16 themes
       │  - Filter by category
       │  - View theme details
       │  - Create from preset
       │
       ├─ /editor/{id} (Theme Editor)
       │  - Customize colors
       │  - Modify typography
       │  - Configure header
       │  - Configure footer
       │  - Manage components
       │  - Live preview
       │  - Save/Publish
       │
       └─ /list (My Themes)
          - List user's themes
          - Activate/Deactivate
          - Edit/Delete
          - Export/Import
```

---

## THEME BUILDER UI STRUCTURE

### Main Sections

#### 1. Presets Gallery Page (`/dashboard/theme-builder/presets`)

**Purpose**: Display all 16 production themes for selection

**Components**:
- **Header**: "Choose a Theme Template"
- **Description**: "Start with a pre-designed theme and customize..."
- **Filter Tabs**: Professional, Luxury, Modern, Casual, Playful, Other
- **Theme Cards**: Grid layout (1-4 columns responsive)
- **Each Card Shows**:
  - Theme name
  - Theme description
  - Color preview (5 color swatches)
  - Component count badge
  - Theme tags/labels
  - "Use Template" overlay on hover

**Layout**:
```
┌─ Presets Gallery Page ─────────────────────────────────────┐
│ Choose a Theme Template                                    │
│ Start with a pre-designed theme and customize...            │
│                                                             │
│ PROFESSIONAL CATEGORY                                       │
│ ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│ │ Modern   │  │ Elegant  │  │ Urban    │  │ Warm     │   │
│ │ Bistro   │  │Simple    │  │ Fresh    │  │Comfort   │   │
│ │ 4 comps  │  │ 3 comps  │  │ 4 comps  │  │ 4 comps  │   │
│ └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│                                                             │
│ LUXURY CATEGORY                                            │
│ ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│ │ Premium  │  │ Purple   │  │ Elegant  │                  │
│ │ Dark     │  │ Luxury   │  │ Dark     │                  │
│ │ 5 comps  │  │ 4 comps  │  │ 4 comps  │                  │
│ └──────────┘  └──────────┘  └──────────┘                  │
│ ... (more categories)                                      │
│                                                             │
│ 💡 Tip: Click on any template to start creating...        │
└─────────────────────────────────────────────────────────────┘
```

#### 2. Theme Editor Page (`/dashboard/theme-builder/editor/{id}`)

**Purpose**: Allow detailed customization of theme

**Layout** (Two-Column):
```
┌─────────────────────────────────────────────────────────────┐
│ Back | Theme Name | Save | Publish                         │
├─────────────────────┬─────────────────────────────────────┤
│                     │                                       │
│  SETTINGS PANEL     │       LIVE PREVIEW                    │
│                     │                                       │
│ ┌─────────────────┐ │ ┌──────────────────────────────────┐ │
│ │ Colors          │ │ │  Theme Preview                   │ │
│ │ - Primary       │ │ │  (Updated in real-time)          │ │
│ │ - Secondary     │ │ │                                  │ │
│ │ - Accent        │ │ │  Header                          │ │
│ │ - Etc.          │ │ │  ┌──────────────────────────┐    │ │
│ │                 │ │ │  │ [Logo] Home Menu About    │    │ │
│ │ Typography      │ │ │  │ Contact         Order Now│    │ │
│ │ - Font Family   │ │ │  └──────────────────────────┘    │ │
│ │ - Base Size     │ │ │                                  │ │
│ │ - Line Height   │ │ │  [Hero Section Preview]         │ │
│ │                 │ │ │                                  │ │
│ │ Header          │ │ │  [Components Preview]           │ │
│ │ - Style         │ │ │                                  │ │
│ │ - Layout        │ │ │  Footer                         │ │
│ │ - Nav Items     │ │ │  ┌──────────────────────────┐    │ │
│ │                 │ │ │  │ Quick Links | Hours | ..│    │ │
│ │ Footer          │ │ │  │ © 2025 Modern Bistro     │    │ │
│ │ - Layout        │ │ │  └──────────────────────────┘    │ │
│ │ - Sections      │ │ │                                  │ │
│ │                 │ │ └──────────────────────────────────┘ │
│ │ Components      │ │                                       │
│ │ - Hero          │ │                                       │
│ │ - Featured      │ │                                       │
│ │ - Why Us        │ │                                       │
│ │ - CTA           │ │                                       │
│ │                 │ │                                       │
│ └─────────────────┘ │                                       │
└─────────────────────┴─────────────────────────────────────────┘
```

---

## PRESETS GALLERY

### Gallery Features

#### Theme Card Display
Each theme card shows:

```
┌─────────────────────────────┐
│  Theme Name                 │  ← Clickable
│                             │
│  ┌┐ ┌┐ ┌┐ ┌┐ ┌┐           │  ← 5 Color Swatches
│  └┘ └┘ └┘ └┘ └┘           │
│                             │
│  Theme description text... │  ← 2-line clamp
│                             │
│  [tag1] [tag2]   [4 comps] │  ← Tags & component count
│                             │
│  ┌─────────────────────────┐│  ← Hover overlay
│  │  Use Template →         ││
│  └─────────────────────────┘│
└─────────────────────────────┘
```

#### Categorization

Themes automatically grouped by category metadata:

```
meta: {
  category: "professional"  // "professional" | "luxury" | "modern" | "casual" | "playful"
}
```

Displays as sections:

```
PROFESSIONAL (2 themes)
LUXURY (1 theme)
MODERN (2 themes)
CASUAL (4 themes)
PLAYFUL (2 themes)
OTHER (4 themes + any uncategorized)
```

#### Responsive Grid

```
Mobile (< 640px):      1 column
Tablet (640-1024px):   2 columns
Desktop (1024-1280px): 3 columns
Large (> 1280px):      4 columns
```

### Gallery Interactions

**Hover** → Shows "Use Template" overlay
**Click** → Calls API to duplicate theme
**Creating** → Shows loading modal
**Success** → Redirects to editor page

---

## THEME EDITOR PAGES

### Top Navigation Bar

```
┌─────────────────────────────────────────────────────────┐
│ ← Back | Theme Name ▼ | [Save] | [Publish] | Menu      │
└─────────────────────────────────────────────────────────┘
```

**Elements**:
- **Back Button**: Return to presets or list
- **Theme Name**: Editable, shows current theme
- **Save Button**: Save all changes to database
- **Publish Button**: Make theme live
- **Menu**: More options (export, duplicate, delete)

### Settings Panel (Left)

The left sidebar contains all customization options:

#### 1. Colors Section

```
┌─ COLORS ───────────────────────┐
│ Click to expand/collapse       │
│                                │
│ Primary       [###] #2563eb    │  ← Click to open color picker
│ Secondary     [###] #059669    │
│ Accent        [###] #db2777    │
│ Background    [###] #ffffff    │
│ Text          [###] #111827    │
│ Border        [###] #e5e7eb    │
│ Shadow        [###] #000000    │
│                                │
│ [Update Color Palette Button]  │
└────────────────────────────────┘
```

**Features**:
- Click color swatch → Open color picker
- Enter hex value directly
- Preview updates in real-time
- Contrast checker (WCAG compliance)

#### 2. Typography Section

```
┌─ TYPOGRAPHY ───────────────────┐
│                                │
│ Font Family  [Select dropdown] │ ← Choose from 50+ fonts
│ Base Size    [16] px           │ ← Slider 10-24px
│ Line Height  [1.6] em          │ ← Slider 1.0-2.0
│ Border Rad.  [8] px            │ ← Slider 0-20px
│                                │
│ Heading 1    [48] px | [700]   │ ← Weight selector
│ Heading 2    [36] px | [600]   │
│ Heading 3    [24] px | [600]   │
│                                │
│ [Preview Typography]           │
└────────────────────────────────┘
```

**Features**:
- Font family selector (Google Fonts)
- Size sliders with live preview
- Heading customization
- Preview shows sample text

#### 3. Header Section

```
┌─ HEADER ───────────────────────┐
│                                │
│ Style        [Select: Modern]  │ ← modern|elegant|classic
│ Layout       [Horizontal]      │ ← horizontal|vertical
│ Position     [Sticky]          │ ← sticky|fixed|static
│ Height       [64] px           │ ← Slider
│                                │
│ Colors                         │
│ Background   [###] #2563eb     │
│ Text         [###] #ffffff     │
│                                │
│ Logo                           │
│ Show Logo    [Toggle: ON]      │
│ Position     [Left]            │ ← left|center|right
│                                │
│ Navigation                     │
│ Item 1: Home / (Order: 1)      │ ← Drag to reorder
│ Item 2: Menu / (Order: 2)      │
│ Item 3: About / (Order: 3)     │
│ Item 4: Contact / (Order: 4)   │
│ [Add Nav Item] [Remove Selected]│
│                                │
│ CTA Button                     │
│ Text: Order Now                │
│ Link: /order                   │
│ Style: Primary                 │
│                                │
│ [Save Header Config]           │
└────────────────────────────────┘
```

**Features**:
- Style preset selector
- Color customization
- Logo upload/toggle
- Navigation item management
- Drag-to-reorder
- CTA button configuration

#### 4. Footer Section

```
┌─ FOOTER ───────────────────────┐
│                                │
│ Style        [Extended]        │ ← extended|compact|minimal
│ Layout       [Multi-column]    │ ← multi-column|single|custom
│ Columns      [4]               │ ← 1-6 columns
│                                │
│ Colors                         │
│ Background   [###] #1f2937     │
│ Text         [###] #ffffff     │
│ Links        [###] #60a5fa     │
│                                │
│ Company Info                   │
│ Name: Modern Bistro            │
│ Address: 123 Main St...        │
│ Phone: (555) 123-4567          │
│ Email: contact@...             │
│                                │
│ Sections (Drag to reorder)     │
│ ☰ Quick Links                  │
│ ☰ Hours                        │
│ ☰ About                        │
│ [Add Section]                  │
│                                │
│ Social Links                   │
│ Facebook: https://...          │
│ Instagram: https://...         │
│ [Add Social]                   │
│                                │
│ Legal Links                    │
│ Privacy Policy: /privacy       │
│ Terms of Service: /terms       │
│ [Add Legal Link]               │
│                                │
│ Copyright Text                 │
│ © 2025 Modern Bistro. All rights│
│ [Save Footer Config]           │
└────────────────────────────────┘
```

**Features**:
- Layout style selector
- Column count adjuster
- Color customization
- Company information fields
- Section management
- Social links management
- Legal links management
- Copyright text customization

#### 5. Components Section

```
┌─ COMPONENTS ───────────────────┐
│ Total: 4 components enabled    │
│                                │
│ ☰ 1 ⚡ Hero Section            │ ← Drag handle + icon
│    [👁] Hero              [−]  │ ← Enable/disable + delete
│                                │
│ ☰ 2 📦 Featured Dishes         │
│    [👁] Featured Items    [−]  │
│                                │
│ ☰ 3 ⭐ Why Choose Us           │
│    [👁] Why Us            [−]  │
│                                │
│ ☰ 4 🎯 Ready to Dine?          │
│    [👁] CTA               [−]  │
│                                │
│ [Add Component]                │
│                                │
│ Component Details (Selected)   │
│ Type: hero                     │
│ Enabled: Yes                   │
│ [Edit Component Config]        │
│                                │
│ [Save Components]              │
└────────────────────────────────┘
```

**Features**:
- Component list with drag-to-reorder
- Enable/disable toggles
- Delete buttons
- Component type icons
- Component count display
- Edit individual components
- Add/remove components

---

## SETTINGS PANELS

### How Settings Work

Each setting panel is **collapsible** and contains related options:

1. **Click header** → Expands/collapses section
2. **Change value** → Instantly updates preview
3. **Save changes** → All updates auto-saved to browser state
4. **Click "Save Theme"** → All changes sent to backend

### Color Picker Modal

```
┌─────────────────────────────────────┐
│ Select Color for: Primary           │
├─────────────────────────────────────┤
│                                     │
│  ┌──────────────────────────────┐   │
│  │                              │   │  ← Large color selector
│  │  ● (Cursor to pick color)    │   │
│  │                              │   │
│  └──────────────────────────────┘   │
│                                     │
│  Saturation                         │
│  │▓▓▓▓░░░░│  75%                    │  ← Saturation slider
│                                     │
│  Brightness                         │
│  │▓▓▓░░░░░│  50%                    │  ← Brightness slider
│                                     │
│  HEX Code                           │
│  [#2563eb]                          │  ← Direct hex input
│                                     │
│  RGB                                │
│  R: [37]  G: [99]  B: [235]         │  ← RGB values
│                                     │
│  Recent Colors                      │
│  [●][●][●][●][●]                    │  ← Previously used colors
│                                     │
│  [Cancel]  [Apply]                  │  ← Action buttons
└─────────────────────────────────────┘
```

### Font Family Selector

```
┌─────────────────────────────────────┐
│ Select Font Family                  │
├─────────────────────────────────────┤
│ Search fonts...        [search box] │
│                                     │
│ Popular Fonts                       │
│ ☐ Inter (Current)                   │
│ ☐ Roboto                            │
│ ☐ Open Sans                         │
│ ☐ Poppins                           │
│ ☐ Playfair Display                  │
│                                     │
│ All Fonts (50+)                     │
│ ☐ Abel                              │
│ ☐ Abril Fatface                     │
│ ... (many more)                     │
│                                     │
│ Preview                             │
│ "The quick brown fox jumps"         │ (in selected font)
│                                     │
│ [Cancel]  [Select]                  │
└─────────────────────────────────────┘
```

---

## HOW EACH THEME LOOKS UNIQUE

### Visual Differences Per Theme

#### Modern Bistro (Professional)

```
HEADER:
┌─────────────────────────────────────┐
│ [Logo] Home Menu About Contact Order│  ← Blue bg, white text
└─────────────────────────────────────┘   Sticky, horizontal

HERO:
┌────────────────────────────────────────┐
│                                        │
│      Welcome to Modern Bistro          │  ← Centered, overlay
│                                        │
│   Contemporary minimalist dining       │
│                                        │
│    [View Menu] [Make Reservation]      │
│                                        │
└────────────────────────────────────────┘

FEATURED DISHES:
┌────────────────────────────────────────┐
│ Our Signature Dishes                   │
│ [Dish 1]  [Dish 2]  [Dish 3]           │  ← 3-column grid
│ [Dish 4]  [Dish 5]  [Dish 6]           │
└────────────────────────────────────────┘

FOOTER:
┌────────────────────────────────────────┐
│ Quick Links | Hours | About | Legal   │  ← 4-column extended
│ Facebook Instagram Twitter LinkedIn    │
│ © 2025 Modern Bistro. All rights...   │  ← Dark bg
└────────────────────────────────────────┘

COLOR SCHEME:
Primary: #2563eb (Blue)
Secondary: #059669 (Green)
Accent: #db2777 (Pink)
```

#### Elegant Simplicity (Luxury)

```
HEADER:
┌─────────────────────────────────────┐
│        [Logo]                        │  ← Centered, elegant
│   Home | Menu | About | Contact      │
└─────────────────────────────────────┘

HERO:
┌────────────────────────────────────────┐
│                                        │
│        Elegant Simplicity              │  ← Large elegant typography
│                                        │
│     Refined Dining Experience         │
│                                        │
│         [Reserve Table]                │
│                                        │
└────────────────────────────────────────┘

FEATURED:
┌────────────────────────────────────────┐
│ The Finest Selection                   │
│ [Item 1]  [Item 2]  [Item 3]           │  ← Elegant layout
└────────────────────────────────────────┘

FOOTER:
┌────────────────────────────────────────┐
│ About Us | Reservations | Contact     │  ← 3-column elegant
│ © 2025 Elegant Simplicity              │
└────────────────────────────────────────┘

COLOR SCHEME:
Primary: #1e3a8a (Navy)
Secondary: #f3f4f6 (Light)
Accent: #d4af37 (Gold)
```

#### Warm Comfort (Casual)

```
HEADER:
┌─────────────────────────────────────┐
│ 🏠 Warm Comfort    Home Menu About  │  ← Warm orange bg
│                            Order    │
└─────────────────────────────────────┘

HERO:
┌────────────────────────────────────────┐
│        Welcome to Our Table            │  ← Warm, inviting
│                                        │
│     Classic comfort food made fresh    │
│                                        │
│      [View Our Menu]                   │
│                                        │
└────────────────────────────────────────┘

FEATURED:
┌────────────────────────────────────────┐
│ Customer Favorites                     │
│ ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│ │ [Image] │ │ [Image] │ │ [Image] │   │ ← Large images
│ │ Classic │ │ Comfort │ │ Family  │   │
│ └─────────┘ └─────────┘ └─────────┘   │
└────────────────────────────────────────┘

FOOTER:
┌────────────────────────────────────────┐
│ Hours | Location | Family-Friendly    │  ← Warm tone
│ Facebook Instagram Yelp               │
│ © 2025 Warm Comfort Restaurant         │
└────────────────────────────────────────┘

COLOR SCHEME:
Primary: #d97706 (Orange)
Secondary: #92400e (Brown)
Accent: #fbbf24 (Warm Gold)
```

### Each Theme's Unique Characteristics

**Modern Bistro**:
- Blue & green color scheme
- Contemporary minimal design
- Sticky horizontal header
- 4-column footer
- 4 pre-configured components

**Elegant Simplicity**:
- Navy & gold color scheme
- Centered, refined layout
- Elegant typography
- 3-column footer
- 3 pre-configured components

**Premium Dark**:
- Dark background with gold accents
- Luxury presentation
- Minimal header
- Extended footer
- 5 pre-configured components

**Urban Fresh**:
- Green color scheme
- Modern, clean design
- Urban aesthetic
- Compact footer
- 4 pre-configured components

... (and so on for all 16 themes)

**Each theme is distinctive because**:
1. ✅ Different color palette (7 unique colors)
2. ✅ Different typography choices
3. ✅ Different header style/layout
4. ✅ Different footer configuration
5. ✅ Different component selection
6. ✅ Different component order

---

## COMPONENT MANAGEMENT

### Component Types

Available component types in theme builder:

```
1. Hero Section
   - Full-screen banner
   - Background image
   - Heading + subtitle
   - CTA buttons

2. Featured Items
   - Grid layout
   - Product/dish display
   - Images and descriptions
   - Price display

3. Why Choose Us
   - 3-column layout
   - Benefits/features
   - Icons
   - Description text

4. Testimonials
   - Quote carousel
   - Customer photos
   - Star ratings
   - Company names

5. Contact Section
   - Contact form
   - Map embed
   - Business hours
   - Contact info

6. CTA (Call to Action)
   - Large heading
   - Description
   - Primary button
   - Newsletter signup

7. Custom
   - User-defined components
   - Any layout
   - Any content
```

### Component Management UI

```
┌─ COMPONENTS ─────────────────────────────┐
│                                          │
│ ☰ 1 ⚡ Hero Section                      │  ← Drag handle
│    [👁 ON] [−]                          │  ← Toggle & delete
│    Component type: hero                 │
│    Status: Enabled                      │
│    [Edit Config]                        │
│                                          │
│ ☰ 2 📦 Featured Items                    │
│    [👁 ON] [−]                          │
│    Component type: featured-items       │
│    Status: Enabled                      │
│    [Edit Config]                        │
│                                          │
│ ☰ 3 ⭐ Why Choose Us                     │
│    [👁 ON] [−]                          │
│    Component type: why_us               │
│    Status: Enabled                      │
│    [Edit Config]                        │
│                                          │
│ ☰ 4 🎯 Call to Action                    │
│    [👁 OFF] [−]                         │  ← Disabled (grayed)
│    Component type: cta                  │
│    Status: Disabled                     │
│    [Edit Config]                        │
│                                          │
│ Total Enabled: 3/4                      │
│ [Add Component]                         │
│                                          │
└──────────────────────────────────────────┘
```

### Reordering Components

**Drag & Drop**:
1. Click and hold "☰" handle
2. Drag to new position
3. Drop to reorder
4. Order updates in preview

**Result**: Components display in custom order on website

---

## LIVE PREVIEW SYSTEM

### Real-Time Updates

When user changes a setting, the preview updates **instantly** (client-side):

```
User changes Primary color: #2563eb → #ff0000
     │
     ├─ Zustand store updates
     │
     ├─ React component re-renders
     │
     ├─ CSS variables updated
     │
     ├─ Live preview section re-renders
     │
     ├─ All elements with primary color now RED
     │
     └─ NO API call (client-side only)
```

### Preview Sections

Live preview shows:

1. **Header** - With logo, navigation, colors
2. **Hero Component** - Banner with image/title
3. **Featured Items** - Product/dish grid
4. **Why Us** - Benefits section
5. **CTA Section** - Call to action
6. **Footer** - With columns, links, social

### Responsive Preview

Preview adjustable to different screen sizes:

```
[📱 Mobile] [📱 Tablet] [💻 Desktop]

Mobile (375px):
┌──────────┐
│ Header   │
│          │
│ Hero     │  ← 1 column layout
│ Featured │
│ Why Us   │
│ CTA      │
│ Footer   │
└──────────┘

Tablet (768px):
┌────────────────┐
│ Header         │
│                │
│ Hero           │
│ [F1][F2][F3]   │  ← 3-column grid
│ [F4][F5][F6]   │
│ Why Us         │
│ Footer         │
└────────────────┘

Desktop (1200px):
┌──────────────────────────┐
│ Header                   │
│                          │
│ Hero                     │
│ [F1][F2][F3][F4]         │  ← 4-column grid
│ [F5][F6][F7][F8]         │
│ Why Us                   │
│ Footer                   │
└──────────────────────────┘
```

---

## FILE LOCATIONS

### Frontend Components

**Theme Builder Pages**:
- `/frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/page.tsx` - Main page (theme list)
- `/frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/presets/page.tsx` - Presets gallery
- `/frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx` - Single editor

**Components**:
- `/frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/components/EditorSidebar.tsx` - Settings panel (left)
- `/frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/components/EditorPreview.tsx` - Live preview (right)
- `/frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/components/ComponentLibrary.tsx` - Component add/reorder/toggle UI
- `/frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/components/ComponentEditor.tsx` - Routes to per-component editors
- `/frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/components/ColorPicker.tsx` - Color picker
- `/frontend/apps/dashboard/src/components/editors/HeaderEditor.tsx` - Header editor
- `/frontend/apps/dashboard/src/components/editors/FooterEditor.tsx` - Footer editor
- `/frontend/apps/dashboard/src/components/theme/ExportDialog.tsx` - Export theme
- `/frontend/apps/dashboard/src/components/theme/ImportDialog.tsx` - Import theme (if enabled in current build)

**State Management**:
- `/frontend/apps/dashboard/src/hooks/useThemeStore.ts` - Active Theme Builder editor store (current editor route)
- `/frontend/apps/dashboard/src/stores/themeBuilderStore.ts` - Older/store utility path still present in repo
- `/frontend/apps/dashboard/src/stores/themeBuilderHooks.ts` - Associated hooks for older store path

**APIs & Utilities**:
- `/frontend/apps/dashboard/src/lib/api/themeApi.ts` - Theme API client
- `/frontend/apps/dashboard/src/types/theme.ts` - Core dashboard theme types
- `/frontend/apps/restaurant-website/lib/api/get-theme-server.ts` - Storefront theme extraction/normalization entry

**Shared Runtime Renderers (Preview + Storefront)**:
- `/frontend/layers/ui-themes/src/components/layout/Header.tsx` - Data-driven header renderer
- `/frontend/layers/ui-themes/src/components/layout/Footer.tsx` - Data-driven footer renderer
- `/frontend/layers/ui-themes/src/components/runtime/SectionRenderer.tsx` - Dynamic section renderer
- `/frontend/layers/ui-themes/src/sections/{marketing,commerce,content,forms,custom}/*.tsx` - Actual section implementations (hero/products/cta/etc.)

### Styling

**Tailwind Classes Used**:
- `grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4` - Responsive grid
- `bg-white dark:bg-gray-800` - Dark mode support
- `border border-gray-200 dark:border-gray-700` - Borders
- `hover:shadow-lg` - Hover effects
- `transition-all` - Smooth animations
- `text-gray-900 dark:text-white` - Text contrast

---

## FEATURES & CAPABILITIES

### ✅ What's Implemented

**Presets Gallery**:
- ✅ Display all 16 themes
- ✅ Organize by category
- ✅ Show theme metadata
- ✅ Responsive grid layout
- ✅ Color preview swatches
- ✅ Component count badges
- ✅ Create from preset

**Theme Editor**:
- ✅ Colors: All 7 colors customizable
- ✅ Typography: Font, size, line-height, border-radius
- ✅ Header: Style, layout, colors, navigation
- ✅ Footer: Style, layout, colors, sections
- ✅ Components: Enable/disable, reorder, delete
- ✅ Live preview: Real-time updates
- ✅ Save: Persist to database
- ✅ Publish: Make theme live

**Component Management**:
- ✅ Component type icons
- ✅ Enable/disable toggles
- ✅ Drag-to-reorder
- ✅ Delete components
- ✅ Component details display
- ✅ Component count
- ✅ Add components

**Bilingual Support**:
- ✅ EN/AR text fields
- ✅ Bilingual utilities
- ✅ Text extraction helpers
- ✅ Direction handling (LTR/RTL)

**Responsive Design**:
- ✅ Mobile (< 640px)
- ✅ Tablet (640-1024px)
- ✅ Desktop (> 1024px)
- ✅ Large displays (> 1280px)

**Dark Mode**:
- ✅ Full dark mode support
- ✅ Dark theme colors
- ✅ Dark component styling

---

## SUMMARY

The **Theme Builder Interface** provides a complete, user-friendly system for:

✅ **Browsing** 16 unique themes with visual previews
✅ **Creating** themes from presets
✅ **Customizing** every aspect (colors, fonts, header, footer, components)
✅ **Managing** components (enable, disable, reorder)
✅ **Previewing** changes in real-time
✅ **Supporting** bilingual content (EN/AR)
✅ **Saving** and publishing themes

All with a **responsive, intuitive interface** that works on mobile, tablet, and desktop.

For system architecture, see: `01-THEMES-SYSTEM-ARCHITECTURE.md`
For production deployment, see: `03-THEME-PRODUCTION-CHECKLIST.md`

---

**Questions? Issues?** Refer to the frontend code:
- `/src/components/theme/` - All UI components
- `/src/store/themeBuilderStore.ts` - State management
- `/src/lib/api/themeApi.ts` - API integration
