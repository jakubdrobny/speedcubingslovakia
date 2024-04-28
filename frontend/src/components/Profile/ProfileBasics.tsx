import { Grid, Stack, Typography } from "@mui/joy";
import { WIN_LG, WIN_SMALL, WIN_VERYSMALL } from "../../constants";
import { useEffect, useState } from "react";

import ProfileBasicsDetailsSmall from "./ProfileBasicsDetailsSmall";
import ProfileBasicsDetailsTable from "./ProfileBasicsDetailsTable";
import { ProfileTypeBasics } from "../../Types";

const ProfileBasics: React.FC<{ basics: ProfileTypeBasics }> = ({ basics }) => {
  const [windowWidth, setWindowWidth] = useState(window.innerWidth);

  useEffect(() => {
    const resizeListener = () => setWindowWidth(window.innerWidth);
    window.addEventListener("resize", resizeListener);
    return () => window.removeEventListener("resize", resizeListener);
  }, []);

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
        {windowWidth < WIN_SMALL ? (
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
