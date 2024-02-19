"""Sphinx configuration."""

project = "model-registry"
author = "Isabella Basso do Amaral"
copyright = f"2023, {author}"
html_theme = "furo"
extensions = [
    "sphinx.ext.autodoc",  # support for automatic documentation
    "sphinx.ext.napoleon",  # support for NumPy and Google style docstrings
    "sphinx.ext.todo",  # support for TODOs on docstrings
    "sphinx.ext.viewcode",  # Add links to highlighted source code
    "myst_parser",  # Markdown support
]
autodoc_typehints = "signature"
napoleon_google_docstring = True
napoleon_use_param = True
napoleon_use_rtype = True

autodoc_default_options = {
    "members": True,
    "show-inheritance": True,
    "inherited-members": True,
    "member-order": "bysource",
}

# https://myst-parser.readthedocs.io/en/latest/syntax/optional.html
myst_enable_extensions = [
    "strikethrough",  # enable ~~strikethrough~~
    "smartquotes",  # converts quotes to “ ” ‘ ’
    "replacements",  # converts symbols like (c), (tm), (r), --, ---, +-
    "colon_fence",  # use ::: to delimit a myst block
    "linkify",  # automatically create links for URLs
    "deflist",  # enables less verbose definition lists
]
myst_heading_anchors = 3  # automatically create anchors for headings up to h3 (###)
myst_linkify_fuzzy_links = False  # only linkify schema-prefixed URLs (e.g. https://*)
