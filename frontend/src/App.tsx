import './App.css'
import SaveFilePage from "@/pages/save-file-page.tsx";
import {Toaster} from "@/components/ui/sonner.tsx";
import {HashRouter, Route, Routes} from "react-router";
import PackDataSidebar from "@/pages/pack-data-sidebar.tsx";
import MainNavMenu from "@/pages/components/main-nav-menu.tsx";

export default function App() {
  return (
    <>
      <HashRouter>
        <MainNavMenu/>
        <Routes>
          <Route path="/" element={<SaveFilePage/>}/>
          <Route path="/pack-data" element={<PackDataSidebar/>}/>
        </Routes>
      </HashRouter>
      <Toaster/>
    </>
  )
}
