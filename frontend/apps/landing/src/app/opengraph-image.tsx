import { ImageResponse } from 'next/og'

export const alt = 'SERVIO — one view of your restaurant’s work'
export const size = { width: 1200, height: 630 }
export const contentType = 'image/png'

export default function OpenGraphImage() {
  return new ImageResponse(
    (
      <div style={{
        width: '100%',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'space-between',
        padding: '64px 76px',
        backgroundColor: '#F7F5FF',
        color: '#171327',
        fontFamily: 'Arial, sans-serif',
        border: '18px solid #6D28D9',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 18, color: '#6D28D9', fontSize: 38, fontWeight: 800 }}>
          <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <rect x="8" y="17" width="7" height="23" fill="#6D28D9" />
            <rect x="21" y="8" width="7" height="32" fill="#6D28D9" />
            <rect x="34" y="22" width="7" height="18" fill="#6D28D9" />
          </svg>
          SERVIO
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', maxWidth: 960, borderTop: '2px solid #6D28D9', paddingTop: 32 }}>
          <div style={{ fontSize: 64, lineHeight: 1.25, fontWeight: 700 }}>Every order, connected.</div>
          <div style={{ marginTop: 20, color: '#5C566F', fontSize: 30 }}>From the counter to the kitchen.</div>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', borderTop: '1px solid #D9D1EB', paddingTop: 20, color: '#6D28D9', fontSize: 22 }}>
          <span>RESTAURANT MANAGEMENT</span>
          <span>SERVIO / RESTAURANT OPERATIONS</span>
        </div>
      </div>
    ),
    size,
  )
}
