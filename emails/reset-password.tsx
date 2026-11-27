import { Button, Heading, Section, Text } from "react-email"
import { Layout } from "./components/layout"
import type { ResetPasswordProps } from "./generated/props"
import { theme } from "./theme"

const headingStyle = {
  color: theme.foreground,
  fontFamily: theme.fontSans,
  fontSize: 30,
  fontWeight: 700,
  lineHeight: "38px",
  margin: 0,
} as const

const paragraphStyle = {
  color: theme.foreground,
  fontFamily: theme.fontSans,
  fontSize: 16,
  lineHeight: "24px",
  margin: 0,
} as const

const noteStyle = {
  color: theme.mutedForeground,
  fontFamily: theme.fontSans,
  fontSize: 14,
  lineHeight: "24px",
  margin: 0,
} as const

export default function ResetPassword({ ResetURL }: ResetPasswordProps) {
  return (
    <Layout
      title="Reset Your Password"
      preview="Reset your password using the link provided."
    >
      <Section
        className="email-px"
        style={{ padding: "48px 48px 20px 48px" }}
      >
        <Heading as="h1" className="email-fg" style={headingStyle}>
          Reset Your Password
        </Heading>
      </Section>
      <Section
        className="email-px"
        style={{ padding: "0 48px 24px 48px" }}
      >
        <Text className="email-fg" style={paragraphStyle}>
          Hi,
        </Text>
        <Text
          className="email-fg"
          style={{ ...paragraphStyle, margin: "16px 0 0 0" }}
        >
          We received a request to reset your password. Click the button below to reset your password:
        </Text>
      </Section>
      <Section
        className="email-px"
        style={{ padding: "0 48px 24px 48px" }}
      >
        <Button
          href={ResetURL}
          style={{
            backgroundColor: theme.primary,
            borderRadius: 6,
            color: theme.primaryForeground,
            display: "inline-block",
            fontFamily: theme.fontSans,
            fontSize: 16,
            fontWeight: 700,
            lineHeight: "24px",
            padding: "12px 24px",
            textAlign: "center",
            textDecoration: "none",
          }}
        >
          Reset My Password
        </Button>
      </Section>
      <Section
        className="email-px"
        style={{ padding: "0 48px 24px 48px" }}
      >
        <Text className="email-muted" style={noteStyle}>
          Or copy and paste this link into your browser:
        </Text>
        <Text
          className="email-primary"
          style={{
            ...noteStyle,
            color: theme.primary,
            margin: "8px 0 0 0",
            wordBreak: "break-all",
          }}
        >
          {ResetURL}
        </Text>
      </Section>
      <Section
        className="email-px"
        style={{ padding: "0 48px 16px 48px" }}
      >
        <Text className="email-muted" style={noteStyle}>
          This link will expire in 1 hour for security reasons.
        </Text>
        <Text
          className="email-muted"
          style={{ ...noteStyle, margin: "16px 0 0 0" }}
        >
          If you didn't request a password reset, please ignore this email. Your password will remain unchanged.
        </Text>
        <Text className="email-muted" style={{ ...noteStyle, margin: "16px 0 0 0" }}>
          Best regards,<br />The Andurel Team
        </Text>
      </Section>
    </Layout>
  )
}
