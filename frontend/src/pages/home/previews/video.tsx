import { Box, Button, HStack, Switch, VStack, Text } from "@hope-ui/solid"
import { createEffect, createMemo, createSignal, onCleanup, onMount, Show } from "solid-js"
import { useLink, useRouter } from "~/hooks"
import { getSettingBool, objStore } from "~/store"
import { ObjType } from "~/types"
import { pathDir, pathJoin } from "~/utils"
import Hls from "hls.js"
import mpegts from "mpegts.js"
import { useNavigate } from "@solidjs/router"

/**
 * Simple video player using native <video> + mpegts.js + hls.js.
 * Supports: MP4, MKV, TS, FLV, M3U8, WebM.
 * Auto-loads subtitles from same directory (same basename, priority: ass > ssa > srt).
 */
const VideoPreview = () => {
  const { pathname } = useRouter()
  const navigate = useNavigate()
  const { rawLink } = useLink()
  let videoEl!: HTMLVideoElement

  const videos = createMemo(() =>
    objStore.objs.filter((obj) => obj.type === ObjType.VIDEO),
  )
  const currentIndex = createMemo(() =>
    videos().findIndex((v) => v.name === objStore.obj.name),
  )

  const hasPrev = () => currentIndex() > 0
  const hasNext = () => currentIndex() < videos().length - 1

  const goTo = (idx: number) => {
    const v = videos()[idx]
    if (v) navigate(pathJoin(pathDir(pathname()), v.name))
  }

  const storedAutoNext = localStorage.getItem("video_auto_next")
  const [autoNext, setAutoNext] = createSignal(
    storedAutoNext === null ? true : storedAutoNext === "true",
  )

  // ── subtitle auto-detection ───────────────────────────────────────
  // Search related files + folder siblings for subtitles matching the video's
  // basename. Priority: .ass > .ssa > .srt
  const subtitleTrack = createMemo(() => {
    const videoName = objStore.obj.name
    const base = videoName.substring(0, videoName.lastIndexOf(".")) || videoName
    const baseLower = base.toLowerCase()

    const candidates = [...(objStore.related || []), ...objStore.objs].filter(
      (o) => o.name !== videoName,
    )

    const priority = [".ass", ".ssa", ".srt"]
    for (const ext of priority) {
      const match = candidates.find(
        (o) => o.name.toLowerCase() === baseLower + ext,
      )
      if (match) return { url: rawLink(match, true), label: match.name, ext }
    }
    return null
  })

  // ── video format handling ─────────────────────────────────────────

  let hls: Hls | undefined
  let flvPlayer: mpegts.Player | undefined

  const ext = () => {
    const name = objStore.obj.name.toLowerCase()
    if (name.endsWith(".flv")) return "flv"
    if (name.endsWith(".ts") || name.endsWith(".mts") || name.endsWith(".m2ts")) return "m2ts"
    if (name.endsWith(".m3u8")) return "m3u8"
    return "native"
  }

  const setupVideo = () => {
    const src = objStore.raw_url
    if (!src || !videoEl) return

    flvPlayer?.destroy()
    flvPlayer = undefined
    hls?.destroy()
    hls = undefined

    switch (ext()) {
      case "m3u8":
        if (videoEl.canPlayType("application/vnd.apple.mpegurl")) {
          videoEl.src = src
        } else {
          hls = new Hls()
          hls.loadSource(src)
          hls.attachMedia(videoEl)
        }
        break
      case "flv":
      case "m2ts":
        flvPlayer = mpegts.createPlayer(
          { type: ext(), url: src },
          { referrerPolicy: "same-origin" },
        )
        flvPlayer.attachMediaElement(videoEl)
        flvPlayer.load()
        break
      default:
        videoEl.src = src
    }
  }

  onMount(() => setupVideo())

  createEffect(() => {
    if (objStore.raw_url) setupVideo()
  })

  onCleanup(() => {
    flvPlayer?.destroy()
    hls?.destroy()
  })

  return (
    <VStack w="$full" spacing="$2">
      <Box w="$full">
        <video
          ref={videoEl}
          controls
          autoplay={getSettingBool("video_autoplay")}
          playsInline
          crossOrigin="anonymous"
          onEnded={() => {
            if (autoNext() && hasNext()) goTo(currentIndex() + 1)
          }}
          style={{ width: "100%", maxHeight: "70vh", background: "#000" }}
        >
          <Show when={subtitleTrack()}>
            {(t) => (
              <track
                kind="subtitles"
                src={t().url}
                label={`${t().label} (字幕)`}
                default
              />
            )}
          </Show>
        </video>
      </Box>

      <Show when={videos().length > 1}>
        <HStack spacing="$2" w="$full" justifyContent="center" flexWrap="wrap">
          <Button size="sm" onClick={() => goTo(currentIndex() - 1)} disabled={!hasPrev()}>
            ← 上一个
          </Button>
          <Text fontSize="$sm" color="$neutral11">
            {currentIndex() + 1} / {videos().length}
          </Text>
          <Button size="sm" onClick={() => goTo(currentIndex() + 1)} disabled={!hasNext()}>
            下一个 →
          </Button>
          <Switch
            checked={autoNext()}
            onChange={(e: { currentTarget: HTMLInputElement }) => {
              const v = e.currentTarget.checked
              setAutoNext(v)
              localStorage.setItem("video_auto_next", v.toString())
            }}
          >
            自动连播
          </Switch>
        </HStack>
      </Show>
    </VStack>
  )
}

export default VideoPreview
