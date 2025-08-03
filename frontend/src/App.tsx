import './App.css'
import SaveFilePage from "@/pages/save-file-page.tsx";
import {registerWailsListeners} from "@/services/game-service.ts";
import {Button} from "@/components/ui/button.tsx";
import {GetTabletopSaves} from "../wailsjs/go/app/CacheManagerApi";
import {useEffect} from "react";

export default function App() {
  useEffect(() => {
    const off = registerWailsListeners();
    return () => {
      off()
    };
  }, []);
  
  return (
    <>
      <SaveFilePage></SaveFilePage> <Button
      onMouseDown={() => GetTabletopSaves().finally(() => console.log("fetch complete"))}>Click Me!</Button>
    </>
  )
}
