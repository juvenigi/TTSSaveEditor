
// listenerSingleton.ts
import {NewFileEvent} from "@/services/game-service.ts";
import {type GameSaveFile, GameSaveFileSchema, useSaveTableStore} from "@/store/save-collection.ts";
import {EventsOn} from "../../wailsjs/runtime";

let alreadyRegistered = false;

export function ensureWailsListenersRegistered() {
  if (alreadyRegistered) return;
  alreadyRegistered = true;

  EventsOn(NewFileEvent, (data: Partial<GameSaveFile>) => {
    const result = GameSaveFileSchema.safeParse(data);
    if (result.success) {
      console.log("adding savefile");
      useSaveTableStore.getState().addFile(result.data);
    }
  });
}
