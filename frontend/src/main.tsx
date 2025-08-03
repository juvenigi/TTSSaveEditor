import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {ensureWailsListenersRegistered} from "@/events/event-registrar.ts";

ensureWailsListenersRegistered()

createRoot(document.getElementById('root')!).render(
  <StrictMode> <App/> </StrictMode>,
)
