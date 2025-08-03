import {create} from "zustand/react";

type DummyStore = {
  files: number;
} & {
  fetchFiles: () => void;
};

export const useDummyStore = create<DummyStore>((set) => ({
  files: 0,
  fetchFiles: () => set((state) => ({files: state.files + 1}))
}))

export type GameSaveFile = {
  filename: string
  objectCount: number
  uncachedResources: number
  cachedRemoteResources: number
  localResources: number
  packedResources: number
}

type LoadingState = "loading" | "done" | "error"

type FileListState = {
  files: GameSaveFile[]
  activeFile?: GameSaveFile
  loading: LoadingState
}

export const ListForTesting: FileListState = {
  files: [
    {
      filename: "project-assets.zip",
      objectCount: 42,
      uncachedResources: 5,
      cachedRemoteResources: 20,
      localResources: 10,
      packedResources: 7,
    },
    {
      filename: "textures.bundle",
      objectCount: 108,
      uncachedResources: 0,
      cachedRemoteResources: 85,
      localResources: 20,
      packedResources: 3,
    },
    {
      filename: "level-dataaa.json",
      objectCount: 12,
      uncachedResources: 2,
      cachedRemoteResources: 5,
      localResources: 3,
      packedResources: 2,
    },
  ],
  activeFile: {
    filename: "textures.bundle",
    objectCount: 108,
    uncachedResources: 0,
    cachedRemoteResources: 85,
    localResources: 20,
    packedResources: 3,
  },
  loading: "done",
};
// todo
// type FileListActions = []

// type FileListStore = FileListState & FileListActions

export const useSaveTableStore = create<FileListState>(() => ({
  ...ListForTesting
}))