import { PortfolioShow } from './Portfolio'
import type { AdminWorkItemProps } from '@/types/payloads'
export default function Show({ item }: AdminWorkItemProps) { return <PortfolioShow kind="Work" item={item} /> }
