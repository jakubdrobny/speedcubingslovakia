import { NavContextType, WindowSizeContextType } from "../../Types";
import { WIN_VERY_LG } from "../../constants";

import { IconBaselineDensityMedium } from "@tabler/icons-react";
import { Link } from "react-router-dom";
import { NavContext } from "../../context/NavContext";
import NavItems from "./NavItems";
import { WindowSizeContext } from "../../context/WindowSizeContext";
import { useContext } from "react";

const NavHorizontal = () => {
  const { navOpen, closeNav, toggleNavOpen } = useContext(
    NavContext,
  ) as NavContextType;
  const { windowSize } = useContext(WindowSizeContext) as WindowSizeContextType;

  return (
    <div className="col-span-12 lg:col-span-10 xl:col-span-8 border-solid border-0 border-b-2 border-gray-200 h-full flex items-center justify-between px-4">
      <Link
        to="/"
        onClick={closeNav}
        className="flex items-center no-underline text-black"
      >
        <img
          src="/speedcubingslovakialogo256.png"
          height="48"
          width="48"
          alt="Speedcubing Slovakia Logo"
        />
        <p className="text-lg font-semibold ml-2 text-nowrap">
          Speedcubing Slovakia
        </p>
      </Link>

      <div className="flex items-center h-full">
        {windowSize.width < WIN_VERY_LG ? (
          <button
            onClick={toggleNavOpen}
            className={`px-3 py-1 border-solid rounded-md ${navOpen
                ? "bg-gray-500 text-white border-gray-500"
                : "bg-transparent border-gray-200"
              } mr-2 flex items-center`}
          >
            <IconBaselineDensityMedium />
          </button>
        ) : (
          <NavItems direction="row" />
        )}
      </div>
    </div>
  );
};

export default NavHorizontal;
