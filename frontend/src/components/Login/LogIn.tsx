import { Alert, CircularProgress, Grid, Typography } from "@mui/joy";
import { AuthContextType, AuthState } from "../../Types";
import { getError, logIn } from "../../utils/utils";
import { useContext, useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

import { AuthContext } from "../../context/AuthContext";
import Cookies from "universal-cookie";

const LogIn = () => {
  const { setAuthState } = useContext(AuthContext) as AuthContextType;
  const [loadingState, setLoadingState] = useState<{
    loading: boolean;
    error: any;
  }>({ loading: false, error: "" });
  const [searchParams, _] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    setLoadingState({ loading: true, error: "" });
    logIn(searchParams)
      .then((res: AuthState) => {
        setAuthState(res);

        const cookies = new Cookies(null, { path: "/" });
        const locationPathname = cookies.get("backlink");
        cookies.remove("backlink");
        navigate(locationPathname || "/", { replace: true });
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
        width: "100%",
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          width: "100%",
        }}
      >
        {loadingState.loading || !loadingState.error ? (
          <>
            <CircularProgress />
            &nbsp; &nbsp; <Typography level="h3">Logging in...</Typography>
          </>
        ) : (
          <>
            <Alert color="danger" sx={{ gap: 0 }}>
              Oops. Something went wrong. Please&nbsp;
              <a href={import.meta.env.VITE_WCA_GET_CODE_URL || ""}>
                try again
              </a>
              .
            </Alert>
          </>
        )}
      </div>
    </Grid>
  );
};

export default LogIn;
