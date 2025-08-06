import {getOsPathSeparator} from "@/events/event-registrar.ts";

export function getAllDirs(mask: string, activeDir: string[], paths: string[]): string[] {
  const osSeparator = getOsPathSeparator()
  const maskSegmentCount = mask.split(osSeparator).length
  const activeDirLen = activeDir.length;

  return [...new Set(paths
    .map(p => p.split(osSeparator))
    .filter(pArr => pArr.length > maskSegmentCount + activeDirLen)
    .filter(pArr => {
      if (activeDirLen === 0) {
        return true
      } else {
        const subPArr = pArr.slice(maskSegmentCount)
        return activeDir.every(segment => subPArr.includes(segment))
      }
    })
    .map(pArr => pArr[maskSegmentCount + activeDirLen])).values()]

}