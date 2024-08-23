import { AuthContextType, WindowSizeContextType } from "./Types";
import { Navigate, Route, Routes } from "react-router-dom";
import { authorizeAdmin, setBearerIfPresent } from "./utils";
import { useContext, useEffect } from "react";

import Announcement from "./components/Announcement/Announcement";
import AnnouncementEdit from "./components/Dashboard/AnnouncementEdit";
import Announcements from "./components/Announcement/Announcements";
import { AuthContext } from "./context/AuthContext";
import Competition from "./components/Competition/Competition";
import CompetitionEdit from "./components/Dashboard/CompetitionEdit";
import Competitions from "./components/Competitions/Competitions";
import Dashboard from "./components/Dashboard/Dashboard";
import Footer from "./components/Footer/Footer";
import { Grid } from "@mui/joy";
import Home from "./components/Home/Home";
import LogIn from "./components/Login/LogIn";
import ManageRoles from "./components/Dashboard/ManageRoles";
import NavHorizontal from "./components/Nav/NavHorizontal";
import NavVertical from "./components/Nav/NavVertical";
import NotFound from "./components/NotFound/NotFound";
import Profile from "./components/Profile/Profile";
import ProtectedRoute from "./components/Login/ProtectedRoute";
import Rankings from "./components/Rankings/Rankings";
import Records from "./components/Records/Records";
import ResultsEdit from "./components/Dashboard/ResultsEdit";
import Users from "./components/Users/Users";
import { WIN_LG } from "./constants";
import { WindowSizeContext } from "./context/WindowSizeContext";
import useState from "react-usestateref";

const App = () => {
  const { authStateRef, setAuthState } = useContext(
    AuthContext
  ) as AuthContextType;
  const { windowSize, setWindowSize } = useContext(
    WindowSizeContext
  ) as WindowSizeContextType;
  useEffect(() => {
    const handleResize = () => {
      setWindowSize({ width: window.innerWidth, height: window.innerHeight });
    };

    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);
  const [authorizationLoadingState, setAuthorizationLoadingState] = useState<{
    loading: boolean;
    error: string;
  }>({ loading: authStateRef.current.token !== "", error: "" });

  useEffect(() => {
    setBearerIfPresent(authStateRef.current.token);

    if (authStateRef.current.token) {
      authorizeAdmin()
        .then((_) => {
          setAuthState({ ...authStateRef.current, isadmin: true });
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
        borderBottom={windowSize.width < WIN_LG ? "" : "2px solid lightgrey"}
        width={"100%"}
      />
      <NavHorizontal />
      <Grid
        xs={0}
        lg={1}
        xl={2}
        borderBottom={windowSize.width < WIN_LG ? "" : "2px solid lightgrey"}
        width={"100%"}
      />
      <NavVertical />
      <Grid xs={0} lg={1} xl={2} />
      <Grid xs={12} lg={10} xl={8}>
        <Routes>
          <Route path="/" Component={Home} />
          <Route path="/competitions" Component={Competitions} />
          <Route path="/competition/:id" Component={Competition} />
          <Route path="/announcements" Component={Announcements} />
          <Route path="/announcement/:id" Component={Announcement} />
          <Route path="/not-found" Component={NotFound} />
          <Route path="/login" Component={LogIn} />
          <Route
            element={
              <ProtectedRoute loadingState={authorizationLoadingState} />
            }
          >
            <Route
              path="/competition/:id/edit"
              element={<CompetitionEdit edit={true} />}
            />
            <Route
              path="/competition/create"
              element={<CompetitionEdit edit={false} />}
            />
            <Route path="/admin/dashboard" Component={Dashboard} />
            <Route path="/admin/manage-roles" Component={ManageRoles} />
            <Route path="/results/edit" Component={ResultsEdit} />
            <Route
              path="/announcement/:id/edit"
              element={<AnnouncementEdit edit={true} />}
            />
            <Route
              path="/announcement/create"
              element={<AnnouncementEdit edit={false} />}
            />
          </Route>
          <Route path="/profile/:id" Component={Profile} />
          <Route path="/results/users" Component={Users} />
          <Route path="/results/records" Component={Records} />
          <Route path="/results/rankings" Component={Rankings} />
          <Route path="*" element={<Navigate to="/not-found" replace />} />
        </Routes>
      </Grid>
      <Grid xs={0} lg={1} xl={2} />
      <Grid xs={12} sx={{ height: "8em" }} /> <Footer />
    </Grid>
  );
};

export default App;
