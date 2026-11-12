export default function SuperAdminLoading() {
  return (
    <main aria-busy="true" className="mx-auto min-h-[60vh] w-full max-w-7xl space-y-8 p-4 sm:p-6">
      <p role="status" className="sr-only">Loading...</p>
      <div aria-hidden="true" className="space-y-4">
        <div className="h-9 w-64 max-w-full animate-pulse rounded-lg bg-gray-200 motion-reduce:animate-none dark:bg-gray-700" />
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {Array.from({ length: 4 }, (_, index) => (
            <div key={index} className="h-28 animate-pulse rounded-xl bg-gray-100 motion-reduce:animate-none dark:bg-gray-800" />
          ))}
        </div>
        <div className="h-80 animate-pulse rounded-xl bg-gray-100 motion-reduce:animate-none dark:bg-gray-800" />
      </div>
    </main>
  );
}
