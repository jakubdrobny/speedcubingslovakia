import { NavContextType, WindowSizeContextType } from "../../Types";

import { Grid } from "@mui/joy";
import { NavContext } from "../../context/NavContext";
import NavItems from "./NavItems";
import { WIN_SMALL } from "../../constants";
import { WindowSizeContext } from "../../context/WindowSizeContext";
import { useContext } from "react";

const NavVertical = () => {
  const { navOpen } = useContext(NavContext) as NavContextType;
  const { windowSize } = useContext(WindowSizeContext) as WindowSizeContextType;
  const navReallyOpen = windowSize.width < WIN_SMALL && navOpen;

  return (
    <Grid
      xs={12}
      flexDirection="column"
      borderBottom={"2px solid lightgrey"}
      padding="0.5em"
      sx={{
        display: navReallyOpen ? "flex" : "none",
        background: "white",
      }}
    >
      <NavItems direction="column" />
    </Grid>
  );
};

export default NavVertical;
