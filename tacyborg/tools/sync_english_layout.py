# /// script
# requires-python = ">=3.10"
# dependencies = ["pymupdf>=1.26,<2"]
# ///
"""Fit English HTML text into the PDF's text regions, reusing text-free artwork.

uv run tacyborg/tools/sync_english_layout.py
Requires the private source PDF and the purchased font archive. No fonts or
artwork are extracted or written by this script; only text geometry is generated.
"""

import argparse
import json
import statistics
from pathlib import Path
from zipfile import ZipFile

import pymupdf


def groups_on(page):
    groups = []
    flags = pymupdf.TEXTFLAGS_DICT & ~pymupdf.TEXT_PRESERVE_IMAGES
    for block in page.get_text("dict", flags=flags)["blocks"]:
        current = None
        for line in block.get("lines", []):
            for span in line["spans"]:
                style = (span["font"], round(span["size"], 1))
                x, baseline = span["origin"]
                # Split separate columns even when the PDF groups them together.
                new_column = current and baseline > current["baselines"][-1] + 1 and abs(x - current["x"]) > span["size"] * 2
                if current is None or current["style"] != style or new_column:
                    current = {"style": style, "size": span["size"], "x": x, "right": span["bbox"][2], "baselines": [baseline]}
                    groups.append(current)
                else:
                    current["right"] = max(current["right"], span["bbox"][2])
                    if baseline > current["baselines"][-1] + 1:
                        current["baselines"].append(baseline)
    return groups


def wrap(text, font, size, width):
    lines = []
    for paragraph in text.split("\n"):
        line = []
        for word in paragraph.split():
            if line and font.text_length(" ".join([*line, word]), fontsize=size) > width + .5:
                lines.append(line)
                line = []
            line.append(word)
        lines.append(line)
    return lines


def main():
    source = Path(__file__).resolve().parents[1]
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--pdf", type=Path, default=source.parent / "portfolio.pdf")
    parser.add_argument("--font-zip", type=Path, default=source / "font files.zip")
    args = parser.parse_args()
    text = {s["ID"]: s for s in json.loads((source / "cmd/api/portfolio.en.json").read_text())}
    layout = json.loads((source / "cmd/api/portfolio-layout.json").read_text())
    by_page = {int(Path(p["Background"]).stem.split("-")[1]): p for p in layout}

    def paragraphs(id):
        return "\n\n".join(text[id]["Paragraphs"])

    def column(id, n, indices=None):
        values = text[id]["Columns"][n]["Paragraphs"]
        return "\n\n".join(values if indices is None else [values[i] for i in indices])

    def caption(id, i, title):
        parts = text[id]["Images"][i]["Caption"].split(" · ")[-3:]
        return f'Title: “{title}”\nDimensions: {parts[0]}\nMedium: {parts[1]}\nYear: {parts[2]}'

    # One entry per source text region. Body copy is shared with the responsive
    # English view, so neither the artwork nor translations are duplicated here.
    regions = {
        "portafolio": ["PORTFOLIO", "Artistic Portfolio 2026", "TACY BORG"],
        "declaracion": ["STATEMENT", paragraphs("declaracion")],
        "semblanza": ["BIOGRAPHY", paragraphs("semblanza")],
        "investigacion": ["RESEARCH\nLINES", "Conceptual research: " + column("investigacion", 0), "Technical research: " + column("investigacion", 1)],
        "curriculum": ["CURRICULUM VITAE", column("curriculum", 0), "Profile", "Education", column("curriculum", 1, [0, 1]), column("curriculum", 1, [2, 3]), column("curriculum", 1, [4]), "Awards", column("curriculum", 2)],
        "exposiciones": [column("exposiciones", 0), "EXHIBITIONS &\nPUBLICATIONS", column("exposiciones", 1) + "\n\n" + column("exposiciones", 2, [0]), column("exposiciones", 2, [1, 2])],
        "experiencia": [*[column("experiencia", 0, [i]) for i in range(4)], "WORK\nEXPERIENCE", column("experiencia", 1, [0]), column("experiencia", 1, [1]), column("experiencia", 2, [0]), column("experiencia", 2, [1])],
        "proyectos": ["PROJECT\nPORTFOLIO", "“" + text["proyectos"]["Quote"] + "”", "Tacy Borg"],
        "mosaicos": [text["mosaicos"]["Title"], "FIRST\nPROJECT"],
        "chakana": [caption("chakana", 0, "Chakana")],
        "cxacxa": [text["cxacxa"]["Paragraphs"][0], text["cxacxa"]["Paragraphs"][1], caption("cxacxa", 0, "Cxacxa"), text["cxacxa"]["Paragraphs"][2]],
        "rupestres": [caption("rupestres", 0, "Plateado rupestre"), caption("rupestres", 1, "Verde rupestre")],
        "hischa": [paragraphs("hischa"), caption("hischa", 0, "Hischa 1")],
        "comunidad": [caption("comunidad", 0, "Comunidad, identidad y rupturas"), caption("comunidad", 1, "Arte deco")],
        "historia-humana": [text["historia-humana"]["Title"], "SECOND\nPROJECT"],
        "amados-monstruos": [text["amados-monstruos"]["Title"], "THIRD\nPROJECT"],
        "genealogia": [text["genealogia"]["Title"], "FOURTH\nPROJECT"],
        "contacto": ["CONTACT", "PHONE", "INSTAGRAM", "+57 301-611-5264", "@tacyborg", "EMAIL", "tacyborg@outlook.com", "Based in\nCajicá\nColombia"],
    }
    questrial = pymupdf.Font(fontfile=str(source / "static/fonts/Questrial-Regular.ttf"))
    fallback = pymupdf.Font(fontfile=str(source / "static/fonts/Silkscreen-Regular.ttf"))
    with ZipFile(args.font_zip) as archive:
        pixelion = pymupdf.Font(fontbuffer=archive.read("font files/TTF/Pixelion-Regular.ttf"))
    result = []
    with pymupdf.open(args.pdf) as pdf:
        for i, page in enumerate(pdf, 1):
            original = by_page[i]
            id = original["ID"]
            groups = groups_on(page)
            if len(groups) != len(regions[id]):
                raise ValueError(f"{id}: {len(groups)} source regions but {len(regions[id])} translations")
            runs = []
            minimum_scale = 1
            for n, (group, value) in enumerate(zip(groups, regions[id])):
                font = pixelion if "Pixelion" in group["style"][0] else questrial
                is_pixel = font is pixelion
                width = group["right"] - group["x"]
                baselines = group["baselines"]
                leading = statistics.median(b - a for a, b in zip(baselines, baselines[1:])) if len(baselines) > 1 else group["size"] * 1.3
                scale = 1.0
                while True:
                    size = group["size"] * scale
                    lines = wrap(value, font, size, width)
                    fits = len(lines) * leading * scale <= len(baselines) * leading + .5
                    fits = fits and all(font.text_length(" ".join(line), fontsize=size) <= width + .5 for line in lines)
                    if fits:
                        break
                    scale -= .025
                    if scale < .5:
                        raise ValueError(f"English text does not fit: {id}, region {n}")
                minimum_scale = min(minimum_scale, scale)
                href = "mailto:tacyborg@outlook.com" if "tacyborg@outlook.com" in value else "tel:+573016115264" if value.startswith("+57") else "https://www.instagram.com/tacyborg/" if value == "@tacyborg" else ""
                for line_no, words in enumerate(lines):
                    x = group["x"]
                    line_width = font.text_length(" ".join(words), fontsize=size)
                    if id == "portafolio" and n in (0, 2):
                        x += (width - line_width) / 2
                    elif id == "portafolio" and n == 1:
                        x += width - line_width
                    # Keep a heading as a single run for the fitted fallback.
                    pieces = [" ".join(words)] if is_pixel else words
                    for word in pieces:
                        if not word:
                            continue
                        advance = font.text_length(word, fontsize=size)
                        fit = advance / fallback.text_length(word, fontsize=size * .768) if is_pixel else 1
                        runs.append({"Text": word + " ", "Font": "pixelion" if is_pixel else "questrial", "Left": round(x / 1440 * 100, 6), "Baseline": round((baselines[0] + line_no * leading * scale) / 810 * 100, 6), "Size": round(size / 1440 * 100, 6), "Width": round(advance / 1440 * 100, 6), "Fit": round(fit, 6), "Color": "#000000", "Href": href})
                        x += advance + font.text_length(" ", fontsize=size)
            result.append({**original, "Runs": runs})
            print(f"{id}: translated {len(groups)} regions; minimum text scale {minimum_scale:.0%}")
    output = source / "cmd/api/portfolio-layout.en.json"
    output.write_text(json.dumps(result, ensure_ascii=False, indent=2) + "\n")


if __name__ == "__main__":
    main()
