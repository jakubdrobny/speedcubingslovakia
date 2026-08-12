import {
  EmojiEvents,
  FormatListNumbered,
  Leaderboard,
  Map,
} from "@mui/icons-material";

import ListItemDropdown from "../ListItemDropdown";

const ResultsListItem = () => {
  return (
    <ListItemDropdown
      content={{
        itemText: "Results",
        itemIcon: FormatListNumbered,
        dropdownItems: [
          {
            text: "Rankings",
            url: "/rankings",
            icon: Leaderboard,
          },
          {
            text: "Records",
            url: "/results/records",
            icon: EmojiEvents,
          },
          {
            text: "User Map",
            url: "/results/users",
            icon: Map,
          },
        ],
      }}
    />
  );
};

export default ResultsListItem;
