import {
  Dropdown,
  ListItemButton,
  ListItemDecorator,
  Menu,
  MenuButton,
  MenuItem,
} from "@mui/joy";

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
    <Dropdown>
      <MenuButton
        slots={{ root: ListItemButton }}
        slotProps={{ root: { sx: { justifyContent: "flex-end" } } }}
      >
        <ListItemDecorator>
          <content.itemIcon />
        </ListItemDecorator>
        {content.itemText}
      </MenuButton>
      <Menu variant="soft" size="sm" placement="bottom-end">
        {content.dropdownItems.map((option: ListItemDropdownOption) => (
          <MenuItem component={Link} to={option.url} onClick={closeNav}>
            <ListItemDecorator>
              <option.icon />
            </ListItemDecorator>
            {option.text}
          </MenuItem>
        ))}
      </Menu>
    </Dropdown>
  );
};

export default ListItemDropdown;
