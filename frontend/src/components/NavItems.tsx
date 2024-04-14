import { AuthContextType, NavContextType } from "../Types";
import { Grid, List, ListItemButton, ListItemDecorator, Stack } from "@mui/joy";
import { Language, ListAlt } from "@mui/icons-material";

import { AuthContext } from "../context/AuthContext";
import LanguageIcon from "@mui/icons-material/Language";
import { Link } from "react-router-dom";
import { NavContext } from "../context/NavContext";
import ProfileListItem from "./Profile/ProfileListItem";
import WCALogoNoText from "../images/WCALogoNoText";
import { useContext } from "react";

const NavItems = () => {
  const { authState } = useContext(AuthContext) as AuthContextType;
  const { navOpen, toggleNavOpen } = useContext(NavContext) as NavContextType;

  return (
    <Stack direction={navOpen ? "column" : "row"}>
      <ListItemButton
        component={Link}
        to="/competitions"
        sx={navOpen ? { justifyContent: "center", mb: 1 } : {}}
        onClick={toggleNavOpen}
      >
        <ListItemDecorator>
          <LanguageIcon />
        </ListItemDecorator>
        Online Competitions
      </ListItemButton>
      {authState.isadmin && (
        <ListItemButton component={Link} to="/admin/dashboard">
          <ListItemDecorator>
            <ListAlt />
          </ListItemDecorator>
          Dashboard
        </ListItemButton>
      )}
      {authState.token ? (
        <ProfileListItem />
      ) : (
        <ListItemButton
          component={Link}
          to={process.env.REACT_APP_WCA_GET_CODE_URL || ""}
          sx={navOpen ? { justifyContent: "center" } : {}}
          onClick={toggleNavOpen}
        >
          <ListItemDecorator>
            <WCALogoNoText />
          </ListItemDecorator>
          Log In
        </ListItemButton>
      )}
    </Stack>
  );
};

export default NavItems;
