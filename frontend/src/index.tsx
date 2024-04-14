import "./styles/index.css";

import App from "./App";
import { AuthProvider } from "./context/AuthContext";
import { BrowserRouter } from "react-router-dom";
import { CompetitionProvider } from "./components/Competition/CompetitionContext";
import { NavProvider } from "./context/NavContext";
import ReactDOM from "react-dom/client";
import { TimerInputProvider } from "./context/TimerInputContext";

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);
root.render(
  <BrowserRouter>
    <AuthProvider>
      <CompetitionProvider>
        <TimerInputProvider>
          <NavProvider>
            <App />
          </NavProvider>
        </TimerInputProvider>
      </CompetitionProvider>
    </AuthProvider>
  </BrowserRouter>
);
