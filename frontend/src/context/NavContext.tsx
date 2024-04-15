import React, { ReactNode, createContext, useState } from "react";

import { NavContextType } from "../Types";

export const NavContext = createContext<NavContextType | null>(null);

export const NavProvider: React.FC<{ children?: ReactNode }> = ({
  children,
}) => {
  const [navOpen, setNavOpen] = useState<boolean>(false);

  const openNav = () => setNavOpen(true);
  const closeNav = () => setNavOpen(false);
  const toggleNavOpen = () => setNavOpen((p) => !p);

  return (
    <NavContext.Provider value={{ navOpen, openNav, closeNav, toggleNavOpen }}>
      {children}
    </NavContext.Provider>
  );
};
