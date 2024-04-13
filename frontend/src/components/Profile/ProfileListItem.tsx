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
import { AuthContextType } from "../../Types";
import { Logout } from "@mui/icons-material";
import { useContext } from "react";

const ProfileListItem = () => {
  const { authState, setAuthState } = useContext(
    AuthContext
  ) as AuthContextType;

  const handleLogOut = () => {
    setAuthState(initialAuthState);
    logOut();
    document.location.href = "/";
    document.location.reload();
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
    >
      <ListItem>
        <Avatar src={authState.avatarUrl} />
        {authState.wcaid}
      </ListItem>
    </Tooltip>
  );
};

export default ProfileListItem;
