'use client'

export default function TestUrlPage() {
  if (typeof window === 'undefined') {
    return <div>Loading...</div>
  }

  return (
    <div className="min-h-screen bg-blue-900 text-white p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">✅ TEST PAGE LOADED!</h1>
        <p className="text-green-400 text-lg mb-6">This page is NOT wrapped by DashboardAuthProvider (no auth protection)</p>

        <div className="bg-blue-800 p-6 rounded mb-6">
          <h2 className="text-xl font-bold mb-3">Current URL:</h2>
          <p className="text-yellow-300 font-mono break-all text-sm">{window.location.href}</p>
        </div>

        <div className="bg-blue-800 p-6 rounded mb-6">
          <h2 className="text-xl font-bold mb-3">Search Params:</h2>
          <p className="text-yellow-300 font-mono break-all text-sm">{window.location.search || '(none)'}</p>
        </div>

        <div className="bg-blue-800 p-6 rounded mb-6">
          <h2 className="text-xl font-bold mb-3">Query Parameters:</h2>
          <div className="space-y-2">
            {Array.from(new URLSearchParams(window.location.search)).map(([key, value]) => (
              <div key={key} className="text-sm border-l-2 border-cyan-400 pl-3">
                <span className="text-cyan-400 font-mono font-bold">{key}</span>
                <span className="text-gray-400">: </span>
                <span className="text-yellow-300 font-mono break-all">{value.substring(0, 100)}{value.length > 100 ? '...' : ''}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-green-900 p-6 rounded mb-6 border-2 border-green-500">
          <h2 className="text-2xl font-bold mb-3 text-green-400">🔍 Debug Check Result:</h2>
          <p className="text-lg mb-3">
            <span className="text-cyan-400">debug param: </span>
            <span className="font-mono text-yellow-300">{new URLSearchParams(window.location.search).get('debug') || 'NOT FOUND'}</span>
          </p>
          <p className="text-3xl font-bold">
            <span className="text-cyan-400">Is Debug Mode: </span>
            <span className={new URLSearchParams(window.location.search).get('debug') === 'true' ? 'text-green-400' : 'text-red-400'}>
              {new URLSearchParams(window.location.search).get('debug') === 'true' ? '✅ YES' : '❌ NO'}
            </span>
          </p>
        </div>

        <p className="text-gray-400 text-sm p-4 bg-blue-950 rounded">
          ℹ️ This page has NO DashboardAuthProvider - it only uses the standard AuthProvider from the locale layout.<br/>
          If you can see this page, it means the URL parameters are reaching your app.<br/>
          If you see this page REDIRECT to login, it means there&apos;s a redirect happening in the locale layout.
        </p>
      </div>
    </div>
  )
}
