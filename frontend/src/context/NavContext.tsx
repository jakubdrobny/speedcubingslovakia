import React, { ReactNode, createContext } from "react";
import useState from "react-usestateref";

import { NavContextType, NavListItemType } from "../Types";

export const NavContext = createContext<NavContextType | null>(null);

export const NavProvider: React.FC<{ children?: ReactNode }> = ({
  children,
}) => {
  const [navOpen, setNavOpen] = useState<boolean>(false);
  const [_, setCompetitionsSublistOpened, compsSublistOpenedRef] =
    useState(false);
  const [__, setAnnouncementsSublistOpened, announcementsSublistOpenedRef] =
    useState(false);
  const [___, setProfileSublistOpened, profileSublistOpenedRef] =
    useState(false);

  const openNav = () => setNavOpen(true);
  const closeNav = () => setNavOpen(false);
  const toggleNavOpen = () => setNavOpen((p) => !p);
  const toggleSublistOpen = (listItemType: NavListItemType | undefined) => {
    if (listItemType === "competitions") {
      setCompetitionsSublistOpened((p) => !p);
      setAnnouncementsSublistOpened(false);
      setProfileSublistOpened(false);
    } else if (listItemType === "results") {
      setAnnouncementsSublistOpened((p) => !p);
      setCompetitionsSublistOpened(false);
      setProfileSublistOpened(false);
    } else {
      setProfileSublistOpened((p) => !p);
      setAnnouncementsSublistOpened(false);
      setCompetitionsSublistOpened(false);
    }
  };
  const sublistOpen = (listItemType: NavListItemType | undefined) => {
    if (listItemType === "competitions") {
      return compsSublistOpenedRef.current;
    } else if (listItemType === "results") {
      return announcementsSublistOpenedRef.current;
    } else if (listItemType === "profile") {
      return profileSublistOpenedRef.current;
    }
    return false;
  };
  const closeSublists = () => {
    setAnnouncementsSublistOpened(false);
    setCompetitionsSublistOpened(false);
    setProfileSublistOpened(false);
  };

  return (
    <NavContext.Provider
      value={{
        navOpen,
        openNav,
        closeNav,
        toggleNavOpen,
        toggleSublistOpen,
        sublistOpen,
        closeSublists,
      }}
    >
      {children}
    </NavContext.Provider>
  );
};
