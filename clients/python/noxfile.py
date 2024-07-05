"""Nox sessions."""

import os
import shutil
import sys
from pathlib import Path
from textwrap import dedent

import nox

try:
    from nox_poetry import Session, session
except ImportError:
    message = f"""\
    Nox failed to import the 'nox-poetry' package.

    Please install it using the following command:

    {sys.executable} -m pip install nox-poetry"""
    raise SystemExit(dedent(message)) from None


package = "model_registry"
python_versions = ["3.12", "3.11", "3.10", "3.9"]
nox.needs_version = ">= 2021.6.6"
nox.options.sessions = (
    "tests",
    "docs-build",
)


@session(python=python_versions)
def lint(session: Session) -> None:
    """Lint using ruff."""
    session.install("ruff")

    # can't check the whole project because of the generated code
    session.run("ruff", "check", "src/model_registry", "tests")


@session(python=python_versions)
def mypy(session: Session) -> None:
    """Type check using mypy."""
    session.install(".")
    session.install("mypy")

    session.run("mypy", "src/model_registry")


@session(python=python_versions)
def tests(session: Session) -> None:
    """Run the test suite."""
    session.install(".")
    session.install(
        "coverage[toml]",
        "pytest",
        "pytest-asyncio",
        "pytest-cov",
        "pygments",
        "huggingface-hub",
    )
    try:
        session.run(
            "pytest",
            "--cov",
            "--cov-config=pyproject.toml",
            *session.posargs,
            env={"COVERAGE_FILE": f".coverage.{session.python}"},
        )
    finally:
        if session.interactive:
            session.notify("coverage", posargs=[])


@session(python=python_versions[0])
def coverage(session: Session) -> None:
    """Produce the coverage report."""
    args = session.posargs or ["report"]

    session.install("coverage[toml]")

    if not session.posargs and any(Path().glob(".coverage.*")):
        session.run("coverage", "combine")

    session.run("coverage", *args)


@session(name="docs-build", python=python_versions[0])
def docs_build(session: Session) -> None:
    """Build the documentation."""
    args = session.posargs or ["docs", "docs/_build"]
    if not session.posargs and "FORCE_COLOR" in os.environ:
        args.insert(0, "--color")

    session.install(".")
    session.install("sphinx", "myst-parser[linkify]", "furo")

    build_dir = Path("docs", "_build")
    if build_dir.exists():
        shutil.rmtree(build_dir)

    session.run("sphinx-build", *args)


@session(python=python_versions[0])
def docs(session: Session) -> None:
    """Build and serve the documentation with live reloading on file changes."""
    args = session.posargs or ["--open-browser", "docs", "docs/_build"]
    session.install(".")
    session.install("sphinx", "myst-parser[linkify]", "furo", "sphinx-autobuild")

    build_dir = Path("docs", "_build")
    if build_dir.exists():
        shutil.rmtree(build_dir)

    session.run("sphinx-autobuild", *args)
