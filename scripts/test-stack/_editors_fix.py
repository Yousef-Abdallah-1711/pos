"""Phase 2A: Header/Footer editor typing. Items loaded without ids get stable ids when the config is read."""
ROOT = "frontend/apps/dashboard/src/"


def rw(rel, pairs):
    p = ROOT + rel
    s = open(p, encoding="utf-8").read()
    crlf = "\r\n" in s
    s = s.replace("\r\n", "\n")
    for a, b in pairs:
        assert a in s, (rel, a[:70])
        s = s.replace(a, b, 1)
    open(p, "w", encoding="utf-8", newline="").write(s.replace("\n", "\r\n") if crlf else s)


# the footer editor offers a 'minimal' layout the shared type did not know about
rw("types/theme.ts", [("  columns?: number\n  layout?: 'expanded' | 'compact'\n}", "  columns?: number\n  layout?: 'expanded' | 'compact' | 'minimal'\n}")])

# ---- HeaderEditor
rw("components/editors/HeaderEditor.tsx", [
    ("interface HeaderConfig extends ThemeHeaderConfig {\n  transparentBackground?: boolean",
     "interface HeaderConfig extends Omit<ThemeHeaderConfig, 'navigationItems'> {\n  navigationItems: NavItem[]\n  shadow?: boolean\n  transparentBackground?: boolean"),
    ("  const initialHeaderConfig: HeaderConfig = useMemo(() => {\n    if (useStore && storeHeader) {\n      return {\n        ...storeHeader,",
     "  const initialHeaderConfig: HeaderConfig = useMemo(() => {\n    if (useStore && storeHeader) {\n      return {\n        ...withNavIds(storeHeader),"),
    ("    return (component?.config as unknown as HeaderConfig) || {\n      navigationItems: [],",
     "    if (component?.config) {\n      return withNavIds(component.config as unknown as ThemeHeaderConfig) as HeaderConfig\n    }\n    return {\n      navigationItems: [],"),
    ("/**\n * HeaderEditor Component", "/**\n * Items saved without ids cannot be edited (the editor looks them up by id), so give them stable ones.\n */\nfunction withNavIds(config: ThemeHeaderConfig): Omit<ThemeHeaderConfig, 'navigationItems'> & { navigationItems: NavItem[] } {\n  return {\n    ...config,\n    navigationItems: (config.navigationItems ?? []).map((item, index) => ({\n      ...item,\n      id: item.id ?? `nav-${index + 1}`,\n      order: item.order ?? index + 1,\n    })),\n  }\n}\n\n/**\n * HeaderEditor Component"),
])

# ---- FooterEditor
rw("components/editors/FooterEditor.tsx", [
    ("interface FooterConfig extends ThemeFooterConfig {\n  layout?: 'compact' | 'expanded' | 'minimal'\n  legalLinks?: LinkItem[]\n}",
     "interface FooterConfig extends Omit<ThemeFooterConfig, 'footerSections' | 'socialLinks'> {\n  footerSections?: FooterSection[]\n  socialLinks: SocialLink[]\n  layout?: 'compact' | 'expanded' | 'minimal'\n  legalLinks?: LinkItem[]\n}\n\n/** Sections and social links saved without ids cannot be edited (lookups are by id): give them stable ones. */\nfunction withFooterIds(config: ThemeFooterConfig): Omit<ThemeFooterConfig, 'footerSections' | 'socialLinks'> & { footerSections: FooterSection[]; socialLinks: SocialLink[] } {\n  return {\n    ...config,\n    footerSections: (config.footerSections ?? []).map((section, index) => ({\n      ...section,\n      id: section.id ?? `section-${index + 1}`,\n      order: section.order ?? index + 1,\n    })),\n    socialLinks: (config.socialLinks ?? []).map((link, index) => ({\n      ...link,\n      id: link.id ?? `social-${index + 1}`,\n      order: link.order ?? index + 1,\n      platform: link.platform as SocialLink['platform'],\n    })),\n  }\n}"),
    ("      return {\n        ...storeFooter,\n        layout: 'expanded',", "      return {\n        ...withFooterIds(storeFooter),\n        layout: 'expanded',"),
    ("    return (component?.config as unknown as FooterConfig) || {\n      footerSections: [],",
     "    if (component?.config) {\n      return withFooterIds(component.config as unknown as ThemeFooterConfig) as FooterConfig\n    }\n    return {\n      footerSections: [],"),
])
print("ok")
