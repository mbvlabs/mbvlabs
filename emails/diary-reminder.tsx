import { Button, Section, Text } from "react-email"
import { Layout } from "./components/layout"
import type { DiaryReminderProps } from "./generated/props"

export default function DiaryReminder({ Title, Period, DiaryURL }: DiaryReminderProps) {
  return (
    <Layout title={Title} preview="Open today's diary entry and capture your thoughts.">
      <Section style={{ padding: "32px 48px", color: "#414552" }}>
        <Text role="heading" aria-level={1} style={{ fontSize: 28, fontWeight: 600 }}>{Title}</Text>
        <Text>Open today's diary entry and write down your {Period} thoughts.</Text>
        <Button href={DiaryURL} style={{ backgroundColor: "#625afa", color: "white", borderRadius: 6, padding: "12px 24px" }}>
          Open Diary
        </Button>
        <Text>Or copy and paste this link into your browser:</Text>
        <Text style={{ color: "#625afa", wordBreak: "break-all" }}>{DiaryURL}</Text>
      </Section>
    </Layout>
  )
}
