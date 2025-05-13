#!/usr/local/bin/bash

ALL="$(doctl compute snapshot list --format CreatedAt,ID,Name)"
ALL_P="$(echo "$ALL" | tail -n +2 | gsed -E -e 's/[\ \t]+/ /g' -e 's/_/ /g')"
KINDS="$(echo "$ALL_P" | cut -d ' ' -f 4 | sort | uniq)"
UPTODATE_IDs="$(echo "$KINDS" | while read KIND; do echo "$ALL_P" | grep "$KIND" | tail -n 1 | cut -d ' ' -f 2; done)"
echo "$ALL" | while read IMAGE; do if echo "$IMAGE" | grep "$UPTODATE_IDs" >/dev/null 2>&1; then echo "$IMAGE" | green; else echo "$IMAGE"; fi; done
