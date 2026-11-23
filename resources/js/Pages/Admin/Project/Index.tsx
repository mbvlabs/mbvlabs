import { PortfolioIndex, type Project } from '../Work/Portfolio'
export default function Index({ items }: { items: Project[] }) { return <PortfolioIndex kind="Project" items={items} /> }
