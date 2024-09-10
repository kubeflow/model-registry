import {
  renderHook as renderHookRTL,
  RenderHookOptions,
  RenderHookResult,
  waitFor,
  waitForOptions,
} from '@testing-library/react';
import { queries, Queries } from '@testing-library/dom';

// eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
export type BooleanValues<T> = T extends boolean | number | null | undefined | Function
  ? boolean | undefined
  : boolean | undefined | { [K in keyof T]?: BooleanValues<T[K]> };

/**
 * Extension of RTL RenderHookResult providing functions used query the current state of the result.
 */
export type RenderHookResultExt<Result, Props> = RenderHookResult<Result, Props> & {
  /**
   * Returns the previous result.
   */
  getPreviousResult: () => Result;

  /**
   * Get the update count for how many times the hook has been rendered.
   * An update occurs initially on render, subsequently when re-rendered, and also whenever the hook itself triggers a re-render.
   * eg. An `useEffect` triggering a state update.
   */
  getUpdateCount: () => number;

  /**
   * Returns a Promise that resolves the next time the hook renders, commonly when state is updated as the result of an asynchronous update.
   *
   * Since `waitForNextUpdate` works using interval checks (backed by `waitFor`), it's possible that multiple updates may occur while waiting.
   */
  waitForNextUpdate: (options?: Pick<waitForOptions, 'interval' | 'timeout'>) => Promise<void>;
};

/**
 * Wrapper on top of RTL `renderHook` returning a result that implements the `RenderHookResultExt` interface.
 *
 * `renderHook` provides full control over the rendering of your hook including the ability to wrap the test component.
 * This is usually used to add context providers from `React.createContext` for the hook to access with `useContext`.
 * `initialProps` and props subsequently set by `rerender` will be provided to the wrapper.
 *
 * ```
 * const renderResult = renderHook(({ who }: { who: string }) => useSayHello(who), { initialProps: { who: 'world' }});
 * expect(renderResult).hookToBe('Hello world!');
 * renderResult.rerender({ who: 'there' });
 * expect(renderResult).hookToBe('Hello there!');
 * ```
 */
export const renderHook = <
  Result,
  Props,
  Q extends Queries = typeof queries,
  Container extends Element | DocumentFragment = HTMLElement,
  BaseElement extends Element | DocumentFragment = Container,
>(
  render: (initialProps: Props) => Result,
  options?: RenderHookOptions<Props, Q, Container, BaseElement>,
): RenderHookResultExt<Result, Props> => {
  let updateCount = 0;
  let prevResult: Result;
  let currentResult: Result;

  const renderResult = renderHookRTL((props) => {
    updateCount++;
    prevResult = currentResult;
    currentResult = render(props);
    return currentResult;
  }, options);

  const renderResultExt: RenderHookResultExt<Result, Props> = {
    ...renderResult,

    getPreviousResult: () => (updateCount > 1 ? prevResult : renderResult.result.current),

    getUpdateCount: () => updateCount,

    waitForNextUpdate: async (currentOptions) => {
      const expected = updateCount;
      try {
        await waitFor(() => expect(updateCount).toBeGreaterThan(expected), currentOptions);
      } catch {
        throw new Error('waitForNextUpdate timed out');
      }
    },
  };

  return renderResultExt;
};

/**
 * Lightweight API for testing a single hook.
 *
 * Prefer this method of testing over `renderHook` for simplicity.
 *
 * ```
 * const renderResult = testHook(useSayHello)('world');
 * expectHook(renderResult).toBe('Hello world!');
 * renderResult.rerender('there');
 * expectHook(renderResult).toBe('Hello there!');
 * ```
 */
export const testHook =
  <Result, P extends unknown[]>(hook: (...params: P) => Result) =>
  // not ideal to nest functions in terms of API but cannot find a better way to infer P from hook and not initialParams
  (
    ...initialParams: P
  ): Omit<RenderHookResultExt<Result, { $params: typeof initialParams }>, 'rerender'> & {
    rerender: (...params: typeof initialParams) => void;
  } => {
    const renderResult = renderHook<Result, { $params: typeof initialParams }>(
      ({ $params }) => hook(...$params),
      {
        initialProps: {
          $params: initialParams,
        },
      },
    );

    return {
      ...renderResult,
      rerender: (...params) => renderResult.rerender({ $params: params }),
    };
  };

/**
 * A helper function for asserting the return value of hooks based on `useFetchState`.
 *
 * eg.
 * ```
 * expectHook(renderResult).isStrictEqual(standardUseFetchState('test value', true))
 * ```
 * is equivalent to:
 * ```
 * expectHook(renderResult).isStrictEqual(['test value', true, undefined, expect.any(Function)])
 * ```
 */
export const standardUseFetchState = <D>(
  data: D,
  loaded = false,
  error?: Error,
): [
  data: D,
  loaded: boolean,
  loadError: Error | undefined,
  refresh: () => Promise<D | undefined>,
] => [data, loaded, error, expect.any(Function)];

// create a new asymmetric matcher that matches everything
const everything = () => {
  const r = expect.anything();
  r.asymmetricMatch = () => true;
  return r;
};

/**
 * Extracts a subset of values from the source that can be used to compare equality.
 *
 * Recursively traverses the `booleanTarget`. For every property or array index equal to `true`,
 * adds the value of the source to the result wrapped in custom matcher `expect.isIdentityEqual`.
 * If the entry is `false` or `undefined`, add an everything matcher to the result.
 */
export const createComparativeValue = <T>(source: T, booleanTarget: BooleanValues<T>): unknown =>
  createComparativeValueRecursive(source, booleanTarget);

const createComparativeValueRecursive = <T>(
  source: unknown,
  // eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
  booleanTarget: boolean | string | number | Function | BooleanValues<T>,
) => {
  if (typeof booleanTarget === 'boolean') {
    return booleanTarget ? expect.isIdentityEqual(source) : everything();
  }
  if (Array.isArray(booleanTarget)) {
    if (Array.isArray(source)) {
      const r = new Array(source.length).fill(everything());
      booleanTarget.forEach((b, i) => {
        if (b != null) {
          r[i] = createComparativeValueRecursive(source[i], b);
        }
      });
      return r;
    }
    return undefined;
  }
  if (
    source == null ||
    typeof source === 'string' ||
    typeof source === 'number' ||
    typeof source === 'function'
  ) {
    return source;
  }
  const obj: { [k: string]: unknown } = {};
  const btObj = booleanTarget as { [k: string]: unknown };
  Object.keys(btObj).forEach((key) => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    obj[key] = createComparativeValueRecursive((source as any)[key] as unknown, btObj[key] as any);
  });
  return expect.objectContaining(obj);
};
