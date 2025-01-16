import { AuthContextType, WindowSizeContextType } from "./Types";
import { Navigate, Route, Routes } from "react-router-dom";
import { authorizeAdmin, setBearerIfPresent } from "./utils/utils";
import { useContext, useEffect, Suspense, lazy } from "react";

import { AuthContext } from "./context/AuthContext";
import { WindowSizeContext } from "./context/WindowSizeContext";
import useState from "react-usestateref";
import LoadingComponent from "./components/Loading/LoadingComponent";

const CompetitionAnnouncements = lazy(
  () =>
    import("./components/CompetitionAnnouncements/CompetitionAnnouncements"),
);
const Announcement = lazy(
  () => import("./components/Announcement/Announcement"),
);
const AnnouncementEdit = lazy(
  () => import("./components/Dashboard/AnnouncementEdit"),
);
const Announcements = lazy(
  () => import("./components/Announcement/Announcements"),
);
const Competition = lazy(() => import("./components/Competition/Competition"));
const CompetitionEdit = lazy(
  () => import("./components/Dashboard/CompetitionEdit"),
);
const Competitions = lazy(
  () => import("./components/Competitions/Competitions"),
);
const Dashboard = lazy(() => import("./components/Dashboard/Dashboard"));
const Footer = lazy(() => import("./components/Footer/Footer"));
const Home = lazy(() => import("./components/Home/Home"));
const LogIn = lazy(() => import("./components/Login/LogIn"));
const ManageRoles = lazy(() => import("./components/Dashboard/ManageRoles"));
const NavHorizontal = lazy(() => import("./components/Nav/NavHorizontal"));
const NavVertical = lazy(() => import("./components/Nav/NavVertical"));
const NotFound = lazy(() => import("./components/NotFound/NotFound"));
const Profile = lazy(() => import("./components/Profile/Profile"));
const ProtectedRoute = lazy(() => import("./components/Login/ProtectedRoute"));
const Rankings = lazy(() => import("./components/Rankings/Rankings"));
const Records = lazy(() => import("./components/Records/Records"));
const ResultsEdit = lazy(() => import("./components/Dashboard/ResultsEdit"));
const Users = lazy(() => import("./components/Users/Users"));
const AdminStats = lazy(() => import("./components/Dashboard/AdminStats"));
const WCACompetitions = lazy(
  () => import("./components/WCACompetitions/WCACompetitions"),
);

const App = () => {
  const { authStateRef, setAuthState } = useContext(
    AuthContext,
  ) as AuthContextType;
  const { setWindowSize } = useContext(
    WindowSizeContext,
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
    <div className="grid grid-cols-12 items-center">
      <Suspense
        fallback={
          <div className="col-span-12 p-4 flex items-center justify-center">
            <LoadingComponent title="Loading..." />
          </div>
        }
      >
        <div className="col-span-0 lg:col-span-1 xl:col-span-2 border-0 border-solid lg:border-b-2 border-gray-200 w-full h-full" />
        <NavHorizontal />
        <div className="col-span-0 lg:col-span-1 xl:col-span-2 border-0 border-solid lg:border-b-2 border-gray-200 w-full h-full" />
        <NavVertical />

        <div className="col-span-0 lg:col-span-1 xl:col-span-2 border-0 border-b-2 lg:border-b-2 border-gray-200 w-full h-full" />
        <div className="col-span-12 lg:col-span-10 xl:col-span-8">
          <Routes>
            <Route path="/" Component={Home} />
            <Route path="/competitions" Component={Competitions} />
            <Route path="/competition/:id" Component={Competition} />
            <Route
              path="/upcoming-wca-competitions"
              Component={() => <Navigate to="/competitions/wca" replace />}
            />
            <Route path="/competitions/wca" Component={WCACompetitions} />
            <Route
              path="/competitions/announcements"
              Component={CompetitionAnnouncements}
            />
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
              <Route path="/admin/stats" Component={AdminStats} />
            </Route>
            <Route path="/profile/:id" Component={Profile} />
            <Route path="/results/users" Component={Users} />
            <Route path="/results/records" Component={Records} />
            <Route path="/results/rankings" Component={Rankings} />
            <Route path="*" element={<Navigate to="/not-found" replace />} />
          </Routes>
        </div>
        <div className="col-span-0 lg:col-span-1 xl:col-span-2" />
        <Footer />
      </Suspense>
    </div>
  );
};

export default App;
