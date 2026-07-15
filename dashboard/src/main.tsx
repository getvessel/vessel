import { RouterProvider } from '@tanstack/react-router';
import React from 'react';
import { createRoot } from 'react-dom/client';
import { getRouter } from './router';

const router = getRouter();

const rootElement = document.getElementById('root');
if (rootElement && !rootElement.innerHTML) {
  const root = createRoot(rootElement);
  root.render(
    <React.StrictMode>
      <RouterProvider router={router} />
    </React.StrictMode>
  );
}
