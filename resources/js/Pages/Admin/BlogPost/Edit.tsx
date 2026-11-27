import Form from './Form'
import type { AdminBlogItemProps } from '@/types/payloads'
export default function Edit({ item }: AdminBlogItemProps) { return <Form item={item} /> }
