#!/bin/bash

# Test script for pdf-compressor service
# This script sends a PDF file to the compression service and saves the compressed output

set -e  # Exit on error

# Configuration
SERVICE_URL="http://localhost:8080/compress"
INPUT_PDF="testing/dummy.pdf"
OUTPUT_PDF="testing/compressed.pdf"

# Check if input file exists
if [ ! -f "$INPUT_PDF" ]; then
    echo "Error: Input file $INPUT_PDF not found"
    exit 1
fi

echo "Testing PDF compression service..."
echo "Input: $INPUT_PDF ($(du -h "$INPUT_PDF" | cut -f1))"

# Send request to compress PDF
curl -s -X POST \
    -F "pdf=@$INPUT_PDF" \
    -o "$OUTPUT_PDF" \
    "$SERVICE_URL"

# Check if compression was successful
if [ -f "$OUTPUT_PDF" ]; then
    echo "Compression successful!"
    echo "Output: $OUTPUT_PDF ($(du -h "$OUTPUT_PDF" | cut -f1))"
    echo "Compression ratio: $(echo "scale=2; $(du -b "$INPUT_PDF" | cut -f1) / $(du -b "$OUTPUT_PDF" | cut -f1)" | bc)"
else
    echo "Error: Compression failed"
    exit 1
fi

echo "Test completed successfully"