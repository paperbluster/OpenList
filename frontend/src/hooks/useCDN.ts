import { joinBase } from "~/utils"
import packageJson from "../../package.json"

export const useCDN = () => {
  const static_path = joinBase("static")

  const npm = (name: string, version: string, path: string) => {
    // Always use local static files, no CDN dependency
    return `${static_path}/${name.split("/").pop()}/${path.split("/").pop()}`
  }

  const monacoPath = () => {
    return `${static_path}/monaco-editor/vs`
  }

  const katexCSSPath = () => {
    return `${static_path}/katex/katex.min.css`
  }

  const mermaidJSPath = () => {
    return `${static_path}/mermaid/mermaid.min.js`
  }

  const libHeifPath = () => {
    return `${static_path}/libheif`
  }

  const libAssPath = () => {
    return `${static_path}/libass-wasm`
  }

  const fontsPath = () => {
    return `${static_path}/fonts`
  }

  return {
    npm,
    monacoPath,
    katexCSSPath,
    mermaidJSPath,
    libHeifPath,
    libAssPath,
    fontsPath,
  }
}
