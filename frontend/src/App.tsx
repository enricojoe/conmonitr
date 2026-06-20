import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { AuthProvider, useAuth } from "./context/AuthContext";
import Dashboard from "./pages/Dashboard";
import ContainerView from "./pages/ContainerView";
import Images from "./pages/Images";
import Login from "./pages/Login";
import Networks from "./pages/Networks";
import Volumes from "./pages/Volumes";

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { token } = useAuth();
  return token ? <>{children}</> : <Navigate to="/login" replace />;
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/containers/:id"
            element={
              <ProtectedRoute>
                <ContainerView />
              </ProtectedRoute>
            }
          />
          <Route
            path="/images"
            element={
              <ProtectedRoute>
                <Images />
              </ProtectedRoute>
            }
          />
          <Route
            path="/volumes"
            element={
              <ProtectedRoute>
                <Volumes />
              </ProtectedRoute>
            }
          />
          <Route
            path="/networks"
            element={
              <ProtectedRoute>
                <Networks />
              </ProtectedRoute>
            }
          />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}
