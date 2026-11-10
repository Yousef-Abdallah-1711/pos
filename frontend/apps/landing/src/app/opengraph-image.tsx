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
          <div style={{ width: 48, height: 48, borderRadius: 13, background: '#6D28D9', color: '#fff', display: 'flex', flexDirection: 'column', justifyContent: 'center', padding: '0 12px', gap: 5 }}>
            <div style={{ height: 2, width: 24, background: 'white' }} />
            <div style={{ height: 2, width: 16, background: 'white' }} />
            <div style={{ height: 2, width: 24, background: 'white' }} />
          </div>
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
