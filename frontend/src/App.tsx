import { AuthContextType, NavContextType } from "./Types";
import { Button, Grid, List, ListItemButton, Stack } from "@mui/joy";
import { Link, Navigate, Route, Routes } from "react-router-dom";
import { authorizeAdmin, setBearerIfPresent } from "./utils";
import { useContext, useEffect, useState } from "react";

import { AuthContext } from "./context/AuthContext";
import Competition from "./components/Competition/Competition";
import CompetitionEdit from "./components/Dashboard/CompetitionEdit";
import Competitions from "./components/Competitions/Competitions";
import Dashboard from "./components/Dashboard/Dashboard";
import { DensityMedium } from "@mui/icons-material";
import Home from "./components/Home/Home";
import LogIn from "./components/Login/LogIn";
import ManageRoles from "./components/Dashboard/ManageRoles";
import { NavContext } from "./context/NavContext";
import NavItems from "./components/NavItems";
import NotFound from "./components/NotFound/NotFound";
import ProtectedRoute from "./components/Login/ProtectedRoute";
import ResultsEdit from "./components/Dashboard/ResultsEdit";

const WIN_SMALL = 900;

const App = () => {
  const { authState, setAuthState } = useContext(
    AuthContext
  ) as AuthContextType;
  const [windowWidth, setWindowWidth] = useState<number>(window.innerWidth);
  const { navOpen, toggleNavOpen, closeNav } = useContext(
    NavContext
  ) as NavContextType;

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
    <Grid container>
      <Grid
        xs={0}
        md={1}
        lg={2}
        borderBottom={windowWidth < WIN_SMALL ? "" : "2px solid lightgrey"}
        width={"100%"}
      />
      <Grid
        xs={12}
        md={10}
        lg={8}
        borderBottom={"2px solid lightgrey"}
        width={"100%"}
      >
        <List
          style={{
            display: "flex",
            flexDirection: "row",
            padding: 20,
            width: "100%",
          }}
        >
          <Grid sx={{ display: "flex", justifyContent: "flex-start" }}>
            <ListItemButton component={Link} to="/" onClick={closeNav}>
              Speedcubing Slovakia
            </ListItemButton>
          </Grid>

          <Grid
            sx={{
              display: "flex",
              justifyContent: "flex-end",
              width: "100%",
            }}
          >
            {windowWidth < WIN_SMALL ? (
              <Button
                onClick={toggleNavOpen}
                variant={navOpen ? "solid" : "outlined"}
                color="neutral"
              >
                <DensityMedium />
              </Button>
            ) : (
              <NavItems />
            )}
          </Grid>
        </List>
      </Grid>
      <Grid
        xs={0}
        md={1}
        lg={2}
        borderBottom={windowWidth < WIN_SMALL ? "" : "2px solid lightgrey"}
        width={"100%"}
      />
      {windowWidth < WIN_SMALL && navOpen && (
        <Grid
          xs={12}
          flexDirection="column"
          borderBottom={"2px solid lightgrey"}
          padding="0.5em"
        >
          <List>
            <NavItems />
          </List>
        </Grid>
      )}
      <Grid xs={0} sm={1} md={2} />
      <Grid xs={12} sm={10} md={8}>
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
      <Grid xs={0} md={1} lg={2} />
    </Grid>
  );
};

export default App;
