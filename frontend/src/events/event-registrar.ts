import {type GameSaveFile, GameSaveFileSchema, useSaveTableStore} from "@/store/save-collection.ts";
import {EventsOn} from "../../wailsjs/runtime";
import {NewFileEvent} from "@/events/event-names.ts";
import {GetProperties} from "../../wailsjs/go/internal/CacheManagerApi";

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

  EventsOn(NewFileEvent, (data: Partial<GameSaveFile>) => {
    const result = GameSaveFileSchema.safeParse(data);
    if (result.success) {
      const file = result.data;
      useSaveTableStore.getState().addFile(file);
    } else {
      console.error("todo: you did not setup toasts!")
    }
  });
}
