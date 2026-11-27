import Form from './Form'
import type { AdminDiaryCreateProps, AdminDiaryData } from '@/types/payloads'
export default function Create({ today, item }: AdminDiaryCreateProps & { item?: AdminDiaryData }) { return <Form today={today} item={item} create /> }
