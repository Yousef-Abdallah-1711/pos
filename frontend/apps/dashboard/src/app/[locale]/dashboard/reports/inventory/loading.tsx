export default function LoadingInventoryReportPage() {
  return (
    <div className="space-y-4 animate-pulse">
      <div className="h-8 w-64 rounded bg-gray-200 dark:bg-gray-700" />
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="h-24 rounded-xl bg-gray-200 dark:bg-gray-700" />
        <div className="h-24 rounded-xl bg-gray-200 dark:bg-gray-700" />
        <div className="h-24 rounded-xl bg-gray-200 dark:bg-gray-700" />
        <div className="h-24 rounded-xl bg-gray-200 dark:bg-gray-700" />
      </div>
      <div className="h-96 rounded-xl bg-gray-200 dark:bg-gray-700" />
    </div>
  )
}
