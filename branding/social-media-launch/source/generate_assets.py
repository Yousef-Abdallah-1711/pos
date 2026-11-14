from __future__ import annotations

import base64
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
BG = ROOT / "assets" / "generated-bg"
FONTS = ROOT / "fonts"

COLORS = {
    "violet": "#6D28D9", "deep": "#4C1D95", "purple": "#8B5CF6",
    "lavender": "#C4B5FD", "paper": "#F7F5FF", "surface": "#F0ECF8",
    "ink": "#171327", "muted": "#575066", "rule": "#D9D1EB", "white": "#FFFFFF",
}
FONT_FILES = {
    "alexandria-arabic": "alexandria-arabic.woff2",
    "alexandria-latin": "alexandria-latin.woff2",
    "plex-400": "ibm-plex-arabic-400.woff2",
    "plex-600": "ibm-plex-arabic-600.woff2",
    "plex-700": "ibm-plex-arabic-700.woff2",
    "manrope": "manrope-latin.woff2",
}


def font_css() -> str:
    faces = []
    for name, filename in FONT_FILES.items():
        data = base64.b64encode((FONTS / filename).read_bytes()).decode("ascii")
        if name.startswith("alexandria"):
            family, weight = "Alexandria", "100 900"
        elif name.startswith("plex"):
            family, weight = "IBM Plex Sans Arabic", name.split("-")[1]
        else:
            family, weight = "Manrope", "200 800"
        faces.append(f"@font-face{{font-family:'{family}';font-weight:{weight};src:url(data:font/woff2;base64,{data}) format('woff2')}}")
    return "".join(faces)


FONT_CSS = font_css()


def svg(width: int, height: int, body: str, title: str, desc: str = "") -> str:
    desc_tag = f"<desc>{desc}</desc>" if desc else ""
    return (
        f'<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" '
        f'width="{width}" height="{height}" viewBox="0 0 {width} {height}" role="img">'
        f"<title>{title}</title>{desc_tag}<defs><style>{FONT_CSS}</style></defs>{body}</svg>"
    )


def image_data(filename: str) -> str:
    data = base64.b64encode((BG / filename).read_bytes()).decode("ascii")
    return f"data:image/png;base64,{data}"


def mark(x: int, y: int, size: int, color: str = "violet", node: str = "purple") -> str:
    return (
        f'<g transform="translate({x} {y}) scale({size / 100:.6f})" fill="none" stroke-linecap="round" stroke-linejoin="round">'
        f'<path d="M62 17C52 11 40 10 29 14 11 21 4 42 12 60c8 19 29 29 48 22 17-6 28-22 28-40 0-7-2-13-6-18" stroke="{COLORS[color]}" stroke-width="10"/>'
        f'<path d="M69 20h10v10" stroke="{COLORS[color]}" stroke-width="7"/>'
        f'<circle cx="62" cy="17" r="6.5" fill="{COLORS[node]}" stroke="none"/>'
        f'<circle cx="79" cy="30" r="5.2" fill="{COLORS["lavender"]}" stroke="none"/></g>'
    )


def wordmark(x: int, y: int, size: int, color: str = "ink", anchor: str = "start") -> str:
    return (
        f'<text x="{x}" y="{y}" text-anchor="{anchor}" fill="{COLORS[color]}" '
        f'font-family="Manrope" font-size="{size}" font-weight="800" letter-spacing="-1.5">ORDERLY</text>'
    )


def arabic(x: int, y: int, size: int, text: str, color: str = "ink", anchor: str = "start", weight: int = 700, family: str = "Alexandria") -> str:
    return (
        f'<text x="{x}" y="{y}" text-anchor="{anchor}" direction="rtl" unicode-bidi="plaintext" '
        f'fill="{COLORS[color]}" font-family="{family}" font-size="{size}" font-weight="{weight}">{text}</text>'
    )


def background(filename: str, width: int, height: int) -> str:
    return (
        f'<image x="0" y="0" width="{width}" height="{height}" preserveAspectRatio="xMidYMid slice" '
        f'href="{image_data(filename)}"/>'
    )


def write_svg(relative_path: str, width: int, height: int, body: str, title: str, desc: str = "") -> None:
    destination = ROOT / relative_path
    destination.parent.mkdir(parents=True, exist_ok=True)
    destination.write_text(svg(width, height, body, title, desc), encoding="utf-8")


def logo_assets() -> None:
    write_svg("logo/symbol/orderly-symbol.svg", 160, 160, mark(20, 20, 120), "ORDERLY symbol", "Open O route mark with connected nodes.")
    write_svg("logo/english/orderly-wordmark-en.svg", 600, 140, wordmark(24, 98, 76), "ORDERLY English wordmark")
    write_svg("logo/arabic/orderly-wordmark-ar.svg", 540, 160, arabic(510, 108, 72, "أوردرلي"), "ORDERLY Arabic wordmark")
    write_svg("logo/primary/orderly-primary.svg", 820, 180, mark(18, 30, 120) + wordmark(164, 105, 72), "ORDERLY horizontal logo")
    write_svg("logo/primary/orderly-primary-arabic.svg", 840, 190, mark(18, 32, 120) + arabic(812, 122, 74, "أوردرلي"), "ORDERLY Arabic horizontal logo")
    write_svg("logo/light/orderly-logo-light.svg", 820, 180, mark(18, 30, 120, "violet", "lavender") + wordmark(164, 105, 72, "white"), "ORDERLY logo for dark surfaces")
    write_svg("logo/dark/orderly-logo-dark.svg", 820, 180, mark(18, 30, 120) + wordmark(164, 105, 72), "ORDERLY logo for light surfaces")
    write_svg("logo/monochrome/orderly-logo-black.svg", 820, 180, mark(18, 30, 120, "ink", "ink") + wordmark(164, 105, 72, "ink"), "ORDERLY monochrome black logo")
    write_svg("logo/monochrome/orderly-logo-white.svg", 820, 180, mark(18, 30, 120, "white", "white") + wordmark(164, 105, 72, "white"), "ORDERLY monochrome white logo")


def avatar(size: int, background_color: str = "paper") -> str:
    center = size / 2
    return (
        f'<rect width="{size}" height="{size}" rx="{size * .24:.0f}" fill="{COLORS[background_color]}"/>'
        f'<circle cx="{center:.0f}" cy="{center:.0f}" r="{size*.37:.0f}" fill="{COLORS["white"]}"/>'
        f'<circle cx="{center:.0f}" cy="{center:.0f}" r="{size*.37:.0f}" fill="none" stroke="{COLORS["rule"]}" stroke-width="{size*.008:.0f}"/>'
        f'{mark(int(size*.29), int(size*.29), int(size*.42))}'
    )


def profile_assets() -> None:
    for platform, size in (("facebook", 1080), ("linkedin", 400), ("instagram", 1080), ("youtube", 800)):
        write_svg(f"{platform}/profile.svg", size, size, avatar(size), f"ORDERLY {platform.title()} profile image")
    write_svg("instagram/highlight-template.svg", 1080, 1080, avatar(1080), "ORDERLY Instagram highlight cover template")
    write_svg("logo/favicon/orderly-favicon.svg", 64, 64, mark(3, 3, 58), "ORDERLY favicon")
    app = f'<rect width="512" height="512" rx="116" fill="{COLORS["violet"]}"/>{mark(102, 102, 308, "white", "lavender")}'
    write_svg("logo/favicon/orderly-app-icon.svg", 512, 512, app, "ORDERLY app icon")


def social_assets() -> None:
    fb = background("facebook-bg.png", 1640, 624)
    fb += mark(1080, 207, 74) + wordmark(1170, 266, 44)
    fb += arabic(1550, 380, 48, "شغل مطعمك في مكان واحد")
    fb += arabic(1550, 435, 28, "من قدّام الضيوف لورا الكواليس", "muted", weight=500, family="IBM Plex Sans Arabic")
    write_svg("facebook/cover.svg", 1640, 624, fb, "Facebook cover", "One platform for connected front- and back-of-house restaurant work.")

    li = background("linkedin-bg.png", 1128, 191)
    li += f'<rect x="448" y="30" width="300" height="132" rx="22" fill="{COLORS["paper"]}" fill-opacity=".94" stroke="{COLORS["rule"]}"/>'
    li += mark(470, 59, 68) + wordmark(550, 108, 32)
    li += f'<path d="M750 42v108" stroke="{COLORS["rule"]}" stroke-width="2"/>'
    li += arabic(1090, 82, 21, "من قدّام الضيوف لورا الكواليس")
    li += arabic(1090, 120, 17, "منصة واحدة لإدارة شغل المطعم", "muted", weight=500, family="IBM Plex Sans Arabic")
    write_svg("linkedin/cover.svg", 1128, 191, li, "LinkedIn company cover", "Dedicated B2B cover with a concise Arabic positioning line.")

    post = background("instagram-post-bg.png", 1080, 1350)
    post += mark(482, 112, 116) + wordmark(432, 304, 46, anchor="middle")
    post += arabic(540, 590, 68, "من قدّام الضيوف", anchor="middle")
    post += arabic(540, 680, 68, "لورا الكواليس", anchor="middle")
    post += arabic(540, 765, 32, "شغل مطعمك في مكان واحد", "muted", anchor="middle", weight=500, family="IBM Plex Sans Arabic")
    post += f'<rect x="340" y="912" width="400" height="3" rx="2" fill="{COLORS["lavender"]}"/>'
    post += f'<text x="540" y="986" text-anchor="middle" fill="{COLORS["deep"]}" font-family="Manrope" font-size="24" font-weight="700">Restaurant management</text>'
    write_svg("instagram/launch-post.svg", 1080, 1350, post, "Instagram launch post", "ORDERLY restaurant management positioning in Egyptian Arabic.")

    story = background("instagram-story-bg.png", 1080, 1920)
    story += mark(482, 350, 116) + wordmark(432, 542, 46, anchor="middle")
    story += arabic(540, 780, 62, "مطعمك شغله متصل", anchor="middle")
    story += arabic(540, 870, 54, "من قدّام الضيوف", anchor="middle")
    story += arabic(540, 940, 54, "لورا الكواليس", anchor="middle")
    story += f'<rect x="335" y="1045" width="410" height="72" rx="36" fill="{COLORS["violet"]}"/>'
    story += f'<text x="540" y="1092" text-anchor="middle" fill="{COLORS["white"]}" font-family="Manrope" font-size="26" font-weight="700">ORDERLY</text>'
    write_svg("instagram/story.svg", 1080, 1920, story, "Instagram Story introduction", "Brand introduction with concise Egyptian Arabic copy.")

    banner = background("youtube-banner-bg.png", 2560, 1440)
    banner += mark(845, 625, 150) + wordmark(1025, 722, 96)
    banner += arabic(1280, 824, 48, "منصة واحدة لشغل المطعم", family="IBM Plex Sans Arabic")
    banner += f'<path d="M920 876h720" stroke="{COLORS["rule"]}" stroke-width="3"/>'
    banner += arabic(1280, 918, 30, "من قدّام الضيوف لورا الكواليس", "muted", weight=500, family="IBM Plex Sans Arabic")
    write_svg("youtube/banner.svg", 2560, 1440, banner, "YouTube channel banner", "Logo and messaging are inside the central 1546 by 423 safe area.")

    thumb = background("youtube-thumbnail-bg.png", 1280, 720)
    thumb += mark(1080, 50, 110) + wordmark(1110, 197, 34, anchor="middle")
    thumb += f'<rect x="64" y="70" width="116" height="38" rx="19" fill="{COLORS["purple"]}"/>'
    thumb += f'<text x="122" y="96" text-anchor="middle" fill="{COLORS["white"]}" font-family="IBM Plex Sans Arabic" font-size="18" font-weight="600" direction="rtl">دليل عملي</text>'
    thumb += arabic(770, 324, 58, "إدارة مطعمك", "white")
    thumb += arabic(770, 408, 58, "في مكان واحد", "white")
    thumb += arabic(1210, 470, 24, "عنوان قابل للتعديل لكل فيديو", "ink", weight=500, family="IBM Plex Sans Arabic")
    write_svg("youtube/thumbnail-template.svg", 1280, 720, thumb, "YouTube thumbnail template", "Reusable title layout; title text is editable in the SVG source.")


def previews() -> None:
    swatches = []
    for i, (name, color) in enumerate(COLORS.items()):
        x, y = 64 + (i % 5) * 220, 430 + (i // 5) * 130
        swatches.append(f'<rect x="{x}" y="{y}" width="180" height="72" rx="18" fill="{color}" stroke="{COLORS["rule"]}"/>')
        swatches.append(f'<text x="{x}" y="{y+102}" fill="{COLORS["ink"]}" font-family="Manrope" font-size="17">{name} · {color}</text>')
    body = f'<rect width="1200" height="720" fill="{COLORS["paper"]}"/>'
    body += mark(64, 56, 116) + wordmark(206, 143, 68)
    body += f'<text x="65" y="236" fill="{COLORS["muted"]}" font-family="Manrope" font-size="24">Restaurant management · Arabic first · front and back of house</text>'
    body += "".join(swatches) + arabic(1135, 670, 32, "من قدّام الضيوف لورا الكواليس", family="IBM Plex Sans Arabic")
    write_svg("previews/brand-overview.svg", 1200, 720, body, "ORDERLY brand overview preview")

    fb_uri = base64.b64encode((ROOT / "facebook" / "cover.svg").read_bytes()).decode("ascii")
    post_uri = base64.b64encode((ROOT / "instagram" / "launch-post.svg").read_bytes()).decode("ascii")
    story_uri = base64.b64encode((ROOT / "instagram" / "story.svg").read_bytes()).decode("ascii")
    body = f'<rect width="1600" height="1000" fill="{COLORS["surface"]}"/>'
    body += f'<rect x="60" y="54" width="1480" height="310" rx="28" fill="{COLORS["white"]}" stroke="{COLORS["rule"]}"/>'
    body += f'<rect x="60" y="396" width="720" height="542" rx="28" fill="{COLORS["white"]}" stroke="{COLORS["rule"]}"/>'
    body += f'<rect x="820" y="396" width="720" height="542" rx="28" fill="{COLORS["ink"]}"/>'
    body += f'<image x="78" y="72" width="1444" height="274" preserveAspectRatio="xMidYMid slice" href="data:image/svg+xml;base64,{fb_uri}"/>'
    body += f'<image x="92" y="426" width="314" height="512" preserveAspectRatio="xMidYMid slice" href="data:image/svg+xml;base64,{post_uri}"/>'
    body += f'<image x="842" y="426" width="314" height="512" preserveAspectRatio="xMidYMid slice" href="data:image/svg+xml;base64,{story_uri}"/>'
    body += f'<text x="440" y="490" fill="{COLORS["ink"]}" font-family="Manrope" font-size="28" font-weight="700">Instagram launch</text>'
    body += f'<text x="1190" y="490" fill="{COLORS["white"]}" font-family="Manrope" font-size="28" font-weight="700">Instagram Story</text>'
    write_svg("previews/social-media-mockups.svg", 1600, 1000, body, "Social media mockups preview")


def main() -> None:
    logo_assets()
    profile_assets()
    social_assets()
    previews()
    print(f"Wrote editable SVG assets under {ROOT}")


if __name__ == "__main__":
    main()
