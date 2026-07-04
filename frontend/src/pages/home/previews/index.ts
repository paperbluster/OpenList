import { Component, lazy } from "solid-js"
import { me, getSettingBool } from "~/store"
import { Obj, ObjType } from "~/types"
import { ext } from "~/utils"
import { useRouter, useT } from "~/hooks"

type Ext = string[] | "*" | ((name: string) => boolean)
type Prior = boolean | (() => boolean)

const extsContains = (exts: Ext | undefined, name: string): boolean => {
  if (exts === undefined) {
    return false
  } else if (exts === "*") {
    return true
  } else if (typeof exts === "function") {
    return (exts as (name: string) => boolean)(name)
  } else {
    return (exts as string[]).includes(ext(name).toLowerCase())
  }
}

const isPrior = (p: Prior): boolean => {
  if (typeof p === "boolean") {
    return p
  }
  return p()
}

export interface Preview {
  key: string
  type?: ObjType
  exts?: Ext
  provider?: RegExp
  component: Component
  prior: Prior
  availableInArchive?: boolean
}

export interface PreviewComponent {
  key: string
  name: string
  component: Component
}

const previews: Preview[] = [
  {
    key: "html",
    exts: ["html"],
    component: lazy(() => import("./html")),
    prior: true,
  },
  {
    key: "image",
    type: ObjType.IMAGE,
    component: lazy(() => import("./image")),
    prior: true,
  },
  {
    key: "video",
    type: ObjType.VIDEO,
    component: lazy(() => import("./video")),
    prior: true,
  },
  {
    key: "audio",
    type: ObjType.AUDIO,
    component: lazy(() => import("./audio")),
    prior: true,
  },
]

export const getPreviews = (
  file: Obj & { provider: string },
): PreviewComponent[] => {
  const { searchParams, isShare } = useRouter()
  const t = useT()
  const typeOverride =
    ObjType[searchParams["type"]?.toUpperCase() as keyof typeof ObjType]
  const res: PreviewComponent[] = []
  const subsequent: PreviewComponent[] = []
  // For image/video/audio files, always prefer preview over download.
  // For other types, respect the user setting.
  const isMediaFile =
    file.type === ObjType.IMAGE ||
    file.type === ObjType.VIDEO ||
    file.type === ObjType.AUDIO
  const downloadPrior = isMediaFile
    ? false
    : (!isShare() && getSettingBool("preview_download_by_default")) ||
      (isShare() && getSettingBool("share_preview_download_by_default"))
  // internal previews
  if (!isShare() || getSettingBool("share_preview")) {
    previews.forEach((preview) => {
      if (preview.provider && !preview.provider.test(file.provider)) {
        return
      }
      if (
        preview.type === file.type ||
        (typeOverride && preview.type === typeOverride) ||
        extsContains(preview.exts, file.name)
      ) {
        const r = {
          key: preview.key,
          name: t(`home.preview.names.${preview.key}`),
          component: preview.component,
        }
        if (!downloadPrior && isPrior(preview.prior)) {
          res.push(r)
        } else {
          subsequent.push(r)
        }
      }
    })
  }

  // download page
  const downloadComponent: PreviewComponent = {
    key: "download",
    name: t("home.preview.names.download"),
    component: lazy(() => import("./download")),
  }
  if (res.length > 0 || subsequent.length > 0) {
    res.push(downloadComponent)
  } else {
    return [downloadComponent]
  }
  return res.concat(subsequent)
}
