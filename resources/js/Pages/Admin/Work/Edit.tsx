import { PortfolioForm } from './Portfolio'
import type { AdminWorkItemProps } from '@/types/payloads'
export default function Edit({ item }: AdminWorkItemProps) { return <PortfolioForm kind="Work" item={item} /> }
