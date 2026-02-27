import sampleCatalogYamlContent from '~/app/pages/modelCatalogSettings/sample-catalog.yaml';

describe('sample-catalog.yaml content', () => {
  it('exports a non-empty string', () => {
    expect(typeof sampleCatalogYamlContent).toBe('string');
    expect(sampleCatalogYamlContent.length).toBeGreaterThan(0);
  });

  it('contains expected YAML structure keys', () => {
    expect(sampleCatalogYamlContent).toContain('source:');
    expect(sampleCatalogYamlContent).toContain('models:');
  });
});
