import { PortfolioIndex } from '../Work/Portfolio'
import type { AdminProjectIndexProps } from '@/types/payloads'
export default function Index({ items }: AdminProjectIndexProps) { return <PortfolioIndex kind="Project" items={items} /> }
