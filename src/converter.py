import argparse
from pdf2image import convert_from_path
import os

def split_pdf_to_images(pdf_path, output_dir, quality=80):
    """Splits a PDF into individual images, saving them in the specified directory."""
    os.makedirs(output_dir, exist_ok=True)
    images = convert_from_path(pdf_path)

    for i, page in enumerate(images):
        output_path = os.path.join(output_dir, f"page_{i + 1}.webp")
        page.save(output_path, "WEBP", quality=quality)
        print(f"Saved: {output_path} with quality {quality}")

    print(f"All pages have been saved in '{output_dir}'")

def main():
    parser = argparse.ArgumentParser(description="Split a PDF into individual page images.")
    parser.add_argument("pdf_path", type=str, help="Path to the PDF file.")
    parser.add_argument("output_dir", type=str, help="Directory to save the images.")
    parser.add_argument(
        "--quality",
        type=int,
        default=80,
        help="Quality of the output WebP images (default: 80).",
    )

    args = parser.parse_args()
    split_pdf_to_images(args.pdf_path, args.output_dir, args.quality)

if __name__ == "__main__":
    main()
