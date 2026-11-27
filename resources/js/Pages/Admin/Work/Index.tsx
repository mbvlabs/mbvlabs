import { PortfolioIndex, type Work } from './Portfolio'
export default function Index({ items }: { items: Work[] }) { return <PortfolioIndex kind="Work" items={items} /> }
