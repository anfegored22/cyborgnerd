import argparse
from pdf2image import convert_from_path
import os

def split_pdf_to_images(pdf_path, output_dir):
    """Splits a PDF into individual images, saving them in the specified directory."""
    os.makedirs(output_dir, exist_ok=True)
    images = convert_from_path(pdf_path)

    for i, page in enumerate(images):
        output_path = os.path.join(output_dir, f"page_{i + 1}.png")
        page.save(output_path, "PNG")
        print(f"Saved: {output_path}")

    print(f"All pages have been saved in '{output_dir}'")

def main():
    parser = argparse.ArgumentParser(description="Split a PDF into individual page images.")
    parser.add_argument("pdf_path", type=str, help="Path to the PDF file.")
    parser.add_argument("output_dir", type=str, help="Directory to save the images.")

    args = parser.parse_args()
    split_pdf_to_images(args.pdf_path, args.output_dir)

if __name__ == "__main__":
    main()
