#!/bin/bash
set -e

# Create _posts/ if it doesn't exist
mkdir -p _posts

# Inject baseurl into all assets
BASEURL="{{ site.baseurl }}"
for filepath in _posts_md/*.md; do
  filename=$(basename "${filepath}")
  sed -E "s#(!\[.*\]\()/assets/#\1${BASEURL}/assets/#g" "${filepath}" > "_posts/${filename}"
done

# Replace mermaid code blocks with HTML
for filepath in _posts/*.md; do
  sed -i -E '/```mermaid/,/```/{
    s/```mermaid/<div class="mermaid">/
    s/```/<\/div>/
  }' "${filepath}"
done

echo "Preprocessed markdown files"
