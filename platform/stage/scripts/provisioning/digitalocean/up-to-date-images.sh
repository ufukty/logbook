#!/usr/local/bin/bash

ALL="$(doctl compute snapshot list --format Name,ID --no-header | grep build | grep -v base | sort -r | gsed -E 's/[ ]+/ /g')"
echo "$ALL" | gsed -E 's/(build_([^_]+)_([0-9_]{17}).*)/\1 \2/g' | sort -r | uniq -f 2 | cut -d ' ' -f 1-2 | gsed -E 's/^([^\ ]+) (.*)/\2 \1/g'
