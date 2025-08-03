import type {FileListState} from "@/store/save-collection.ts";

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