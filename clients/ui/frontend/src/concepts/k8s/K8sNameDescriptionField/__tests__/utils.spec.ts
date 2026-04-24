import { translateDisplayNameForK8s } from '~/concepts/k8s/K8sNameDescriptionField/utils';

describe('translateDisplayNameForK8s', () => {
  it('converts spaces to dashes and lowercases', () => {
    expect(translateDisplayNameForK8s('My Model')).toBe('my-model');
  });

  it('collapses multiple dashes (e.g. "My -- Model (v2)" → "my-model-v2")', () => {
    expect(translateDisplayNameForK8s('My -- Model (v2)')).toBe('my-model-v2');
  });

  it('strips leading and trailing dashes', () => {
    expect(translateDisplayNameForK8s('--hello--')).toBe('hello');
  });

  it('removes non-alphanumeric non-dash characters', () => {
    expect(translateDisplayNameForK8s('hello!@#world')).toBe('helloworld');
  });

  it('generates a stable name when input normalizes to nothing (dashes only)', () => {
    const result = translateDisplayNameForK8s('----');
    expect(result).toMatch(/^gen-[a-z0-9]+$/);
    expect(translateDisplayNameForK8s('----')).toBe(result);
  });

  it('generates a stable name when input is punctuation only', () => {
    const result = translateDisplayNameForK8s('!!!');
    expect(result).toMatch(/^gen-[a-z0-9]+$/);
    expect(translateDisplayNameForK8s('!!!')).toBe(result);
  });

  it('uses different stable gen suffixes for different punctuation-only inputs', () => {
    expect(translateDisplayNameForK8s('----')).not.toBe(translateDisplayNameForK8s('####'));
  });

  it('returns empty string when input is empty', () => {
    expect(translateDisplayNameForK8s('')).toBe('');
  });

  it('prepends safePrefix when provided', () => {
    expect(translateDisplayNameForK8s('my model', 'prefix-')).toBe('prefix-my-model');
  });

  it('prepends safePrefix when name normalizes to empty', () => {
    expect(translateDisplayNameForK8s('!!!', 'safe-')).toBe('safe-');
  });
});
