"""Regenerate web images and text from the extracted portfolio assets.

Requires ImageMagick (magick) and Poppler (pdftotext), not Python packages.
Run from the repository root: python3 tacyborg/tools/prepare_portfolio.py
"""

import argparse
import json
from pathlib import Path
import shutil
import subprocess


def main():
    source = Path(__file__).resolve().parents[1]
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--pdf", type=Path, default=source.parent / "portfolio.pdf")
    parser.add_argument("--originals", type=Path, default=source.parent / "source-assets/portfolio")
    parser.add_argument("--additional-originals", type=Path, default=source.parent / "source-assets/genealogia")
    args = parser.parse_args()
    for command in ("magick", "pdftotext"):
        if not shutil.which(command):
            parser.error(f"Required command is not installed: {command}")
    if not args.pdf.is_file():
        parser.error(f"PDF not found: {args.pdf}")

    assets = source / "static" / "portfolio"
    manifest = json.loads((assets / "manifest.json").read_text())
    content = json.loads((source / "cmd/api/portfolio.json").read_text())
    originals = {Path(image["file"]).stem: image["file"] for image in manifest["images"]}
    needed = dict.fromkeys(image["File"] for slide in content for image in slide.get("Images", []))
    for stem in needed:
        original = (args.originals / originals[stem]) if stem in originals else (args.additional_originals / (stem + ".jpg"))
        if not original.is_file():
            parser.error(f"Original image missing for {stem}; restore source-assets or use --originals PATH")
        needed[stem] = original
    web = assets / "web"
    web.mkdir(exist_ok=True)
    for stem, original in needed.items():
        target = web / (stem + ".webp")
        subprocess.run([
            "magick", str(original), "-auto-orient", "-resize", "1600x1600>", "-strip",
            "-quality", "82", "-define", "webp:method=6", str(target),
        ], check=True)
    subprocess.run([
        "pdftotext", "-layout", str(args.pdf), str(assets / "text.txt"),
    ], check=True)
    print(f"Web images and extracted text saved to {assets}")


if __name__ == "__main__":
    main()
