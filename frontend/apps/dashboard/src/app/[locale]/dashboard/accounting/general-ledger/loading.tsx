export default function Loading() {
  return (
    <div className="animate-pulse space-y-4 p-4 md:p-6">
      <div className="h-7 w-56 rounded bg-gray-200 dark:bg-gray-700" />
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        <div className="h-28 rounded-xl bg-gray-200 dark:bg-gray-700" />
        <div className="h-28 rounded-xl bg-gray-200 dark:bg-gray-700" />
        <div className="h-28 rounded-xl bg-gray-200 dark:bg-gray-700" />
      </div>
      <div className="h-64 rounded-xl bg-gray-200 dark:bg-gray-700" />
    </div>
  )
}
