"""Render every editable SVG master to a dimension-matched PNG with Edge."""
from __future__ import annotations

from pathlib import Path

from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[1]
EDGE = Path(r"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe")


def render_svg(page, source: Path) -> None:
    page.goto(source.as_uri(), wait_until="load")
    page.locator("svg").wait_for()
    page.evaluate("document.fonts.ready")
    size = page.locator("svg").evaluate(
        "svg => ({width: Number(svg.getAttribute('width')), height: Number(svg.getAttribute('height'))})"
    )
    output = source.with_suffix(".png")
    page.locator("svg").screenshot(path=str(output), animations="disabled", omit_background=True)
    from PIL import Image

    with Image.open(output) as image:
        if image.size != (size["width"], size["height"]):
            raise ValueError(f"{source}: expected {size}, got {image.size}")
    print(f"{output.relative_to(ROOT)} {size['width']}x{size['height']}")


def main() -> None:
    sources = sorted(path for path in ROOT.rglob("*.svg") if "assets/generated-bg" not in str(path))
    with sync_playwright() as playwright:
        browser = playwright.chromium.launch(headless=True, executable_path=str(EDGE))
        page = browser.new_page(viewport={"width": 2560, "height": 1440}, device_scale_factor=1)
        for source in sources:
            render_svg(page, source)
        browser.close()


if __name__ == "__main__":
    main()
