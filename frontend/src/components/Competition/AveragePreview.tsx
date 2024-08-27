import {
  AverageInfo,
  CompetitionContextType,
  ResponseError,
} from "../../Types";
import { Card, Chip, Stack, Typography } from "@mui/joy";
import {
  GetAverageInfo,
  getError,
  initialAverageInfo,
  isObjectEmpty,
  renderResponseError,
} from "../../utils";
import { useContext, useEffect, useState } from "react";

import { CompetitionContext } from "../../context/CompetitionContext";

const AveragePreview = () => {
  const [averageInfo, setAverageInfo] =
    useState<AverageInfo>(initialAverageInfo);
  const [error, setError] = useState<ResponseError>({});
  const { currentResults } = useContext(
    CompetitionContext
  ) as CompetitionContextType;

  useEffect(() => {
    GetAverageInfo(currentResults)
      .then((res) => setAverageInfo(res))
      .catch((err) => setError(getError(err)));
  }, []);

  return (
    <Stack
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        overflow: "auto",
      }}
      component={Card}
    >
      <h3 style={{ padding: 0, margin: 0, borderBottom: "1px solid #CDD7E1" }}>
        Average preview:
      </h3>
      {!isObjectEmpty(error) ? (
        renderResponseError(error)
      ) : (
        <>
          <Stack
            spacing={2}
            direction="row"
            display="flex"
            justifyContent="center"
          >
            <div
              style={
                averageInfo.showPossibleAverage
                  ? {
                      display: "flex",
                      justifyContent: "flex-end",
                    }
                  : {}
              }
            >
              <Typography sx={{ display: "flex", alignItems: "center" }}>
                <b>Single:</b>&nbsp;
                <Chip variant="soft" color="primary">
                  {averageInfo.single}
                </Chip>
              </Typography>
            </div>
            {averageInfo.finished && (
              <div>
                <Typography sx={{ display: "flex", alignItems: "center" }}>
                  <b>Average:</b>&nbsp;
                  <Chip variant="soft" color="primary">
                    {averageInfo.average}
                  </Chip>
                </Typography>
              </div>
            )}
          </Stack>
          <div
            style={{
              display: "flex",
              justifyContent: "center",
            }}
          >
            <Typography sx={{ display: "flex", alignItems: "center" }}>
              <b>Times:</b>&nbsp;
              <Chip variant="soft" color="warning">
                <Stack spacing={1} direction="row">
                  {averageInfo.times.map((solveTime, idx) => (
                    <div key={idx.toString() + "#" + solveTime}>
                      {solveTime}
                    </div>
                  ))}
                </Stack>
              </Chip>
            </Typography>
          </div>
          {averageInfo.showPossibleAverage && (
            <Stack
              spacing={2}
              direction="row"
              display="flex"
              justifyContent="center"
            >
              <div style={{ display: "flex", justifyContent: "flex-end" }}>
                <Typography sx={{ display: "flex", alignItems: "center" }}>
                  <b>BPA:</b>&nbsp;
                  <Chip variant="soft" color="success">
                    {averageInfo.bpa}
                  </Chip>
                </Typography>
              </div>
              <div>
                <Typography sx={{ display: "flex", alignItems: "center" }}>
                  <b>WPA:</b>&nbsp;
                  <Chip variant="soft" color="danger">
                    {averageInfo.wpa}
                  </Chip>
                </Typography>
              </div>
            </Stack>
          )}
        </>
      )}
    </Stack>
  );
};

export default AveragePreview;
