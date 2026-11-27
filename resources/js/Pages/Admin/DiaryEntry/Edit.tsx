import Form from './Form'
import type { AdminDiaryEditProps } from '@/types/payloads'
export default function Edit({ item, focus }: AdminDiaryEditProps) { return <Form item={item} focus={focus} /> }
