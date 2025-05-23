import { Button, Grid, List, ListItemButton, Typography } from "@mui/joy";
import { NavContextType, WindowSizeContextType } from "../../Types";
import { WIN_SMALL, WIN_VERY_LG } from "../../constants";

import { DensityMedium } from "@mui/icons-material";
import { Link } from "react-router-dom";
import { NavContext } from "../../context/NavContext";
import NavItems from "./NavItems";
import { WindowSizeContext } from "../../context/WindowSizeContext";
import { useContext } from "react";

const NavHorizontal = () => {
  const { navOpen, closeNav, toggleNavOpen } = useContext(
    NavContext
  ) as NavContextType;
  const { windowSize } = useContext(WindowSizeContext) as WindowSizeContextType;

  return (
    <Grid
      xs={12}
      lg={10}
      xl={8}
      borderBottom={"2px solid lightgrey"}
      width={"100%"}
    >
      <List
        style={{
          display: "flex",
          flexDirection: "row",
          padding: WIN_SMALL ? 10 : 20,
          width: "100%",
        }}
      >
        <Grid
          sx={{
            display: "flex",
            justifyContent: "flex-start",
          }}
        >
          <ListItemButton
            component={Link}
            to="/"
            onClick={closeNav}
            sx={WIN_SMALL ? { ml: -1 } : {}}
          >
            <img src="/speedcubingslovakialogo256.png" height="48" width="48" alt="SpeedcubingSlovakia Logo"></img>
            <Typography level="h4">Speedcubing Slovakia</Typography>
          </ListItemButton>
        </Grid>

        <Grid
          sx={{
            display: "flex",
            justifyContent: "flex-end",
            alignItems: "center",
            width: "100%",
          }}
        >
          {windowSize.width < WIN_VERY_LG ? (
            <Button
              onClick={toggleNavOpen}
              variant={navOpen ? "solid" : "outlined"}
              color="neutral"
              sx={{ mr: 1 }}
            >
              <DensityMedium />
            </Button>
          ) : (
            <NavItems direction="row" />
          )}
        </Grid>
      </List>
    </Grid>
  );
};

export default NavHorizontal;
