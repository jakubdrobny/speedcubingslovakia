import React, { ReactNode, createContext, useState } from "react";
import { WindowSize, WindowSizeContextType } from "../Types";

export const WindowSizeContext = createContext<WindowSizeContextType | null>(
  null
);

export const WindowSizeProvider: React.FC<{ children?: ReactNode }> = ({
  children,
}) => {
  const [windowSize, setWindowSize] = useState<WindowSize>({
    width: 0,
    height: 0,
  });

  return (
    <WindowSizeContext.Provider value={{ windowSize, setWindowSize }}>
      {children}
    </WindowSizeContext.Provider>
  );
};
