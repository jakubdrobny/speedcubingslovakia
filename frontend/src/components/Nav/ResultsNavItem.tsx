import {
  IconChartBar,
  IconListNumbers,
  IconSearch,
  IconTrophy,
} from "@tabler/icons-react";
import NavItem from "./NavItem";
import { useContext } from "react";
import { NavContext } from "../../context/NavContext";
import { NavContextType } from "../../Types";

const ResultsNavItem = () => {
  const { closeNav } = useContext(NavContext) as NavContextType;

  return (
    <NavItem
      Title={<div>Results</div>}
      onClick={closeNav}
      TitleIcon={<IconListNumbers />}
      sublistItems={[
        {
          title: "Rankings",
          icon: <IconChartBar />,
          to: "/results/rankings",
        },
        {
          title: "Records",
          icon: <IconTrophy />,
          to: "/results/records",
        },
        {
          title: "Users",
          icon: <IconSearch />,
          to: "/results/users",
        },
      ]}
    />
  );
};

export default ResultsNavItem;
