import { Email, EmojiEvents, Language } from "@mui/icons-material";

import WCALogoNoText from "../../images/WCALogoNoText";
import ListItemDropdown from "../ListItemDropdown";

const CompetitionsListItem = () => {
  return (
    <ListItemDropdown
      content={{
        itemText: "Competitions",
        itemIcon: EmojiEvents,
        dropdownItems: [
          {
            text: "Online Competitions",
            url: "/competitions",
            icon: Language,
          },
          {
            text: "WCA Competitions",
            url: "/competitions/wca",
            icon: WCALogoNoText,
          },
          {
            text: "Announcements",
            url: "/competitions/announcements",
            icon: Email,
          },
        ],
      }}
    />
  );
};

export default CompetitionsListItem;
