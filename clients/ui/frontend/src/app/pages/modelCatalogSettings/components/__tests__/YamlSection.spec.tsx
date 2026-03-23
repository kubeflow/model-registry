import React from 'react';
import { screen, render } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import '@testing-library/jest-dom';
import YamlSection from '~/app/pages/modelCatalogSettings/components/YamlSection';
import { ManageSourceFormData } from '~/app/pages/modelCatalogSettings/useManageSourceData';
import { CatalogSourceType } from '~/app/modelCatalogTypes';
import {
  FORM_LABELS,
  EXPECTED_YAML_FORMAT_LABEL,
} from '~/app/pages/modelCatalogSettings/constants';

const createFormData = (overrides: Partial<ManageSourceFormData> = {}): ManageSourceFormData => ({
  name: '',
  id: '',
  sourceType: CatalogSourceType.YAML,
  accessToken: '',
  organization: '',
  yamlContent: '',
  allowedModels: '',
  excludedModels: '',
  enabled: false,
  isDefault: false,
  ...overrides,
});

describe('YamlSection', () => {
  const setData = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders YAML label and required indicator', () => {
    render(<YamlSection formData={createFormData()} setData={setData} />);
    expect(screen.getByText(FORM_LABELS.YAML_CONTENT)).toBeInTheDocument();
    expect(screen.getByTestId('yaml-section')).toBeInTheDocument();
    expect(screen.getByTestId('yaml-content-input')).toBeInTheDocument();
  });

  it('shows "View expected file format" link when onToggleExpectedFormatDrawer is provided', () => {
    render(
      <YamlSection
        formData={createFormData()}
        setData={setData}
        onToggleExpectedFormatDrawer={jest.fn()}
      />,
    );
    const link = screen.getByTestId('view-expected-yaml-format-link');
    expect(link).toBeInTheDocument();
    expect(link).toHaveTextContent(EXPECTED_YAML_FORMAT_LABEL);
  });

  it('does not show the link when onToggleExpectedFormatDrawer is not provided', () => {
    render(<YamlSection formData={createFormData()} setData={setData} />);
    expect(screen.queryByTestId('view-expected-yaml-format-link')).not.toBeInTheDocument();
  });

  it('calls onToggleExpectedFormatDrawer when the link is clicked', async () => {
    const user = userEvent.setup();
    const onToggleExpectedFormatDrawer = jest.fn();
    render(
      <YamlSection
        formData={createFormData()}
        setData={setData}
        onToggleExpectedFormatDrawer={onToggleExpectedFormatDrawer}
      />,
    );
    const link = screen.getByTestId('view-expected-yaml-format-link');
    await user.click(link);
    expect(onToggleExpectedFormatDrawer).toHaveBeenCalledTimes(1);
  });
});
