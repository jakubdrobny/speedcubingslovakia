import { Grid, Typography } from "@mui/joy";

import ProfileBasicsDetailsSmall from "./ProfileBasicsDetailsSmall";
import ProfileBasicsDetailsTable from "./ProfileBasicsDetailsTable";
import { ProfileTypeBasics } from "../../Types";
import { WIN_SMALL } from "../../constants";
import { WindowSizeContext } from "../../context/WindowSizeContext";
import { WindowSizeContextType } from "../../Types";
import { useContext } from "react";

const ProfileBasics: React.FC<{ basics: ProfileTypeBasics }> = ({ basics }) => {
  const { windowSize, setWindowSize } = useContext(
    WindowSizeContext
  ) as WindowSizeContextType;

  return (
    <div
      style={{
        flexDirection: "column",
      }}
    >
      <Grid
        xs={12}
        sx={{
          display: "flex",
          justifyContent: "center",
        }}
      >
        <Typography level="h2">{basics.name}</Typography>
      </Grid>
      <Grid
        xs={12}
        sx={{
          display: "flex",
          justifyContent: "center",
          marginTop: "2em",
        }}
      >
        <img src={basics.imageurl} alt="profile_image" />
      </Grid>
      <Grid
        xs={12}
        sx={{
          display: "flex",
          justifyContent: "center",
          marginTop: "2em",
        }}
      >
        {windowSize.width < WIN_SMALL ? (
          <ProfileBasicsDetailsSmall
            basics={basics}
          ></ProfileBasicsDetailsSmall>
        ) : (
          <ProfileBasicsDetailsTable
            basics={basics}
          ></ProfileBasicsDetailsTable>
        )}
      </Grid>
    </div>
  );
};

export default ProfileBasics;
