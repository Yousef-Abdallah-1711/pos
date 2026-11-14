import { RootDocument, rootMetadata } from '@/components/RootDocument'
export const metadata = rootMetadata
export default function EntryLayout({ children }: { children: React.ReactNode }) {
  return <RootDocument locale="en">{children}</RootDocument>
}
