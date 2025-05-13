#!/usr/local/bin/bash

find image -name 'dr.yml' | while read FOLDER; do
  CHILD="$(basename "$(dirname "$FOLDER")")"
  PARENT="$(basename "$(cat "$FOLDER" | yq -r '.depends_on.folder')")"
  echo "$PARENT -> $CHILD"
done | sort
