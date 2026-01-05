"""Nox sessions for Model Catalog Python Client."""

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


python_versions = ["3.12", "3.11", "3.10"]
nox.needs_version = ">= 2021.6.6"
nox.options.sessions = (
    "lint",
    "e2e",
)


@session(python=python_versions)
def lint(session: Session) -> None:
    """Lint using ruff."""
    session.install("ruff")
    session.run("ruff", "check", ".")


@session(python=python_versions)
def mypy(session: Session) -> None:
    """Type check using mypy."""
    session.install(
        ".",
        "mypy",
        "types-python-dateutil",
        "types-pyyaml",
        "types-requests",
    )
    session.run("mypy", ".")


@session(name="e2e", python=python_versions)
def e2e_tests(session: Session) -> None:
    """Run E2E tests (requires running catalog service)."""
    session.install(
        ".",
        "pytest",
        "pytest-asyncio",
        "pytest-mock",
        "pytest-timeout",
        "pytest-xdist",
        "coverage[toml]",
        "pytest-cov",
    )
    try:
        session.run(
            "pytest",
            "tests/",
            "--e2e",
            "-v",
            "-rA",
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
