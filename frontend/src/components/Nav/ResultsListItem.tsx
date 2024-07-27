import { AuthContextType, NavContextType } from "../../Types";
import {
  Avatar,
  List,
  ListItemButton,
  ListItemDecorator,
  Tooltip,
} from "@mui/joy";
import {
  EmojiEvents,
  FormatListNumbered,
  Leaderboard,
  Logout,
  Person,
  Search,
} from "@mui/icons-material";
import { Link, useNavigate } from "react-router-dom";
import { initialAuthState, logOut } from "../../utils";

import { AuthContext } from "../../context/AuthContext";
import { NavContext } from "../../context/NavContext";
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
              <Search />
            </ListItemDecorator>
            Users
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
