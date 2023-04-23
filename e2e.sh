#!/usr/bin/env bash

# Script for end-to-end testing evernote2md
# It requires a path to the current version and the candidate version
#
# The script processes the same input file with 2 different evernote2md versions
# and fails if the results are not 100% equal.

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
