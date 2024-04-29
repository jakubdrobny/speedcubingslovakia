import { AuthContextType, NavContextType } from "../../Types";
import {
  Grid,
  List,
  ListItem,
  ListItemButton,
  ListItemDecorator,
  Stack,
} from "@mui/joy";
import { Language, ListAlt, Search } from "@mui/icons-material";

import { AuthContext } from "../../context/AuthContext";
import LanguageIcon from "@mui/icons-material/Language";
import { Link } from "react-router-dom";
import { NavContext } from "../../context/NavContext";
import ProfileListItem from "../Profile/ProfileListItem";
import WCALogoNoText from "../../images/WCALogoNoText";
import { useContext } from "react";

const NavItems: React.FC<{
  direction: "row" | "row-reverse" | "column" | "column-reverse";
  windowWidth: number;
}> = ({ direction, windowWidth }) => {
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
      <ListItemButton
        component={Link}
        to="/users"
        onClick={closeNav}
        sx={{ justifyContent: "flex-end" }}
      >
        <ListItemDecorator>
          <Search />
        </ListItemDecorator>
        Users
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
          onClick={closeNav}
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
