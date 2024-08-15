import { AuthContextType, NavContextType } from "../../Types";
import {
  Avatar,
  List,
  ListItemButton,
  ListItemDecorator,
  Tooltip,
} from "@mui/joy";
import { Logout, Person } from "@mui/icons-material";
import { initialAuthState, logOut } from "../../utils";

import { AuthContext } from "../../context/AuthContext";
import { NavContext } from "../../context/NavContext";
import { useContext } from "react";
import { useNavigate } from "react-router-dom";

const ProfileListItem = () => {
  const { authState, setAuthState } = useContext(
    AuthContext
  ) as AuthContextType;
  const { closeNav } = useContext(NavContext) as NavContextType;
  const navigate = useNavigate();

  const handleLogOut = () => {
    setAuthState(initialAuthState);
    closeNav();
    logOut();
  };

  const goToProfile = () => {
    closeNav();
    navigate(`/profile/${authState.wcaid}`);
  };

  return (
    <Tooltip
      variant="soft"
      title={
        <List size="sm">
          <ListItemButton onClick={goToProfile}>
            <ListItemDecorator>
              <Person />
            </ListItemDecorator>
            My profile
          </ListItemButton>
          <ListItemButton onClick={handleLogOut}>
            <ListItemDecorator>
              <Logout />
            </ListItemDecorator>
            Logout
          </ListItemButton>
        </List>
      }
      arrow
      enterTouchDelay={0}
    >
      <ListItemButton sx={{ justifyContent: "flex-end" }}>
        <Avatar src={authState.avatarUrl} />
        {authState.wcaid}
      </ListItemButton>
    </Tooltip>
  );
};

export default ProfileListItem;
