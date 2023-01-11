
// ArrayBuffer to URLBase64
export function bufferEncode(value: number) {
    return btoa(
      String.fromCharCode.apply(
        null,
        new Uint8Array(value) as unknown as number[]
      )
    )
      .replace(/\+/g, "-")
      .replace(/\//g, "_")
      .replace(/=/g, "")
  }
  
  // Base64 to ArrayBuffer
  export function bufferDecode(value: string) {
    return Uint8Array.from(atob(value), (c) => c.charCodeAt(0))
  }