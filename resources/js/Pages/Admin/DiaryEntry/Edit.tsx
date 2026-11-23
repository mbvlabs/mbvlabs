import Form, { type DiaryEntry } from './Form'
export default function Edit({ item, focus }: { item: DiaryEntry; focus?: 'morning' | 'evening' | '' }) { return <Form item={item} focus={focus} /> }
