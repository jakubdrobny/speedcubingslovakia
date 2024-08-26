import { Card, Chip, Grid, Stack, Typography } from "@mui/joy";
import { useEffect, useState } from "react";

import { AverageInfo } from "../../Types";
import { initialAverageInfo } from "../../utils";

const AveragePreview = () => {
  const [averageInfo, setAverageInfo] =
    useState<AverageInfo>(initialAverageInfo);

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
      <Stack spacing={2} direction="row" display="flex" justifyContent="center">
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
                <div key={idx.toString() + "#" + solveTime}>{solveTime}</div>
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
    </Stack>
  );
};

export default AveragePreview;
