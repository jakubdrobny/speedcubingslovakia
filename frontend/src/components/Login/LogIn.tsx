import { Alert, CircularProgress, Grid, Typography } from "@mui/joy";
import { AuthContextType, AuthState } from "../../Types";
import { getError, logIn } from "../../utils";
import { useContext, useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

import { AuthContext } from "../../context/AuthContext";

const LogIn = () => {
  const { setAuthState } = useContext(AuthContext) as AuthContextType;
  const [loadingState, setLoadingState] = useState<{
    loading: boolean;
    error: string;
  }>({ loading: false, error: "" });
  const [searchParams, _] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    setLoadingState({ loading: true, error: "" });
    logIn(searchParams)
      .then((res: AuthState) => {
        setAuthState(res);
        navigate("/", { replace: true });
      })
      .catch((err) => {
        setLoadingState({ loading: false, error: getError(err) });
      });
  }, []);

  return (
    <Grid
      container
      sx={{
        position: "absolute",
        top: "50%",
        left: "50%",
        transform: "translate(-50%, -50%)",
      }}
    >
      {loadingState.loading ? (
        <>
          <Typography level="h3" sx={{ display: "flex", alignItems: "center" }}>
            <CircularProgress /> &nbsp; Logging in...
          </Typography>
        </>
      ) : (
        <>
          <Alert color="danger" sx={{ gap: 0 }}>
            Oops. Something went wrong. Please&nbsp;
            <a href={process.env.REACT_APP_WCA_GET_CODE_URL || ""}>try again</a>
            .
          </Alert>
        </>
      )}
    </Grid>
  );
};

export default LogIn;
