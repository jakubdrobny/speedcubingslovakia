import { CircularProgress } from "@mui/joy";

import React from "react";

const LoadingComponent: React.FC<{ title: string }> = ({ title }) => {
  return (
    <div className="flex gap-4 text-2xl font-semibold justify-center items-center h-full">
      <CircularProgress />
      {title}
    </div>
  );
};
export default LoadingComponent;
