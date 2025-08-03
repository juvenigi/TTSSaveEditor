import {create} from "zustand/react";
import {z} from "zod";

export const GameSaveFileSchema = z.object({
  filename: z.string().nonempty("Filename cannot be empty"),
  objectCount: z.number().min(0),
  uncachedResources: z.number().min(0),
  cachedRemoteResources: z.number().min(0),
  localResources: z.number().min(0),
  packedResources: z.number().min(0),
});

export type GameSaveFile = z.infer<typeof GameSaveFileSchema>;


export type LoadingState = "loading" | "done" | "error"

export type FileListState = {
  files: GameSaveFile[]
  activeFile?: GameSaveFile
  loading: LoadingState
}


export type FileListActions = {
  addFile: (file: GameSaveFile) => void
}

type FileListStore = FileListState & FileListActions

export const useSaveTableStore = create<FileListStore>((set) => ({
  loading: "done",
  files: [],
  addFile: (file) =>
    set((state) => ({
      files: [file, ...state.files],
    })),
}))