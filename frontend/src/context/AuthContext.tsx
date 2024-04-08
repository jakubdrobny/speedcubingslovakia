import { AuthContextType, AuthState } from "../Types";
import React, { ReactNode, createContext, useState } from "react";

import { initialAuthState } from "../utils";

export const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider: React.FC<{ children?: ReactNode }> = ({
  children,
}) => {
  const [authState, setAuthState] = useState<AuthState>(initialAuthState);

  const updateAuthToken = (newToken: string) =>
    setAuthState({ ...authState, token: newToken });

  return (
    <AuthContext.Provider value={{ authState, updateAuthToken, setAuthState }}>
      {children}
    </AuthContext.Provider>
  );
};
