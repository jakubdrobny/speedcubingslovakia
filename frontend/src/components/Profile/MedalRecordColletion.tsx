import { Card, Stack, Table, Typography } from "@mui/joy";

import React from "react";

const MedalRecordColletion: React.FC<{
  title: string;
  headers: string[];
  values: number[];
}> = ({ title, headers, values }) => {
  const center: React.CSSProperties = { textAlign: "center", height: "1em" };

  return (
    <Stack spacing={2}>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <Typography level="h3">{title}</Typography>
      </div>
      <Card sx={{ padding: "0.25em 0.5em" }}>
        <Table>
          <thead>
            <tr>
              {headers.map((columnTitle) => (
                <th style={center}>
                  <b>{columnTitle}</b>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            <tr>
              {values.map((columnContent) => (
                <td style={center}>{columnContent}</td>
              ))}
            </tr>
          </tbody>
        </Table>
      </Card>
    </Stack>
  );
};

export default MedalRecordColletion;
