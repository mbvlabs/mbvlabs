import { PortfolioShow, type Work } from './Portfolio'
export default function Show({ item }: { item: Work }) { return <PortfolioShow kind="Work" item={item} /> }
