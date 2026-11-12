// Minimal typings for the parts of `qrcode` used by the dashboard (the package ships no types and
// @types/qrcode is not a dependency).
declare module 'qrcode' {
  interface QRCodeRenderOptions {
    errorCorrectionLevel?: 'low' | 'medium' | 'quartile' | 'high' | 'L' | 'M' | 'Q' | 'H'
    margin?: number
    width?: number
    scale?: number
    color?: { dark?: string; light?: string }
  }

  export function toDataURL(text: string, options?: QRCodeRenderOptions): Promise<string>
  export function toString(text: string, options?: QRCodeRenderOptions & { type?: 'svg' | 'utf8' | 'terminal' }): Promise<string>

  const QRCode: {
    toDataURL: typeof toDataURL
    toString: typeof toString
  }
  export default QRCode
}
