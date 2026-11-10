import { ImageResponse } from 'next/og'

export const alt = 'POOS — one view of your restaurant’s work'
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
            <path d="M11 39V9h14a10 10 0 0 1 0 20H11m14 0h12" stroke="#6D28D9" strokeWidth="7" strokeLinecap="round" strokeLinejoin="round" />
            <circle cx="37" cy="29" r="3.5" fill="#8B5CF6" />
          </svg>
          POOS
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', maxWidth: 960, borderTop: '2px solid #6D28D9', paddingTop: 32 }}>
          <div style={{ fontSize: 64, lineHeight: 1.25, fontWeight: 700 }}>One view of your restaurant’s work.</div>
          <div style={{ marginTop: 20, color: '#5C566F', fontSize: 30 }}>One platform for front- and back-of-house work.</div>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', borderTop: '1px solid #D9D1EB', paddingTop: 20, color: '#6D28D9', fontSize: 22 }}>
          <span>RESTAURANT MANAGEMENT</span>
          <span>POOS / RESTAURANT OPERATIONS</span>
        </div>
      </div>
    ),
    size,
  )
}
