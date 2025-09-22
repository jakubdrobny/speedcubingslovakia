import {
  EmojiEvents,
  FormatListNumbered,
  Leaderboard,
  Map,
} from "@mui/icons-material";
import { List, ListItemButton, ListItemDecorator, Tooltip } from "@mui/joy";

import { Link } from "react-router-dom";
import { NavContext } from "../../context/NavContext";
import { NavContextType } from "../../Types";
import { useContext } from "react";

const ResultsListItem = () => {
  const { closeNav } = useContext(NavContext) as NavContextType;

  return (
    <Tooltip
      variant="soft"
      title={
        <List size="sm">
          <ListItemButton
            component={Link}
            to="/results/rankings"
            onClick={closeNav}
          >
            <ListItemDecorator>
              <Leaderboard />
            </ListItemDecorator>
            Rankings
          </ListItemButton>
          <ListItemButton
            component={Link}
            to="/results/records"
            onClick={closeNav}
          >
            <ListItemDecorator>
              <EmojiEvents />
            </ListItemDecorator>
            Records
          </ListItemButton>
          <ListItemButton
            component={Link}
            to="/results/users"
            onClick={closeNav}
          >
            <ListItemDecorator>
              <Map />
            </ListItemDecorator>
            User map
          </ListItemButton>
        </List>
      }
      arrow
      enterTouchDelay={0}
    >
      <ListItemButton sx={{ justifyContent: "flex-end" }}>
        <ListItemDecorator>
          <FormatListNumbered />
        </ListItemDecorator>
        Results
      </ListItemButton>
    </Tooltip>
  );
};

export default ResultsListItem;
