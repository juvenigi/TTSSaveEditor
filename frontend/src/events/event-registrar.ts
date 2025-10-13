import {type GameSaveFile, GameSaveFileSchema, useSaveTableStore} from "@/store/save-collection.ts";
import {EventsOn} from "../../wailsjs/runtime";
import {NewFileEvent, WailsPanicEvent} from "@/events/event-names.ts";
import {toast} from "sonner";

let alreadyRegistered = false; // only works because js is single-threaded

export async function ensureWailsInitialized() {
  if (alreadyRegistered) {
    return;
  } else {
    alreadyRegistered = true;
  }
  setupListeners();
  // todo: finalize
}

function setupListeners() {
  EventsOn(NewFileEvent, (data: Partial<GameSaveFile>) => {
    const result = GameSaveFileSchema.safeParse(data);
    if (result.success) {
      const file = result.data;

      if (file.directory === "todo") { //todo
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

