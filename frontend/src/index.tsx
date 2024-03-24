import './styles/index.css';

import App from './App';
import { BrowserRouter } from 'react-router-dom'
import { CompetitionProvider } from './components/Competition/CompetitionContext';
import React from 'react';
import ReactDOM from 'react-dom/client';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <BrowserRouter>
        <CompetitionProvider>
            <App />
        </CompetitionProvider>
    </BrowserRouter>
  </React.StrictMode>
);