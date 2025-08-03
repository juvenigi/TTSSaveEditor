import {type GameSaveFile, GameSaveFileSchema, useSaveTableStore} from "@/store/save-collection.ts";
import {EventsOn} from "../../wailsjs/runtime";
import {NewFileEvent} from "@/events/event-names.ts";


let alreadyRegistered = false; // only works because js is single-threaded

export function ensureWailsListenersRegistered() {
  if (alreadyRegistered) return;
  alreadyRegistered = true;

  EventsOn(NewFileEvent, (data: Partial<GameSaveFile>) => {
    const result = GameSaveFileSchema.safeParse(data);
    if (result.success) {
      useSaveTableStore.getState().addFile(result.data);
    } else {
      console.error("todo: you did not setup toasts!")
    }
  });
}
