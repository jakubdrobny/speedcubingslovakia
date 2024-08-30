import { CircularProgress, Typography } from "@mui/joy";

import React from "react";

const LoadingComponent: React.FC<{ title: string }> = ({ title }) => {
  return (
    <Typography
      level="h3"
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100%",
      }}
    >
      <CircularProgress /> &nbsp; {title}
    </Typography>
  );
};
export default LoadingComponent;
