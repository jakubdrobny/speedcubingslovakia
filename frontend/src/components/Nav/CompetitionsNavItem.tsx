import { IconMail, IconTrophy, IconWorld } from "@tabler/icons-react";
import NavItem from "./NavItem";
import WCALogoNoText from "../../images/WCALogoNoText";
import { useContext } from "react";
import { NavContext } from "../../context/NavContext";
import { NavContextType } from "../../Types";

const CompetitionsNavItem = () => {
  const { closeNav } = useContext(NavContext) as NavContextType;

  return (
    <NavItem
      Title={<div>Competitions</div>}
      TitleIcon={<IconTrophy />}
      sublistItems={[
        {
          title: "Online Competitions",
          icon: <IconWorld />,
          to: "/competitions",
        },
        {
          title: "WCA Competitions",
          icon: <WCALogoNoText />,
          to: "/competitions/wca",
        },
        {
          title: "Announcements",
          icon: <IconMail />,
          to: "/competitions/announcements",
        },
      ]}
      onClick={closeNav}
    />
  );
};

export default CompetitionsNavItem;
