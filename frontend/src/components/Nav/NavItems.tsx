import { AuthContextType, NavContextType } from "../../Types";
import { Badge, List, ListItemButton, ListItemDecorator } from "@mui/joy";
import { Campaign, ListAlt } from "@mui/icons-material";
import { GetNoOfNewAnnouncements, saveCurrentLocation } from "../../utils";
import { useContext, useEffect, useState } from "react";

import { AuthContext } from "../../context/AuthContext";
import LanguageIcon from "@mui/icons-material/Language";
import { Link } from "react-router-dom";
import { NavContext } from "../../context/NavContext";
import ProfileListItem from "../Profile/ProfileListItem";
import ResultsListItem from "./ResultsListItem";
import WCALogoNoText from "../../images/WCALogoNoText";

const NavItems: React.FC<{
  direction: "row" | "row-reverse" | "column" | "column-reverse";
}> = ({ direction }) => {
  const { authStateRef } = useContext(AuthContext) as AuthContextType;
  const { closeNav } = useContext(NavContext) as NavContextType;
  const [newAnnouncements, setNewAnnouncements] = useState<number>(0);

  useEffect(() => {
    if (authStateRef.current.token) {
      GetNoOfNewAnnouncements()
        .then((res) => setNewAnnouncements(res))
        .catch((err) => {});
    }
  }, []);

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
        {authStateRef.current.token && newAnnouncements > 0 ? (
          <Badge badgeContent={newAnnouncements.toString()} color="danger">
            <ListItemDecorator>
              <Campaign />
            </ListItemDecorator>
            Announcements
          </Badge>
        ) : (
          <>
            <ListItemDecorator>
              <Campaign />
            </ListItemDecorator>
            Announcements
          </>
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
