import { PortfolioIndex } from './Portfolio'
import type { AdminWorkIndexProps } from '@/types/payloads'
export default function Index({ items }: AdminWorkIndexProps) { return <PortfolioIndex kind="Work" items={items} /> }
