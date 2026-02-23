import { EXPECTED_YAML_FORMAT_CONTENT } from '~/app/pages/modelCatalogSettings/expectedYamlFormatContent';

describe('expectedYamlFormatContent', () => {
  it('exports a non-empty string', () => {
    expect(typeof EXPECTED_YAML_FORMAT_CONTENT).toBe('string');
    expect(EXPECTED_YAML_FORMAT_CONTENT.length).toBeGreaterThan(0);
  });

  it('contains expected YAML structure keys', () => {
    expect(EXPECTED_YAML_FORMAT_CONTENT).toContain('all:');
    expect(EXPECTED_YAML_FORMAT_CONTENT).toContain('children:');
    expect(EXPECTED_YAML_FORMAT_CONTENT).toContain('control_nodes:');
    expect(EXPECTED_YAML_FORMAT_CONTENT).toContain('hosts:');
  });
});
