import { Card, Grid, Stack, Typography } from "@mui/joy";

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
            <Typography key={idx}>{content}</Typography>
          ))}
        </Stack>
      </Grid>
    </Card>
  );
};

export default ProfileBasicsDetailsSmall;
