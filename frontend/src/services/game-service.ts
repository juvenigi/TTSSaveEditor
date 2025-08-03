import {EventsOn} from "../../wailsjs/runtime";
import {type GameSaveFile, GameSaveFileSchema, useSaveTableStore} from "@/store/save-collection.ts";

export const NewFileEvent = "ttsc:newFile"

export function registerWailsListeners() {
  const unsubscribe = EventsOn(NewFileEvent, (data: Partial<GameSaveFile>) => {
    const result = GameSaveFileSchema.safeParse(data)
    if (result.success) {
      console.log("adding savefile")
      useSaveTableStore.getState().addFile(result.data);
      return;
    }

    console.error("todo: you did not setup toasts yet")
  });
  return () => {
    unsubscribe();
  }
}