.PHONY: test
test:
    poetry run pytest -s -rA

test-with-cov:
    poetry run pytest --cov=job --cov-report=term-missing tests/
