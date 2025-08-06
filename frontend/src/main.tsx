import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {ensureWailsInitialized} from "@/events/event-registrar.ts";

await ensureWailsInitialized()

createRoot(document.getElementById('root')!).render(
  <StrictMode> <App/> </StrictMode>,
)
