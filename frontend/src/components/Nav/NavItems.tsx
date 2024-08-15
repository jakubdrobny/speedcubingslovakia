import { AuthContextType, NavContextType } from "../../Types";
import { Leaderboard, ListAlt, Search } from "@mui/icons-material";
import { Link, useLocation } from "react-router-dom";
import { List, ListItemButton, ListItemDecorator } from "@mui/joy";
import { useContext, useEffect } from "react";

import { AuthContext } from "../../context/AuthContext";
import LanguageIcon from "@mui/icons-material/Language";
import { NavContext } from "../../context/NavContext";
import ProfileListItem from "../Profile/ProfileListItem";
import ResultsListItem from "./ResultsListItem";
import WCALogoNoText from "../../images/WCALogoNoText";
import { saveCurrentLocation } from "../../utils";

const NavItems: React.FC<{
  direction: "row" | "row-reverse" | "column" | "column-reverse";
}> = ({ direction }) => {
  const { authState } = useContext(AuthContext) as AuthContextType;
  const { closeNav } = useContext(NavContext) as NavContextType;

  let location = useLocation();

  return (
    <List
      sx={{
        flexDirection: direction,
        justifyContent: "flex-end",
      }}
    >
      <ListItemButton
        component={Link}
        to="/competitions"
        onClick={closeNav}
        sx={{ justifyContent: "flex-end" }}
      >
        <ListItemDecorator>
          <LanguageIcon />
        </ListItemDecorator>
        Online Competitions
      </ListItemButton>
      <ResultsListItem />
      {authState.isadmin && (
        <ListItemButton
          component={Link}
          to="/admin/dashboard"
          onClick={closeNav}
          sx={{ justifyContent: "flex-end" }}
        >
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
          onClick={() => {
            saveCurrentLocation(location.pathname);
            closeNav();
          }}
          sx={{ justifyContent: "flex-end" }}
        >
          <ListItemDecorator>
            <WCALogoNoText />
          </ListItemDecorator>
          Log In
        </ListItemButton>
      )}
    </List>
  );
};

export default NavItems;
