import { AuthContextType, AuthState } from "../Types";
import React, { ReactNode, createContext } from "react";

import { initialAuthState } from "../utils/utils";
import useState from "react-usestateref";

export const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider: React.FC<{ children?: ReactNode }> = ({
  children,
}) => {
  const [authState, setAuthState, authStateRef] =
    useState<AuthState>(initialAuthState);

  const updateAuthToken = (newToken: string) =>
    setAuthState({ ...authState, token: newToken });

  return (
    <AuthContext.Provider
      value={{ authState, updateAuthToken, setAuthState, authStateRef }}
    >
      {children}
    </AuthContext.Provider>
  );
};
