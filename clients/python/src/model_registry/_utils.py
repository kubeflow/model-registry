from __future__ import annotations

import functools
import inspect
import typing as t
from collections.abc import Sequence
from dataclasses import dataclass

from .exceptions import StoreError

CallableT = t.TypeVar("CallableT", bound=t.Callable[..., t.Any])


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
def required_args(*variants: Sequence[str]) -> t.Callable[[CallableT], CallableT]:  # noqa: C901
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


T = t.TypeVar("T")

E = t.TypeVar("E", bound=Exception)


@dataclass
class Result(t.Generic[T, E]):
    value: T | None
    error: E | None

    @property
    def ok(self) -> bool:
        return self.error is None

    @property
    def has_value(self) -> bool:
        return self.value is not None


class Address(t.NamedTuple):
    protocol: str
    host: str
    port: int


@dataclass
class Kube:
    user_token: str | None = None
    from kubernetes import client, config

    DEFAULT_NS = "kubeflow"
    DSC_CRD = "datasciencecluster.opendatahub.io/v1"
    DSC_NS_CONFIG = "registriesNamespace"
    EXTERNAL_ADDR_ANNOTATION = "routing.opendatahub.io/external-address-rest"

    def __post_init__(self):
        self.config.load_incluster_config()
        client = Kube.client.ApiClient()
        self.sa_token = client.configuration.api_key["authorization"]
        self.api_client = client

    def __enter__(self) -> Kube:
        return self

    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        self.api_client.close()

    def try_get(
        self, op: t.Callable[[], t.Any], as_user: bool = False
    ) -> Result[t.Any, client.ApiException]:
        if as_user and self.user_token is not None:
            # NOTE: even though this config is consumed by the RESTClient, auth is refreshed on every request: https://github.com/kubernetes-client/python/blob/b7ccf179f1b0194a0ed18e39fb063ef8a963fc6b/kubernetes/client/api_client.py#L166
            self.api_client.configuration.api_key["authorization"] = self.user_token
        try:
            return Result(op(), None)
        except Kube.client.ApiException as e:
            if e.status != 403:
                raise e
            return Result(None, e)
        finally:
            self.api_client.configuration.api_key["authorization"] = self.sa_token

    def try_get_with_any_token(
        self, op: t.Callable[[], t.Any]
    ) -> Result[t.Any, client.ApiException]:
        res = self.try_get(op)
        if res.error is not None and self.user_token:
            res = self.try_get(op, as_user=True)
        return res

    def get_default_dsc(self) -> Result[dict[str, t.Any], StoreError]:
        kcustom = Kube.client.CustomObjectsApi(self.api_client)

        g, v = Kube.DSC_CRD.split("/")
        p = f"{g.split('.')[0]}s"

        def list_dscs() -> t.Any:
            return kcustom.list_cluster_custom_object(
                group=g,
                version=v,
                plural=p,
            )

        res = self.try_get_with_any_token(list_dscs)
        if dscs := res.value:
            return Result(
                t.cast(
                    dict[str, t.Any],
                    dscs["items"][0],
                ),
                None,
            )
        return Result(None, StoreError(f"Failed to list {p}: {res.error}"))

    def get_mr_ns(self) -> Result[str, StoreError]:
        res = self.get_default_dsc()
        if dsc_raw := res.value:
            return Result(
                dsc_raw["status"]["components"]["modelregistry"][Kube.DSC_NS_CONFIG],
                None,
            )
        return Result(Kube.DEFAULT_NS, res.error)

    def get_namespaced_service(
        self, name: str, ns: str
    ) -> Result[client.V1Service, StoreError]:
        kcore = self.client.CoreV1Api(self.api_client)

        def get_service() -> t.Any:
            return kcore.read_namespaced_service(name, ns)

        res = self.try_get_with_any_token(get_service)
        if serv := res.value:
            return Result(t.cast(Kube.client.V1Service, serv), None)
        return Result(None, StoreError(f"Failed to get service {name}: {res.error}"))

    def get_service_addr(self, name: str, ns: str) -> Result[Address, StoreError]:
        res = self.get_namespaced_service(name, ns)
        if res.error:
            return Result(None, res.error)

        serv = res.value
        assert serv is not None
        meta = t.cast(Kube.client.V1ObjectMeta, serv.metadata)
        ext_addr = t.cast(dict[str, str], meta.annotations).get(
            Kube.EXTERNAL_ADDR_ANNOTATION
        )
        err = None
        if not ext_addr:
            host = str(meta.name)
            port_by_protocol = {
                port.app_protocol: port
                for port in t.cast(
                    list[Kube.client.V1ServicePort],
                    t.cast(Kube.client.V1ServiceSpec, serv.spec).ports,
                )
                if port.app_protocol in ("http", "https")
            }
            if p := port_by_protocol.get("https"):
                port = int(str(p.port))
                protocol = "https"
            elif p := port_by_protocol.get("http"):
                port = int(str(p.port))
                protocol = "http"
            else:
                err = StoreError(f"Service {name} has no http(s) ports")
                port = 8080
                protocol = "http"
        else:
            from urllib.parse import urlparse

            parsed = urlparse(ext_addr)
            protocol = parsed.scheme
            host, port = parsed.netloc.split(":")
            port = int(port)

        return Result(Address(protocol, host, port), err)
