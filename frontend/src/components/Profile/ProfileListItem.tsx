import { AuthContextType, NavContextType } from "../../Types";
import {
  Avatar,
  Box,
  List,
  ListItem,
  ListItemButton,
  ListItemDecorator,
  Tooltip,
} from "@mui/joy";
import { initialAuthState, logOut } from "../../utils";

import { AuthContext } from "../../context/AuthContext";
import { Logout } from "@mui/icons-material";
import { NavContext } from "../../context/NavContext";
import { useContext } from "react";

const ProfileListItem = () => {
  const { authState, setAuthState } = useContext(
    AuthContext
  ) as AuthContextType;
  const { navOpen, setNavOpen } = useContext(NavContext) as NavContextType;

  const handleLogOut = () => {
    setAuthState(initialAuthState);
    setNavOpen(false);
    logOut();
    document.location.reload();
    document.location.href = "/";
  };

  return (
    <Tooltip
      variant="soft"
      title={
        <List size="sm">
          <ListItemButton onClick={handleLogOut}>
            <ListItemDecorator>
              <Logout />
            </ListItemDecorator>
            Logout
          </ListItemButton>
        </List>
      }
      enterTouchDelay={0}
    >
      <ListItem sx={navOpen ? { justifyContent: "center" } : {}}>
        <Avatar src={authState.avatarUrl} />
        {authState.wcaid}
      </ListItem>
    </Tooltip>
  );
};

export default ProfileListItem;
