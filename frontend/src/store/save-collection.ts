import {create} from "zustand/react";
import {z} from "zod";
import {getAllDirs} from "@/lib/paths.ts";
import {getOsPathSeparator} from "@/events/event-registrar.ts";

export const GameSaveFileSchema = z.object({
  filename: z.string().nonempty("Filename cannot be empty"),
  savename: z.string(),
  pack_id: z.string(),
  pack_revision: z.int().min(0),
  directory: z.string().nonempty("Directory cannot be empty"),
  objectCount: z.number().min(0),
  uncachedResources: z.number().min(0),
  cachedRemoteResources: z.number().min(0),
  localResources: z.number().min(0),
  packedResources: z.number().min(0),
});

export type GameSaveFile = z.infer<typeof GameSaveFileSchema>;


export type LoadingState = "loading" | "done" | "error"

export type FileListState = {
  rootSaveDir: string
  activeDir: string[]
  childSegments: string[]
  allSavefiles: GameSaveFile[]
  activeSavefiles: GameSaveFile[]
  activeFile?: GameSaveFile
  loading: LoadingState
}


export type FileListActions = {
  addFile: (file: GameSaveFile) => void
  setActiveSegment: (segment: string | "..") => void
  setChildSegments: () => void
  retToParent: (times: number) => void
}

type FileListStore = FileListState & FileListActions

export const useSaveTableStore = create<FileListStore>((set) => ({
  loading: "done",
  rootSaveDir: "",
  activeDir: [],
  childSegments: [],
  allSavefiles: [],
  activeSavefiles: [],
  retToParent: (times: number) => {
    set((state) => ({
      activeDir: state.activeDir.slice(0, Math.max(0, state.activeDir.length - times)),
    }))
    set(updateChildSegments);
    set(updateActiveSavefiles);
  },
  setChildSegments: () => set(updateChildSegments),
  setActiveSegment: (segmentUpdate) => {
    set((state) => {
      let updatedActiveDir: string[];
      if (segmentUpdate === "..") {
        updatedActiveDir = state.activeDir.slice(0, state.activeDir.length - 1);
      } else {
        updatedActiveDir = [...state.activeDir, segmentUpdate];
      }

      return ({
        activeDir: updatedActiveDir,
      });

    });
    set(updateChildSegments);
    set(updateActiveSavefiles);
  },
  addFile: (file) => {

    set((state) => {
      const fileDirSegs = file.directory.split(getOsPathSeparator())
      let activeFilesUpdate = state.activeSavefiles;
      if (state.activeDir.every(seg => fileDirSegs.includes(seg))) {
        activeFilesUpdate = [...activeFilesUpdate, file];
      }
      return ({
        activeSavefiles: activeFilesUpdate,
        allSavefiles: [file, ...state.allSavefiles],
      });
    })
  },
}))

function updateChildSegments(state: FileListStore) {
  const allSegments = getAllDirs(state.rootSaveDir, state.activeDir, state.allSavefiles.map(f => f.directory));

  return ({childSegments: allSegments})
}

function updateActiveSavefiles(state: FileListStore) {
  const osPathSeparator = getOsPathSeparator();
  const activeUpdate = state.allSavefiles
    .filter(sf => {
      if (state.activeDir.length == 0) {
        return true;
      }
      return state.activeDir.every(seg => sf.directory.split(osPathSeparator).includes(seg));
    })

  return ({activeSavefiles: activeUpdate})
}
