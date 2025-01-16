import { AuthContextType, NavContextType } from "../../Types";
import { initialAuthState, logOut } from "../../utils/utils";

import { AuthContext } from "../../context/AuthContext";
import { NavContext } from "../../context/NavContext";
import { useContext } from "react";
import { useNavigate } from "react-router-dom";
import NavItem from "./NavItem";
import { IconLogout, IconUser } from "@tabler/icons-react";

const ProfileListItem = () => {
  const { authState, setAuthState } = useContext(
    AuthContext,
  ) as AuthContextType;
  const { closeNav } = useContext(NavContext) as NavContextType;
  const navigate = useNavigate();

  const handleLogOut = () => {
    setAuthState(initialAuthState);
    closeNav();
    logOut();
    document.location.reload();
  };

  const goToProfile = () => {
    closeNav();
    console.log("navigating");
    navigate(`profile/${authState.wcaid}`);
  };

  return (
    <NavItem
      listItemType="profile"
      Title={<div>{authState.wcaid}</div>}
      TitleIcon={
        <img
          height="24"
          width="24"
          src={authState.avatarUrl}
          alt="avatar"
          className="rounded-full"
        />
      }
      onClick={closeNav}
      sublistItems={[
        {
          title: "My profile",
          to: "",
          onClick: goToProfile,
          icon: <IconUser />,
        },
        {
          title: "Logout",
          to: "",
          onClick: handleLogOut,
          icon: <IconLogout />,
        },
      ]}
    />
  );
};

export default ProfileListItem;
