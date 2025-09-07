import './App.css'
import SaveFilePage from "@/pages/save-file-page.tsx";
import {Toaster} from "@/components/ui/sonner.tsx";

export default function App() {
  return (
    <>
      <SaveFilePage></SaveFilePage>
      <Toaster/>
    </>
  )
}
