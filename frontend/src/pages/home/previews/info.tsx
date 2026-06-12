import type { ParentProps } from "solid-js"
import { Center, VStack } from "@hope-ui/solid"
import { useT } from "~/hooks"
import { objStore } from "~/store"
import { formatDate, objType } from "~/utils"

export const FileInfo = (props: ParentProps) => {
  const t = useT()
  return (
    <VStack w="$full" spacing="$2" alignItems="start">
      <Center w="$full" py="$4">
        <VStack w="$full" spacing="$2">
          {props.children}
        </VStack>
      </Center>
    </VStack>
  )
}
