import axios from "axios"
import {
  appendObjs,
  password,
  ObjStore,
  State,
  getPagination,
  objStore,
  me,
  shouldKeepState,
} from "~/store"
import { fsGet, fsList, notify } from "~/utils"
import { useRouter } from "./useRouter"

let first_fetch = true
let globalPage = 1

export const getGlobalPage = () => globalPage
export const setGlobalPage = (p: number) => { globalPage = p }
export const resetGlobalPage = () => { setGlobalPage(1) }

// ── Simple file browser hook ──────────────────────────────────────────
// Design: URL pathname() is the sole source of truth.
// The server response tells us if a path is a file or folder.
// No IsDirRecord cache. No pushHref pre-marking. No stale async races.
// ──────────────────────────────────────────────────────────────────────

export const usePath = () => {
  const { pathname, to, searchParams } = useRouter()

  // Each new navigation bumps this; async callbacks discard stale results.
  let opId = 0

  // Cancellers for in-flight axios requests
  let cancelGet: axios.Canceler | undefined
  let cancelList: axios.Canceler | undefined

  const pagination = getPagination()
  if (pagination.type === "pagination") {
    setGlobalPage(parseInt(searchParams["page"]) || 1)
  }

  // ── api helpers ────────────────────────────────────────────────────

  const apiGet = (path: string) =>
    fsGet(path, password(), new axios.CancelToken(c => { cancelGet = c }))

  const apiList = (path: string, index?: number, size?: number, force?: boolean) =>
    fsList(path, password(), index, size, force, new axios.CancelToken(c => { cancelList = c }))

  // ── error handling ─────────────────────────────────────────────────

  let retryPass = false

  const handleErr = (msg: string, code?: number) => {
    if (code === 403) {
      ObjStore.setState(State.NeedPassword)
      if (retryPass) notify.error(msg)
      return
    }
    // base_path redirect on first "not found"
    const bp = me().base_path
    if (first_fetch && bp !== "/" && pathname().includes(bp) && msg.endsWith("object not found")) {
      first_fetch = false
      to(pathname().replace(bp, ""))
      return
    }
    if (code === undefined || code >= 0) {
      ObjStore.setErr(msg)
    }
  }

  // ── load folder contents ───────────────────────────────────────────

  const loadFolder = async (
    path: string,
    index?: number,
    size?: number,
    append = false,
    force?: boolean,
  ) => {
    const id = opId
    if (!size) size = pagination.size
    if (size !== undefined && pagination.type === "all") size = undefined

    if (!append && !shouldKeepState())
      ObjStore.setState(State.FetchingObjs)

    let resp
    try {
      resp = await apiList(path, index, size, force)
    } catch { return } // cancelled or network error
    if (id !== opId) return // stale

    if (resp.code !== 200) {
      handleErr(resp.message, resp.code)
      return
    }

    const data = resp.data
    setGlobalPage(index ?? 1)
    if (append) {
      appendObjs(data.content)
    } else {
      // Only overwrite the main view if the URL hasn't changed
      if (pathname() !== path) return
      ObjStore.setObjs(data.content ?? [])
      ObjStore.setTotal(data.total)
    }
    ObjStore.setReadme(data.readme || "")
    ObjStore.setHeader(data.header || "")
    ObjStore.setWrite(data.write)
    ObjStore.setWriteContentBypass(data.write_content_bypass)
    ObjStore.setProvider(data.provider)
    ObjStore.setDirectUploadTools(data.direct_upload_tools)
    if (!shouldKeepState()) ObjStore.setState(State.Folder)
  }

  // ── load single object (file or folder) ────────────────────────────

  const loadObj = async (path: string, index?: number) => {
    const id = opId
    if (!shouldKeepState()) ObjStore.setState(State.FetchingObj)

    let resp
    try {
      resp = await apiGet(path)
    } catch { return }
    if (id !== opId) return

    if (resp.code === 403) {
      ObjStore.setState(State.NeedPassword)
      if (retryPass) notify.error(resp.message)
      return
    }
    if (resp.code !== 200) {
      handleErr(resp.message, resp.code)
      return
    }

    const data = resp.data
    if (pathname() !== path) return // URL changed, discard

    ObjStore.setObj(data)
    ObjStore.setProvider(data.provider)

    if (data.is_dir) {
      await loadFolder(path, index, undefined, false)
    } else {
      ObjStore.setReadme(data.readme || "")
      ObjStore.setHeader(data.header || "")
      ObjStore.setRelated(data.related ?? [])
      ObjStore.setRawUrl(data.raw_url)
      if (!shouldKeepState()) ObjStore.setState(State.File)
    }
  }

  // ── public API ─────────────────────────────────────────────────────

  const handlePathChange = (path: string, index?: number, rp?: boolean, force?: boolean) => {
    // Cancel pending requests and invalidate all in-flight callbacks
    cancelGet?.()
    cancelList?.()
    ++opId
    retryPass = rp ?? false
    ObjStore.setErr("")
    loadObj(path, index)
  }

  const refresh = async (rp?: boolean, force?: boolean) => {
    const path = pathname()
    const scroll = window.scrollY
    if (pagination.type === "load_more" || pagination.type === "auto_load_more") {
      const page = globalPage
      resetGlobalPage()
      handlePathChange(path, globalPage, rp, force)
      while (globalPage < page) {
        await loadFolder(pathname(), globalPage + 1, undefined, true)
      }
    } else {
      handlePathChange(path, globalPage, rp, force)
    }
    if (pathname() === path) {
      window.scroll({ top: scroll, behavior: "smooth" })
    }
  }

  const loadMore = () => loadFolder(pathname(), globalPage + 1, undefined, true)

  return {
    handlePathChange,
    handleFolder: loadFolder,
    // setPathAs is a no-op — the old IsDirRecord pre-marking has been removed.
    // The server response is the only authority on file-vs-directory.
    setPathAs: (_p?: string, _d?: boolean, _push?: boolean) => {},
    refresh,
    loadMore,
    allLoaded: () => globalPage >= Math.ceil(objStore.total / pagination.size),
  }
}
