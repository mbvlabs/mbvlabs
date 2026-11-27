import Form, { type DiaryEntry } from './Form'
export default function Create({ today, item }: { today: string; item?: DiaryEntry }) { return <Form today={today} item={item} create /> }
