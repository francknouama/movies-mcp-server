#!/bin/bash

echo "=== Checking for missing step definitions ==="
echo

# Extract all unique step patterns from feature files
echo "Extracting steps from feature files..."
grep -hE "^\s*(Given|When|Then|And)" features/*.feature | \
  sed 's/^\s*\(Given\|When\|Then\|And\)\s*//' | \
  sed 's/"[^"]*"/"{param}"/g' | \
  sed 's/\d+/{number}/g' | \
  sort | uniq > /tmp/feature_steps.txt

# Extract implemented step patterns from Go files
echo "Extracting implemented steps..."
grep -h 'ctx.Step(' steps/*.go | \
  sed 's/.*ctx.Step(`\^\(.*\)\$`.*/\1/' | \
  sed 's/\([^\\]\)\$/\1/g' | \
  sed 's/\\d+/{number}/g' | \
  sed 's/([^)]*)/{param}/g' | \
  sed 's/"[^"]*"/{param}/g' | \
  sort | uniq > /tmp/implemented_steps.txt

echo
echo "=== Steps in features but not implemented ==="
while IFS= read -r step; do
  # Clean up the step pattern for matching
  clean_step=$(echo "$step" | sed 's/[.*+?{}()|[\]\\]/\\&/g')
  if ! grep -qF "$clean_step" /tmp/implemented_steps.txt; then
    echo "- $step"
  fi
done < /tmp/feature_steps.txt

echo
echo "=== Summary ==="
echo "Total steps in features: $(wc -l < /tmp/feature_steps.txt)"
echo "Total implemented steps: $(wc -l < /tmp/implemented_steps.txt)"