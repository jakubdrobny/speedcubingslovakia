import { EmojiEvents, Language } from "@mui/icons-material";
import { List, ListItemButton, ListItemDecorator, Tooltip } from "@mui/joy";

import { Link } from "react-router-dom";
import { NavContext } from "../../context/NavContext";
import { NavContextType } from "../../Types";
import { useContext } from "react";
import WCALogoNoText from "../../images/WCALogoNoText";

const CompetitionsListItem = () => {
  const { closeNav } = useContext(NavContext) as NavContextType;

  return (
    <Tooltip
      variant="soft"
      title={
        <List size="sm">
          <ListItemButton
            component={Link}
            to="/competitions"
            onClick={closeNav}
          >
            <ListItemDecorator>
              <Language />
            </ListItemDecorator>
            Online Competitions
          </ListItemButton>
          <ListItemButton
            component={Link}
            to="/upcoming-wca-competitions"
            onClick={closeNav}
          >
            <ListItemDecorator>
              <WCALogoNoText />
            </ListItemDecorator>
            WCA Competitions
          </ListItemButton>
        </List>
      }
      arrow
      enterTouchDelay={0}
    >
      <ListItemButton sx={{ justifyContent: "flex-end" }}>
        <ListItemDecorator>
          <EmojiEvents />
        </ListItemDecorator>
        Competitions
      </ListItemButton>
    </Tooltip>
  );
};

export default CompetitionsListItem;
