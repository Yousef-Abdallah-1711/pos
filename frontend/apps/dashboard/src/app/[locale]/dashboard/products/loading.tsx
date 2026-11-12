export default function Loading() {
  return (
    <div className="space-y-4 animate-pulse">
      <div className="h-8 w-56 rounded bg-gray-200 dark:bg-gray-700" />
      <div className="h-4 w-80 rounded bg-gray-200 dark:bg-gray-700" />
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        <div className="h-36 rounded bg-gray-200 dark:bg-gray-700" />
        <div className="h-36 rounded bg-gray-200 dark:bg-gray-700" />
        <div className="h-36 rounded bg-gray-200 dark:bg-gray-700" />
      </div>
    </div>
  )
}
