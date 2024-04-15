import { Grid } from "@mui/joy";
import { NavContext } from "../../context/NavContext";
import { NavContextType } from "../../Types";
import NavItems from "./NavItems";
import { WIN_SMALL } from "../../constants";
import { useContext } from "react";

const NavVertical: React.FC<{ windowWidth: number }> = ({ windowWidth }) => {
  const { navOpen } = useContext(NavContext) as NavContextType;
  const navReallyOpen = windowWidth < WIN_SMALL && navOpen;

  return (
    <Grid
      xs={12}
      flexDirection="column"
      borderBottom={"2px solid lightgrey"}
      padding="0.5em"
      sx={{
        display: navReallyOpen ? "flex" : "none",
      }}
    >
      <NavItems direction="column" windowWidth={windowWidth} />
    </Grid>
  );
};

export default NavVertical;
