import * as React from 'react';
import { createComparativeValue, renderHook, standardUseFetchState, testHook } from './hooks';

const useSayHello = (who: string, showCount = false) => {
  const countRef = React.useRef(0);
  countRef.current++;
  return `Hello ${who}!${showCount && countRef.current > 1 ? ` x${countRef.current}` : ''}`;
};

const useSayHelloDelayed = (who: string, delay = 0) => {
  const [speech, setSpeech] = React.useState('');
  React.useEffect(() => {
    const handle = setTimeout(() => setSpeech(`Hello ${who}!`), delay);
    return () => clearTimeout(handle);
  }, [who, delay]);
  return speech;
};

describe('hook test utils', () => {
  it('simple testHook', () => {
    const renderResult = testHook((who: string) => `Hello ${who}!`)('world');
    expect(renderResult).hookToBe('Hello world!');
    expect(renderResult).hookToHaveUpdateCount(1);
    renderResult.rerender('world');
    expect(renderResult).hookToBe('Hello world!');
    expect(renderResult).hookToBeStable();
    expect(renderResult).hookToHaveUpdateCount(2);
  });

  it('use testHook for rendering', () => {
    const renderResult = testHook(useSayHello)('world');
    expect(renderResult).hookToHaveUpdateCount(1);
    expect(renderResult).hookToBe('Hello world!');
    expect(renderResult).hookToStrictEqual('Hello world!');

    renderResult.rerender('world', false);

    expect(renderResult).hookToHaveUpdateCount(2);
    expect(renderResult).hookToBe('Hello world!');
    expect(renderResult).hookToStrictEqual('Hello world!');
    expect(renderResult).hookToBeStable();

    renderResult.rerender('world', true);

    expect(renderResult).hookToHaveUpdateCount(3);
    expect(renderResult).hookToBe('Hello world! x3');
    expect(renderResult).hookToStrictEqual('Hello world! x3');
  });

  it('use renderHook for rendering', () => {
    type Props = {
      who: string;
      showCount?: boolean;
    };
    const renderResult = renderHook(({ who, showCount }: Props) => useSayHello(who, showCount), {
      initialProps: {
        who: 'world',
      },
    });

    expect(renderResult).hookToHaveUpdateCount(1);
    expect(renderResult).hookToBe('Hello world!');
    expect(renderResult).hookToStrictEqual('Hello world!');

    renderResult.rerender({
      who: 'world',
    });

    expect(renderResult).hookToHaveUpdateCount(2);
    expect(renderResult).hookToBe('Hello world!');
    expect(renderResult).hookToStrictEqual('Hello world!');

    renderResult.rerender({ who: 'world' });

    expect(renderResult).hookToHaveUpdateCount(3);
    expect(renderResult).hookToBe('Hello world!');
    expect(renderResult).hookToStrictEqual('Hello world!');
  });

  it('should use waitForNextUpdate for async update testing', async () => {
    const renderResult = testHook(useSayHelloDelayed)('world');
    expect(renderResult).hookToHaveUpdateCount(1);
    expect(renderResult).hookToBe('');

    await renderResult.waitForNextUpdate();
    expect(renderResult).hookToHaveUpdateCount(2);
    expect(renderResult).hookToBe('Hello world!');
  });

  it('should throw error if waitForNextUpdate times out', async () => {
    const renderResult = renderHook(() => useSayHelloDelayed('', 20));

    await expect(renderResult.waitForNextUpdate({ timeout: 10, interval: 5 })).rejects.toThrow();
    expect(renderResult).hookToHaveUpdateCount(1);

    // unmount to test waiting for an update that will never happen
    renderResult.unmount();

    await expect(renderResult.waitForNextUpdate({ timeout: 500, interval: 10 })).rejects.toThrow();

    expect(renderResult).hookToHaveUpdateCount(1);
  });

  it('should not throw if waitForNextUpdate timeout is sufficient', async () => {
    const renderResult = renderHook(() => useSayHelloDelayed('', 20));

    await expect(
      renderResult.waitForNextUpdate({ timeout: 500, interval: 10 }),
    ).resolves.not.toThrow();

    expect(renderResult).hookToHaveUpdateCount(2);
  });

  it('should assert stability of results using isStable', () => {
    let testValue = 'test';
    const renderResult = renderHook(() => testValue);
    expect(renderResult).hookToHaveUpdateCount(1);

    renderResult.rerender();
    expect(renderResult).hookToHaveUpdateCount(2);
    expect(renderResult).hookToBeStable();

    testValue = 'new';
    renderResult.rerender();
    expect(renderResult).hookToHaveUpdateCount(3);

    renderResult.rerender();
    expect(renderResult).hookToHaveUpdateCount(4);
    expect(renderResult).hookToBeStable();
  });

  it(`should assert stability of result using isStable 'array'`, () => {
    let testValue = ['test'];
    // explicitly returns a new array each render to show the difference between `isStable` and `isStableArray`
    const renderResult = renderHook(() => testValue);
    expect(renderResult).hookToHaveUpdateCount(1);

    renderResult.rerender();
    expect(renderResult).hookToHaveUpdateCount(2);
    expect(renderResult).hookToBeStable();
    expect(renderResult).hookToBeStable([true]);

    testValue = ['new'];
    renderResult.rerender();
    expect(renderResult).hookToHaveUpdateCount(3);
    expect(renderResult).hookToBeStable([false]);

    renderResult.rerender();
    expect(renderResult).hookToHaveUpdateCount(4);
    expect(renderResult).hookToBeStable();
    expect(renderResult).hookToBeStable([true]);
  });

  it('standardUseFetchState should return an array matching the state of useFetchState', () => {
    expect(['test', false, undefined, () => null]).toStrictEqual(standardUseFetchState('test'));
    expect(['test', true, undefined, () => null]).toStrictEqual(
      standardUseFetchState('test', true),
    );
    expect(['test', false, new Error('error'), () => null]).toStrictEqual(
      standardUseFetchState('test', false, new Error('error')),
    );
  });

  describe('createComparativeValue', () => {
    it('should extract array values according to the boolean object', () => {
      expect([1, 2, 3]).toStrictEqual(createComparativeValue([1, 2, 3], [true, true, true]));
      expect([1, 2, 3]).toStrictEqual(createComparativeValue([1, 2, 3], [true, true, false]));
      expect([1, 2, 3]).toStrictEqual(createComparativeValue([1, 2, 4], [true, true, false]));
      expect([1, 2, 3]).toStrictEqual(createComparativeValue([1, 2, 4], [true, true]));
      expect([1, 2, 3]).not.toStrictEqual(createComparativeValue([1, 4, 3], [true, true, true]));
      expect([3, 2, 1]).not.toStrictEqual(createComparativeValue([1, 2, 3], [true, true, true]));
      expect([true, false]).not.toStrictEqual(createComparativeValue([false, true], [true, true]));
      // array comparison must have the same length, however the stability array may have a lesser length
      expect([1, 2, 3, 4]).toStrictEqual(createComparativeValue([1, 2, 3, 5], [true, true, true]));
    });

    it('should extract object values according to the boolean object', () => {
      expect({ a: 1, b: 2, c: 3 }).toStrictEqual(
        createComparativeValue({ a: 1, b: 2, c: 3 }, { a: true, b: true, c: true }),
      );
      expect({ a: 1, b: 2, c: 3 }).toStrictEqual(
        createComparativeValue({ a: 1, b: 2, c: 3 }, { a: true, b: true, c: true }),
      );
      expect({ a: 1, b: 2, c: 3 }).toStrictEqual(
        createComparativeValue({ a: 1, b: 2, c: 4 }, { a: true, b: true, c: false }),
      );
      expect({ a: 1, b: 2, c: 3 }).toStrictEqual(
        createComparativeValue({ a: 1, b: 2, c: 4 }, { a: true, b: true }),
      );
      expect({ a: 1, b: 2, c: 3 }).not.toStrictEqual(
        createComparativeValue({ a: 1, b: 4, c: 3 }, { a: true, b: true, c: true }),
      );
    });

    it('should extract nested values', () => {
      const testValue = {
        a: 1,
        b: {
          c: 2,
          d: [{ e: 3 }, 'f', {}],
        },
      };
      expect(testValue).toStrictEqual(
        createComparativeValue(
          { a: 10, b: { c: 2, d: [null, 'f', null] } },
          {
            b: {
              c: true,
              d: [false, true],
            },
          },
        ),
      );
    });

    it('should extract objects for identity comparisons', () => {
      const obj = {};
      const array: string[] = [];
      const testValue = {
        a: obj,
        b: array,
        c: {
          d: obj,
          e: array,
        },
      };

      expect(testValue).not.toStrictEqual(
        createComparativeValue(
          {
            a: {},
            b: [],
            c: {
              d: {},
              e: [],
            },
          },
          {
            a: true,
            b: true,
            c: { d: true, e: true },
          },
        ),
      );

      expect(testValue).toStrictEqual(
        createComparativeValue(
          {
            a: obj,
            b: array,
            c: {
              d: obj,
              e: array,
            },
          },
          {
            a: true,
            b: true,
            c: { d: true, e: true },
          },
        ),
      );
    });
  });
});
