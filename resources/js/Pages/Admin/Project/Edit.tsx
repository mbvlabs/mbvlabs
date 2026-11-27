import { PortfolioForm } from '../Work/Portfolio'
import type { AdminProjectItemProps } from '@/types/payloads'
export default function Edit({ item }: AdminProjectItemProps) { return <PortfolioForm kind="Project" item={item} /> }
