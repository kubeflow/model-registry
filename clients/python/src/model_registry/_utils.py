from __future__ import annotations

import functools
import inspect
from collections.abc import Sequence
from typing import Any, Callable, TypeVar

CallableT = TypeVar("CallableT", bound=Callable[..., Any])


# copied from https://github.com/Rapptz/RoboDanny
def human_join(seq: Sequence[str], *, delim: str = ", ", final: str = "or") -> str:
    size = len(seq)
    if size == 0:
        return ""

    if size == 1:
        return seq[0]

    if size == 2:
        return f"{seq[0]} {final} {seq[1]}"

    return delim.join(seq[:-1]) + f" {final} {seq[-1]}"


def quote(string: str) -> str:
    """Add single quotation marks around the given string. Does *not* do any escaping."""
    return f"'{string}'"


# copied from https://github.com/openai/openai-python
def required_args(*variants: Sequence[str]) -> Callable[[CallableT], CallableT]:  # noqa: C901
    """Decorator to enforce a given set of arguments or variants of arguments are passed to the decorated function.

    Useful for enforcing runtime validation of overloaded functions.

    Example usage:
    ```py
    @overload
    def foo(*, a: str) -> str:
        ...


    @overload
    def foo(*, b: bool) -> str:
        ...


    # This enforces the same constraints that a static type checker would
    # i.e. that either a or b must be passed to the function
    @required_args(["a"], ["b"])
    def foo(*, a: str | None = None, b: bool | None = None) -> str:
        ...
    ```
    """

    def inner(func: CallableT) -> CallableT:  # noqa: C901
        params = inspect.signature(func).parameters
        positional = [
            name
            for name, param in params.items()
            if param.kind
            in {
                param.POSITIONAL_ONLY,
                param.POSITIONAL_OR_KEYWORD,
            }
        ]

        @functools.wraps(func)
        def wrapper(*args: object, **kwargs: object) -> object:
            given_params: set[str] = set()
            for i, _ in enumerate(args):
                try:
                    given_params.add(positional[i])
                except IndexError:
                    msg = f"{func.__name__}() takes {len(positional)} argument(s) but {len(args)} were given"
                    raise TypeError(msg) from None

            for key in kwargs:
                given_params.add(key)

            for variant in variants:
                matches = all(param in given_params for param in variant)
                if matches:
                    break
            else:  # no break
                if len(variants) > 1:
                    variations = human_join(
                        [
                            "("
                            + human_join([quote(arg) for arg in variant], final="and")
                            + ")"
                            for variant in variants
                        ]
                    )
                    msg = f"Missing required arguments; Expected either {variations} arguments to be given"
                else:
                    # TODO: this error message is not deterministic
                    missing = list(set(variants[0]) - given_params)
                    if len(missing) > 1:
                        msg = f"Missing required arguments: {human_join([quote(arg) for arg in missing])}"
                    else:
                        msg = f"Missing required argument: {quote(missing[0])}"
                raise TypeError(msg)
            return func(*args, **kwargs)

        return wrapper  # type: ignore

    return inner
