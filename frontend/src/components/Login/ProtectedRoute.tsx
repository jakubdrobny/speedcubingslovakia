import { Navigate, Outlet } from "react-router-dom";
import { useContext, useEffect } from "react";

import { AuthContext } from "../../context/AuthContext";
import { AuthContextType } from "../../Types";

const ProtectedRoute = () => {
  const { authState } = useContext(AuthContext) as AuthContextType;

  return !authState.isadmin ? <Navigate to="/" /> : <Outlet />;
};

export default ProtectedRoute;
