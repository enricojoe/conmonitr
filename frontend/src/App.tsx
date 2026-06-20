import { BrowserRouter, Route, Routes } from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import ContainerView from "./pages/ContainerView";
import Images from "./pages/Images";
import Volumes from "./pages/Volumes";
import Networks from "./pages/Networks";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/containers/:id" element={<ContainerView />} />
        <Route path="/images" element={<Images />} />
        <Route path="/volumes" element={<Volumes />} />
        <Route path="/networks" element={<Networks />} />
      </Routes>
    </BrowserRouter>
  );
}
