import React, { ReactNode, createContext, useState } from "react";

import { NavContextType } from "../Types";

export const NavContext = createContext<NavContextType | null>(null);

export const NavProvider: React.FC<{ children?: ReactNode }> = ({
  children,
}) => {
  const [navOpen, setNavOpen] = useState<boolean>(false);

  return (
    <NavContext.Provider value={{ navOpen, setNavOpen }}>
      {children}
    </NavContext.Provider>
  );
};
