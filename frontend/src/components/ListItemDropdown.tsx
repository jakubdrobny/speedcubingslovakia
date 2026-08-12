import { List, ListItemButton, ListItemDecorator, Tooltip } from "@mui/joy";

import { Link } from "react-router-dom";
import { NavContext } from "../context/NavContext";
import { NavContextType } from "../Types";
import { useContext } from "react";
import { ListItemDropdownContent, ListItemDropdownOption } from "../Types";

const ListItemDropdown: React.FC<{ content: ListItemDropdownContent }> = ({
  content,
}) => {
  const { closeNav } = useContext(NavContext) as NavContextType;

  return (
    <Tooltip
      variant="soft"
      arrow
      enterTouchDelay={0}
      leaveTouchDelay={5000}
      sx={{ pointerEvents: "auto" }}
      title={
        <List
          size="sm"
          onTouchStart={(e) => e.stopPropagation()}
          onTouchEnd={(e) => e.stopPropagation()}
        >
          {content.dropdownItems.map((option: ListItemDropdownOption) => (
            <ListItemButton
              component={Link}
              to={option.url}
              key={option.url}
              onClick={closeNav}
            >
              <ListItemDecorator>
                <option.icon />
              </ListItemDecorator>
              {option.text}
            </ListItemButton>
          ))}
        </List>
      }
    >
      <ListItemButton sx={{ justifyContent: "flex-end" }}>
        <ListItemDecorator>
          <content.itemIcon />
        </ListItemDecorator>
        {content.itemText}
      </ListItemButton>
    </Tooltip>
  );
};

export default ListItemDropdown;
