import { Box, Option, Select } from "@mui/joy";
import { RegionSelectGroup } from "../Types";

const RegionGroupSelect: React.FC<{
  regionValue: string;
  handleRegionChange: (newRegionValue: string) => void;
  regionGroups: RegionSelectGroup[];
  disabled: boolean;
}> = ({ regionValue, handleRegionChange, regionGroups, disabled }) => {
  return (
    <Select
      value={regionValue}
      onChange={(_, val) => {
        handleRegionChange(val || "");
      }}
      renderValue={(sel) => <Box sx={{ pl: 1 }}>{sel?.label}</Box>}
      sx={{ minWidth: "200px" }}
      disabled={disabled}
    >
      {regionGroups.map((regionGroup: RegionSelectGroup, idx: number) => (
        <div key={idx}>
          <Option value={regionGroup.groupName} disabled sx={{ pl: 2 }}>
            <b style={{ color: "black" }}>{regionGroup.groupName}</b>
          </Option>
          {regionGroup.groupMembers.map((groupMember: string, idx2: number) => (
            <Option
              key={idx2}
              value={regionGroup.groupName + "+" + groupMember}
              label={groupMember}
              sx={{ pl: 4 }}
              color="neutral"
            >
              {groupMember}
            </Option>
          ))}
        </div>
      ))}
    </Select>
  );
};

export default RegionGroupSelect;
