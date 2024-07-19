import * as React from 'react';
import App from '@app/index';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';

describe('App tests', () => {
  test('should render default App component', () => {
    const { asFragment } = render(<App />);

    expect(asFragment()).toMatchSnapshot();
  });

  it('should render a nav-toggle button', () => {
    render(<App />);

    expect(screen.getByRole('button', { name: 'Global navigation' })).toBeVisible();
  });
});
