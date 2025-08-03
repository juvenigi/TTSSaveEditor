import {create} from "zustand/react";
type FileState = {
  files: number;
};

type FileStateActions = {
  fetchFiles: () => void;
};

type FileStore = FileState & FileStateActions;

export const useFileStore = create<FileStore>((set) => ({
  files: 0,
  fetchFiles: () => set((state) => ({files: state.files + 1}))
}))