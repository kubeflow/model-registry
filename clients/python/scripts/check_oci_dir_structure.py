#!/usr/bin/env python3
"""Test that save_to_oci_registry preserves nested directory structure

Ref: test_save_to_oci_registry_with_skopeo in tests/test_utils.py:86

Usage:

  python scripts/check_oci_dir_structure.py

  # or

  poetry run python scripts/check_oci_dir_structure.py
"""

import json
import subprocess
import sys
import tarfile
import tempfile
from pathlib import Path

from model_registry.utils import save_to_oci_registry

OCI_REF = "quay.io/jburdo/test1:dir-structure-test"
BASE_IMAGE = "quay.io/mmortari/hello-world-wait:latest"

EXPECTED_PATHS = {
    "README.md",
    "onnx/model.onnx",
    "onnx/weights/quantized.bin",
    "tokenizer/vocab.txt",
    "tokenizer/config.json",
}

with tempfile.TemporaryDirectory() as tmpdir:
    tmp = Path(tmpdir)
    model_dir = tmp / "my-model"
    model_dir.mkdir()

    # Create nested model structure
    (model_dir / "README.md").write_text("# Test Model\n")
    (model_dir / "onnx").mkdir()
    (model_dir / "onnx" / "model.onnx").write_bytes(b"\x00" * 64)
    (model_dir / "onnx" / "weights").mkdir()
    (model_dir / "onnx" / "weights" / "quantized.bin").write_bytes(b"\x01" * 32)
    (model_dir / "tokenizer").mkdir()
    (model_dir / "tokenizer" / "vocab.txt").write_text("hello\nworld\n")
    (model_dir / "tokenizer" / "config.json").write_text('{"type": "test"}\n')

    # Push
    print(f"Pushing to {OCI_REF}...")
    save_to_oci_registry(
        base_image=BASE_IMAGE,
        oci_ref=OCI_REF,
        model_files_path=model_dir,
        dest_dir=tmp / "oci-dest",
        backend="skopeo",
    )

    # Pull back
    print("Pulling back...")
    pull_dir = tmp / "oci-pull"
    pull_dir.mkdir()
    subprocess.run(
        ["skopeo", "copy", f"docker://{OCI_REF}", f"oci:{pull_dir}:latest"],
        check=True,
        capture_output=True,
    )

    # Collect all file paths from tar layers (skip base image layer 0)
    index = json.loads((pull_dir / "index.json").read_text())
    algo, digest = index["manifests"][0]["digest"].split(":")
    manifest = json.loads((pull_dir / "blobs" / algo / digest).read_text())

    found_paths = set()
    layer_num = 0
    print("\nModel layers:")
    for layer in manifest["layers"]:
        algo, digest = layer["digest"].split(":")
        blob = pull_dir / "blobs" / algo / digest
        try:
            with tarfile.open(blob, "r:*") as tar:
                model_entries = [
                    m for m in tar.getmembers() if m.name.startswith("models/")
                ]
                if not model_entries:
                    continue
                print(f"  layer {layer_num} ({digest[:12]})")
                for m in sorted(model_entries, key=lambda m: m.name):
                    path = m.name.removeprefix("models/")
                    kind = "dir" if m.isdir() else "file"
                    print(f"    {'└── ' if m == model_entries[-1] else '├── '}{path}  ({kind})")
                    found_paths.add(path)
                layer_num += 1
        except tarfile.TarError:
            continue

    # Report
    print(f"\nExpected: {sorted(EXPECTED_PATHS)}")
    print(f"Found:    {sorted(found_paths)}")
    missing = EXPECTED_PATHS - found_paths
    if missing:
        print(f"FAIL - missing paths: {sorted(missing)}")
        sys.exit(1)
    else:
        print("PASS - all expected paths preserved")
