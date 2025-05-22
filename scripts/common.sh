#!/usr/bin/env bash

realpath_fallback() {
  if command -v realpath >/dev/null 2>&1; then
    realpath "$@"
  elif command -v grealpath >/dev/null 2>&1; then
    grealpath "$@"
  else
    python3 -c "import os,sys; print(os.path.realpath(sys.argv[1]))" "$1"
  fi
} 