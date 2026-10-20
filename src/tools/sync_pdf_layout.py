# /// script
# requires-python = ">=3.10"
# dependencies = ["pymupdf>=1.26,<2", "pillow>=12,<13"]
# ///
"""Export PDF artwork and positioned, selectable text without exporting fonts.

Run with: uv run src/tools/sync_pdf_layout.py
The PDF is read only. Rendered backdrops contain no text; all text is HTML.
"""

import argparse
import json
from pathlib import Path

import pymupdf
from PIL import Image


def text_runs(page, fallback_font):
    runs = []
    # Do not decode the PDF's large images just to obtain text positions.
    flags = pymupdf.TEXTFLAGS_RAWDICT & ~pymupdf.TEXT_PRESERVE_IMAGES
    for block in page.get_text("rawdict", flags=flags)["blocks"]:
        for line in block.get("lines", []):
            for span in line["spans"]:
                start = len(runs)
                pixelion = "Pixelion" in span["font"]
                span_text = "".join(char["c"] for char in span["chars"])
                href = ""
                if "tacyborg@outlook.com" in span_text:
                    href = "mailto:tacyborg@outlook.com"
                elif "+57 301" in span_text:
                    href = "tel:+573016115264"
                elif "@tacyborg" in span_text:
                    href = "https://www.instagram.com/tacyborg/"
                # Keep headings together so the substitute can be fitted to the
                # source width. Position body words separately to retain the
                # PDF's justification, instead of guessing CSS word spacing.
                groups = []
                group = []
                for char in span["chars"]:
                    group.append(char)
                    if char["c"].isspace() and not pixelion:
                        groups.append(group)
                        group = []
                if group:
                    groups.append(group)
                for group in groups:
                    text = "".join(char["c"] for char in group)
                    if not text.strip():
                        continue
                    x, baseline = group[0]["origin"]
                    right = max(char["bbox"][2] for char in group)
                    fit = 1
                    if pixelion:
                        natural_width = fallback_font.text_length(text, fontsize=span["size"] * .768)
                        if natural_width:
                            fit = (right - x) / natural_width
                    runs.append({
                        "Text": text,
                        "Font": "pixelion" if pixelion else "questrial",
                        "Left": round(x / page.rect.width * 100, 6),
                        "Baseline": round(baseline / page.rect.height * 100, 6),
                        "Size": round(span["size"] / page.rect.width * 100, 6),
                        "Width": round((right - x) / page.rect.width * 100, 6),
                        "Fit": round(fit, 6),
                        "Color": f"#{span['color']:06x}",
                        "Href": href,
                    })
                if len(runs) > start:
                    # Keep selection/copy boundaries without introducing a
                    # second visual line in each absolutely positioned span.
                    if not runs[-1]["Text"].endswith(" "):
                        runs[-1]["Text"] += " "
    return runs


def main():
    source = Path(__file__).resolve().parents[1]
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--pdf", type=Path, default=source.parent / "portfolio.pdf")
    args = parser.parse_args()
    if not args.pdf.is_file():
        parser.error(f"PDF not found: {args.pdf}")
    content = json.loads((source / "cmd/api/portfolio.json").read_text())
    backdrops = source / "static/portfolio/pages"
    backdrops.mkdir(exist_ok=True)
    layout = []
    fallback_font = pymupdf.Font(fontfile=str(source / "static/fonts/Silkscreen-Regular.ttf"))
    with pymupdf.open(args.pdf) as pdf:
        if len(pdf) != len(content):
            parser.error("The PDF and portfolio content must have the same page count")
        for index, page in enumerate(pdf):
            runs = text_runs(page, fallback_font)
            # Remove text ONLY. Retain raster artwork, gradients, and all vector
            # ornaments, even when they overlap a text bounding box.
            flags = pymupdf.TEXTFLAGS_DICT & ~pymupdf.TEXT_PRESERVE_IMAGES
            for block in page.get_text("dict", flags=flags)["blocks"]:
                for line in block.get("lines", []):
                    for span in line["spans"]:
                        page.add_redact_annot(span["bbox"], fill=False, cross_out=False)
            page.apply_redactions(images=0, graphics=0, text=0)
            if page.get_text().strip():
                raise RuntimeError(f"Page {index + 1} still contains text")
            scale = 1920 / page.rect.width
            pixmap = page.get_pixmap(matrix=pymupdf.Matrix(scale, scale), alpha=False)
            image = Image.frombytes("RGB", (pixmap.width, pixmap.height), pixmap.samples)
            name = f"page-{index + 1:02d}.webp"
            image.save(backdrops / name, "WEBP", quality=90, method=6)
            layout.append({
                "ID": content[index]["ID"],
                "Background": name,
                "Width": page.rect.width,
                "Height": page.rect.height,
                "Runs": runs,
            })
            print(f"Page {index + 1:02d}: {len(runs)} positioned text runs, text-free backdrop")
    target = source / "cmd/api/portfolio-layout.json"
    target.write_text(json.dumps(layout, ensure_ascii=False, indent=2) + "\n")
    print(f"Wrote {target}")


if __name__ == "__main__":
    main()
