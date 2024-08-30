import { Grid, Stack } from "@mui/joy";
import { NavContextType, WindowSizeContextType } from "../../Types";

import { NavContext } from "../../context/NavContext";
import NavItems from "./NavItems";
import { WIN_VERY_LG } from "../../constants";
import { WindowSizeContext } from "../../context/WindowSizeContext";
import { useContext } from "react";

const NavVertical = () => {
  const { navOpen } = useContext(NavContext) as NavContextType;
  const { windowSize } = useContext(WindowSizeContext) as WindowSizeContextType;
  const navReallyOpen = windowSize.width < WIN_VERY_LG && navOpen;

  return (
    <Stack direction="row" style={{ width: "100%" }}>
      <Grid
        xs={0}
        lg={1}
        xl={2}
        borderBottom={"2px solid lightgrey"}
        sx={{
          display: navReallyOpen ? "flex" : "none",
          background: "white",
        }}
      />
      <Grid
        xs={12}
        lg={10}
        xl={8}
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
      <Grid
        xs={0}
        lg={1}
        xl={2}
        borderBottom={"2px solid lightgrey"}
        sx={{
          display: navReallyOpen ? "flex" : "none",
          background: "white",
        }}
      />
    </Stack>
  );
};

export default NavVertical;
