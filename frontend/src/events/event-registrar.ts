import {type GameSaveFile, GameSaveFileSchema, useSaveTableStore} from "@/store/save-collection.ts";
import {EventsOn} from "../../wailsjs/runtime";
import {NewFileEvent, WailsPanicEvent} from "@/events/event-names.ts";
import {GetProperties, GetTabletopSaves} from "../../wailsjs/go/internal/CacheManagerApi";
import {properties} from "../../wailsjs/go/models.ts";
import {toast} from "sonner";
import ApplicationPropertiesView = properties.ApplicationPropertiesView;

let osSeparator = "/";

let alreadyRegistered = false; // only works because js is single-threaded

export function getOsPathSeparator() {
  return osSeparator;
}

async function setupProperties() {
  const properties = await GetProperties();
  osSeparator = properties.OsPathSeparator;
  useSaveTableStore.getState().setPackDataDir(properties.PackDataDir.replaceAll(properties.OsPathSeparator, "/"));
  useSaveTableStore.getState().rootSaveDir = `${properties.GameDir}${osSeparator}Saves`;
  return properties;
}

export async function ensureWailsInitialized() {
  if (alreadyRegistered) {
    return;
  } else {
    alreadyRegistered = true;
  }
  const properties = await setupProperties();
  setupListeners(properties);

  await GetTabletopSaves(properties.GameDir + osSeparator + 'Saves');
  await GetTabletopSaves(properties.PackDataDir);
  useSaveTableStore.getState().setChildSegments()
}

function setupListeners(properties: ApplicationPropertiesView) {
  EventsOn(NewFileEvent, (data: Partial<GameSaveFile>) => {
    const result = GameSaveFileSchema.safeParse(data);
    if (result.success) {
      const file = result.data;

      if (file.directory === properties.PackDataDir) {
        useSaveTableStore.getState().addPackDataFile(file);
      } else {
        useSaveTableStore.getState().addFile(file);
      }
    } else {
      console.error("todo: you did not setup toasts!")
    }
  });

  EventsOn(WailsPanicEvent, (error: unknown) => {
    console.error("Wails backend panic caught:", error);
    // If the error object contains structured fields, you can inspect them:

    if (isWailsPanicPayload(error)) {
      toast(`A backend panic occurred:\n\n${error.message}`);
    } else {
      toast("A backend panic occurred. Check console for details.");
    }
  });
}

interface WailsPanicPayload {
  message?: string;
  stack?: string;

  [key: string]: unknown;
}

function isWailsPanicPayload(value: unknown): value is WailsPanicPayload {
  return (
    typeof value === "object" &&
    value !== null &&
    "message" in value
  );
}

