import { PortfolioShow } from '../Work/Portfolio'
import type { AdminProjectItemProps } from '@/types/payloads'
export default function Show({ item }: AdminProjectItemProps) { return <PortfolioShow kind="Project" item={item} /> }
