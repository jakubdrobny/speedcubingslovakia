import { Button, Grid, List, ListItemButton, Typography } from "@mui/joy";
import { useContext, useEffect, useState } from "react";

import { DensityMedium } from "@mui/icons-material";
import { Link } from "react-router-dom";
import { NavContext } from "../../context/NavContext";
import { NavContextType } from "../../Types";
import NavItems from "./NavItems";
import { WIN_SMALL } from "../../constants";

const mainLogo = require("../../images/speedcubingslovakialogo256.png");

const NavHorizontal: React.FC<{ windowWidth: number }> = ({ windowWidth }) => {
  const { navOpen, closeNav, toggleNavOpen } = useContext(
    NavContext
  ) as NavContextType;

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
            sx={WIN_SMALL ? { pl: 1 } : {}}
          >
            <img src={mainLogo} height="48" width="48"></img>
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
          {windowWidth < WIN_SMALL ? (
            <Button
              onClick={toggleNavOpen}
              variant={navOpen ? "solid" : "outlined"}
              color="neutral"
            >
              <DensityMedium />
            </Button>
          ) : (
            <NavItems direction="row" windowWidth={windowWidth} />
          )}
        </Grid>
      </List>
    </Grid>
  );
};

export default NavHorizontal;
