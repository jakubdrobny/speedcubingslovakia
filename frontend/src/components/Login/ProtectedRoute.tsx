import { CircularProgress, Grid, Typography } from "@mui/joy";
import { Navigate, Outlet } from "react-router-dom";

import { AuthContext } from "../../context/AuthContext";
import { AuthContextType } from "../../Types";
import { useContext } from "react";

const ProtectedRoute: React.FC<{
  loadingState: {
    loading: boolean;
    error: string;
  };
}> = ({ loadingState }) => {
  const { authState } = useContext(AuthContext) as AuthContextType;

  return loadingState.loading ? (
    <Grid
      container
      sx={{
        position: "absolute",
        top: "50%",
        left: "50%",
        transform: "translate(-50%, -50%)",
      }}
    >
      <Typography level="h3" sx={{ display: "flex", alignItems: "center" }}>
        <CircularProgress /> &nbsp; Authorizing...
      </Typography>
    </Grid>
  ) : !authState.isadmin ? (
    <Navigate to="/" />
  ) : (
    <Outlet />
  );
};

export default ProtectedRoute;
