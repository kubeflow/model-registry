import '@testing-library/jest-dom';
import React from 'react';
import { act, render } from '@testing-library/react';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import RegisterAndStoreFields from '~/app/pages/modelRegistry/screens/RegisterModel/RegisterAndStoreFields';
import {
  ModelLocationType,
  RegistrationCommonFormData,
} from '~/app/pages/modelRegistry/screens/RegisterModel/useRegisterModelData';
import { K8sNameDescriptionFieldData } from '~/concepts/k8s/K8sNameDescriptionField/types';

type MockK8sNameDescriptionFieldProps = {
  data: K8sNameDescriptionFieldData;
  onDataChange?: (key: keyof K8sNameDescriptionFieldData, value: string) => void;
};

let latestK8sNameDescriptionFieldProps: MockK8sNameDescriptionFieldProps | undefined;

jest.mock('~/concepts/k8s/K8sNameDescriptionField/K8sNameDescriptionField', () => ({
  __esModule: true,
  default: (props: MockK8sNameDescriptionFieldProps) => {
    latestK8sNameDescriptionFieldProps = props;
    return <div data-testid="mock-k8s-name-description-field" />;
  },
}));

jest.mock('~/concepts/k8s/NamespaceSelectorField/NamespaceSelectorField', () => ({
  __esModule: true,
  default: () => <div data-testid="mock-namespace-selector-field" />,
}));

jest.mock('~/app/pages/modelRegistry/components/pf-overrides/FormSection', () => ({
  __esModule: true,
  default: ({ children }: { children?: React.ReactNode }) => <>{children}</>,
}));

jest.mock(
  '~/app/pages/modelRegistry/screens/RegisterModel/RegistrationModelLocationFields',
  () => ({
    __esModule: true,
    default: () => <div data-testid="mock-model-location-fields" />,
  }),
);

jest.mock(
  '~/app/pages/modelRegistry/screens/RegisterModel/RegistrationDestinationLocationFields',
  () => ({
    __esModule: true,
    default: () => <div data-testid="mock-destination-location-fields" />,
  }),
);

const createFormData = (
  overrides: Partial<RegistrationCommonFormData> = {},
): RegistrationCommonFormData => ({
  versionName: '',
  versionDescription: '',
  sourceModelFormat: '',
  sourceModelFormatVersion: '',
  modelLocationType: ModelLocationType.ObjectStorage,
  modelLocationEndpoint: '',
  modelLocationBucket: '',
  modelLocationRegion: '',
  modelLocationPath: '',
  modelLocationURI: '',
  modelLocationS3AccessKeyId: '',
  modelLocationS3SecretAccessKey: '',
  destinationOciRegistry: '',
  destinationOciUsername: '',
  destinationOciPassword: '',
  destinationOciUri: '',
  namespace: '',
  registrationMode: undefined,
  jobName: '',
  jobResourceName: '',
  ...overrides,
});

const Harness: React.FC<{ initialData?: Partial<RegistrationCommonFormData> }> = ({
  initialData,
}) => {
  const [formData, setFormData] = React.useState<RegistrationCommonFormData>(
    createFormData(initialData),
  );
  const setData = React.useCallback(
    (key: keyof RegistrationCommonFormData, value: string) => {
      setFormData((previousData) => ({ ...previousData, [key]: value }));
    },
    [setFormData],
  );

  return (
    <RegisterAndStoreFields
      formData={formData}
      setData={setData as unknown as UpdateObjectAtPropAndValue<RegistrationCommonFormData>}
    />
  );
};

describe('RegisterAndStoreFields k8s name wiring', () => {
  beforeEach(() => {
    latestK8sNameDescriptionFieldProps = undefined;
  });

  it('uses the 63-char cap for resource name validity', () => {
    render(<Harness initialData={{ jobResourceName: 'a'.repeat(64) }} />);

    expect(latestK8sNameDescriptionFieldProps).toBeDefined();
    expect(latestK8sNameDescriptionFieldProps?.data.k8sName.state.maxLength).toBe(63);
    expect(latestK8sNameDescriptionFieldProps?.data.k8sName.state.invalidLength).toBe(true);
  });

  it('keeps manual resource name after touched, even when name changes', () => {
    render(<Harness />);

    expect(latestK8sNameDescriptionFieldProps).toBeDefined();

    act(() => {
      latestK8sNameDescriptionFieldProps?.onDataChange?.('name', '!!!');
    });

    const generatedName = latestK8sNameDescriptionFieldProps?.data.k8sName.value;
    expect(generatedName).toMatch(/^gen-[a-z0-9]+$/);
    expect(latestK8sNameDescriptionFieldProps?.data.k8sName.state.touched).toBe(false);

    act(() => {
      latestK8sNameDescriptionFieldProps?.onDataChange?.('k8sName', 'manual-name');
    });

    expect(latestK8sNameDescriptionFieldProps?.data.k8sName.value).toBe('manual-name');
    expect(latestK8sNameDescriptionFieldProps?.data.k8sName.state.touched).toBe(true);

    act(() => {
      latestK8sNameDescriptionFieldProps?.onDataChange?.('name', 'new name');
    });

    expect(latestK8sNameDescriptionFieldProps?.data.k8sName.value).toBe('manual-name');
    expect(latestK8sNameDescriptionFieldProps?.data.k8sName.state.touched).toBe(true);
  });
});
