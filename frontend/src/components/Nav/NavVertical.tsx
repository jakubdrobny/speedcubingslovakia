import { useContext } from "react";
import { WIN_VERY_LG } from "../../constants";
import { NavContext } from "../../context/NavContext";
import { WindowSizeContext } from "../../context/WindowSizeContext";
import { NavContextType, WindowSizeContextType } from "../../Types";
import NavItems from "./NavItems";

const NavVertical = () => {
  const { navOpen } = useContext(NavContext) as NavContextType;
  const { windowSize } = useContext(WindowSizeContext) as WindowSizeContextType;
  const navReallyOpen = windowSize.width < WIN_VERY_LG && navOpen;

  return [
    <div
      key="navVertigridleft"
      className={`col-span-0 lg:col-span-1 xl:col-span-2 border-solid border-0 border-b-2 border-b-gray-100 ${navReallyOpen ? "flex" : "hidden"} bg-white h-full mt-1`}
    />,
    <div
      key="navVertigridcenter"
      className={`col-span-12 lg:col-span-10 xl:col-span-8 border-solid border-0 border-b-2 border-b-gray-100 ${navReallyOpen ? "flex flex-col" : "hidden"} bg-white h-full mt-1`}
    >
      <NavItems direction="col" />
    </div>,
    <div
      key="navVertigridright"
      className={`col-span-0 lg:col-span-1 xl:col-span-2 border-solid border-0 border-b-2 border-b-gray-100 ${navReallyOpen ? "flex" : "hidden"} bg-white h-full mt-1`}
    />,
  ];
};

export default NavVertical;
