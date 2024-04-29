import { Navigate, Route, Routes } from "react-router-dom";
import { authorizeAdmin, setBearerIfPresent } from "./utils";
import { useContext, useEffect } from "react";

import { AuthContext } from "./context/AuthContext";
import { AuthContextType } from "./Types";
import Competition from "./components/Competition/Competition";
import CompetitionEdit from "./components/Dashboard/CompetitionEdit";
import Competitions from "./components/Competitions/Competitions";
import Dashboard from "./components/Dashboard/Dashboard";
import { Grid } from "@mui/joy";
import Home from "./components/Home/Home";
import LogIn from "./components/Login/LogIn";
import ManageRoles from "./components/Dashboard/ManageRoles";
import NavHorizontal from "./components/Nav/NavHorizontal";
import NavVertical from "./components/Nav/NavVertical";
import NotFound from "./components/NotFound/NotFound";
import Profile from "./components/Profile/Profile";
import ProtectedRoute from "./components/Login/ProtectedRoute";
import ResultsEdit from "./components/Dashboard/ResultsEdit";
import Users from "./components/Users/Users";
import { WIN_LG } from "./constants";
import useState from "react-usestateref";

const App = () => {
  const { authState, setAuthState } = useContext(
    AuthContext
  ) as AuthContextType;
  const [windowWidth, setWindowWidth] = useState<number>(window.innerWidth);
  useEffect(() => {
    const handleResize = () => {
      setWindowWidth(window.innerWidth);
    };

    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);
  const [authorizationLoadingState, setAuthorizationLoadingState] = useState<{
    loading: boolean;
    error: string;
  }>({ loading: authState.token != "", error: "" });

  useEffect(() => {
    setBearerIfPresent(authState.token);

    if (authState.token) {
      authorizeAdmin()
        .then((_) => {
          setAuthState({ ...authState, isadmin: true });
          setAuthorizationLoadingState((ps) => ({ ...ps, loading: false }));
        })
        .catch((_) => {
          setAuthorizationLoadingState({ loading: false, error: "" });
        });
    } else {
      setAuthorizationLoadingState({ loading: false, error: "" });
    }
  }, []);

  return (
    <Grid container>
      <Grid
        xs={0}
        lg={1}
        xl={2}
        borderBottom={windowWidth < WIN_LG ? "" : "2px solid lightgrey"}
        width={"100%"}
      />
      <NavHorizontal windowWidth={windowWidth} />
      <Grid
        xs={0}
        lg={1}
        xl={2}
        borderBottom={windowWidth < WIN_LG ? "" : "2px solid lightgrey"}
        width={"100%"}
      />
      <NavVertical windowWidth={windowWidth} />
      <Grid xs={0} lg={1} xl={2} />
      <Grid xs={12} lg={10} xl={8}>
        <Routes>
          <Route path="/" Component={Home} />
          <Route path="/competitions" Component={Competitions} />
          <Route path="/competition/:id" Component={Competition} />
          <Route path="/not-found" Component={NotFound} />
          <Route path="/login" Component={LogIn} />
          <Route
            Component={() => (
              <ProtectedRoute loadingState={authorizationLoadingState} />
            )}
          >
            <Route
              path="/competition/:id/edit"
              Component={() => <CompetitionEdit edit={true} />}
            />
            <Route
              path="/competition/create"
              Component={() => <CompetitionEdit edit={false} />}
            />
            <Route path="/admin/dashboard" Component={Dashboard} />
            <Route path="/admin/manage-roles" Component={ManageRoles} />
            <Route path="/results/edit" Component={ResultsEdit} />
          </Route>
          <Route path="/profile/:id" Component={Profile} />
          <Route path="/users" Component={Users} />
          <Route path="*" element={<Navigate to="/not-found" replace />} />
        </Routes>
      </Grid>
      <Grid xs={0} lg={1} xl={2} />
      <Grid xs={12} sx={{ height: "3em" }} />{" "}
      {/* padding at the bottom of page */}
    </Grid>
  );
};

export default App;
