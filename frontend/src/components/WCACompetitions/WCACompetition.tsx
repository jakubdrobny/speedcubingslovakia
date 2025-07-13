import { Button, Card, Chip, Divider, Stack, Typography } from "@mui/joy";
import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { CompetitionEvent, WCACompetitionType } from "../../Types";
import {
  getCubingIconClassName,
  RenderUpcomingWCACompetitionDateRange,
} from "../../utils/utils";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
dayjs.extend(relativeTime);

const wcaLiveAPIBody = (name: string) => {
  return {
    operationName: "Competitions",
    variables: {
      from: "2025-06-13",
      name: name,
    },
    query: `query Competitions($from: Date!, $name: String!) {\n  competitions(from: $from, filter: $name) {\n    id\n    name\n    startDate\n    endDate\n    startTime\n    endTime\n    venues { \n      id\n      country { \n        iso2\n        __typename\n }\n
__typename\n    }\n    __typename\n  }\n  recentRecords { \n    ...records\n    __typename\n } \n}\n\nfragment records on Record {
\n  id\n
  tag\n  type\n  attemptResult\n  result {
  \n    id\n    person {
  \n      id\n      name\n      country {
  \n        iso2\n        name\n
    __typename\n
  } \n      __typename\n
  } \n    round {
    \n      id\n      competitionEvent {
    \n        id\n        event {
    \n
      id\n          name\n          __typename\n
    } \n        competition { \n          id\n          __typename\n } \n        __typename\n
    } \n      __typename\n
    } \n    __typename\n
  } \n  __typename\n
}`
  }
};

const WCACompetition: React.FC<{ comp: WCACompetitionType }> = ({ comp }) => {
  const [wcaLiveId, setWcaLiveId] = useState("");
  const isLive = dayjs(comp.startdate).isBefore(dayjs());

  useEffect(() => {
    const fetchWCALiveCompInfo = async () => {
      fetch('https://live.worldcubeassociation.org/api', {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(wcaLiveAPIBody(comp.name)),
      }).then(response => {
        return response.json();
      }).then(data => {
        if (!data || !data.data) return;
        if (!data.data.competitions) {
          console.log('competitions are not in the response');
          return;
        }
        if (data.data.competitions.length === 0) {
          console.log('competition not found: ' + comp.name)
          return;
        }
        const wcaLiveCompetition = data.data.competitions[0];
        setWcaLiveId(wcaLiveCompetition.id)
      }).catch(err => {
        console.log(err)
      });

    };

    fetchWCALiveCompInfo();
  }, []);

  return (
    <Stack component={Card} direction="column">
      <Stack spacing={1} direction="row">
        {isLive && (
          <Chip color="success">Live!</Chip>
        )}
        <Typography level="h3">{comp.name}</Typography>
      </Stack>
      <Divider />
      <Typography>
        <b>Place:</b>&nbsp;{comp.venueAddress}
      </Typography>
      <Typography>
        <b>Date:</b>&nbsp;
        {RenderUpcomingWCACompetitionDateRange(comp.startdate, comp.enddate)}
      </Typography>
      {dayjs().isBefore(dayjs(comp.registrationOpen)) ? (
        <>
          <Typography>
            <b>Competitor limit:</b>&nbsp;{comp.competitorLimit}
          </Typography>
          <Stack spacing={1} direction="row">
            <Typography>
              <b>Registration opens:</b>
            </Typography>
            <Typography>
              {new Date(comp.registrationOpen).toLocaleDateString() +
                " " +
                new Date(comp.registrationOpen).toLocaleTimeString()}
            </Typography>
            <Chip color="warning">
              {dayjs(comp.registrationOpen).fromNow()}
            </Chip>
          </Stack>
        </>
      ) : (
        <>
          {comp.competitorLimit !== 0 && (
            <Stack spacing={1} direction="row">
              <Typography>
                <b>Competitors:</b>
              </Typography>
              <Typography>
                {comp.registered + "/" + comp.competitorLimit}
              </Typography>
              {dayjs().isBefore(dayjs(comp.registrationClose)) && (
                <Chip
                  color={
                    comp.registered === comp.competitorLimit
                      ? "danger"
                      : "success"
                  }
                >
                  {comp.registered === comp.competitorLimit
                    ? "Full"
                    : (comp.competitorLimit - comp.registered).toString() +
                    " spot" +
                    (comp.competitorLimit - comp.registered > 1 ? "s" : "") +
                    " remaining"}
                </Chip>
              )}
            </Stack>
          )}
          <Stack spacing={1} direction="row">
            <Typography>
              <b>Registration: </b>
            </Typography>
            {dayjs().isBefore(dayjs(comp.registrationClose)) ? (
              <Chip color="warning">
                {"Closes " + dayjs(comp.registrationClose).fromNow()}
              </Chip>
            ) : (
              <Chip color="danger">Closed</Chip>
            )}
          </Stack>
        </>
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
      {isLive && (
        <Typography>
          <b>Group assignments:</b>&nbsp;
          {
            <Link
              to={"https://competitiongroups.com/competitions/" + comp.id}
              style={{ color: "#0B6BCB", textDecoration: "none" }}
            >
              Link
            </Link>
          }
        </Typography>
      )}
      {isLive && (
        <Typography>
          <b>Live results:</b>&nbsp;
          {
            <Link
              to={"https://live.worldcubeassociation.com/competitions/" + wcaLiveId}
              style={{ color: "#0B6BCB", textDecoration: "none" }}
            >
              Link
            </Link>
          }
        </Typography>
      )}
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
