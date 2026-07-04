import { Box, VStack, Text } from "@hope-ui/solid"
import { createMemo, Show } from "solid-js"
import { objStore } from "~/store"
import { ObjType } from "~/types"

/** Simple HTML5 audio player — replaces the heavy APlayer dependency. */
const AudioPreview = () => {
  const audios = createMemo(() =>
    objStore.objs.filter((obj) => obj.type === ObjType.AUDIO),
  )

  return (
    <VStack w="$full" spacing="$2" alignItems="center">
      <Box
        w="$full"
        maxW="600px"
        css={{
          "& audio": { width: "100%" },
        }}
      >
        <audio
          src={objStore.raw_url}
          controls
          autoplay={false}
          preload="auto"
          style={{ width: "100%" }}
        />
      </Box>
      <Show when={audios().length > 1}>
        <Text fontSize="$sm" color="$neutral11">
          当前文件夹共有 {audios().length} 个音频文件
        </Text>
      </Show>
    </VStack>
  )
}

export default AudioPreview
