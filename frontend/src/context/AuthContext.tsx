import { AuthContextType, AuthState } from "../Types";
import React, { ReactNode, createContext, useState } from "react";

export const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider: React.FC<{ children?: ReactNode }> = ({ children }) => {
    const [authState, setAuthState] = useState<AuthState>(initialState);

    const updateAuthToken = (newToken: string) => setAuthState({...authState, token: newToken});

    return (
        <AuthContext.Provider value={{authState, updateAuthToken}}>
            {children}
        </AuthContext.Provider>
    );
}

const initialState: AuthState = {
    token: '',
    authenticated: true,
    admin: true
};