declare namespace jest {
  interface Expect {
    isIdentityEqual: <T>(expected: T) => T;
  }

  interface Matchers<R, T> {
    hookToBe: (expected: unknown) => R;
    hookToStrictEqual: (expected: unknown) => R;
    hookToHaveUpdateCount: (expected: number) => R;
    hookToBeStable: <
      V extends T extends Pick<
        import('~/__tests__/unit/testUtils/hooks').RenderHookResultExt<
          infer Result,
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          any
        >,
        'result'
      >
        ? import('~/__tests__/unit/testUtils/hooks').BooleanValues<Result>
        : never,
    >(
      expected?: V,
    ) => R;
  }

  interface Expect {
    isIdentityEqual: (expected: unknown) => AsymmetricMatcher;
  }
}
