import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {ensureWailsInitialized} from "@/events/event-registrar.ts";
import {GetTabletopSaves} from "../wailsjs/go/internal/CacheManagerApi";
import {useSaveTableStore} from "@/store/save-collection.ts";

await ensureWailsInitialized()
await GetTabletopSaves()
useSaveTableStore.getState().setChildSegments()

createRoot(document.getElementById('root')!).render(
  <StrictMode> <App/> </StrictMode>,
)
