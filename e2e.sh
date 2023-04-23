#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

# Absolute path to the latest released version
current_bin="/usr/local/bin/evernote2md"

# Absolute path to the new release candidate
candidate_bin="$(pwd)/bin/evernote2md"

# Sample Evernote export file
input_file=".trash/My\ Notes.enex"

rm -rf "notes1" "notes2"

$current_bin "$input_file" "notes1"

$candidate_bin "$input_file" "notes2"

if ! diff -qr "notes1" "notes2"; then
    echo "End-to-end test failed"
    exit 1
else
    rm -r "notes1" "notes2"
fi
