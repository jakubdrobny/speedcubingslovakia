import { Alert, CircularProgress, Grid, Typography } from "@mui/joy";
import { AuthContextType, AuthState } from "../../Types";
import { Link, redirect, useNavigate, useSearchParams } from "react-router-dom";
import { useContext, useEffect, useState } from "react";

import { AuthContext } from "../../context/AuthContext";
import { logIn } from "../../utils";

const LogIn = () => {
  const { setAuthState } = useContext(AuthContext) as AuthContextType;
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState("");
  const [searchParams, _] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    setIsLoading(true);
    logIn(searchParams)
      .then((res: AuthState) => {
        setAuthState(res);
        navigate("/");
      })
      .catch((err) => {
        setIsLoading(false);
        setError(err.message);
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
      {isLoading ? (
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
