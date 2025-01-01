import { AuthContextType, NavContextType } from "../../Types";
import { Badge, List, ListItemButton, ListItemDecorator } from "@mui/joy";
import { Campaign, ListAlt } from "@mui/icons-material";
import {
  GetNoOfNewAnnouncements,
  saveCurrentLocation,
} from "../../utils/utils";
import { Link, useLocation } from "react-router-dom";
import { useContext, useEffect, useState } from "react";

import { AuthContext } from "../../context/AuthContext";
import LanguageIcon from "@mui/icons-material/Language";
import { NavContext } from "../../context/NavContext";
import ProfileListItem from "../Profile/ProfileListItem";
import ResultsListItem from "./ResultsListItem";
import WCALogoNoText from "../../images/WCALogoNoText";

const NavItems: React.FC<{
  direction: "row" | "row-reverse" | "column" | "column-reverse";
}> = ({ direction }) => {
  const { authStateRef } = useContext(AuthContext) as AuthContextType;
  const { closeNav } = useContext(NavContext) as NavContextType;
  const [newAnnouncements, setNewAnnouncements] = useState<number>(-1);
  const location = useLocation();

  useEffect(() => {
    if (authStateRef.current.token) {
      GetNoOfNewAnnouncements()
        .then((res) => setNewAnnouncements(res))
        .catch((_) => {});
    }
  }, [location.pathname]);

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
        {authStateRef.current.token && newAnnouncements > 0 ? (
          <Badge
            badgeContent={newAnnouncements.toString()}
            color="danger"
            variant="soft"
          >
            Announcements
          </Badge>
        ) : (
          <>Announcements</>
        )}
      </ListItemButton>
      {authStateRef.current.isadmin && (
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
      {authStateRef.current.token ? (
        <ProfileListItem />
      ) : (
        <ListItemButton
          component={Link}
          to={import.meta.env.VITE_WCA_GET_CODE_URL || ""}
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
