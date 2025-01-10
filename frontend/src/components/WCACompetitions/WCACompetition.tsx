import { Button, Card, Chip, Divider, Stack, Typography } from "@mui/joy";
import React from "react";
import { Link } from "react-router-dom";
import { CompetitionEvent, WCACompetitionType } from "../../Types";
import {
  getCubingIconClassName,
  renderUpcomingWCACompetitionDateRange,
} from "../../utils/utils";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
dayjs.extend(relativeTime);

const WCACompetition: React.FC<{ comp: WCACompetitionType }> = ({ comp }) => {
  return (
    <Stack component={Card} direction="column">
      <Typography level="h3">{comp.name}</Typography>
      <Divider />
      <Typography>
        <b>Place:</b>&nbsp;{comp.venueAddress}
      </Typography>
      <Typography>
        <b>Date:</b>&nbsp;
        {renderUpcomingWCACompetitionDateRange(comp.startdate, comp.enddate)}
      </Typography>
      {dayjs().isBefore(dayjs(comp.registrationOpen)) ? (
        <Stack spacing={1} direction="row">
          <Typography>
            <b>Registration opens:</b>
          </Typography>
          <Typography>
            {new Date(comp.registrationOpen).toLocaleDateString() +
              " " +
              new Date(comp.registrationOpen).toLocaleTimeString()}
          </Typography>
          <Chip color="warning">{dayjs(comp.registrationOpen).fromNow()}</Chip>
        </Stack>
      ) : (
        <Stack spacing={1} direction="row">
          <Typography>
            <b>Competitors:</b>
          </Typography>
          <Typography>
            {comp.registered + "/" + comp.competitorLimit}
          </Typography>
          <Chip
            color={
              comp.registered === comp.competitorLimit ? "danger" : "success"
            }
          >
            {comp.registered === comp.competitorLimit
              ? "Full"
              : (comp.competitorLimit - comp.registered).toString() +
              " spot" +
              (comp.competitorLimit - comp.registered > 1 ? "s" : "") +
              " remaining"}
          </Chip>
        </Stack>
      )}
      <Stack direction="row" alignItems="center" flexWrap="wrap" spacing={1}>
        <Typography>
          <b>Events:</b>
        </Typography>
        {comp.events.map((event: CompetitionEvent, idx2: number) => (
          <span
            key={idx2 + 100000}
            className={`${getCubingIconClassName(
              event.iconcode,
            )} profile-cubing-icon-mock`}
          />
        ))}
      </Stack>
      <Divider />
      <Button
        sx={{ float: "right" }}
        variant="outlined"
        component={Link}
        to={comp.url}
      >
        More info!
      </Button>
    </Stack>
  );
};

export default WCACompetition;
