import { AuthContextType, NavContextType } from "./Types";
import { Navigate, Route, Routes } from "react-router-dom";
import { WIN_LG, WIN_SMALL } from "./constants";
import { authorizeAdmin, setBearerIfPresent } from "./utils";
import { useContext, useEffect, useState } from "react";

import { AuthContext } from "./context/AuthContext";
import Competition from "./components/Competition/Competition";
import CompetitionEdit from "./components/Dashboard/CompetitionEdit";
import Competitions from "./components/Competitions/Competitions";
import Dashboard from "./components/Dashboard/Dashboard";
import { Grid } from "@mui/joy";
import Home from "./components/Home/Home";
import LogIn from "./components/Login/LogIn";
import ManageRoles from "./components/Dashboard/ManageRoles";
import { NavContext } from "./context/NavContext";
import NavHorizontal from "./components/Nav/NavHorizontal";
import NavVertical from "./components/Nav/NavVertical";
import NotFound from "./components/NotFound/NotFound";
import ProtectedRoute from "./components/Login/ProtectedRoute";
import ResultsEdit from "./components/Dashboard/ResultsEdit";

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

  useEffect(() => {
    setBearerIfPresent(authState.token);

    if (authState.token) {
      authorizeAdmin()
        .then((_) => setAuthState({ ...authState, isadmin: true }))
        .catch((_) => null);
    }
  }, []);

  return (
    <Grid container sx={{ background: "white" }}>
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
          <Route Component={ProtectedRoute}>
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
          <Route path="*" element={<Navigate to="/not-found" replace />} />
        </Routes>
      </Grid>
      <Grid xs={0} lg={1} xl={2} />
    </Grid>
  );
};

export default App;
