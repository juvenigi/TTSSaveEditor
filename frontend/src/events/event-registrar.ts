import {type GameSaveFile, GameSaveFileSchema, useSaveTableStore} from "@/store/save-collection.ts";
import {EventsOn} from "../../wailsjs/runtime";
import {NewFileEvent} from "@/events/event-names.ts";
import {GetProperties, GetTabletopSaves} from "../../wailsjs/go/internal/CacheManagerApi";
import {properties} from "../../wailsjs/go/models.ts";
import ApplicationPropertiesView = properties.ApplicationPropertiesView;

let osSeparator = "/";

let alreadyRegistered = false; // only works because js is single-threaded

export function getOsPathSeparator() {
  return osSeparator;
}


export async function ensureWailsInitialized() {
  if (alreadyRegistered) {
    return;
  } else {
    alreadyRegistered = true;
  }

  const properties = await GetProperties();
  osSeparator = properties.OsPathSeparator;
  useSaveTableStore.getState().rootSaveDir = `${properties.GameDir}${osSeparator}Saves`;
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
}

