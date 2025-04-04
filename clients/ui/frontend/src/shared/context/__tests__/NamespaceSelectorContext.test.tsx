import '@testing-library/jest-dom';
import * as React from 'react';
import { render, screen, act } from '@testing-library/react';
import useNamespaces from '~/shared/hooks/useNamespaces';
import * as constUtils from '~/shared/utilities/const';
import {
  NamespaceSelectorContext,
  NamespaceSelectorContextProvider,
} from '~/shared/context/NamespaceSelectorContext';

// Mock the hooks and utilities
jest.mock('~/shared/hooks/useNamespaces');
jest.mock('~/shared/utilities/const');

const mockNamespaces = [{ name: 'namespace-2' }, { name: 'namespace-3' }, { name: 'namespace-1' }];

describe('NamespaceSelectorContext', () => {
  const TestConsumer = () => {
    const {
      namespaces,
      namespacesLoaded,
      namespacesLoadError,
      preferredNamespace,
      initializationError,
    } = React.useContext(NamespaceSelectorContext);

    return (
      <div>
        <div data-testid="loading-state">{namespacesLoaded.toString()}</div>
        <div data-testid="error-state">{namespacesLoadError?.message || 'no-error'}</div>
        <div data-testid="init-error">{initializationError?.message || 'no-init-error'}</div>
        <div data-testid="namespaces">{namespaces.map((ns) => ns.name).join(',')}</div>
        <div data-testid="preferred">{preferredNamespace?.name || 'none'}</div>
      </div>
    );
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should provide initial empty state', () => {
    (useNamespaces as jest.Mock).mockReturnValue([[], true, undefined]);
    (constUtils.isIntegrated as jest.Mock).mockReturnValue(false);

    render(
      <NamespaceSelectorContextProvider>
        <TestConsumer />
      </NamespaceSelectorContextProvider>,
    );

    expect(screen.getByTestId('loading-state')).toHaveTextContent('true');
    expect(screen.getByTestId('error-state')).toHaveTextContent('no-error');
    expect(screen.getByTestId('namespaces')).toHaveTextContent('');
    expect(screen.getByTestId('preferred')).toHaveTextContent('none');
  });

  it('should load namespaces and set first namespace as preferred when not integrated', () => {
    (useNamespaces as jest.Mock).mockReturnValue([mockNamespaces, true, undefined]);
    (constUtils.isIntegrated as jest.Mock).mockReturnValue(false);

    render(
      <NamespaceSelectorContextProvider>
        <TestConsumer />
      </NamespaceSelectorContextProvider>,
    );

    expect(screen.getByTestId('loading-state')).toHaveTextContent('true');
    expect(screen.getByTestId('error-state')).toHaveTextContent('no-error');
    expect(screen.getByTestId('namespaces')).toHaveTextContent(
      'namespace-1,namespace-2,namespace-3',
    );
    expect(screen.getByTestId('preferred')).toHaveTextContent('namespace-1');
  });

  it('should handle errors during namespace loading', () => {
    const error = new Error('Failed to load namespaces');
    (useNamespaces as jest.Mock).mockReturnValue([[], true, error]);
    (constUtils.isIntegrated as jest.Mock).mockReturnValue(false);

    render(
      <NamespaceSelectorContextProvider>
        <TestConsumer />
      </NamespaceSelectorContextProvider>,
    );

    expect(screen.getByTestId('loading-state')).toHaveTextContent('true');
    expect(screen.getByTestId('error-state')).toHaveTextContent('Failed to load namespaces');
    expect(screen.getByTestId('namespaces')).toHaveTextContent('');
    expect(screen.getByTestId('preferred')).toHaveTextContent('none');
  });

  it('should initialize central dashboard client when integrated', () => {
    (useNamespaces as jest.Mock).mockReturnValue([mockNamespaces, true, undefined]);
    (constUtils.isIntegrated as jest.Mock).mockReturnValue(true);

    // Mock window.centraldashboard
    const mockInit = jest.fn();
    global.window.centraldashboard = {
      CentralDashboardEventHandler: {
        init: mockInit,
      },
    };

    render(
      <NamespaceSelectorContextProvider>
        <TestConsumer />
      </NamespaceSelectorContextProvider>,
    );

    expect(mockInit).toHaveBeenCalled();
    expect(screen.getByTestId('preferred')).toHaveTextContent('namespace-1');
  });

  it('should update preferred namespace when selected', () => {
    (useNamespaces as jest.Mock).mockReturnValue([mockNamespaces, true, undefined]);
    (constUtils.isIntegrated as jest.Mock).mockReturnValue(false);

    const TestUpdater = () => {
      const { updatePreferredNamespace } = React.useContext(NamespaceSelectorContext);
      React.useEffect(() => {
        act(() => {
          updatePreferredNamespace({ name: 'namespace-2' });
        });
      }, [updatePreferredNamespace]);
      return null;
    };

    render(
      <NamespaceSelectorContextProvider>
        <TestConsumer />
        <TestUpdater />
      </NamespaceSelectorContextProvider>,
    );

    expect(screen.getByTestId('preferred')).toHaveTextContent('namespace-2');
  });

  it('should handle central dashboard initialization failure gracefully', () => {
    (useNamespaces as jest.Mock).mockReturnValue([mockNamespaces, true, undefined]);
    (constUtils.isIntegrated as jest.Mock).mockReturnValue(true);

    // Mock console.error to avoid test output noise
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation();

    // Mock window.centraldashboard with an init function that throws an error
    const error = new Error('Central dashboard initialization failed');
    global.window.centraldashboard = {
      CentralDashboardEventHandler: {
        init: () => {
          throw error;
        },
      },
    };

    render(
      <NamespaceSelectorContextProvider>
        <TestConsumer />
      </NamespaceSelectorContextProvider>,
    );

    expect(consoleSpy).toHaveBeenCalledWith('Failed to initialize central dashboard client', error);
    expect(screen.getByTestId('init-error')).toHaveTextContent(
      'Central dashboard initialization failed',
    );
    expect(screen.getByTestId('preferred')).toHaveTextContent('namespace-1');

    consoleSpy.mockRestore();
  });

  it('should reflect loading state changes', async () => {
    // Start with loading state
    (useNamespaces as jest.Mock).mockReturnValue([[], false, undefined]);
    (constUtils.isIntegrated as jest.Mock).mockReturnValue(false);

    const { rerender } = render(
      <NamespaceSelectorContextProvider>
        <TestConsumer />
      </NamespaceSelectorContextProvider>,
    );

    expect(screen.getByTestId('loading-state')).toHaveTextContent('false');

    // Update to loaded state
    (useNamespaces as jest.Mock).mockReturnValue([mockNamespaces, true, undefined]);

    rerender(
      <NamespaceSelectorContextProvider>
        <TestConsumer />
      </NamespaceSelectorContextProvider>,
    );

    expect(screen.getByTestId('loading-state')).toHaveTextContent('true');
    expect(screen.getByTestId('namespaces')).toHaveTextContent(
      'namespace-1,namespace-2,namespace-3',
    );
  });

  it('should handle multiple namespace updates correctly', () => {
    (useNamespaces as jest.Mock).mockReturnValue([mockNamespaces, true, undefined]);
    (constUtils.isIntegrated as jest.Mock).mockReturnValue(false);

    const TestMultipleUpdates = () => {
      const { updatePreferredNamespace } = React.useContext(NamespaceSelectorContext);

      React.useEffect(() => {
        act(() => {
          updatePreferredNamespace({ name: 'namespace-2' });
          updatePreferredNamespace({ name: 'namespace-3' });
          updatePreferredNamespace({ name: 'namespace-1' });
        });
      }, [updatePreferredNamespace]);

      return null;
    };

    render(
      <NamespaceSelectorContextProvider>
        <TestConsumer />
        <TestMultipleUpdates />
      </NamespaceSelectorContextProvider>,
    );

    expect(screen.getByTestId('preferred')).toHaveTextContent('namespace-1');
  });
});
