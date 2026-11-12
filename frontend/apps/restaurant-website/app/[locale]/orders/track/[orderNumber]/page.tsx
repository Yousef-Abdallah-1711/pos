import TrackPageContent from '../track-page-content'

interface PageProps {
  params: Promise<{ locale: 'en' | 'ar'; orderNumber: string }>
}

export default async function TrackOrderNumberPage({ params }: PageProps) {
  const { orderNumber } = await params
  const initialOrderNumber = decodeURIComponent(orderNumber)

  return (
    <main className="min-h-screen w-full">
      <TrackPageContent initialOrderNumber={initialOrderNumber} />
    </main>
  )
}
