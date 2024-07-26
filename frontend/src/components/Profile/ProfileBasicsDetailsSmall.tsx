import { Card, Grid, Stack, Typography } from "@mui/joy";
import { ProfileTypeBasics, WindowSizeContextType } from "../../Types";

import { Link } from "react-router-dom";
import { WindowSizeContext } from "../../context/WindowSizeContext";
import { useContext } from "react";

const ProfileBasicsDetailsSmall: React.FC<{ basics: ProfileTypeBasics }> = ({
  basics,
}) => {
  const { windowSize } = useContext(WindowSizeContext) as WindowSizeContextType;
  const tooSmall = windowSize.width < 397;

  return (
    <Card
      sx={{
        width: "100%",
        display: "flex",
        flexDirection: "row",
        overflowX: "auto",
      }}
    >
      <Grid xs={tooSmall ? 7 : 6}>
        <Stack spacing={1}>
          {["Region", "WCA ID", "Sex", "Competitions", "Completed solves"].map(
            (title, idx) => (
              <Typography key={idx} style={{ whiteSpace: "nowrap" }}>
                <b>{title}</b>
              </Typography>
            )
          )}
        </Stack>
      </Grid>
      <Grid xs={tooSmall ? 5 : 6}>
        <Stack spacing={1}>
          <div style={{ whiteSpace: "nowrap" }}>
            <span className={`fi fi-${basics.region.iso2.toLowerCase()}`} />
            &nbsp;&nbsp;{basics.region.name}
          </div>
          {[
            basics.wcaid,
            basics.sex,
            basics.noOfCompetitions,
            basics.completedSolves,
          ].map((content, idx) => (
            <Typography key={idx} sx={{ whiteSpace: "nowrap" }}>
              {idx === 0 ? (
                <Link
                  to={`https://worldcubeassociation.org/persons/${content}`}
                  style={{ color: "#0B6BCB", textDecoration: "none" }}
                >
                  {content}
                </Link>
              ) : (
                content
              )}
            </Typography>
          ))}
        </Stack>
      </Grid>
    </Card>
  );
};

export default ProfileBasicsDetailsSmall;
