import { AuthContextType, NavContextType } from "../../Types";
import {
  GetNoOfNewAnnouncements,
  saveCurrentLocation,
} from "../../utils/utils";
import { useLocation } from "react-router-dom";
import { useContext, useEffect } from "react";

import { AuthContext } from "../../context/AuthContext";
import { NavContext } from "../../context/NavContext";
import ProfileListItem from "./ProfileListItem";
import WCALogoNoText from "../../images/WCALogoNoText";
import NavItem from "./NavItem";
import { IconList } from "@tabler/icons-react";
import CompetitionsNavItem from "./CompetitionsNavItem";
import ResultsNavItem from "./ResultsNavItem";
import AnnouncementsNavItem from "./AnnouncementsNavItem";
import useState from "react-usestateref";

const NavItems: React.FC<{
  direction: "row" | "col";
}> = ({ direction }) => {
  const { authStateRef } = useContext(AuthContext) as AuthContextType;
  const { closeNav } = useContext(NavContext) as NavContextType;
  const [newAnnouncements, setNewAnnouncements] = useState<number>(-1);
  const location = useLocation();

  useEffect(() => {
    if (authStateRef.current.token) {
      GetNoOfNewAnnouncements()
        .then((res) => setNewAnnouncements(res))
        .catch((_) => { });
    }
  }, [location.pathname]);

  return (
    <div
      className={`flex flex-${direction} h-full ${direction === "col" ? "gap-2 my-2" : ""}`}
    >
      <CompetitionsNavItem />
      <ResultsNavItem />
      <AnnouncementsNavItem
        isAuthenticated={
          authStateRef.current.token !== undefined &&
          authStateRef.current.token !== ""
        }
        newAnnouncements={newAnnouncements}
      />
      {authStateRef.current.isadmin && (
        <NavItem
          Title={<div>Dashboard</div>}
          to="/admin/dashboard"
          TitleIcon={<IconList />}
          onClick={closeNav}
          sublistItems={[]}
        />
      )}
      {authStateRef.current.token ? (
        <ProfileListItem />
      ) : (
        <NavItem
          Title={<div>Log In</div>}
          TitleIcon={<WCALogoNoText />}
          to={import.meta.env.VITE_WCA_GET_CODE_URL || ""}
          onClick={() => {
            saveCurrentLocation(window.location.pathname);
            closeNav();
          }}
          sublistItems={[]}
        />
      )}
    </div>
  );
};

export default NavItems;
