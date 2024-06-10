import { Card, Grid, Stack, Typography } from "@mui/joy";

import { Link } from "react-router-dom";
import { ProfileTypeBasics } from "../../Types";

const ProfileBasicsDetailsSmall: React.FC<{ basics: ProfileTypeBasics }> = ({
  basics,
}) => {
  return (
    <Card sx={{ width: "100%", display: "flex", flexDirection: "row" }}>
      <Grid xs={6}>
        <Stack spacing={1}>
          {["Region", "WCA ID", "Sex", "Competitions", "Completed solves"].map(
            (title, idx) => (
              <Typography key={idx}>
                <b>{title}</b>
              </Typography>
            )
          )}
        </Stack>
      </Grid>
      <Grid xs={6}>
        <Stack spacing={1}>
          <div>
            <span className={`fi fi-${basics.region.iso2.toLowerCase()}`} />
            &nbsp;&nbsp;{basics.region.name}
          </div>
          {[
            basics.wcaid,
            basics.sex,
            basics.noOfCompetitions,
            basics.completedSolves,
          ].map((content, idx) => (
            <Typography key={idx}>
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
