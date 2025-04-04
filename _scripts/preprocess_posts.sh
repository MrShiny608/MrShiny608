#!/bin/bash
set -e

# Set the baseurl
BASEURL="{{ site.baseurl }}"

# Create _posts/ if it doesn't exist
mkdir -p _posts

# Process each Markdown file
for filepath in _posts_md/*.md; do
  filename=$(basename "${filepath}")
  sed -E "s#(!\[.*\]\()\/assets/#\1${BASEURL}/assets/#g" "${filepath}" > "_posts/${filename}"
done

echo "Preprocessed markdown with baseurl to _posts/"
