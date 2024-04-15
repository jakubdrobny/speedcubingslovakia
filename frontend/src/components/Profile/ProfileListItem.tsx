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
import { WIN_SMALL } from "../../constants";
import { useContext } from "react";

const ProfileListItem: React.FC<{ windowWidth: number }> = ({
  windowWidth,
}) => {
  const { authState, setAuthState } = useContext(
    AuthContext
  ) as AuthContextType;
  const { closeNav } = useContext(NavContext) as NavContextType;

  const handleLogOut = () => {
    setAuthState(initialAuthState);
    closeNav();
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
      placement={windowWidth < WIN_SMALL ? "bottom-start" : "bottom"}
      enterTouchDelay={0}
    >
      <ListItem>
        <Avatar src={authState.avatarUrl} />
        {authState.wcaid}
      </ListItem>
    </Tooltip>
  );
};

export default ProfileListItem;
