import { AuthContextType, NavContextType } from "../../Types";
import { Campaign, ListAlt } from "@mui/icons-material";
import { List, ListItemButton, ListItemDecorator } from "@mui/joy";

import { AuthContext } from "../../context/AuthContext";
import LanguageIcon from "@mui/icons-material/Language";
import { Link } from "react-router-dom";
import { NavContext } from "../../context/NavContext";
import ProfileListItem from "../Profile/ProfileListItem";
import ResultsListItem from "./ResultsListItem";
import WCALogoNoText from "../../images/WCALogoNoText";
import { saveCurrentLocation } from "../../utils";
import { useContext } from "react";

const NavItems: React.FC<{
  direction: "row" | "row-reverse" | "column" | "column-reverse";
}> = ({ direction }) => {
  const { authState } = useContext(AuthContext) as AuthContextType;
  const { closeNav } = useContext(NavContext) as NavContextType;

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
      <ListItemButton
        component={Link}
        to="/announcements"
        onClick={closeNav}
        sx={{ justifyContent: "flex-end" }}
      >
        <ListItemDecorator>
          <Campaign />
        </ListItemDecorator>
        Announcements
      </ListItemButton>
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
            saveCurrentLocation(window.location.pathname);
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
